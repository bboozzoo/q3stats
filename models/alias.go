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

const (
	NoUser = 0
)

// Alias is a nickname used during the match.
type Alias struct {
	gorm.Model

	// alias, ex. 'honeybunny', 'Klesk'
	Alias string

	// player claiming this alias
	PlayerID uint
}

func GetAliases(store store.DB, user uint) []Alias {
	var aliases []Alias
	store.Conn().
		Model(&Alias{}).
		Where("player_id = ?", user).
		Find(&aliases)

	return aliases
}

func ClaimAliasesByPlayer(store store.DB, player uint, aliases []string) {
	// update aliases set player_id = `player` where aliases in
	// `aliases`?

	db := store.Conn()

	for _, a := range aliases {
		db.Model(&Alias{}).
			Where(&Alias{Alias: a}).
			Update("player_id", player)
	}
}
