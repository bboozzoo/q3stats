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
	"github.com/bboozzoo/q3stats/models/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGlobalStats(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	NewMatch(store, Match{Map: "q3dm17", Duration: 10})
	NewMatch(store, Match{Map: "q3dm17", Duration: 100})
	mid := NewMatch(store, Match{Map: "q3dm4", Duration: 100})

	NewAliasOrCurrent(store, Alias{Alias: "foo"})
	abarid := NewAliasOrCurrent(store, Alias{Alias: "bar"})

	pmsid := NewPlayerMatchStat(store, PlayerMatchStat{
		Score:    10,
		Kills:    12,
		Suicides: 2,
		MatchID:  mid,
		AliasID:  abarid,
	})

	NewWeaponStat(store, WeaponStat{
		Type:              RocketLauncher,
		Shots:             100,
		Hits:              10,
		PlayerMatchStatID: pmsid,
	})

	gs := GetGlobalStats(store)

	assert.Equal(t, uint(12), gs.Frags)
	assert.Equal(t, uint(2), gs.Suicides)
	assert.Equal(t, uint(2), gs.Maps)
	assert.Equal(t, uint(3), gs.Matches)
	assert.Equal(t, uint(100), gs.RocketsLaunched)
	assert.Equal(t, uint(2), gs.Players)
}
