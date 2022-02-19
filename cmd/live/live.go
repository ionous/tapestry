package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net/http"
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

	http.HandleFunc("/io/", web.HandleResource(&web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			if name == "io" {
				ret = &web.Wrapper{
					Gets: func(ctx context.Context, w http.ResponseWriter) (err error) {
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
						// write msgs back:
						w.Header().Set("Content-Type", "application/json")
						_, err = io.WriteString(w, out.String())
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
