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
package player

import (
	"github.com/bboozzoo/q3stats/models"
	"log"
)

func (p *PlayerController) ListUnclaimedAliases() []models.Alias {
	return models.GetAliases(p.db, models.NoUser)
}

func (p *PlayerController) ListPlayerAliases(player uint) []models.Alias {
	return models.GetAliases(p.db, player)
}

// update aliases so that the player owns only the aliases passed as
// `aliases` parameter
func (p *PlayerController) UpdatePlayerAliases(player uint, aliases []string) error {
	// repack aliases to map, value true indicates that alias was
	// acted upon
	am := make(map[string]bool)
	for _, a := range aliases {
		am[a] = false
	}

	// grab the list of currently owned aliases
	current := models.GetAliases(p.db, player)

	to_claim := make([]string, 0)
	to_unclaim := make([]string, 0)

	for _, a := range current {
		// if currently owned alias is in map, things don't change
		_, ok := am[a.Alias]
		if ok == false {
			// but if it's not, we unclaim it
			to_unclaim = append(to_unclaim, a.Alias)
		}
		// alias already looked at
		am[a.Alias] = true
	}

	// scan remaining aliases, if alias was not acted upon, means
	// we must claim it
	for k, v := range am {
		if v == false {
			// alias was not acted upon, since we're
			// passing a list of aliases that a player
			// should own, it means that this alias will
			// be claimed by this player
			to_claim = append(to_claim, k)
		}
	}

	// claim new ones
	log.Printf("claim new aliases: %v", to_claim)
	models.ClaimAliasesByPlayer(p.db, player, to_claim)
	// clear old ones
	log.Printf("release aliases: %v", to_unclaim)
	models.ClaimAliasesByPlayer(p.db, models.NoUser, to_unclaim)

	return nil
}
