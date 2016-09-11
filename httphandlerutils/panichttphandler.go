/*
Copyright 2016 Euan James Hunter

panichttphandler.go: Panic HTTP Handler utility

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
