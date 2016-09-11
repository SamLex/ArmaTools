package httphandlerutils

import (
	"fmt"
	"net/http"
)

type panicHTTPHandler struct {
	nested         http.Handler
	revealInternal bool
}

func PanicHandler(nested http.Handler) http.Handler {
	return &panicHTTPHandler{
		nested:         nested,
		revealInternal: true,
	}
}

func (p *panicHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
			if p.revealInternal {
				fmt.Fprintln(w, r)
			}
		}
	}()

	p.nested.ServeHTTP(w, r)
}
