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
	"github.com/jinzhu/gorm"
	"log"
)

type PlayerGlobalStats struct {
	Kills         uint
	Deaths        uint
	Suicides      uint
	Matches       uint
	KillsPerMatch float32
	// frags per life
	Efficiency float32
	Wins       uint
	// global weapon stats
	Weapons []WeaponStat
	// global item stats
	Items []ItemStat
}

func getPlayerKillDeathSuicide(db *gorm.DB, player uint) (uint, uint, uint) {
	var result = struct {
		Kills    uint
		Deaths   uint
		Suicides uint
	}{}
	db.Model(&PlayerMatchStat{}).
		Select("sum(kills) as kills, sum(deaths) as deaths, sum(suicides)as suicides").
		Joins("join aliases on player_match_stats.alias_id = aliases.id").
		Where("aliases.player_id = ?", player).
		Scan(&result)
	log.Printf("stats: %+v", result)

	return result.Kills, result.Deaths, result.Suicides
}

func getPlayerMatches(db *gorm.DB, player uint) uint {
	var result = struct {
		Matches uint
	}{}
	db.Model(&PlayerMatchStat{}).
		Select("count(distinct player_match_stats.match_id) as matches").
		Joins("join aliases on player_match_stats.alias_id = aliases.id").
		Where("aliases.player_id = ?", player).
		Scan(&result)
	log.Printf("stats: %+v", result)

	return result.Matches
}

func GetPlayerGlobaStats(store store.DB, player uint) *PlayerGlobalStats {
	db := store.Conn()

	pgs := PlayerGlobalStats{}
	pgs.Kills, pgs.Deaths, pgs.Suicides = getPlayerKillDeathSuicide(db,
		player)
	pgs.Matches = getPlayerMatches(db, player)
	if pgs.Matches != 0 {
		pgs.KillsPerMatch = float32(pgs.Kills) / float32(pgs.Matches)
	}

	if (pgs.Deaths + pgs.Suicides) != 0 {
		pgs.Efficiency = float32(pgs.Kills) / float32(pgs.Deaths+pgs.Suicides)
	}
	return &pgs
}
