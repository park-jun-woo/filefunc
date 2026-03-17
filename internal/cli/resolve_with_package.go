//ff:func feature=cli type=util control=iteration dimension=1
//ff:what 매칭된 qualified name 목록에서 지정 패키지에 해당하는 것을 선택
package cli

import (
	"fmt"
	"strings"
)

func resolveWithPackage(matches []string, target string, pkg string) (string, error) {
	for _, m := range matches {
		if pkgFromQualified(m) == pkg {
			return m, nil
		}
	}
	return "", fmt.Errorf("func %q not found in package %q (available: %s)", target, pkg, strings.Join(matches, ", "))
}
