package main

import (
	"flag"
	"go/build"

	"git.sr.ht/~ionous/tapestry/composer"
	"git.sr.ht/~ionous/tapestry/web"
	"git.sr.ht/~ionous/tapestry/web/support"
	"github.com/ionous/errutil"
)

// ex. go run compose.go -open -in /Users/ionous/Documents/Tapestry
// the specified directory needs two sub-directories:
// 1. "stories" - containing .if files
// 2. "ifspec"  - containing .ifspec files
func main() {
	var dir string
	var open bool
	flag.StringVar(&dir, "in", "", "directory for processing if files.")
	flag.BoolVar(&open, "open", false, "open a new browser window.")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error")
	flag.Parse()
	//
	cfg := web.DevConfig(build.Default.GOPATH, dir)
	if open {
		support.OpenBrowser(web.Endpoint(3000, "localhost", "compose"))
	}
	// by design, this never returns.
	composer.Compose(cfg, 3000)
}
