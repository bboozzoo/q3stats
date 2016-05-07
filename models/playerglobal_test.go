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
	"time"
)

func TestGetPlayerGlobalStats(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	m1 := NewMatch(store, Match{Map: "q3dm17", Duration: 10})
	m2 := NewMatch(store, Match{Map: "q3dm4", Duration: 100})

	a1 := NewAliasOrCurrent(store, Alias{Alias: "foo"})
	p1 := NewPlayer(store, "fooplayer", "")

	// claim aliases
	ClaimAliasesByPlayer(store, p1, []string{"foo"})

	pms1 := NewPlayerMatchStat(store, PlayerMatchStat{
		Score:    10,
		Kills:    12,
		Deaths:   33,
		Suicides: 2,
		MatchID:  m1,
		AliasID:  a1,
	})
	NewWeaponStat(store, WeaponStat{
		Type:              RocketLauncher,
		Shots:             100,
		Hits:              10,
		PlayerMatchStatID: pms1,
	})

	NewWeaponStat(store, WeaponStat{
		Type:              GrenadeLauncher,
		Shots:             25,
		Hits:              5,
		PlayerMatchStatID: pms1,
	})

	NewItemStat(store, ItemStat{
		Type:              MegaHealth,
		Pickups:           12,
		PlayerMatchStatID: pms1,
	})

	NewItemStat(store, ItemStat{
		Type:              RedArmor,
		Pickups:           5,
		PlayerMatchStatID: pms1,
	})
	NewItemStat(store, ItemStat{
		Type:              QuadDamage,
		Pickups:           5,
		Time:              time.Duration(30) * time.Second,
		PlayerMatchStatID: pms1,
	})

	pms2 := NewPlayerMatchStat(store, PlayerMatchStat{
		Score:    5,
		Kills:    5,
		Deaths:   3,
		Suicides: 2,
		MatchID:  m2,
		AliasID:  a1,
	})
	NewWeaponStat(store, WeaponStat{
		Type:              RocketLauncher,
		Shots:             10,
		Hits:              5,
		PlayerMatchStatID: pms2,
	})
	NewWeaponStat(store, WeaponStat{
		Type:              Shotgun,
		Shots:             10,
		Hits:              5,
		PlayerMatchStatID: pms2,
	})
	NewItemStat(store, ItemStat{
		Type:              GreenArmor,
		Pickups:           3,
		PlayerMatchStatID: pms2,
	})
	NewItemStat(store, ItemStat{
		Type:              QuadDamage,
		Pickups:           1,
		Time:              time.Duration(5) * time.Second,
		PlayerMatchStatID: pms2,
	})

	pgs := GetPlayerGlobaStats(store, p1)

	assert.Equal(t, uint(12+5), pgs.Kills)
	assert.Equal(t, uint(2+2), pgs.Suicides)
	assert.Equal(t, uint(33+3), pgs.Deaths)
	assert.Equal(t, uint(2), pgs.Matches)

	assert.Len(t, pgs.Weapons, 3)
	assert.Len(t, pgs.Items, 4)

	// these should be aggregate values from all weapons
	for _, w := range pgs.Weapons {
		switch w.Type {
		case RocketLauncher:
			assert.Equal(t, 110, w.Shots)
			assert.Equal(t, 15, w.Hits)
		case GrenadeLauncher:
			assert.Equal(t, 25, w.Shots)
		case Shotgun:
			assert.Equal(t, 10, w.Shots)
			assert.Equal(t, 5, w.Hits)
		default:
			assert.Fail(t, "unexpected weapon type: %s", w.Type)
		}
	}

	// aggregate values for all items
	for _, i := range pgs.Items {
		switch i.Type {
		case MegaHealth:
			assert.Equal(t, 12, i.Pickups)
		case RedArmor:
			assert.Equal(t, 5, i.Pickups)
		case GreenArmor:
			assert.Equal(t, 3, i.Pickups)
		case QuadDamage:
			assert.Equal(t, 5+1, i.Pickups)
			assert.Equal(t, time.Duration(30+5)*time.Second,
				i.Time)
		}
	}
}
