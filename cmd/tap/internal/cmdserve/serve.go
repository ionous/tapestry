package cmdserve

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/support/shuttle"
	"git.sr.ht/~ionous/tapestry/web"
)

func serveWithOptions(inFile string, opts qna.Options, listenTo, requestFrom int) (ret int, err error) {
	if ctx, e := shuttle.NewShuttle(inFile, opts); e != nil {
		err = e
	} else {
		defer ctx.Close()
		mux := http.NewServeMux()
		// our main command service:
		mux.HandleFunc("/shuttle/", newServer("shuttle", ctx))
		// create a proxy for the web apps ( does nothing if requestFrom port is zero )
		if requestFrom != 0 {
			proxyToVite(mux, requestFrom)
		}
		// block forever ish.
		// note: on windows the localhost is required in order to avoid the windows firewall popup
		where := "localhost:" + strconv.Itoa(listenTo)
		err = http.ListenAndServe(where, mux)
	}
	return
}

// create a reverse proxy for the web console app ( to handle CORS issues )
// the user can access the console app via their web browser at localhost:listenTo/play/
// amd this forwards those requests to a local vita web server
// which serves the vue console app at localhost:requestFrom/play/
func proxyToVite(mux *http.ServeMux, port int) {
	vite := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:" + strconv.Itoa(port),
	})
	// longer paths take precedence over shorter ones
	// therefore this handles everything that isnt handled.
	mux.Handle("/", web.MethodHandler{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Println(req.Method, req.RequestURI)
			vite.ServeHTTP(w, req)
		}),
	})
}

func newServer(path string, ctx shuttle.Shuttle) http.HandlerFunc {
	return web.HandleResource(&web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			if name == path {
				ret = &web.Wrapper{
					Finds: func(endpoint string) (ret web.Resource) {
						return &web.Wrapper{
							// client sent a command
							Posts: func(_ context.Context, r io.Reader, w http.ResponseWriter) (err error) {
								w.Header().Set("Content-Type", "application/json")
								if raw, e := io.ReadAll(r); e != nil {
									err = e
								} else {
									err = ctx.Post(w, endpoint, raw)
								}
								return
							},
						}
					},
				}
			}
			return
		}})
}
