//ff:func feature=chain type=util control=sequence
//ff:what qualified name(pkg.FuncName)에서 패키지명을 추출
package chain

import "strings"

// PkgFromQualified extracts the package name from a qualified name (pkg.FuncName).
func PkgFromQualified(qname string) string {
	if i := strings.LastIndex(qname, "."); i >= 0 {
		return qname[:i]
	}
	return ""
}
