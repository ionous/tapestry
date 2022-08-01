// Run the blockly editor.
// Note: currently requires running the npm content server in tapestry/www. ( > npm run dev )
package main

import (
	"flag"
	"go/build"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/composer"
	"git.sr.ht/~ionous/tapestry/web"
	"git.sr.ht/~ionous/tapestry/web/support"
	"github.com/ionous/errutil"
)

func main() {
	var dir string
	var open, dist bool
	flag.StringVar(&dir, "in", "", "directory for if files.")
	flag.BoolVar(&dist, "dist", false, "use the www production files built by vite .")
	flag.BoolVar(&open, "open", false, "open a new browser window.")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error.")
	flag.Parse()
	if len(dir) == 0 {
		flag.Usage()
		return
	}

	if open {
		support.OpenBrowser(web.Endpoint(8080, "localhost", "mosaic"))
	}
	// by design, this never returns.
	var cfg *web.Config
	if !dist {
		cfg = web.DevConfig(build.Default.GOPATH, dir)
	} else {
		cfg = web.DevDist(build.Default.GOPATH, dir, filepath.Clean("../../www/dist"))
	}
	composer.RunMosaic(cfg, 8080)
}

const Description = //
`Mosaic: starts the backend for Tapestry's blockly editor providing shape definitions, save/load functionality, etc.
Currently requires separately running the npm content server in tapestry/www. ( > npm run dev )

Requires a directory containing two sub-directories:
	1. "stories" - containing .if files ( the target for save/load )
	2. "ifspec"  - containing .ifspecs ( these define how to present the story content )
`
const Example = "go build mosaic.go && mosaic.exe -in /Users/ionous/Documents/Tapestry -dist -open"

func init() {
	flag.Usage = func() {
		println(Description)
		flag.PrintDefaults()
		println("\nex.", Example)
	}
}
