package main

import (
	"io/fs"
	"log"
	"net/http"
	"strconv"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// wails application framing
type Host struct{}

// startup is called when the app starts.
// saves context in case we need to call wails runtime methods

// not expected to return
func startBackend(listenTo int, mux http.Handler) {
	// start the thing that listens to browser requests
	// explicitly specifying localhost quiets the windows firewall.
	if e := http.ListenAndServe("localhost:"+strconv.Itoa(listenTo), mux); e != nil {
		log.Fatal(e)
	}
}

// not expected to return
func runWails(ws *Workspace, width, height int, mux http.Handler) {
	var wailsFakeout fs.FS
	if BuildConfig == Prod {
		wailsFakeout = ErrFS{}
	}

	var host Host
	if e := wails.Run(&options.App{
		Title:  "Tapestry",
		Width:  width,
		Height: height,
		// in dev: this is nil, so that all files are served via the mux "AssetsHandler".
		// in production: all assets are built in, and so anything else is considered an error.
		Assets: wailsFakeout,
		// this serves dynamic data from specific endpoints:
		// in dev: any paths not handled are passed to vite for processing.
		// in production: it excludes the fallback handler
		AssetsHandler: mux,
		//BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: ws.Startup,
		// OnShutdown: host.shutdown,
		//Menu:       host.menu(), -> https://wails.io/docs/reference/menus
		Bind: []interface{}{
			&host, // all public methods are exposed to wails as javascript commands
		},
	}); e != nil {
		log.Fatal()
	}
}
