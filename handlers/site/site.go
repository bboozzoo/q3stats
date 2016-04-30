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
	"github.com/bboozzoo/q3stats/controllers/match"
	"github.com/bboozzoo/q3stats/models"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path"
)

type Site struct {
	m    *match.MatchController
	tdir string
	r    *mux.Router
}

func NewSite(m *match.MatchController, webroot string) *Site {
	return &Site{
		m:    m,
		tdir: path.Join(webroot, "templates"),
	}
}

func (s *Site) SetupHandlers(r *mux.Router) {
	r.HandleFunc("/", s.siteHomeHandler).
		Methods("GET")
	r.HandleFunc("/matches", s.matchesViewHandler).
		Methods("GET").Name("matches")
	r.HandleFunc("/matches/{id}", s.matchViewHandler).
		Methods("GET")

	// keep track of router
	s.r = r
}

func (s *Site) siteHomeHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("site / handler")

	url, _ := s.r.Get("matches").URL()
	http.Redirect(w, req, url.String(), http.StatusFound)
}

func (s *Site) matchesViewHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("site matches view handler")

	matches := s.m.ListMatches()
	data := struct {
		Matches []models.Match
	}{
		matches,
	}

	s.loadRenderOrError(w, data, "matches.tmpl", "base.tmpl")
}

func (s *Site) matchViewHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("site match view handler")
	id := mux.Vars(req)["id"]

	log.Printf("match ID: %s", id)

	m := s.m.GetMatchInfo(id)
	if m == nil {
		http.Error(w, "match not found", http.StatusNotFound)
	}

	data := struct {
		models.MatchInfo
	}{
		*m,
	}

	s.loadRenderOrError(w, data, "match.tmpl", "base.tmpl")

}

func (s *Site) loadRenderOrError(w http.ResponseWriter, data interface{},
	templates ...string) {

	mt := s.loadTemplates("match.tmpl", "base.tmpl")
	// got match info
	err := renderTemplate(w, mt, data)
	if err != nil {
		log.Printf("failed to execute template: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Site) loadTemplates(names ...string) *template.Template {
	paths := make([]string, 0, len(names))
	for _, n := range names {
		tpath := path.Join(s.tdir, n)
		paths = append(paths, tpath)
	}

	log.Printf("loading templates: %s", paths)

	t, err := template.ParseFiles(paths...)
	if err != nil {
		log.Printf("failed to parse templates: %s", err)
		return nil
	}
	return t
}

func renderTemplate(w http.ResponseWriter, t *template.Template,
	data interface{}) error {

	// set header?
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	return t.ExecuteTemplate(w, "base", data)
}
