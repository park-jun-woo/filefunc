//ff:func feature=parse type=util control=sequence
//ff:what test: writeTestFile
package parse

import (
	"os"
)

func writeTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
