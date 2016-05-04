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