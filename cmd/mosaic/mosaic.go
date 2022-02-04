package main

import (
	"flag"
	"go/build"

	"git.sr.ht/~ionous/tapestry/composer"
	"git.sr.ht/~ionous/tapestry/web/support"
	"github.com/ionous/errutil"
)

// ex. go run mosaic.go -in /Users/ionous/Documents/Tapestry -open
// the specified directory needs two sub-directories:
// 1. "stories" - containing .if files
// 2. "ifspec"  - containing .ifspecs
func main() {
	var dir string
	var open bool
	flag.StringVar(&dir, "in", "", "directory for processing if files.")
	flag.BoolVar(&open, "open", false, "open a new browser window.")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error")
	flag.Parse()
	//
	cfg := composer.DevConfig(build.Default.GOPATH)
	cfg.Port = 8080
	cfg.Mosaic = "http://localhost:3000/"
	if len(dir) > 0 {
		cfg.Root = dir
	}
	if open {
		support.OpenBrowser("http://localhost" + cfg.PortString())
	}
	// by design, this never returns.
	composer.Mosaic(cfg)
}
