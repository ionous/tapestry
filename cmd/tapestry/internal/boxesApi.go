package tap

import (
	"context"
	"io"
	"net/http"

	"git.sr.ht/~ionous/tapestry/blockly/box"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/web"
)

func BoxesApi(*Config) web.Resource {
	return &web.Wrapper{
		Finds: func(str string) (ret web.Resource) {
			if len(str) == 0 {
				ret = &web.Wrapper{
					Gets: func(ctx context.Context, w http.ResponseWriter) (err error) {
						if shapes, e := box.FromSpecs(idl.Specs); e != nil {
							err = e
						} else {
							w.Header().Set("Content-Type", "application/json")
							_, err = io.WriteString(w, shapes)
						}
						return
					},
				}
			}
			return
		},
	}
}
