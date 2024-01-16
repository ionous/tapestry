package serve

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/play"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/web"
)

// read a request from the client
// see live.Vue... example: io.send({cmd: txt});
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

func Marshal(m play.PlayMessage) (ret string, err error) {
	var enc encode.Encoder
	if d, e := enc.Encode(m.(jsn.Marshalee)); e != nil {
		err = e
	} else {
		var str strings.Builder
		if e := files.JsonEncoder(&str, files.EscapeHtml).Encode(d); e != nil {
			err = e
		} else {
			ret = str.String()
		}
	}
	return
}

func marshal(m jsn.Marshalee) (ret string, err error) {
	var enc encode.Encoder
	if d, e := enc.Encode(m); e != nil {
		err = e
	} else {
		var str strings.Builder
		if e := files.JsonEncoder(&str, files.EscapeHtml).Encode(d); e != nil {
			err = e
		} else {
			ret = str.String()
		}
	}
	return
}

func ListenAndServe(endpoint string, cs *Channels) error {
	log.Println("Listening to", endpoint, "...")

	// our micro-server:
	u, _ := url.Parse("http://localhost:3000/")
	p := httputil.NewSingleHostReverseProxy(u)
	http.Handle("/", p)

	// input/output processing:
	http.HandleFunc("/io/", web.HandleResource(&web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			if name == "io" {
				ret = &web.Wrapper{
					// client polled for data.
					Gets: func(ctx context.Context, w http.ResponseWriter) (err error) {
						w.Header().Set("Content-Type", "application/json")
						cs.WriteMessages(w, false)
						return
					},
					// client sent a command
					// ( uses post because the same command multiple times can produce different results )
					Posts: func(ctx context.Context, r io.Reader, w http.ResponseWriter) (err error) {
						if cmd, e := Decode(r); e != nil {
							err = e
						} else {
							select {
							case cs.input <- cmd:
								log.Println("pushed:", cmd)
							default:
								log.Println("ignored:", cmd)
							}
							w.Header().Set("Content-Type", "application/json")
							cs.WriteMessages(w, true)
						}
						return
					},
				}
			}
			return
		}}))

	return http.ListenAndServe(endpoint, nil)
}
