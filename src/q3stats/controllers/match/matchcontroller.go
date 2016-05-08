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
	"q3stats/loader"
	"q3stats/models"
	"q3stats/store"
	"github.com/pkg/errors"
	"log"
)

type MatchController struct {
	db store.DB
}

func NewController(db store.DB) *MatchController {

	// move this somewhere else?
	models.CreateSchema(db)

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

	match := models.FindMatchByHash(m.db, rmatch.DataHash)
	if match == nil {
		return false
	}
	return true
}
