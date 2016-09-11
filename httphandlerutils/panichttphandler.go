package httphandlerutils

import "net/http"

type panicHTTPHandler struct {
	nested http.Handler
}

func PanicHandler(nested http.Handler) http.Handler {
	return &panicHTTPHandler{
		nested: nested,
	}
}

func (p *panicHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	p.nested.ServeHTTP(w, r)
}
