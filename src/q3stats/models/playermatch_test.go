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

func TestNewPlayerMatchStat(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	pms := PlayerMatchStat{
		Score:   10,
		Kills:   10,
		AliasID: 2,
		MatchID: 3,
	}

	pmsid := NewPlayerMatchStat(store, pms)

	var fpms PlayerMatchStat
	nf := db.Find(&fpms, pmsid).RecordNotFound()
	assert.False(t, nf, "expected to find player-match stat of ID %u", pmsid)

	assert.Equal(t, pms.AliasID, fpms.AliasID)
	assert.Equal(t, pms.MatchID, fpms.MatchID)
	assert.Equal(t, pms.Score, fpms.Score)
}

func TestListPlayerMatchStat(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	NewPlayerMatchStat(store, PlayerMatchStat{
		Score:   1,
		AliasID: 1,
		MatchID: 2,
	})
	NewPlayerMatchStat(store, PlayerMatchStat{
		Score:   -1,
		AliasID: 2,
		MatchID: 2,
	})
	NewPlayerMatchStat(store, PlayerMatchStat{
		Score:   -1,
		AliasID: 2,
		MatchID: 1,
	})

	pls := ListPlayerMatchStat(store, 2)
	assert.Len(t, pls, 2)

	pls = ListPlayerMatchStat(store, 1)
	assert.Len(t, pls, 1)

	pls = ListPlayerMatchStat(store, 5)
	assert.Len(t, pls, 0)
}
