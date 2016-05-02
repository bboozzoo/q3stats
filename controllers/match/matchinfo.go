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
package match

import (
	"github.com/bboozzoo/q3stats/models"
	"log"
)

func (m *MatchController) GetMatchInfo(id string) *models.MatchInfo {

	log.Printf("get match infor for %s", id)

	db := m.db.Conn()

	var match models.Match
	// grab match data first
	notfound := db.Where(&models.Match{DataHash: id}).
		First(&match).
		RecordNotFound()
	if notfound == true {
		return nil
	}

	// filling up
	mi := models.MatchInfo{}

	// already know map
	mi.Match = match

	// locate all players in this match
	var pls []models.PlayerMatchStat
	db.Where(&models.PlayerMatchStat{MatchID: match.ID}).
		Find(&pls)
	log.Printf("found %u players", len(pls))

	mi.PlayerData = make([]models.PlayerDataItem, len(pls))
	// start populating PlayerData
	for i, p := range pls {
		log.Printf("processing player %u", p.ID)
		// already have player stats
		mi.PlayerData[i].Stats = p

		// fill alias info for this player
		db.First(&mi.PlayerData[i].Alias, p.AliasID)
		log.Printf("   found alias: %s", mi.PlayerData[i].Alias.Alias)

		// fill weapon statistics
		db.Where(&models.WeaponStat{
			PlayerMatchStatID: p.ID,
		}).Find(&mi.PlayerData[i].Weapons)

		// fill items
		db.Where(&models.ItemStat{
			PlayerMatchStatID: p.ID,
		}).Find(&mi.PlayerData[i].Items)
	}

	return &mi
}
