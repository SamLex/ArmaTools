package httphandlerutils

import (
	"fmt"
	"net/http"
)

type panicHTTPHandler struct {
	nested         http.Handler
	revealInternal bool
}

// PanicHandler recovers a panic from the nested Handler
// Causes a HTTP 500 and an appropriate message on a panic
func PanicHandler(nested http.Handler) http.Handler {
	return &panicHTTPHandler{
		nested:         nested,
		revealInternal: true,
	}
}

func (p *panicHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Actual panic handler
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
