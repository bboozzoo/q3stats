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
	"time"
)

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
	PlayerMatchStatID uint
}

// true if effect of given item is time bounded, such as Quad Damage,
// Regeneration, Battle Suit
func (i ItemStat) HasDuration() bool {
	return ItemHasDuration(i.Type)
}

func (i ItemStat) DurationDesc() string {
	return (i.Time).String()
}

// create new item stat and return its ID
func NewItemStat(store store.DBConn, is ItemStat) uint {
	db := store.Conn()

	db.Create(&is)

	return is.ID
}

// list weapon statistics for given player match stat ID
func ListItemStats(store store.DBConn, pmsID uint) []ItemStat {
	db := store.Conn()

	var is []ItemStat
	db.Where(&ItemStat{
		PlayerMatchStatID: pmsID,
	}).Find(&is)

	return is
}
