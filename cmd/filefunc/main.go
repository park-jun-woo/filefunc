//ff:func feature=cli type=command control=sequence
//ff:what filefunc CLI 엔트리포인트
package main

import (
	"os"

	"github.com/park-jun-woo/filefunc/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
