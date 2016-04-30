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
	"fmt"
	"github.com/bboozzoo/q3stats/loader"
	"github.com/bboozzoo/q3stats/models"
	"github.com/bboozzoo/q3stats/store"
	"github.com/pkg/errors"
	"log"
	"time"
)

type MatchController struct {
	db store.DB
}

func NewController(db store.DB) *MatchController {

	// move this somewhere else?
	models.CreateSchema(db.Conn())

	return &MatchController{db}
}

func makePlayerMatchStat(stats []loader.RawStat) models.PlayerMatchStat {
	var pms models.PlayerMatchStat
	mapping := map[string]*int{
		"Score":       &pms.Score,
		"Kills":       &pms.Kills,
		"Deaths":      &pms.Deaths,
		"Suicides":    &pms.Suicides,
		"Net":         &pms.Net,
		"DamageGiven": &pms.DamageGiven,
		"DamageTaken": &pms.DamageTaken,
		"HealthTotal": &pms.HealthTotal,
		"ArmorTotal":  &pms.ArmorTotal,
	}

	for _, stat := range stats {
		ival, ok := mapping[stat.Name]
		if ok == false {
			log.Printf("skipping unknown stat: %s: %d",
				stat.Name, stat.Value)
			continue
		}

		*ival = stat.Value
	}
	return pms
}

// Add match from raw XML data in `data`. Returns match hash or error
func (m *MatchController) AddFromData(data []byte) (string, error) {
	rmatch, err := loader.LoadMatchData(data)
	if err != nil {
		return "", errors.Wrap(err, "failed to load match data")
	}

	if m.isMatchIndexed(rmatch) == true {
		return "", fmt.Errorf("already present")
	}

	if err := m.storeMatch(rmatch); err != nil {
		return "", errors.Wrap(err, "failed to store match data")
	}

	return rmatch.DataHash, nil
}

func (m *MatchController) isMatchIndexed(rmatch *loader.RawMatch) bool {
	db := m.db.Conn()

	var mfound models.Match
	notfound := db.Where("data_hash = ?", rmatch.DataHash).
		Find(&mfound).
		RecordNotFound()
	if notfound == false {
		log.Printf("record found: %+v", mfound)
		return true
	}

	return false
}

func (m *MatchController) storeMatch(rmatch *loader.RawMatch) error {

	match := models.Match{
		DataHash: rmatch.DataHash,
		DateTime: time.Now(),
		Duration: rmatch.Duration,
		Map:      rmatch.Map,
		Type:     rmatch.Type,
	}

	tx := m.db.Conn().Begin()

	tx.Create(&match)

	log.Printf("got match: %+v", match)

	for _, player := range rmatch.Players {
		var alias models.Alias
		tx.FirstOrCreate(&alias, models.Alias{Alias: player.Name})
		log.Printf("got alias: %+v", alias)

		pms := makePlayerMatchStat(player.Stats)
		pms.MatchID = match.ID
		pms.AliasID = alias.ID

		tx.Create(&pms)

		log.Printf("created PMS: %+v", pms)

		for _, wep := range player.Weapons {
			tx.Create(&models.WeaponStat{
				Type:              wep.Name,
				Hits:              wep.Hits,
				Shots:             wep.Shots,
				Kills:             wep.Kills,
				PlayerMatchStatID: pms.ID,
			})
		}

		// join items and powerups
		itms := []loader.RawItem{}
		itms = append(itms, player.Items...)
		itms = append(itms, player.Powerups...)

		for _, itm := range itms {
			tx.Create(&models.ItemStat{
				Type:              itm.Name,
				Pickups:           itm.Pickups,
				PlayerMatchStatID: pms.ID,
				// fill duration
			})
		}
	}

	tx.Commit()

	return nil
}

func (m *MatchController) ListMatches() []models.Match {
	db := m.db.Conn()

	var matches []models.Match
	db.Find(&matches)

	log.Printf("matches found: %u", len(matches))

	return matches
}

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
	mi.Map = match.Map

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
