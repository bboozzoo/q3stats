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
	"github.com/bboozzoo/q3stats/models/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayer(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	id := NewPlayer(store, "foo", "123")

	var p Player
	nf := db.First(&p, id).RecordNotFound()
	assert.False(t, nf)

	assert.Equal(t, "foo", p.Name)
	assert.Equal(t, "123", p.PasswordHash)
}

func TestHasPlayer(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	NewPlayer(store, "foo", "123")

	assert.True(t, HasPlayer(store, "foo"))
	assert.False(t, HasPlayer(store, "bar"))
}

func TestListPlayers(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	NewPlayer(store, "foo", "123")
	NewPlayer(store, "bar", "123")
	NewPlayer(store, "baz", "123")

	players := ListPlayers(store)

	assert.Len(t, players, 3)

	d := map[string]int{
		"foo": 0,
		"bar": 0,
		"baz": 0,
	}
	for _, p := range players {
		d[p.Name] = 1
	}
	for k, v := range d {
		if v != 1 {
			t.Fatalf("player %s not seen", k)
		}
	}
}

func TestHasPlayerId(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	p1 := NewPlayer(store, "foo", "123")
	p2 := NewPlayer(store, "bar", "123")

	assert.True(t, HasPlayerId(store, p1))
	assert.True(t, HasPlayerId(store, p2))
	assert.False(t, HasPlayerId(store, p1+p2))
}

func TestGetPlayer(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	pid1 := NewPlayer(store, "foo", "123")

	p := GetPlayer(store, pid1)

	assert.Equal(t, "foo", p.Name)
	assert.Equal(t, "123", p.PasswordHash)
	assert.Equal(t, pid1, p.ID)
}
