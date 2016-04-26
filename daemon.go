// The MIT License (MIT)

// Copyright (c) 2016 Maciej Borzecki

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	defaultListenPort = 9090

	uriApi    = "/api"
	uriStatic = "/static/"
	uriIndex  = "/"
)

var (
	defaultListenAddr = fmt.Sprintf("localhost:%d",
		defaultListenPort)
)

func setupHandlers() {
	r := mux.NewRouter()

	SetupSiteHandlers(r)

	apir := r.PathPrefix(uriApi).Subrouter()
	SetupApiHandlers(apir)

	// static files
	staticroot := path.Join(C.webroot, "static")
	log.Printf("serving static files from %s", staticroot)

	filehandler := http.FileServer(http.Dir(staticroot))
	r.PathPrefix(uriStatic).
		Handler(http.StripPrefix(uriStatic, filehandler))

	// setup logging for all handlers
	lr := handlers.LoggingHandler(os.Stdout, r)

	http.Handle("/", lr)
}

func daemonMain() error {

	setupHandlers()

	return http.ListenAndServe(fmt.Sprintf(":%d", C.port), nil)
}

func runDaemon() error {
	err := LoadConfig()
	if err != nil {
		return errors.Wrap(err, "daemon startup failed")
	}

	log.Printf("listen port: %d", C.port)

	return daemonMain()
}
