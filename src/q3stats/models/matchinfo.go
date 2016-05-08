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
package models

import (
	"github.com/bboozzoo/q3stats/store"
	"log"
)

// aggregate player data from particular match
type PlayerDataItem struct {
	// use Alias to access PlayerID
	Alias   Alias
	Weapons []WeaponStat
	Items   []ItemStat
	Stats   PlayerMatchStat
}

// aggregate match info
type MatchInfo struct {
	Match Match
	// stats of all players who attended this match
	PlayerData []PlayerDataItem
}

// load user provided matchinfo into store, returns match ID and hash
func NewMatchInfo(store store.DB, m *MatchInfo) (uint, string) {
	tx := store.Begin()

	mid := NewMatch(tx, m.Match)

	log.Printf("got match: %u", mid)

	for _, player := range m.PlayerData {
		aid := NewAliasOrCurrent(tx, player.Alias)

		log.Printf("got alias: %u", aid)

		player.Stats.AliasID = aid
		player.Stats.MatchID = mid
		pmsid := NewPlayerMatchStat(tx, player.Stats)

		log.Printf("created PMS: %u", pmsid)

		for _, wep := range player.Weapons {
			wep.PlayerMatchStatID = pmsid
			NewWeaponStat(tx, wep)
		}

		for _, itm := range player.Items {
			itm.PlayerMatchStatID = pmsid
			NewItemStat(tx, itm)
		}
	}

	tx.Commit()

	return mid, m.Match.DataHash
}

// load match data and produce a match info data wrapper, returns nil
// if match was not found
func GetMatchInfo(store store.DB, hash string) *MatchInfo {
	db := store.Conn()

	var match Match
	// grab match data first
	notfound := db.Where(&Match{DataHash: hash}).
		First(&match).
		RecordNotFound()
	if notfound == true {
		return nil
	}

	// filling up
	mi := MatchInfo{}

	// already know basic match data
	mi.Match = match

	// locate all players in this match
	pls := ListPlayerMatchStat(store, match.ID)

	mi.PlayerData = make([]PlayerDataItem, len(pls))
	// start populating PlayerData
	for i, p := range pls {
		log.Printf("processing player %v", p.ID)
		// already have player stats
		mi.PlayerData[i].Stats = p

		// fill alias info for this player
		mi.PlayerData[i].Alias = *GetAlias(store, p.AliasID)
		log.Printf("   found alias: %s", mi.PlayerData[i].Alias.Alias)

		// fill weapon statistics
		mi.PlayerData[i].Weapons = ListWeaponStats(store, p.ID)
		// fill items
		mi.PlayerData[i].Items = ListItemStats(store, p.ID)
	}

	return &mi

}
