//ff:func feature=cli type=util
//ff:what 환경변수 값이 있으면 반환하고 없으면 기본값 반환
//ff:checked llm=gpt-oss:20b hash=a842ab03
package cli

import "os"

// EnvOrDefault returns the environment variable value or the default.
func EnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
