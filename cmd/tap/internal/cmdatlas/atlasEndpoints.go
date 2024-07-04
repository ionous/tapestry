package cmdatlas

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/support/dump"
)

type atlasContext struct {
	db  *sql.DB
	dec *query.QueryDecoder
}

func newAtlas(mux *http.ServeMux, a *atlasContext) {
	mux.HandleFunc("/atlas/scenes/", func(w http.ResponseWriter, req *http.Request) {
		if val, e := dump.QueryAllScenes(a.db); e != nil {
			handleError(w, e)
		} else {
			writeJson(w, val)
		}
	})
	mux.HandleFunc("/atlas/scenes/{scene}/patterns/", func(w http.ResponseWriter, req *http.Request) {
		scene := req.PathValue("scene")
		if val, e := dump.QueryPatterns(a.db, a.dec, scene); e != nil {
			handleError(w, e)
		} else {
			writeJson(w, val)
		}
	})
}

// call a function which returns some object,
// write that object as json to the output
func write[E any](w http.ResponseWriter, a *atlasContext, q func(db *sql.DB) (E, error)) (err error) {
	if val, e := q(a.db); e != nil {
		err = e
	} else {
		writeJson(w, val)
	}
	return
}

func writeJson(w http.ResponseWriter, val any) {
	w.Header().Set("Content-Type", "application/json")
	js := json.NewEncoder(w)
	js.SetEscapeHTML(false)
	if e := js.Encode(val); e != nil {
		handleError(w, e)
	}
}

func handleError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(e)
}
