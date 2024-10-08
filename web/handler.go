package web

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/ionous/errutil"
)

// HandleResource turns a Resource into an http.HandlerFunc;
// providing responses to http get and post requests.
func HandleResource(root Resource) http.HandlerFunc {
	return HandleResourceWithContext(root, func(ctx context.Context) context.Context {
		return ctx
	})
}

// HandleResource turns a Resource into an http.HandlerFunc;
// providing responses to http get and post requests.
func HandleResourceWithContext(root Resource, xform func(context.Context) context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Println("( handling", r.URL.Path, r.Method, ")")
		if e := handleResponse(w, r, root, xform); e != nil {
			// note: handle response already sets http.Error;
			// so this is just for local logging purposes.
			log.Println("error handling", r.URL.Path, r.Method, e)
		}
	}
}

// NOTE: the error, if any, is automatically passed to http.Error
func handleResponse(w http.ResponseWriter, r *http.Request,
	root Resource, xform func(context.Context) context.Context,
) (err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// chop off the leading and trailing slash. wise? i dont know.
	if res, e := FindResource(root, strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/")); e != nil {
		http.NotFound(w, r)
		err = e
	} else {
		ctx := xform(r.Context())
		switch r.Method {
		case "GET":
			if e := res.Get(ctx, w); e != nil {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				err = e
			}
		case "POST":
			if e := res.Post(ctx, r.Body, w); e != nil {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				err = e
			}
		case "PUT":
			if e := res.Put(ctx, r.Body, w); e != nil {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				err = e
			}
		default:
			http.Error(w, r.Method, http.StatusMethodNotAllowed)
			err = errutil.Fmt("method %s not allowed", r.Method)
		}
	}
	return
}

// FindResource expands the passed resource, using each element of the passed path in turn.
// Returns an error, PathError, describing the extent of the matched path.
func FindResource(res Resource, path string) (ret Resource, err error) {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if sub := res.Find(part); sub != nil {
			res = sub // set for next iteration of the loop
		} else {
			err = errutil.Fmt("failed to find resource %d(%s) in %s", i, part, path)
			break
		}
	}
	if err == nil {
		ret = res
	}
	return
}
