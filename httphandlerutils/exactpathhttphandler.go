package httphandlerutils

import (
	"fmt"
	"net/http"
)

type exactPathHTTPHandler struct {
	path   string
	nested http.Handler
}

// ExactPath only accepts a request if it for a certain path exactly
// Returns HTTP 404 when the request path does not match
func ExactPath(path string, nested http.Handler) http.Handler {
	return &exactPathHTTPHandler{
		path:   path,
		nested: nested,
	}
}

func (exact *exactPathHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == exact.path {
		exact.nested.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, http.StatusText(http.StatusNotFound))
	}
}
