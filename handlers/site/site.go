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
package site

import (
	"github.com/bboozzoo/q3stats/controllers"
	"github.com/bboozzoo/q3stats/controllers/match"
	"github.com/bboozzoo/q3stats/controllers/player"
	"github.com/gorilla/mux"
	"path"
)

type Site struct {
	m    *match.MatchController
	p    *player.PlayerController
	tdir string
	r    *mux.Router
}

func NewSite(c controllers.Controllers, webroot string) *Site {
	return &Site{
		m:    c.M,
		p:    c.P,
		tdir: path.Join(webroot, "templates"),
	}
}

func (s *Site) SetupHandlers(r *mux.Router) {
	r.HandleFunc("/", s.siteHomeHandler).
		Methods("GET").Name("home")
	r.HandleFunc("/matches", s.matchesViewHandler).
		Methods("GET").Name("matches")
	r.HandleFunc("/matches/{id}", s.matchViewHandler).
		Methods("GET")
	r.HandleFunc("/players", s.playersViewHandler).
		Methods("GET").Name("players")
	r.HandleFunc("/players/new", s.newPlayerViewHandler).
		Methods("GET")
	r.HandleFunc("/players/new", s.createNewPlayerViewHandler).
		Methods("POST")
	r.HandleFunc("/players/{id}/edit", s.editPlayerViewHandler).
		Methods("GET").Name("playeredit")
	r.HandleFunc("/players/{id}/edit", s.editSubmitPlayerViewHandler).
		Methods("POST")
	r.HandleFunc("/players/{id}", s.playerViewHandler).
		Methods("GET").Name("player")

	// keep track of router
	s.r = r
}
