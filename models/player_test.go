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
	if nf == true {
		t.Fatalf("expected to find player with ID %u", id)
	}

	if p.Name != "foo" {
		t.Fatalf("player of ID %u has unexpected name %s",
			id, p.Name)
	}
}

func TestHasPlayer(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	NewPlayer(store, "foo", "123")

	if HasPlayer(store, "foo") != true {
		t.Fatalf("expected to find player foo")
	}

	if HasPlayer(store, "bar") != false {
		t.Fatalf("expected player bar not to be found")
	}
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

	if len(players) != 3 {
		t.Fatalf("expected 3 players, got %u", len(players))
	}

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
