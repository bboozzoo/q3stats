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
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type matchViewWeaponStat struct {
	models.WeaponStat
}

type matchViewPlayerData struct {
	PlayerID     uint
	PlayerUrl    string
	Alias        string
	GeneralStats models.PlayerMatchStat
	Weapons      map[string]matchViewWeaponStat
	Items        map[string]models.ItemStat
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
		// data of all players
		Players []matchViewPlayerData
		// user presentable data descriptor
		Desc *models.DescData
		// match information
		Match models.Match
	}{
		make([]matchViewPlayerData, len(m.PlayerData)),
		models.GetDesc(),
		m.Match,
	}

	// repack to a format expected by the template
	for i, pd := range m.PlayerData {
		mvpd := matchViewPlayerData{
			pd.Alias.PlayerID,
			s.playerViewURL(pd.Alias.PlayerID),
			pd.Alias.Alias,
			pd.Stats,
			make(map[string]matchViewWeaponStat),
			make(map[string]models.ItemStat),
		}
		for _, wd := range pd.Weapons {
			mvpd.Weapons[wd.Type] = matchViewWeaponStat{
				wd,
			}
		}

		for _, itd := range pd.Items {
			mvpd.Items[itd.Type] = itd
		}
		data.Players[i] = mvpd
	}

	p := s.NewPage(catMatches, data)
	s.loadRenderOrError(w, p, "match.tmpl", "base.tmpl")
}
