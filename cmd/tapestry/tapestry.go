package main

import (
	"context"
	"flag"
	"go/build"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	tap "git.sr.ht/~ionous/tapestry/cmd/tapestry/internal"

	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func main() {
	var folder tap.Folder
	var listen, request tap.Port
	var width, height int
	//
	flag.Var(&folder, "in", "directory for if files.")
	flag.IntVar(&width, "width", 1024, "width of the application window.")
	flag.IntVar(&height, "height", 768, "height of the application window.")
	if tap.BuildConfig != tap.Prod {
		flag.Var(&listen, "listen", "location for everything happening via the browser. specify a port number; or, 'true' for the default (8080).")
		if tap.BuildConfig == tap.Web {
			flag.Var(&request, "www", "location of the tapestry webapps. specify a port number; or, 'true' to use the default port (3000).")
		}
	}
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error.")
	flag.Parse()

	if dir, e := folder.GetFolder(); e != nil {
		log.Println(e)
		flag.Usage()
	} else {
		mux := http.NewServeMux()

		// redirect index.html to the mosaic app.
		// [ doesnt work for wails -- it looks for a built in index.html ]
		// composer.RedirectIndex(mux, "mosaic")

		// FIX: remove the "cmdDir"
		cfg := tap.DevConfig(build.Default.GOPATH, dir)

		// raw story files ( because why not )
		mux.Handle("/stories/", web.HandleResource(tap.FilesApi(cfg)))

		// blockly blocks ( from .if )
		mux.Handle("/blocks/", web.HandleResource(tap.FilesApi(cfg)))

		// blockly shape files ( from .ifspecs )
		mux.Handle("/shapes/", http.StripPrefix("/shapes/", web.HandleResource(tap.ShapesApi(cfg))))

		// blockly shape files ( from .ifspecs )
		mux.Handle("/boxes/", http.StripPrefix("/boxes/", web.HandleResource(tap.BoxesApi(cfg))))

		// fix: serve .ifspecs from package idl?
		// below is the older .ifspec endpoint --
		// not clear if that's desirable at this pint...

		if tap.BuildConfig == tap.Prod {
			// prod redirects unknown url requests to our embedded assets
			// note: package http will tack on index.html redirects for bare directories automatically.
			// and, for good or ill, it will serve directories of files as actual directory listings
			mux.Handle("/", web.MethodHandler{
				http.MethodGet:  http.FileServer(http.FS(tap.Frontend)),
				http.MethodPost: tap.HandleCommands(cfg),
			})

		} else {
			// web and dev start a custom server to listen for incoming requests
			// and send unknown url requests to the vite backend
			listenTo, _ := listen.GetPort(8080)
			requestFrom, _ := request.GetPort(3000)
			log.Println("using story files from:", dir)
			log.Println("listening to:", listenTo)
			log.Println("requesting from:", requestFrom)
			if viteBackend, e := url.Parse(web.Endpoint(requestFrom, "localhost")); e != nil {
				log.Fatal(e)
			} else {
				// anything not handled by "mux" gets sent to the vite backend.
				mux.Handle("/", web.MethodHandler{
					http.MethodGet:  httputil.NewSingleHostReverseProxy(viteBackend),
					http.MethodPost: tap.HandleCommands(cfg),
				})

			}

			// web mode stops here
			if tap.BuildConfig != tap.Dev {
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
		runWails(width, height, mux)
	}
}

// not expected to return
func startBackend(listenTo int, mux http.Handler) {
	// start the thing that listens to browser requests
	if e := http.ListenAndServe(":"+strconv.Itoa(listenTo), mux); e != nil {
		log.Fatal(e)
	}
}

// not expected to return
func runWails(width, height int, mux http.Handler) {
	var wailsFakeout fs.FS
	if tap.BuildConfig == tap.Prod {
		wailsFakeout = tap.ErrFS{}
	}

	var host Host
	if e := wails.Run(&options.App{
		Title:  "Tapestry",
		Width:  width,
		Height: height,
		// in dev: we pass nil, and all files are served via the "AssetsHandler"
		// in production: we pass a special filesystem that always errors out.
		Assets: wailsFakeout,
		// a way to serve dynamic data from specific endpoints
		// in dev: any paths not handled are passed to vite for processing.
		// in production: it excludes the fallback handler
		AssetsHandler: mux,
		//BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: host.startup,
		// OnShutdown: host.shutdown,
		//Menu:       host.menu(), -> https://wails.io/docs/reference/menus
		Bind: []interface{}{
			&host, // all public methods are exposed to wails as javascript commands
		},
	}); e != nil {
		log.Fatal()
	}
}

// wails application framing
type Host struct{ ctx context.Context }

// startup is called when the app starts.
// saves context in case we need to call wails runtime methods
func (w *Host) startup(ctx context.Context) {
	w.ctx = ctx
}

const Description = //
`Start the Tapestry editor.

Requires a directory containing two sub-directories:
	1. "stories" - containing .if files ( the target for save/load )
	2. "ifspec"  - containing .ifspecs ( these define how to present the story content )

By default, attempts to use a directory called Tapestry in your Documents folder.
`

var Examples = []string{
	`go build -o web.exe && web.exe`,
	`go build -tags dev -o dev.exe && dev.exe`,
	`go build -tags production,desktop -ldflags "-w -s -H windowsgui" && tapestry.exe`,
}

func init() {
	flag.Usage = func() {
		println(Description)
		flag.PrintDefaults()
		for _, ex := range Examples {
			println("\nex.", ex)
		}
	}
}
