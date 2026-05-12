#!/usr/bin/env node
/**
 * Parse TypeScript source files via Compiler API and output JSON structure.
 *
 * Usage:
 *   node scripts/ts_ast.js file.ts
 *   node scripts/ts_ast.js --batch file1.ts file2.ts ...
 */
"use strict";

const ts = require("typescript");
const crypto = require("crypto");
const fs = require("fs");
const path = require("path");

function isUpperCase(name) {
  return /^[A-Z][A-Z0-9_]*$/.test(name);
}

function stmtDepth(node, current) {
  let max = current;
  ts.forEachChild(node, function visit(child) {
    let d = current;
    if (ts.isIfStatement(child)) {
      d = current + 1;
      if (d > max) max = d;
      if (child.thenStatement) {
        let inner = stmtDepth(child.thenStatement, d);
        if (inner > max) max = inner;
      }
      if (child.elseStatement) {
        if (ts.isIfStatement(child.elseStatement)) {
          // else if: same depth
          let inner = stmtDepth(child.elseStatement, current);
          if (inner > max) max = inner;
        } else {
          let inner = stmtDepth(child.elseStatement, d);
          if (inner > max) max = inner;
        }
      }
      return;
    }
    if (ts.isSwitchStatement(child)) {
      d = current + 1;
      if (d > max) max = d;
      for (const clause of child.caseBlock.clauses) {
        let inner = stmtDepth(clause, d);
        if (inner > max) max = inner;
      }
      return;
    }
    if (
      ts.isForStatement(child) ||
      ts.isForInStatement(child) ||
      ts.isForOfStatement(child) ||
      ts.isWhileStatement(child) ||
      ts.isDoStatement(child)
    ) {
      d = current + 1;
      if (d > max) max = d;
      let inner = stmtDepth(child.statement, d);
      if (inner > max) max = inner;
      return;
    }
    // recurse into blocks, arrow bodies, etc.
    let inner = stmtDepth(child, current);
    if (inner > max) max = inner;
  });
  return max;
}

function countLines(node, sourceFile) {
  const startLine = sourceFile.getLineAndCharacterOfPosition(node.getStart(sourceFile)).line;
  const endLine = sourceFile.getLineAndCharacterOfPosition(node.getEnd()).line;
  return endLine - startLine + 1;
}

function pureLines(bodyNode, sourceFile) {
  // Total lines in body minus lines belonging to inner control statements
  const bodyStart = sourceFile.getLineAndCharacterOfPosition(bodyNode.getStart(sourceFile)).line;
  const bodyEnd = sourceFile.getLineAndCharacterOfPosition(bodyNode.getEnd()).line;
  const total = bodyEnd - bodyStart + 1;

  let innerLines = 0;
  ts.forEachChild(bodyNode, function (child) {
    if (
      ts.isIfStatement(child) ||
      ts.isSwitchStatement(child) ||
      ts.isForStatement(child) ||
      ts.isForInStatement(child) ||
      ts.isForOfStatement(child) ||
      ts.isWhileStatement(child) ||
      ts.isDoStatement(child)
    ) {
      const s = sourceFile.getLineAndCharacterOfPosition(child.getStart(sourceFile)).line;
      const e = sourceFile.getLineAndCharacterOfPosition(child.getEnd()).line;
      innerLines += e - s + 1;
    }
  });

  return total - innerLines;
}

function getCallName(expr) {
  if (ts.isIdentifier(expr)) {
    return expr.text;
  }
  if (ts.isPropertyAccessExpression(expr)) {
    const obj = getCallName(expr.expression);
    if (obj) return obj + "." + expr.name.text;
    return expr.name.text;
  }
  return null;
}

function collectCalls(node, calls) {
  if (ts.isCallExpression(node)) {
    const name = getCallName(node.expression);
    if (name) calls.push(name);
  }
  ts.forEachChild(node, function (child) {
    collectCalls(child, calls);
  });
}

function parseFile(filePath) {
  const result = {
    path: filePath,
    functions: [],
    classes: [],
    interfaces: [],
    type_aliases: [],
    methods: [],
    has_constructor: false,
    vars: [],
    lines: 0,
    max_depth: 0,
    control: "sequence",
    has_loop_at_depth1: false,
    has_switch_at_depth1: false,
    func_lines: {},
    q4_results: [],
    calls: [],
    imports: [],
    body_hash: "",
  };

  let src;
  try {
    src = fs.readFileSync(filePath, "utf8");
  } catch (e) {
    result.error = e.message;
    return result;
  }

  if (src.trim() === "") {
    return result;
  }

  const sourceFile = ts.createSourceFile(
    path.basename(filePath),
    src,
    ts.ScriptTarget.Latest,
    true,
    filePath.endsWith(".tsx") ? ts.ScriptKind.TSX : ts.ScriptKind.TS
  );

  const lineCount = sourceFile.getLineAndCharacterOfPosition(sourceFile.getEnd()).line;
  const lastChar = src.length > 0 ? src[src.length - 1] : "";
  result.lines = lastChar === "\n" || lastChar === "\r" ? lineCount : lineCount + 1;
  if (src.trim() === "") {
    result.lines = 0;
  }

  let firstFuncBody = null;

  for (const stmt of sourceFile.statements) {
    // Functions
    if (ts.isFunctionDeclaration(stmt) && stmt.name) {
      const name = stmt.name.text;
      result.functions.push(name);
      const ln = countLines(stmt, sourceFile);
      result.func_lines[name] = ln;

      if (!firstFuncBody && stmt.body) {
        firstFuncBody = stmt;
      }

      // depth inside this function
      if (stmt.body) {
        const d = stmtDepth(stmt.body, 0);
        if (d > result.max_depth) result.max_depth = d;

        // depth-1 control detection
        for (const s of stmt.body.statements) {
          checkDepth1Control(s, result);
          checkQ4(s, name, sourceFile, result);
        }
      }

      // calls
      if (stmt.body) {
        collectCalls(stmt.body, result.calls);
      }
    }

    // Classes
    if (ts.isClassDeclaration(stmt) && stmt.name) {
      const className = stmt.name.text;
      result.classes.push(className);

      for (const member of stmt.members) {
        if (ts.isConstructorDeclaration(member)) {
          result.has_constructor = true;
          const ctorName = className + ".constructor";
          const ln = countLines(member, sourceFile);
          result.func_lines[ctorName] = ln;

          if (member.body) {
            const d = stmtDepth(member.body, 0);
            if (d > result.max_depth) result.max_depth = d;
            collectCalls(member.body, result.calls);
          }
        }

        if (ts.isMethodDeclaration(member) && member.name) {
          const methodName = className + "." + member.name.getText(sourceFile);
          result.methods.push(methodName);
          const ln = countLines(member, sourceFile);
          result.func_lines[methodName] = ln;

          if (!firstFuncBody && member.body) {
            firstFuncBody = member;
          }

          if (member.body) {
            const d = stmtDepth(member.body, 0);
            if (d > result.max_depth) result.max_depth = d;

            for (const s of member.body.statements) {
              checkDepth1Control(s, result);
              checkQ4(s, methodName, sourceFile, result);
            }

            collectCalls(member.body, result.calls);
          }
        }
      }
    }

    // Interfaces
    if (ts.isInterfaceDeclaration(stmt) && stmt.name) {
      result.interfaces.push(stmt.name.text);
    }

    // Type aliases
    if (ts.isTypeAliasDeclaration(stmt) && stmt.name) {
      result.type_aliases.push(stmt.name.text);
    }

    // Variable statements — UPPER_CASE const only
    if (ts.isVariableStatement(stmt)) {
      const isConst =
        stmt.declarationList &&
        (stmt.declarationList.flags & ts.NodeFlags.Const) !== 0;
      if (isConst) {
        for (const decl of stmt.declarationList.declarations) {
          if (ts.isIdentifier(decl.name) && isUpperCase(decl.name.text)) {
            result.vars.push(decl.name.text);
          }
        }
      }
    }

    // Imports
    if (ts.isImportDeclaration(stmt) && stmt.moduleSpecifier) {
      const mod = stmt.moduleSpecifier.getText(sourceFile).replace(/['"]/g, "");
      const names = [];
      if (stmt.importClause) {
        if (stmt.importClause.name) {
          names.push(stmt.importClause.name.text);
        }
        if (stmt.importClause.namedBindings) {
          if (ts.isNamedImports(stmt.importClause.namedBindings)) {
            for (const el of stmt.importClause.namedBindings.elements) {
              names.push(el.name.text);
            }
          } else if (ts.isNamespaceImport(stmt.importClause.namedBindings)) {
            names.push(stmt.importClause.namedBindings.name.text);
          }
        }
      }
      result.imports.push({ module: mod, names: names });
    }
  }

  // control field
  if (result.has_switch_at_depth1) {
    result.control = "selection";
  } else if (result.has_loop_at_depth1) {
    result.control = "iteration";
  } else {
    result.control = "sequence";
  }

  // body hash: first func/method source SHA-256 first 8 chars
  if (firstFuncBody) {
    const bodyText = firstFuncBody.getText(sourceFile);
    const hash = crypto.createHash("sha256").update(bodyText).digest("hex");
    result.body_hash = hash.substring(0, 8);
  }

  // deduplicate calls
  result.calls = [...new Set(result.calls)];

  return result;
}

function checkDepth1Control(stmt, result) {
  if (ts.isSwitchStatement(stmt)) {
    result.has_switch_at_depth1 = true;
  }
  if (
    ts.isForStatement(stmt) ||
    ts.isForInStatement(stmt) ||
    ts.isForOfStatement(stmt) ||
    ts.isWhileStatement(stmt) ||
    ts.isDoStatement(stmt)
  ) {
    result.has_loop_at_depth1 = true;
  }
}

function checkQ4(stmt, funcName, sourceFile, result) {
  let stmtType = null;
  let bodyNode = null;

  if (ts.isIfStatement(stmt)) {
    stmtType = "if";
    bodyNode = stmt.thenStatement;
  } else if (ts.isSwitchStatement(stmt)) {
    stmtType = "switch";
    bodyNode = stmt.caseBlock;
  } else if (ts.isForStatement(stmt) || ts.isForInStatement(stmt) || ts.isForOfStatement(stmt)) {
    stmtType = "for";
    bodyNode = stmt.statement;
  } else if (ts.isWhileStatement(stmt)) {
    stmtType = "while";
    bodyNode = stmt.statement;
  } else if (ts.isDoStatement(stmt)) {
    stmtType = "do";
    bodyNode = stmt.statement;
  }

  if (stmtType && bodyNode) {
    const pure = pureLines(bodyNode, sourceFile);
    if (pure > 10) {
      const line = sourceFile.getLineAndCharacterOfPosition(stmt.getStart(sourceFile)).line + 1;
      result.q4_results.push({
        func_name: funcName,
        stmt_type: stmtType,
        pure_lines: pure,
        line: line,
      });
    }
  }
}

// Main
const args = process.argv.slice(2);
if (args.length === 0) {
  process.stderr.write("Usage: ts_ast.js [--batch] file.ts ...\n");
  process.exit(1);
}

let files;
let batch = false;
if (args[0] === "--batch") {
  batch = true;
  files = args.slice(1);
} else {
  files = [args[0]];
}

const results = files.map(parseFile);

if (batch) {
  process.stdout.write(JSON.stringify(results));
} else {
  process.stdout.write(JSON.stringify(results[0]));
}
