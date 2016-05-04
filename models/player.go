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

// Registered players
type Player struct {
	gorm.Model

	// user defined name, ex. joe
	Name string

	// password hash
	PasswordHash string
}

// create new player returning its ID
func NewPlayer(store store.DBConn, name string, passwordhash string) uint {
	player := Player{
		Name:         name,
		PasswordHash: passwordhash,
	}

	db := store.Conn()

	db.Create(&player)

	return player.ID
}

func HasPlayer(store store.DBConn, name string) bool {
	db := store.Conn()

	var player Player

	notfound := db.Where(&Player{Name: name}).
		First(&player).
		RecordNotFound()

	if notfound == true {
		return false
	}
	return true
}

func ListPlayers(store store.DBConn) []Player {
	db := store.Conn()

	var players []Player
	db.Find(&players)

	return players
}
