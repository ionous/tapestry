package web

import "net/http"

// implements a handler per method.
// ( doesnt attempt to validate method names, and doesn't notice duplicate registration. last one wins )
type MethodHandler map[string]http.Handler

func (m MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := m[r.Method]; !ok {
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
	} else {
		h.ServeHTTP(w, r)
	}
}
