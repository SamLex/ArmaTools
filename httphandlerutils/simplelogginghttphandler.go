package httphandlerutils

import (
	"log"
	"net/http"
)

type simpleLoggingHTTPHandler struct {
	nested http.Handler
	logger *log.Logger
}

func SimpleLogging(nested http.Handler) http.Handler {
	return &simpleLoggingHTTPHandler{
		nested: nested,
		logger: nil,
	}
}

func (simple *simpleLoggingHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	simple.nested.ServeHTTP(w, r)

	if simple.logger == nil {
		log.Println(r.RemoteAddr, r.Method, r.Host, r.RequestURI, r.Proto, r.ContentLength, w.Header().Get("Content-Length"), r.UserAgent())
	} else {
		simple.logger.Println(r.RemoteAddr, r.Method, r.Host, r.RequestURI, r.Proto, r.ContentLength, w.Header().Get("Content-Length"), r.UserAgent())
	}
}
