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
	"q3stats/models/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWeaponAccuracy(t *testing.T) {
	w := WeaponStat{
		Hits:  10,
		Shots: 100,
	}

	if w.Accuracy() != 10 {
		t.Fatalf("incorrect accuracy, expected 10, got %u",
			w.Accuracy())
	}
}

func TestNewWeaponStat(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	w := WeaponStat{
		Hits:              10,
		Shots:             100,
		PlayerMatchStatID: 2,
	}

	wid := NewWeaponStat(store, w)

	var fw WeaponStat
	nf := db.Find(&fw, wid).RecordNotFound()
	assert.False(t, nf)

	assert.Equal(t, w.Hits, fw.Hits)
	assert.Equal(t, w.Shots, fw.Shots)
	assert.Equal(t, w.PlayerMatchStatID, fw.PlayerMatchStatID)
}

func TestListWeaponStat(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	w := WeaponStat{
		Hits:              10,
		Shots:             100,
		PlayerMatchStatID: 2,
	}

	NewWeaponStat(store, w)

	wfs := ListWeaponStats(store, 2)
	assert.Len(t, wfs, 1)

	assert.Len(t, ListWeaponStats(store, 9999), 0)
}
