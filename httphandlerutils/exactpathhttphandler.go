/*
Copyright 2016 Euan James Hunter

exactpathhttphandler.go: ExactPath HTTP Handler utility

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
