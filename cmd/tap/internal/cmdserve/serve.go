package cmdserve

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
)

func serveWithOptions(inFile string, opts qna.Options, listenTo, requestFrom int) (ret int, err error) {
	if serverContext, e := newServerContext(inFile, opts); e != nil {
		err = e
	} else {
		defer serverContext.Close()
		mux := http.NewServeMux()
		// our main command service:
		mux.HandleFunc("/shuttle/", newServer("shuttle", serverContext))
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

func newServer(path string, serverContext serverContext) http.HandlerFunc {
	var state State
	return web.HandleResource(&web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			if name == path {
				ret = &web.Wrapper{
					// client sent a command
					Posts: func(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
						if msg, e := Decode(r); e != nil {
							err = e
						} else {
							switch msg := msg.(type) {
							case string:
								if h := state.HandleInput; h == nil {
									err = errutil.New("invalid input state", state.Name)
								} else if next, e := h(w, msg); e != nil {
									err = e
								} else if len(next.Name) > 0 {
									state = next
								}
							case map[string]any:
								// maybe also "$weave", "$shutdown" ...
								if scene, ok := msg["$restart"].(string); ok {
									if next, e := restart(w, serverContext, scene); e != nil {
										err = e
									} else {
										state = next
									}
								} else {
									// not a system level command, so pass to the current state (if any)
									if h := state.HandleCommand; h == nil {
										err = errutil.New("invalid command state", state.Name)
									} else if next, e := h(w, msg); e != nil {
										err = e
									} else if len(next.Name) > 0 {
										state = next
									}
								}
							}
						}
						return
					},
				}
			}
			return
		}})
}

// read a request from the client
// see Play.vue... example: io.send({in: txt});
func Decode(r io.Reader) (ret any, err error) {
	var msg struct {
		Input   string         `json:"in"`
		Command map[string]any `json:"cmd"`
	}
	dec := json.NewDecoder(r)
	// decode an array value (Message)
	if e := dec.Decode(&msg); e != nil {
		err = e
	} else if len(msg.Command) > 0 {
		ret = msg.Command
	} else if len(msg.Input) > 0 {
		ret = msg.Input
	} else {
		err = errutil.New("unknown or empty input")
	}
	return
}
