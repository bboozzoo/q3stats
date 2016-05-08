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
	"q3stats/store"
	"github.com/jinzhu/gorm"
)

// Overall statistics achieved by player during a match
type PlayerMatchStat struct {
	gorm.Model

	// score
	Score int
	// kills count
	Kills int
	// number of deaths
	Deaths int
	// self-kills count
	Suicides int
	// net (kills - deaths - suicides)
	Net int
	// damage given in HP
	DamageGiven int
	// damage taken in HP
	DamageTaken int
	// total health taken in HP
	HealthTotal int
	// total armor points
	ArmorTotal int

	// alias used
	AliasID uint

	// match
	MatchID uint
}

func NewPlayerMatchStat(store store.DBConn, pms PlayerMatchStat) uint {
	db := store.Conn()

	db.Create(&pms)

	return pms.ID
}

func ListPlayerMatchStat(store store.DBConn, matchID uint) []PlayerMatchStat {
	db := store.Conn()

	// locate all players in this match
	var pls []PlayerMatchStat

	db.Where(&PlayerMatchStat{MatchID: matchID}).
		Find(&pls)

	return pls
}
