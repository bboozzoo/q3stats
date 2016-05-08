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
package api

import (
	"github.com/bboozzoo/q3stats/controllers/match"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	UriAddMatch = "/matches/new"
)

type Api struct {
	mc *match.MatchController
}

func NewApi(mc *match.MatchController) *Api {
	return &Api{mc}
}

func (a *Api) apiAddMatch(w http.ResponseWriter, req *http.Request) {
	// add new match
	log.Printf("add new match")

	matchdata, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("failed to receive match data: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(matchdata) == 0 {
		log.Printf("no data provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("match data: %s", matchdata)

	hash, err := a.mc.AddFromData(matchdata)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("match hash: %s", hash)
	w.Write([]byte(hash))
}

func (a *Api) SetupHandlers(r *mux.Router) {
	// matches only come through POST
	r.HandleFunc(UriAddMatch, a.apiAddMatch).
		Methods("POST")
}
