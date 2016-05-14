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
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Site) loadRenderOrError(w http.ResponseWriter, data interface{},
	templates ...string) {

	mt := s.loadTemplates(templates...)
	// got match info
	err := renderTemplate(w, mt, data)
	if err != nil {
		log.Printf("failed to execute template: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Site) loadTemplates(names ...string) *template.Template {
	t := template.New("site")
	for _, n := range names {
		//
		f, err := s.fs.Open("/" + n)
		if err != nil {
			log.Printf("failed to load template file %s: %s",
				n, err)
			return nil
		}

		tdata, _ := ioutil.ReadAll(f)
		f.Close()
		_, err = t.Parse(string(tdata))
		if err != nil {
			log.Printf("failed to parse template from file %s: %s",
				n, err)
		}
	}

	return t
}

func renderTemplate(w http.ResponseWriter, t *template.Template,
	data interface{}) error {

	// set header?
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	return t.ExecuteTemplate(w, "base", data)
}
