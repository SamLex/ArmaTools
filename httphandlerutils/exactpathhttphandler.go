package httphandlerutils

import "net/http"

type exactPathHTTPHandler struct {
	path    string
	nested  http.Handler
	invalid http.Handler
}

func ExactPath(path string, nested http.Handler) http.Handler {
	return &exactPathHTTPHandler{
		path:    path,
		nested:  nested,
		invalid: http.NotFoundHandler(),
	}
}

func (exact *exactPathHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == exact.path {
		exact.nested.ServeHTTP(w, r)
	} else {
		exact.invalid.ServeHTTP(w, r)
	}
}
