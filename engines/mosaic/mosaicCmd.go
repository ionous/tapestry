// Copyright (C) 2022 - Simon Travis. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"git.sr.ht/~ionous/tapestry/web"
)

// exported to package main in cmd/tap
// var CmdMosaic = &base.Command{
// 	Run:       runMosaic,
// 	Flag:      buildFlags(),
// 	UsageLine: "tap edit [-in <directory>] [mosaic flags]",
// 	Short:     "run the tapestry story editor",
// 	Long: `Start the Tapestry story editor.

// The 'in' directory should contain two sub-directories:
// 	1. "stories" - containing story files ( the target for save/load )
// 	2. "shared"  - containing shared libraries ( defaults will be used if this folder is missing. )

// By default, attempts to use a directory called Tapestry in your Documents folder.
// `,
// }

func main() {
	if e := runMosaic(); e != nil {
		fmt.Println("error", e)
	}
}
func runMosaic() (err error) {
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if e := cmdLine.Parse(os.Args[1:]); e != nil {
		err = e
	} else if dir, e := mosaicFlags.folder.GetFolder(); e != nil {
		err = e
	} else {
		var ws Workspace
		mux := http.NewServeMux()

		// FIX: remove the "cmdDir"
		// everything should be using tap internals at this point i think.
		cfg := Configure(build.Default.GOPATH, dir)

		// raw story files ( because why not )
		mux.Handle("/stories/", web.HandleResource(FilesApi(cfg)))

		// blockly blocks ( from story files )
		mux.Handle("/blocks/", web.HandleResource(FilesApi(cfg)))

		// blockly shape files ( from spec files )
		mux.Handle("/shapes/", http.StripPrefix("/shapes/", web.HandleResource(ShapesApi(cfg))))

		// blockly toolbox files ( from spec files )
		mux.Handle("/boxes/", http.StripPrefix("/boxes/", web.HandleResource(BoxesApi(cfg))))

		// ui actions
		mux.Handle("/actions/", http.StripPrefix("/actions/", web.HandleResource(ActionsApi(cfg, &ws))))

		// fix: serve specs from package idl?
		// below is the file endpoint

		if BuildConfig == Prod {
			// prod redirects unknown url requests to our embedded assets
			// note: package http will tack on index.html redirects for bare directories automatically.
			// and, for good or ill, it will serve directories of files as actual directory listings
			mux.Handle("/", web.MethodHandler{
				http.MethodGet:  http.FileServer(http.FS(Frontend)),
				http.MethodPost: HandleCommands(cfg),
			})

		} else {
			// web and dev start a custom server to listen for incoming requests
			// and send unknown url requests to the vite server
			listenTo, _ := mosaicFlags.listen.GetPort(8080)
			requestFrom, _ := mosaicFlags.request.GetPort(3000)
			log.Println("using story files from:", dir)
			log.Println("listening to:", listenTo)
			log.Println("requesting from:", requestFrom)
			log.Printf("browse to: http://localhost:%d/\n", listenTo)

			// anything not handled by "mux" gets sent to the vite server.
			vite := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",
				Host:   "localhost:" + strconv.Itoa(requestFrom),
			})
			mux.Handle("/", web.MethodHandler{
				http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					log.Println(req.Method, req.RequestURI)
					vite.ServeHTTP(w, req)
				}),
				http.MethodPost: HandleCommands(cfg),
			})

			if BuildConfig != Prod {
				log.Println("don't forget to run the vite web server")
				log.Println("in the directory tapestry/www type 'npm run dev'.")
			}

			// NOTE: web mode stops here
			if BuildConfig != Dev {
				startBackend(listenTo, mux)
				// doesn't return.
			}

			// dev mode starts the tapestry server; then continues on to start the wails webkit browser.
			// the server has to tell the browser that its okay to make requests from the webkit origin.
			// wails webkit requests pages from "wails://" ( used to be from "http://wails.localhost" )
			// the webapps ask for ( and post/put data at ) "http://localhost:8080"
			go startBackend(listenTo, web.HandleCors("wails://wails", mux))

			// fix? having these commented out skips activating the wails dev server.
			// i haven't quite figured out how to integrate it with successfully with backend
			// not sure if its really needed, but it really *wants* to serve content.
			//os.Setenv("devserver", "http://localhost:34115")
			//os.Setenv("frontenddevserverurl", "http://localhost:3000")
		}
		// doesn't return.
		runWails(&ws, mosaicFlags.width, mosaicFlags.height, mux)
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
