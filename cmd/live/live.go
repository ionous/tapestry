package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"git.sr.ht/~ionous/tapestry/web"
	"git.sr.ht/~ionous/tapestry/web/js"
)

func main() {
	r := bufio.NewReader(os.Stdin)

	messages := make(chan string, 1024)

	// read from stdin
	go func() {
		for {
			if in, e := r.ReadString('\n'); e != nil {
				break
			} else if cnt := len(in); cnt > 0 {
				messages <- in[:cnt-1] // trim the newline.
			}
		}
	}()

	buildMsgs := func() string {
		var out js.Builder // fix: really need to make this streamable.
		out.Brace(js.Array, func(a *js.Builder) {
			// some arbitrary maximum size
			for i := 0; i < 512; i++ {
				select {
				case msg := <-messages:
					if i > 0 {
						a.R(js.Comma)
					}
					a.Q(msg) // json quote each string.
				default:
					break // done
				}
			}
		})
		return out.String()
	}
	writeMsgs := func(w http.ResponseWriter) (err error) {
		w.Header().Set("Content-Type", "application/json")
		msgs := buildMsgs()
		_, err = io.WriteString(w, msgs)
		return
	}

	// backend
	u, _ := url.Parse("http://localhost:3000/")
	p := httputil.NewSingleHostReverseProxy(u)
	http.Handle("/", p)

	// input/output processing
	http.HandleFunc("/io/", web.HandleResource(&web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			if name == "io" {
				ret = &web.Wrapper{

					// polling for data.
					Gets: func(ctx context.Context, w http.ResponseWriter) (err error) {
						return writeMsgs(w)
					},
					// sending a command.
					Posts: func(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
						if msg, e := Decode(r); e != nil {
							err = e
						} else {
							select {
							case messages <- msg:
								log.Println("posted:", msg)
								err = writeMsgs(w)
							default:
								log.Println("ignored:", msg)
								_, err = io.WriteString(w, "(server busy)")
							}
						}
						return
					},
				}
			}
			return
		}}))
	log.Println("Listening on port", ":8088", "...")
	if e := http.ListenAndServe(":8088", nil); e != nil {
		log.Fatal(e)
	}
}

func Decode(r io.Reader) (ret string, err error) {
	var cmd struct {
		Msg string `json:"cmd"`
	}
	dec := json.NewDecoder(r)
	// decode an array value (Message)
	if e := dec.Decode(&cmd); e != nil {
		err = e
	} else {
		ret = cmd.Msg
	}
	return
}
