package cmdserve

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/support/play"
	"git.sr.ht/~ionous/tapestry/support/player"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web"
)

func serveWithOptions(inFile string, opts qna.Options, listenTo, requestFrom int) (ret int, err error) {
	if db, e := tables.CreateRunTime(inFile); e != nil {
		err = e
	} else if ctx, e := makeShuttle(db, opts); e != nil {
		err = e
	} else {
		defer db.Close()
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

func makeShuttle(db *sql.DB, opts qna.Options) (ret *frame.Shuttle, err error) {
	// fix: merge with others and put in package player?
	decoder := query.NewDecoder(tapestry.AllSignatures)
	if grammar, e := player.MakeGrammar(db); e != nil {
		err = e
	} else if q, e := qdb.NewQueries(db, decoder); e != nil {
		err = e
	} else {
		run := qna.NewRuntimeOptions(q, opts)
		survey := play.MakeDefaultSurveyor(run)
		pt := play.NewPlaytime(run, survey, grammar)
		play.CaptureInput(pt)
		ret = frame.NewShuttle(pt, decoder)
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

func newServer(path string, ctx *frame.Shuttle) http.HandlerFunc {
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
