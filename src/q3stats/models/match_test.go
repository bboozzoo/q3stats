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
	"testing"
	"time"
)

func TestFindMatchByHash(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	hash := "123123foo"
	m := Match{
		Map:      "q3dm17",
		DataHash: hash,
	}

	db.Create(&m)

	mf := FindMatchByHash(store, hash)
	if mf == nil {
		t.Fatalf("expected to find a match with hash %s", hash)
	}

	if mf.DataHash != hash {
		t.Fatalf("expected match with hash %s, got %s",
			hash, mf.DataHash)
	}
}

func TestNewMatch(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	m := Match{
		Map:  "foo",
		Type: "FFA",
	}

	mid := NewMatch(store, m)

	var fm Match
	nf := db.Find(&fm, mid).RecordNotFound()
	if nf == true {
		t.Fatalf("expected to find match of ID %u", mid)
	}

	if m.Map != fm.Map || m.Type != fm.Type ||
		mid != fm.ID {
		t.Fatalf("found data mismatch in player stat data")
	}
}

func TestListMatches(t *testing.T) {
	store := test.GetStore(t)

	db := store.Conn()
	defer db.Close()

	CreateSchema(store)

	NewMatch(store, Match{Map: "foo", Type: "FFA",
		DateTime: time.Now().Add(time.Duration(1) * time.Hour),
	})
	NewMatch(store, Match{Map: "bar", Type: "TDM",
		DateTime: time.Now(),
	})
	NewMatch(store, Match{Map: "baz", Type: "FFA",
		DateTime: time.Now().Add(time.Duration(2) * time.Hour),
	})
	NewMatch(store, Match{Map: "zab", Type: "FFA",
		DateTime: time.Now().Add(time.Duration(3) * time.Hour),
	})
	// time order (newest to oldest):
	// - zab
	// - baz
	// - foo
	// - bar

	m := ListMatches(store, MatchListParams{})
	if len(m) != 4 {
		t.Fatalf("expected 4 matches, got %u", len(m))
	}

	m = ListMatches(store, MatchListParams{Limit: 2})
	if len(m) != 2 {
		t.Fatalf("expected 2 matches, got %u", len(m))
	}

	m = ListMatches(store, MatchListParams{TimeSort: true})
	exp := []string{"bar", "foo", "baz", "zab"}
	for i := range m {
		if exp[i] != m[i].Map {
			t.Fatalf("mismatch in time ordered matches, idx %u, got %s, expected %s",
				i, m[i].Map, exp[i])
		}
	}

	m = ListMatches(store, MatchListParams{TimeSort: true,
		SortDesc: true})
	exp = []string{"zab", "baz", "foo", "bar"}
	for i := range m {
		if exp[i] != m[i].Map {
			t.Fatalf("mismatch in time ordered matches, idx %u, got %s, expected %s",
				i, m[i].Map, exp[i])
		}
	}
}
