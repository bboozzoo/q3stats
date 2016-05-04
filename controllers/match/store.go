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

	mid := models.NewMatch(m.db, match)

	log.Printf("got match: %u", mid)

	for _, player := range rmatch.Players {
		aid := models.NewAliasOrCurrent(m.db,
			models.Alias{Alias: player.Name})

		log.Printf("got alias: %u", aid)

		pms := makePlayerMatchStat(player.Stats)
		pms.MatchID = mid
		pms.AliasID = aid

		pmsid := models.NewPlayerMatchStat(m.db, pms)

		log.Printf("created PMS: %u", pmsid)

		for _, wep := range player.Weapons {
			models.NewWeaponStat(m.db,
				models.WeaponStat{
					Type:              wep.Name,
					Hits:              wep.Hits,
					Shots:             wep.Shots,
					Kills:             wep.Kills,
					PlayerMatchStatID: pmsid,
				})
		}

		// join items and powerups
		itms := []loader.RawItem{}
		itms = append(itms, player.Items...)
		itms = append(itms, player.Powerups...)

		for _, itm := range itms {
			models.NewItemStat(m.db,
				models.ItemStat{
					Type:              itm.Name,
					Pickups:           itm.Pickups,
					PlayerMatchStatID: pmsid,
					// fill duration, rawItem time is in ms
					Time: time.Duration(itm.Time) * time.Millisecond,
				})
		}
	}

	return nil
}
