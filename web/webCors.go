package web

import (
	"net/http"
)

// Allow some remote origin access to the specified handler.
func HandleCors(allowedOrigin string, h http.Handler) http.Handler {
	return &corsHandler{allowedOrigin, h}
}

type corsHandler struct {
	s string
	h http.Handler
}

// tbd: could this expose a nicer set of things?
func (c corsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := w.Header()
	h.Set("Access-Control-Allow-Origin", c.s)
	//log.Println("processing", req.Host, req.Method, req.URL.String())
	if req.Method == http.MethodOptions {
		h.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
		h.Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	} else {
		c.h.ServeHTTP(w, req)
	}
}
