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

func (m Match) DurationDesc() string {
	return (time.Duration(m.Duration) * time.Second).String()
}

func (m Match) NiceDateTime() string {
	return m.DateTime.Format("02-01-2006 15:04:05")
}

func FindMatchByHash(store store.DB, hash string) *Match {

	db := store.Conn()

	var mfound Match
	notfound := db.Where("data_hash = ?", hash).
		Find(&mfound).
		RecordNotFound()
	if notfound == true {
		return nil
	}
	return &mfound
}
