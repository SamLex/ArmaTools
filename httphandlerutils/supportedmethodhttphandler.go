/*
Copyright 2016 Euan James Hunter

supportedmethodshttphandler.go: SupportedMethods HTTP Handler utility

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
