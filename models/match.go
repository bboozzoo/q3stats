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
	"time"
)

// Match data
type Match struct {
	gorm.Model

	// match data hash
	DataHash string `sql:"not null,unique"`
	// when the match was played
	DateTime time.Time
	// duration in seconds
	Duration uint
	// map, ex. q3dm17
	Map string
	// type of match, 1v1, FFA, CTF
	Type string
}

// Alias is a nickname used during the match.
type Alias struct {
	gorm.Model

	// alias, ex. 'honeybunny', 'Klesk'
	Alias string

	// player claiming this alias
	PlayerID int
}

// Registered players
type Player struct {
	gorm.Model

	// user defined name, ex. joe
	Name string
}

// Per weapon statistics achieved during the match
type WeaponStat struct {
	gorm.Model

	// weapon type
	Type string
	// number of hits
	Hits uint
	// number of shots
	Shots uint
	// number of kills
	Kills uint

	// player's match stats this weapon is part of
	PlayerMatchStatID int
}

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
	AliasID int
}

// Per item statistics achieved during the match
type ItemStat struct {
	gorm.Model

	// item type (RA, YA, GA, MH)
	Type string
	// pickups count
	Pickups uint
	// time held (when applicable, ex. for quad)
	Time time.Duration

	// player's match stats this item is part of
	PlayerMatchStatID int
}
