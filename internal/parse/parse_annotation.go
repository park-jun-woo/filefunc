//ff:func feature=parse type=parser
//ff:what Go 소스 파일 상단의 //ff: 어노테이션을 파싱
//ff:checked llm=gpt-oss:20b hash=119a8578
package parse

import (
	"bufio"
	"os"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// ParseAnnotation parses //ff: annotations from the top of a Go source file.
func ParseAnnotation(path string) (*model.Annotation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ann := &model.Annotation{
		Func:    make(map[string]string),
		Type:    make(map[string]string),
		Checked: make(map[string]string),
	}
	found := false
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !strings.HasPrefix(line, "//ff:") && (found || !IsSkippableLine(line)) {
			break
		}
		if !strings.HasPrefix(line, "//ff:") {
			continue
		}

		found = true
		rest := line[len("//ff:"):]
		spaceIdx := strings.IndexByte(rest, ' ')
		if spaceIdx == -1 {
			continue
		}
		ApplyAnnotationLine(ann, rest[:spaceIdx], strings.TrimSpace(rest[spaceIdx+1:]))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return ann, nil
}
