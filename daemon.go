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
	astatic "github.com/bboozzoo/q3stats/assets/static"
	atemplates "github.com/bboozzoo/q3stats/assets/templates"
	"github.com/bboozzoo/q3stats/controllers"
	"github.com/bboozzoo/q3stats/controllers/match"
	"github.com/bboozzoo/q3stats/controllers/player"
	"github.com/bboozzoo/q3stats/handlers"
	"github.com/bboozzoo/q3stats/handlers/api"
	"github.com/bboozzoo/q3stats/handlers/site"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const (
	defaultListenPort = 9090

	uriApi    = "/api"
	uriStatic = "/static/"
	uriSite   = "/site/"
)

var (
	defaultListenAddr = fmt.Sprintf("localhost:%d",
		defaultListenPort)
)

type handlerRouting struct {
	prefix  string
	handler handlers.Handler
}

func setupHandlers(handlers []handlerRouting) {
	r := mux.NewRouter()

	for _, h := range handlers {
		subr := r.PathPrefix(h.prefix).Subrouter()
		h.handler.SetupHandlers(subr)
	}

	// redirect to site by default
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, uriSite, http.StatusFound)
	})

	filehandler := http.FileServer(astatic.FS(false))
	r.PathPrefix(uriStatic).
		Handler(http.StripPrefix(uriStatic, filehandler))

	// setup logging for all handlers
	lr := ghandlers.LoggingHandler(os.Stdout, r)

	http.Handle("/", lr)
}

func daemonMain() error {
	db := NewDB()
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	matchctrl := match.NewController(db)
	userctrl := player.NewController(db)
	api := api.NewApi(matchctrl)
	ctrls := controllers.Controllers{
		matchctrl,
		userctrl,
	}
	site := site.NewSite(ctrls, atemplates.FS(false))

	hrouting := []handlerRouting{
		{uriApi, api},
		{uriSite, site},
	}
	setupHandlers(hrouting)

	return http.ListenAndServe(fmt.Sprintf(":%d", C.Port), nil)
}

func runDaemon() error {
	log.Printf("listen port: %d", C.Port)

	return daemonMain()
}
