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
package test

import (
	"github.com/bboozzoo/q3stats/store"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"testing"
)

type DB struct {
	db *gorm.DB
}

func (d *DB) Open() error {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return errors.Wrap(err, "failed to open DB")
	}

	d.db = db

	return nil
}

func (d *DB) Conn() *gorm.DB {
	return d.db
}

func (d *DB) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func GetStore(t *testing.T) store.DB {
	db := &DB{}

	if err := db.Open(); err != nil {
		t.Fatalf("failed to open DB: %s", err)
	}
	return db
}

func (d *DB) Begin() store.DBTransaction {
	return &dbTx{
		db: d.db.Begin(),
	}
}

type dbTx struct {
	db *gorm.DB
}

func (dbt *dbTx) Conn() *gorm.DB {
	return dbt.db
}

func (dbt *dbTx) Commit() {
	dbt.db.Commit()
}
