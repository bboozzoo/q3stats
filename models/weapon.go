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
)

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
	PlayerMatchStatID uint
}

func (w WeaponStat) Accuracy() uint {
	var acc uint = 0
	if w.Shots != 0 {
		acc = (100 * w.Hits) / w.Shots
	}
	return acc
}

// create new weapon stat and return its ID
func NewWeaponStat(store store.DBConn, ws WeaponStat) uint {
	db := store.Conn()

	db.Create(&ws)

	return ws.ID
}

// list weapon statistics for given player match stat ID
func ListWeaponStats(store store.DBConn, pmsID uint) []WeaponStat {
	db := store.Conn()

	var ws []WeaponStat
	db.Where(&WeaponStat{
		PlayerMatchStatID: pmsID,
	}).Find(&ws)

	return ws
}
