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
	"strconv"
)

// wrap weapon stat
type playerViewWeapon struct {
	models.WeaponStat
	// weapon description
	Desc models.ItemDesc
}

// wrap item stat
type playerViewItem struct {
	models.ItemStat
	// item descripion
	Desc models.ItemDesc
}

func (s *Site) playerViewHandler(w http.ResponseWriter, req *http.Request) {
	// player view handler
	log.Printf("player view handler")
	id := mux.Vars(req)["id"]

	log.Printf("player ID: %s", id)

	pid := strToUint(id)

	pl := s.p.GetPlayer(pid)
	pgs := s.p.GetPlayerGlobaStats(pid)

	data := struct {
		// palayer global stats
		PlayerGlobalStats *models.PlayerGlobalStats
		Name              string
		Weapons           map[string]playerViewWeapon
		Items             map[string]playerViewItem
	}{
		pgs,
		pl.Name,
		make(map[string]playerViewWeapon),
		make(map[string]playerViewItem),
	}
	for _, w := range pgs.Weapons {
		data.Weapons[w.Type] = playerViewWeapon{
			w,
			models.GetItemDesc(w.Type),
		}
	}
	for _, i := range pgs.Items {
		data.Items[i.Type] = playerViewItem{
			i,
			models.GetItemDesc(i.Type),
		}
	}

	s.loadRenderOrError(w, data, "player.tmpl", "base.tmpl")
}

func (s *Site) playerViewURL(id uint) string {
	u, _ := s.r.Get("player").URL("id", strconv.FormatUint(uint64(id), 10))
	return u.String()
}
