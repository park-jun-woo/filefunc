//ff:func feature=chain type=util control=sequence
//ff:what qualified name(pkg.FuncName)에서 함수명만 추출
package chain

import "strings"

func nameFromQualified(qname string) string {
	if i := strings.LastIndex(qname, "."); i >= 0 {
		return qname[i+1:]
	}
	return qname
}
