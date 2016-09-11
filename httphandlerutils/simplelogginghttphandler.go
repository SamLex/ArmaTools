/*
Copyright 2016 Euan James Hunter

simplelogginghttphandler.go: SimpleLogging HTTP Handler utility

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
	"log"
	"net/http"
)

type simpleLoggingHTTPHandler struct {
	nested http.Handler
	logger *log.Logger
}

// SimpleLogging logs every request in a simple format:
// "$remoteAddress $httpMethod $requestedHost $requestURI $httpProtocol $requestLength $responseLength $requestUserAgent"
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
