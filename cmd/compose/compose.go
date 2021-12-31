package main

import (
	"flag"
	"go/build"

	"git.sr.ht/~ionous/tapestry/composer"
	"git.sr.ht/~ionous/tapestry/web/support"
)

// ex. go run compose.go -open -dir /Users/ionous/Documents/tapestry
// the specified directory needs two sub-directories:
// 1. "stories" - containing .if files
// 2. "ifspec"  - containing .ifspec files
func main() {
	var dir string
	var open bool
	flag.StringVar(&dir, "dir", "", "directory for processing if files.")
	flag.BoolVar(&open, "open", false, "open a new browser window.")
	flag.Parse()
	//
	cfg := composer.DevConfig(build.Default.GOPATH)
	if len(dir) > 0 {
		cfg.Root = dir
	}
	if open {
		support.OpenBrowser("http://localhost:3000/compose/")
	}
	// by design, this never returns.
	composer.Compose(cfg)
}
