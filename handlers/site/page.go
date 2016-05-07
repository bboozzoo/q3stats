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
// "github.com/gorilla/mux"
)

type ActiveMap map[string]bool
type UrlsMap map[string]string

type Page struct {
	// active page
	Active ActiveMap
	// main urls
	Urls UrlsMap

	ContentData interface{}
}

const (
	catHome    = "Home"
	catPlayers = "Players"
	catMatches = "Matches"
)

var (
	categories = []string{
		catHome,
		catPlayers,
		catMatches,
	}
)

func (s *Site) NewPage(category string, content interface{}) Page {
	p := Page{
		Active: make(ActiveMap),
		Urls:   make(UrlsMap),
	}

	for _, c := range categories {
		if category == c {
			p.Active[c] = true
		} else {
			p.Active[c] = false
		}
	}

	p.Urls[catHome] = s.getSimpleURL("home")
	// p.Urls[catPlayers] = s.getSimpleURL("players")
	p.Urls[catMatches] = s.getSimpleURL("matches")

	p.ContentData = content

	return p
}

func (s *Site) getSimpleURL(name string) string {
	u, _ := s.r.Get(name).URL()
	return u.String()
}
