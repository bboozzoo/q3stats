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
	"testing"
	"time"
)

func TestHasDuration(t *testing.T) {
	d := []struct {
		item string
		exp  bool
	}{
		{RocketLauncher, false},
		{RedArmor, false},
		{QuadDamage, true},
		{BattleSuit, true},
		{Regeneration, true},
		{Haste, true},
		{MegaHealth, false},
	}

	for _, v := range d {
		i := ItemStat{Type: v.item}
		if i.HasDuration() != v.exp {
			t.Fatalf("expected %s for item %s, got %s",
				v.exp, v.item, i.HasDuration())
		}
	}
}

func TestDurationDesc(t *testing.T) {

	i := ItemStat{
		Time: time.Duration(1)*time.Minute + time.Duration(1)*time.Second,
	}

	dd := i.DurationDesc()
	if dd != "1m1s" {
		t.Fatalf("expected duration 1m1s, got %s", dd)
	}
}
