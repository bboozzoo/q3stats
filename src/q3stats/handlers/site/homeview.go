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
	"q3stats/models"
	"log"
	"net/http"
)

func (s *Site) siteHomeHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("site / handler")

	// show only 10 last matches
	matches := s.m.ListMatches(models.MatchListParams{
		TimeSort: true,
		SortDesc: true,
		Limit:    10,
	})

	data := struct {
		RecentMatches []matchesViewMatchData
		Global        models.GlobalStats
	}{
		make([]matchesViewMatchData, len(matches)),
		models.GlobalStats{},
	}
	for i, m := range matches {
		data.RecentMatches[i] = matchesViewMatchData{
			m,
		}
	}

	data.Global = s.m.GetGlobalStats()

	p := s.NewPage(catHome, data)
	s.loadRenderOrError(w, p,
		"home.tmpl", "matchlist.tmpl", "base.tmpl")
}
