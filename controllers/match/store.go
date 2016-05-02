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
	"github.com/bboozzoo/q3stats/loader"
	"github.com/bboozzoo/q3stats/models"
	"github.com/pkg/errors"
	"log"
	"time"
)

func (m *MatchController) storeMatch(rmatch *loader.RawMatch) error {

	match := models.Match{
		DataHash: rmatch.DataHash,
		Duration: rmatch.Duration,
		Map:      rmatch.Map,
		Type:     rmatch.Type,
	}

	tm, err := time.Parse("2006/01/02 15:04:05",
		rmatch.Datetime)
	if err != nil {
		return errors.Wrap(err, "failed to parse time")
	}
	match.DateTime = tm

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
				// fill duration, rawItem time is in ms
				Time: time.Duration(itm.Time) * time.Millisecond,
			})
		}
	}

	tx.Commit()

	return nil
}
