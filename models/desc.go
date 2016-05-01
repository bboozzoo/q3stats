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

// item (such as rocket launcher, railgun, red armor) descriptor
type ItemDesc struct {
	// short name, ex. RA (red armor), RL (rocket launcher)
	Short string
	// long name - Rocket Launcher
	Name string
}

// descriptor of items, mainly for presentation purposes
type DescData struct {
	// mapping of short name to long names
	ItemMap map[string]ItemDesc
}

var (
	descData = &DescData{
		ItemMap: map[string]ItemDesc{
			"G":   ItemDesc{"G", "Gauntlet"},
			"MG":  ItemDesc{"MG", "Machine Gun"},
			"SG":  ItemDesc{"SG", "Shotgun"},
			"PG":  ItemDesc{"PG", "Plasma Gun"},
			"GL":  ItemDesc{"GL", "Grenade Launcher"},
			"RL":  ItemDesc{"RL", "Rocket Launcher"},
			"RG":  ItemDesc{"RG", "Railgun"},
			"LG":  ItemDesc{"LG", "Lightning Gun"},
			"BFG": ItemDesc{"BFG", "BFG"},
			"Q":   ItemDesc{"Q", "Quad Damage"},
			"MH":  ItemDesc{"MH", "Mega Health"},
			"BS":  ItemDesc{"BS", "Battle Suit"},
			"RA":  ItemDesc{"RA", "Red Armor"},
			"YA":  ItemDesc{"YA", "Yellow Armor"},
			"GA":  ItemDesc{"GA", "Green Armor"},
		},
	}
)

// obtain pointer to item descriptor
func GetItemDesc() *DescData {
	return descData
}

// returns true if effects of given item are time bounded
func ItemHasDuration(itype string) bool {
	switch itype {
	case "Q", "BS":
		return true
	}
	return false
}
