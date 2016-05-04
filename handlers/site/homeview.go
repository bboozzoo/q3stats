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

// wrapper for player data with URL to player's profile
type homeViewPlayerData struct {
	models.Player
	URL string
}

func (s *Site) siteHomeHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("site / handler")

	// show only 10 last matches
	matches := s.m.ListMatches(models.MatchListParams{
		TimeSort: true,
		SortDesc: true,
		Limit:    10,
	})

	players := s.p.ListPlayers()

	data := struct {
		RecentMatches []matchesViewMatchData
		Global        models.GlobalStats
		Players       []homeViewPlayerData
	}{
		make([]matchesViewMatchData, len(matches)),
		models.GlobalStats{},
		make([]homeViewPlayerData, len(players)),
	}
	for i, m := range matches {
		data.RecentMatches[i] = matchesViewMatchData{
			m,
		}
	}

	for i, p := range players {
		data.Players[i] = homeViewPlayerData{
			p,
			s.playerViewURL(p.ID),
		}
	}

	data.Global = s.m.GetGlobalStats()

	s.loadRenderOrError(w, data,
		"home.tmpl", "matchlist.tmpl", "base.tmpl")
}
