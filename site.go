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
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"path"
)

type Site struct {
	m    *MatchController
	tdir string
}

func NewSite(m *MatchController) *Site {
	return &Site{
		m:    m,
		tdir: path.Join(C.Webroot, "templates"),
	}
}

func (s *Site) SetupHandlers(r *mux.Router) {
	r.HandleFunc(uriIndex, s.siteHomeHandler).
		Methods("GET")
}

func (s *Site) siteHomeHandler(w http.ResponseWriter, req *http.Request) {

	log.Printf("site home handler")

	t := s.loadTemplate("base.tmpl")
	err := t.Execute(w, nil)
	if err != nil {
		log.Printf("failed to execute template: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Site) loadTemplate(name string) *template.Template {
	tpath := path.Join(s.tdir, "base.tmpl")
	log.Printf("template path: %s", tpath)
	// homepage
	t, err := template.ParseFiles(tpath)
	if err != nil {
		log.Printf("failed to parse template: %s", err)
		return nil
	}
	return t
}
