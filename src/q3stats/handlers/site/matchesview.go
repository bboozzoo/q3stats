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
	"github.com/bboozzoo/q3stats/models"
	"log"
	"net/http"
)

type matchesViewMatchData struct {
	models.Match
}

func (s *Site) matchesViewHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("site matches view handler")

	matches := s.m.ListMatches(models.MatchListParams{
		TimeSort: true,
		SortDesc: true,
	})

	data := struct {
		Matches []matchesViewMatchData
	}{
		make([]matchesViewMatchData, len(matches)),
	}
	for i, m := range matches {
		data.Matches[i] = matchesViewMatchData{
			m,
		}
	}

	p := s.NewPage(catMatches, data)
	s.loadRenderOrError(w, p,
		"matches.tmpl", "matchlist.tmpl", "base.tmpl")
}
