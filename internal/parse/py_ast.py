#!/usr/bin/env python3
"""Parse Python source files via AST and output JSON structure.

Usage:
    python3 scripts/py_ast.py file.py
    python3 scripts/py_ast.py --batch file1.py file2.py ...
"""
import ast
import hashlib
import json
import sys


def _depth_of_stmts(stmts, current):
    """Compute max depth across a list of statements."""
    max_d = current
    for stmt in stmts:
        d = _stmt_depth(stmt, current)
        if d > max_d:
            max_d = d
    return max_d


def _stmt_depth(node, current):
    """Compute max depth starting from a single statement node.

    Counted: if, for, while, match.
    elif does NOT increase depth.
    with/try/except/comprehension are NOT counted.
    """
    if isinstance(node, (ast.For, ast.While, ast.AsyncFor)):
        return _depth_of_stmts(node.body, current + 1)
    if isinstance(node, ast.If):
        max_d = _depth_of_stmts(node.body, current + 1)
        # elif/else chain
        orelse = node.orelse
        while orelse:
            if len(orelse) == 1 and isinstance(orelse[0], ast.If):
                # elif: same depth as parent if
                elif_node = orelse[0]
                d = _depth_of_stmts(elif_node.body, current + 1)
                if d > max_d:
                    max_d = d
                orelse = elif_node.orelse
            else:
                # else block
                d = _depth_of_stmts(orelse, current + 1)
                if d > max_d:
                    max_d = d
                break
        return max_d
    if hasattr(ast, "Match") and isinstance(node, ast.Match):
        max_d = current + 1
        for case_node in node.cases:
            d = _depth_of_stmts(case_node.body, current + 1)
            if d > max_d:
                max_d = d
        return max_d
    # Non-control: scan children for nested control stmts
    max_d = current
    for child in ast.iter_child_nodes(node):
        if isinstance(child, ast.stmt):
            d = _stmt_depth(child, current)
            if d > max_d:
                max_d = d
    return max_d


def _func_max_depth(func_node):
    """Compute max nesting depth of a function body."""
    return _depth_of_stmts(func_node.body, 0)


def _has_control_at_depth1(func_node, kind):
    """Check if func body (depth 1) has a specific control type.

    kind: 'loop' (for/while), 'match' (match/case).
    """
    for stmt in func_node.body:
        if kind == "loop" and isinstance(stmt, (ast.For, ast.While, ast.AsyncFor)):
            return True
        if kind == "match" and hasattr(ast, "Match") and isinstance(stmt, ast.Match):
            return True
    return False


def _detect_control(func_node):
    """Detect control type from depth-1 statements."""
    has_match = _has_control_at_depth1(func_node, "match")
    has_loop = _has_control_at_depth1(func_node, "loop")
    if has_match:
        return "selection"
    if has_loop:
        return "iteration"
    return "sequence"


def _count_lines(node):
    """Count the number of source lines a node spans."""
    if not hasattr(node, "end_lineno") or node.end_lineno is None:
        return 0
    return node.end_lineno - node.lineno + 1


def _is_control_stmt(node):
    """Check if node is a control statement (if/for/while/match)."""
    if isinstance(node, (ast.If, ast.For, ast.While, ast.AsyncFor)):
        return True
    if hasattr(ast, "Match") and isinstance(node, ast.Match):
        return True
    return False


def _inner_control_lines(stmts):
    """Sum lines of inner control statements."""
    total = 0
    for stmt in stmts:
        if _is_control_stmt(stmt):
            total += _count_lines(stmt)
    return total


def _q4_check_body(func_name, stmts, results):
    """Check depth-1 control statement bodies for PURE line violations (>10)."""
    for stmt in stmts:
        if isinstance(stmt, (ast.For, ast.While, ast.AsyncFor)):
            body_lines = _count_lines(stmt) - 1  # exclude the statement line itself
            inner = _inner_control_lines(stmt.body)
            pure = body_lines - inner
            stmt_type = "for" if isinstance(stmt, (ast.For, ast.AsyncFor)) else "while"
            if pure > 10:
                results.append({
                    "func_name": func_name,
                    "stmt_type": stmt_type,
                    "pure_lines": pure,
                    "line": stmt.lineno,
                })
        elif isinstance(stmt, ast.If):
            body_lines = 0
            if stmt.body:
                last_body = stmt.body[-1]
                if hasattr(last_body, "end_lineno") and last_body.end_lineno is not None:
                    body_lines = last_body.end_lineno - stmt.lineno
            inner = _inner_control_lines(stmt.body)
            pure = body_lines - inner
            if pure > 10:
                results.append({
                    "func_name": func_name,
                    "stmt_type": "if",
                    "pure_lines": pure,
                    "line": stmt.lineno,
                })
        elif hasattr(ast, "Match") and isinstance(stmt, ast.Match):
            for case_node in stmt.cases:
                case_body_lines = _count_lines(case_node) - 1
                inner = _inner_control_lines(case_node.body)
                pure = case_body_lines - inner
                if pure > 10:
                    results.append({
                        "func_name": func_name,
                        "stmt_type": "match",
                        "pure_lines": pure,
                        "line": case_node.pattern.lineno if hasattr(case_node, "pattern") else stmt.lineno,
                    })


def _extract_calls(node):
    """Extract function call names from AST (simple static analysis)."""
    calls = []
    for child in ast.walk(node):
        if not isinstance(child, ast.Call):
            continue
        func = child.func
        if isinstance(func, ast.Name):
            calls.append(func.id)
        elif isinstance(func, ast.Attribute):
            parts = []
            obj = func
            while isinstance(obj, ast.Attribute):
                parts.append(obj.attr)
                obj = obj.value
            if isinstance(obj, ast.Name):
                parts.append(obj.id)
            parts.reverse()
            calls.append(".".join(parts))
    return calls


def _extract_imports(tree):
    """Extract import information."""
    imports = []
    for node in ast.iter_child_nodes(tree):
        if isinstance(node, ast.Import):
            for alias in node.names:
                imports.append({"module": alias.name, "names": []})
        elif isinstance(node, ast.ImportFrom):
            module = node.module or ""
            names = [alias.name for alias in node.names]
            imports.append({"module": module, "names": names})
    return imports


def _body_hash(source, func_node):
    """Compute SHA-256 hash of first func/method signature+body (first 8 hex chars)."""
    lines = source.splitlines(True)
    start = func_node.lineno - 1
    end = func_node.end_lineno if hasattr(func_node, "end_lineno") and func_node.end_lineno else len(lines)
    chunk = "".join(lines[start:end])
    h = hashlib.sha256(chunk.encode("utf-8")).hexdigest()
    return h[:8]


def parse_file(path):
    """Parse a single Python file and return its structure as a dict."""
    try:
        with open(path, "r", encoding="utf-8") as f:
            source = f.read()
    except (OSError, UnicodeDecodeError) as e:
        return {"path": path, "error": str(e)}

    try:
        tree = ast.parse(source, filename=path)
    except SyntaxError as e:
        return {"path": path, "error": str(e)}

    functions = []
    classes = []
    methods = []
    has_init_method = False
    module_vars = []
    func_lines = {}
    all_calls = []
    q4_results = []
    first_func_node = None

    has_loop_d1 = False
    has_match_d1 = False

    for node in ast.iter_child_nodes(tree):
        if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)):
            functions.append(node.name)
            lines = _count_lines(node)
            func_lines[node.name] = lines
            if first_func_node is None:
                first_func_node = node
            if _has_control_at_depth1(node, "loop"):
                has_loop_d1 = True
            if _has_control_at_depth1(node, "match"):
                has_match_d1 = True
            all_calls.extend(_extract_calls(node))
            _q4_check_body(node.name, node.body, q4_results)

        elif isinstance(node, ast.ClassDef):
            classes.append(node.name)
            for item in node.body:
                if isinstance(item, (ast.FunctionDef, ast.AsyncFunctionDef)):
                    if item.name == "__init__":
                        has_init_method = True
                    else:
                        method_name = node.name + "." + item.name
                        methods.append(method_name)
                        lines = _count_lines(item)
                        func_lines[method_name] = lines
                        if first_func_node is None:
                            first_func_node = item
                        if _has_control_at_depth1(item, "loop"):
                            has_loop_d1 = True
                        if _has_control_at_depth1(item, "match"):
                            has_match_d1 = True
                        all_calls.extend(_extract_calls(item))
                        _q4_check_body(method_name, item.body, q4_results)

        elif isinstance(node, ast.Assign):
            for target in node.targets:
                if isinstance(target, ast.Name) and target.id.isupper():
                    module_vars.append(target.id)

    # Determine control
    if has_match_d1:
        control = "selection"
    elif has_loop_d1:
        control = "iteration"
    else:
        control = "sequence"

    # max_depth over all funcs/methods
    max_depth = 0
    for node in ast.iter_child_nodes(tree):
        if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)):
            d = _func_max_depth(node)
            if d > max_depth:
                max_depth = d
        elif isinstance(node, ast.ClassDef):
            for item in node.body:
                if isinstance(item, (ast.FunctionDef, ast.AsyncFunctionDef)):
                    d = _func_max_depth(item)
                    if d > max_depth:
                        max_depth = d

    # total lines
    total_lines = len(source.splitlines())

    # body hash
    bhash = ""
    if first_func_node is not None:
        bhash = _body_hash(source, first_func_node)

    # imports
    imports = _extract_imports(tree)

    # deduplicate calls
    seen = set()
    unique_calls = []
    for c in all_calls:
        if c not in seen:
            seen.add(c)
            unique_calls.append(c)

    return {
        "path": path,
        "functions": functions,
        "classes": classes,
        "methods": methods,
        "has_init_method": has_init_method,
        "vars": module_vars,
        "lines": total_lines,
        "max_depth": max_depth,
        "control": control,
        "has_loop_at_depth1": has_loop_d1,
        "has_match_at_depth1": has_match_d1,
        "func_lines": func_lines,
        "q4_results": q4_results,
        "calls": unique_calls,
        "imports": imports,
        "body_hash": bhash,
    }


def main():
    if len(sys.argv) < 2:
        print("Usage: py_ast.py [--batch] file.py [file2.py ...]", file=sys.stderr)
        sys.exit(1)

    args = sys.argv[1:]
    batch = False
    if args[0] == "--batch":
        batch = True
        args = args[1:]

    if not args:
        print("No files specified", file=sys.stderr)
        sys.exit(1)

    if batch:
        results = [parse_file(p) for p in args]
        json.dump(results, sys.stdout, ensure_ascii=False)
    else:
        result = parse_file(args[0])
        json.dump(result, sys.stdout, ensure_ascii=False)


if __name__ == "__main__":
    main()
