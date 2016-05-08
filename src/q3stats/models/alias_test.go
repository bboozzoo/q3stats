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
	"github.com/stretchr/testify/assert"
	"q3stats/models/test"
	"testing"
)

func TestGetAliases(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	db.Create(&Alias{Alias: "foo"})
	db.Create(&Alias{Alias: "bar"})

	aliases := GetAliases(store, NoUser)

	if len(aliases) != 2 {
		t.Fatalf("returned aliases %u expected 2",
			len(aliases))
	}

	// create another one, assigned to player 1
	db.Create(&Alias{Alias: "baz", PlayerID: 1})

	aliases = GetAliases(store, NoUser)
	// still expecting 2
	if len(aliases) != 2 {
		t.Fatalf("returned aliases %u expected 2",
			len(aliases))
	}

	aliases = GetAliases(store, 1)
	// still expecting 2
	if len(aliases) != 1 {
		t.Fatalf("returned aliases %u expected 1",
			len(aliases))
	}
	al := aliases[0]
	if al.Alias != "baz" {
		t.Fatalf("expected baz, got %s", al.Alias)
	}

}

func TestClaimAlias(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	db.Create(&Alias{Alias: "foo"})
	db.Create(&Alias{Alias: "bar"})
	db.Create(&Alias{Alias: "baz"})

	ClaimAliasesByPlayer(store, 1, []string{"foo", "bar"})

	player1_aliases := GetAliases(store, 1)
	if len(player1_aliases) != 2 {
		t.Fatalf("expected 2 aliases, got %u",
			len(player1_aliases))
	}

	for _, a := range player1_aliases {
		if a.Alias != "foo" && a.Alias != "bar" {
			t.Fatalf("unexpected alias: %s", a.Alias)
		}
		if a.PlayerID != 1 {
			t.Fatalf("alias with incorrect player ID %u",
				a.PlayerID)
		}
	}

	unclaimed := GetAliases(store, NoUser)
	if len(unclaimed) != 1 {
		t.Fatal("expected 1 unclaimed alias, got %u",
			len(unclaimed))
	}
	ua := unclaimed[0]
	if ua.Alias != "baz" {
		t.Fatalf("expected 'baz' to be unclaimed, got %s",
			ua.Alias)
	}
}

func TestNewAlias(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	// add bogus alias first
	aidfoo := NewAliasOrCurrent(store, Alias{Alias: "foo"})
	// now the one we're going to track
	aidbar := NewAliasOrCurrent(store, Alias{Alias: "bar"})

	if aidfoo == aidbar {
		t.Fatalf("expected IDs to be different, got %u %u",
			aidfoo, aidbar)
	}

	// try again, we should get the same alias
	aidbar2 := NewAliasOrCurrent(store, Alias{Alias: "bar"})
	if aidbar2 != aidbar {
		t.Fatalf("expected the same ID, got %u %u",
			aidbar, aidbar2)
	}
}

func TestGetAlias(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	// add bogus alias first
	aidfoo := NewAliasOrCurrent(store, Alias{Alias: "foo"})

	afoo := GetAlias(store, aidfoo)
	assert.NotNil(t, afoo)

	assert.Nil(t, GetAlias(store, 999))
}
