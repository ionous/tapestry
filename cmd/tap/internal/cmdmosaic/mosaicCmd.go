// Copyright (C) 2022 - Simon Travis. All rights reserved.
// Use of this source code is governed by the Hippocratic 2.1
// license that can be found in the LICENSE file.
package cmdmosaic

import (
	"context"
	"go/build"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/mosaic"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
)

var CmdMosaic = &base.Command{
	Run:       runMosaic,
	Flag:      buildFlags(),
	UsageLine: "tap mosaic [-in <directory>] [mosaic flags]",
	Short:     "tapestry story editor",
	Long: `Start the Tapestry story editor.

The 'in' directory should contain two sub-directories:
	1. "stories" - containing .if files ( the target for save/load )
	2. "ifspec"  - containing .ifspecs ( these define how to present the story content )

By default, attempts to use a directory called Tapestry in your Documents folder.
`,
}

func runMosaic(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if dir, e := mosaicFlags.folder.GetFolder(); e != nil {
		err = e
	} else if types, e := rs.FromSpecs(idl.Specs); e != nil {
		err = errutil.New("fatal error:", e)
	} else {
		mux := http.NewServeMux()

		// redirect index.html to the mosaic app.
		// [ doesnt work for wails -- it looks for a built in index.html ]
		// composer.RedirectIndex(mux, "mosaic")

		// FIX: remove the "cmdDir"
		cfg := mosaic.Configure(types, build.Default.GOPATH, dir)

		// raw story files ( because why not )
		mux.Handle("/stories/", web.HandleResource(mosaic.FilesApi(cfg)))

		// blockly blocks ( from .if )
		mux.Handle("/blocks/", web.HandleResource(mosaic.FilesApi(cfg)))

		// blockly shape files ( from .ifspecs )
		mux.Handle("/shapes/", http.StripPrefix("/shapes/", web.HandleResource(mosaic.ShapesApi(cfg))))

		// blockly shape files ( from .ifspecs )
		mux.Handle("/boxes/", http.StripPrefix("/boxes/", web.HandleResource(mosaic.BoxesApi(cfg))))

		// fix: serve .ifspecs from package idl?
		// below is the older .ifspec endpoint --
		// not clear if that's desirable at this pint...

		if mosaic.BuildConfig == mosaic.Prod {
			// prod redirects unknown url requests to our embedded assets
			// note: package http will tack on index.html redirects for bare directories automatically.
			// and, for good or ill, it will serve directories of files as actual directory listings
			mux.Handle("/", web.MethodHandler{
				http.MethodGet:  http.FileServer(http.FS(mosaic.Frontend)),
				http.MethodPost: mosaic.HandleCommands(cfg),
			})

		} else {
			// web and dev start a custom server to listen for incoming requests
			// and send unknown url requests to the vite backend
			listenTo, _ := mosaicFlags.listen.GetPort(8080)
			requestFrom, _ := mosaicFlags.request.GetPort(3000)
			log.Println("using story files from:", dir)
			log.Println("listening to:", listenTo)
			log.Println("requesting from:", requestFrom)
			log.Printf("browse to: http://localhost:%d/mosaic/\n", listenTo)
			if viteBackend, e := url.Parse(web.Endpoint(requestFrom, "localhost")); e != nil {
				log.Fatal(e)
			} else {
				// anything not handled by "mux" gets sent to the vite backend.
				mux.Handle("/", web.MethodHandler{
					http.MethodGet:  httputil.NewSingleHostReverseProxy(viteBackend),
					http.MethodPost: mosaic.HandleCommands(cfg),
				})

			}

			// web mode stops here
			if mosaic.BuildConfig != mosaic.Dev {
				startBackend(listenTo, mux) // doesn't return.
			}

			// dev mode starts the tapestry server; then continues on to start the wails webkit browser.
			// the server has to tell the browser that its okay to make requests from the webkit origin.
			// wails webkit requests pages from: "http://wails.localhost"
			// the webapps ask for ( and post/put data at ) "http://localhost:8080"
			go startBackend(listenTo, web.HandleCors("http://wails.localhost", mux))

			// fix? having these commented out skips activating the wails dev server.
			// i haven't quite figured out how to integrate it with succesfully with backend
			// not sure if its really needed, but it really *wants* to serve content.
			//os.Setenv("devserver", "http://localhost:34115")
			//os.Setenv("frontenddevserverurl", "http://localhost:3000")
		}
		// doesn't return.
		runWails(mosaicFlags.width, mosaicFlags.height, mux)
	}
	return
}

// TODO:
// `go build -o web.exe && web.exe`,
// `go build -tags dev -o dev.exe && dev.exe`,
// `go build -tags production,desktop -ldflags "-w -s -H windowsgui" && tapestry.exe`,

// func init() {
// 	flag.Usage = func() {
// 		println(Description)
// 		flag.PrintDefaults()
// 		for _, ex := range Examples {
// 			println("\nex.", ex)
// 		}
// 	}
// }
