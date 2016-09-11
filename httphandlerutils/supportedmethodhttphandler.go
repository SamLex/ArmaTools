package httphandlerutils

import "net/http"

type supportedMethodHTTPHandler struct {
	nested           http.Handler
	supportedMethods []string
}

// SupportedMethods only accepts if the request method is one of those specified
// Returns HTTP 405 on an unsupported method
func SupportedMethods(nested http.Handler, methods ...string) http.Handler {
	return &supportedMethodHTTPHandler{
		nested:           nested,
		supportedMethods: methods,
	}
}

func (support *supportedMethodHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, method := range support.supportedMethods {
		if r.Method == method {
			support.nested.ServeHTTP(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
