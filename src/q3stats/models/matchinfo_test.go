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
	"q3stats/models/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetMatchinfo(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	m := GetMatchInfo(store, "fooo")
	assert.Nil(t, m)

	// start with match
	mid := NewMatch(store, Match{
		DataHash: "foobar",
		Map:      "q3dm17",
		Duration: 100,
	})

	// alias foo
	aidfoo := NewAliasOrCurrent(store, Alias{
		Alias: "foo",
	})
	// match stats for foo
	pmsfoo := NewPlayerMatchStat(store, PlayerMatchStat{
		Score:    1,
		Kills:    2,
		Suicides: 1,
		AliasID:  aidfoo,
		MatchID:  mid,
	})
	// weapon stats for foo's stats in match
	NewWeaponStat(store, WeaponStat{
		Type:              RocketLauncher,
		Shots:             1,
		PlayerMatchStatID: pmsfoo,
	})
	// item stats for foo's stats in match
	NewItemStat(store, ItemStat{
		Type:              MegaHealth,
		Pickups:           1,
		PlayerMatchStatID: pmsfoo,
	})

	aidbar := NewAliasOrCurrent(store, Alias{
		Alias: "bar",
	})
	pmsbar := NewPlayerMatchStat(store, PlayerMatchStat{
		Score:    -2,
		Kills:    2,
		Suicides: 1,
		AliasID:  aidbar,
		MatchID:  mid,
	})
	NewWeaponStat(store, WeaponStat{
		Type:              Plasmagun,
		Shots:             10,
		PlayerMatchStatID: pmsbar,
	})

	mi := GetMatchInfo(store, "foobar")
	assert.NotNil(t, mi)
	assert.Equal(t, "q3dm17", mi.Match.Map)
	assert.Equal(t, "foobar", mi.Match.DataHash)
	assert.Equal(t, uint(100), mi.Match.Duration)
	assert.Len(t, mi.PlayerData, 2)

	for _, p := range mi.PlayerData {
		if p.Alias.Alias == "foo" {
			assert.Equal(t, 1, p.Stats.Score)
			assert.Equal(t, aidfoo, p.Stats.AliasID)
			assert.Equal(t, mid, p.Stats.MatchID)

			assert.Len(t, p.Weapons, 1)
			assert.Len(t, p.Items, 1)

			w := p.Weapons[0]
			assert.Equal(t, RocketLauncher, w.Type)
			assert.Equal(t, uint(1), w.Shots)
			assert.Equal(t, pmsfoo, w.PlayerMatchStatID)

			i := p.Items[0]
			assert.Equal(t, MegaHealth, i.Type)
			assert.Equal(t, uint(1), i.Pickups)
			assert.Equal(t, pmsfoo, i.PlayerMatchStatID)
		} else if p.Alias.Alias == "bar" {
			assert.Len(t, p.Weapons, 1)
			assert.Len(t, p.Items, 0)
		} else {
			assert.Fail(t, "unexpected alias: %s",
				p.Alias.Alias)
		}
	}
}

func TestNewMatchInfo(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	ws := WeaponStat{
		Type:  RocketLauncher,
		Shots: 10,
		Hits:  10,
		Kills: 1,
	}
	is := ItemStat{
		Type:    MegaHealth,
		Pickups: 10,
	}
	pmsfoo := PlayerMatchStat{
		Kills:    10,
		Deaths:   10,
		Suicides: 1,
		Score:    -2,
	}
	pmsbar := PlayerMatchStat{
		Kills:    10,
		Deaths:   10,
		Suicides: 1,
		Score:    2,
	}

	matchtime := time.Now()
	t.Logf("time: %v", matchtime.Format(time.RFC822Z))
	m := MatchInfo{
		Match: Match{
			DataHash: "foobar",
			Map:      "q3dm17",
			Duration: 100,
			DateTime: matchtime,
		},
		PlayerData: []PlayerDataItem{
			{
				Alias{Alias: "foo"},
				[]WeaponStat{
					ws,
				},
				[]ItemStat{
					is,
				},
				pmsfoo,
			},
			{
				Alias{Alias: "bar"},
				[]WeaponStat{
					ws,
				},
				[]ItemStat{},
				pmsbar,
			},
		},
	}

	id, hash := NewMatchInfo(store, &m)

	assert.Equal(t, "foobar", hash, "should be equal")

	match := FindMatchByHash(store, "foobar")
	assert.NotNil(t, match)
	assert.Equal(t, id, match.ID, "different match ID")
	assert.Equal(t, "q3dm17", match.Map, "incorrect map")
	assert.Equal(t, uint(100), match.Duration, "incorrect duration")
	assert.Equal(t, matchtime.Format(time.RFC822Z),
		match.DateTime.Format(time.RFC822Z), "incorrect match datetime")

	aliases := GetAliases(store, NoUser)
	assert.Len(t, aliases, 2)

	var afooid, abarid uint
	for _, a := range aliases {
		if a.Alias != "foo" && a.Alias != "bar" {
			assert.Fail(t, "unexpected alias: %s", a.Alias)
		}
		if a.Alias == "foo" {
			afooid = a.ID
		} else if a.Alias == "bar" {
			abarid = a.ID
		}
	}

	// find player match stats for this match
	pms := ListPlayerMatchStat(store, id)
	assert.Len(t, pms, 2)

	for _, p := range pms {
		if p.AliasID == afooid {
			// player foo
			// 1 weapon stat
			// 1 item stats
			assert.Equal(t, pmsfoo.Kills, p.Kills)
			assert.Len(t, ListWeaponStats(store, p.ID), 1)
			assert.Len(t, ListItemStats(store, p.ID), 1)
		} else if p.AliasID == abarid {
			// player bar
			// 1 weapon stat
			// 0 item stats
			assert.Equal(t, pmsbar.Kills, p.Kills)
			assert.Len(t, ListWeaponStats(store, p.ID), 1)
			assert.Len(t, ListItemStats(store, p.ID), 0)
		}
	}

}
