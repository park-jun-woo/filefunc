//ff:func feature=parse type=util control=sequence
//ff:what 디렉토리 내 파일 존재 여부 확인
package walk

import "os"

// fileExists checks if a file with the given name exists in the root directory.
func fileExists(root, name string) bool {
	_, err := os.Stat(root + "/" + name)
	return err == nil
}
