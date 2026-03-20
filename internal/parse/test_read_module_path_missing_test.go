//ff:func feature=parse type=util control=sequence
//ff:what test: TestReadModulePath_Missing
package parse

import (
	"testing"
)

func TestReadModulePath_Missing(t *testing.T) {
	_, err := ReadModulePath("/nonexistent/go.mod")
	if err == nil {
		t.Error("expected error for nonexistent go.mod")
	}
}
