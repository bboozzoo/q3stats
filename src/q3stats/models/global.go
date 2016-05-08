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
	"github.com/jinzhu/gorm"
	"log"
	"q3stats/store"
)

type GlobalStats struct {
	// total number of frags/kills
	Frags uint
	// total number of suicides
	Suicides uint
	// total number of matches
	Matches uint
	// total number of unique maps
	Maps uint
	// number of rockets launched
	RocketsLaunched uint
	// players
	Players uint
}

// obtain kills/suicidees/deaths counts
func getKillSuicideDeathCount(db *gorm.DB) (uint, uint, uint) {
	var result = struct {
		Kills    uint
		Deaths   uint
		Suicides uint
	}{}
	db.Model(&PlayerMatchStat{}).
		Select("sum(kills) as kills, sum(deaths) as deaths, sum(suicides)as suicides").
		Scan(&result)
	log.Printf("kills: %u deaths: %u suicides: %u",
		result.Kills, result.Deaths, result.Suicides)

	return result.Kills, result.Deaths, result.Suicides
}

func getUniqueMapsCount(db *gorm.DB) uint {
	var result = struct {
		Maps uint
	}{}
	db.Raw("select count(map) as maps from (select distinct map from matches)").
		Scan(&result)
	log.Printf("unique maps: %u", result.Maps)
	return result.Maps
}

func getUniquePlayersCount(db *gorm.DB) uint {
	var count uint
	db.Model(&Alias{}).Count(&count)
	log.Printf("unique players: %u", count)
	return count
}

func getRocketsLaunchedCount(db *gorm.DB) uint {
	var result = struct {
		Launched uint
	}{}

	db.Model(&WeaponStat{}).
		Where(&WeaponStat{Type: RocketLauncher}).
		Select("sum(shots) as launched").
		Scan(&result)
	log.Printf("launched rockets: %u", result.Launched)

	return result.Launched
}

func getMatchCount(db *gorm.DB) uint {
	var count uint
	db.Model(&Match{}).Count(&count)
	return count
}

func GetGlobalStats(store store.DB) GlobalStats {
	db := store.Conn()

	gs := GlobalStats{}

	gs.Matches = getMatchCount(db)
	k, _, s := getKillSuicideDeathCount(db)

	gs.Frags = k
	gs.Suicides = s
	gs.Players = getUniquePlayersCount(db)
	gs.Maps = getUniqueMapsCount(db)
	gs.RocketsLaunched = getRocketsLaunchedCount(db)

	return gs

}
