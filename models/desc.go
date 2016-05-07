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

const (
	// for sake of simplicity define item types using names
	// matching data coming from CPMA Q3 server
	RedArmor     = "RA"
	YellowArmor  = "YA"
	GreenArmor   = "GA"
	MegaHealth   = "MH"
	QuadDamage   = "Quad"
	BattleSuit   = "BattleSuit"
	Regeneration = "Regen"
	Haste        = "Haste"
	Flight       = "Flight"
	Invisibility = "Invis"

	// weapons
	Gauntlet        = "G"
	Machinegun      = "MG"
	Shotgun         = "SG"
	Plasmagun       = "PG"
	GrenadeLauncher = "GL"
	RocketLauncher  = "RL"
	Railgun         = "RG"
	BFG             = "BFG"
	LightningGun    = "LG"
)

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
			// weapons
			Gauntlet:        ItemDesc{"G", "Gauntlet"},
			Machinegun:      ItemDesc{"MG", "Machine Gun"},
			Shotgun:         ItemDesc{"SG", "Shotgun"},
			Plasmagun:       ItemDesc{"PG", "Plasma Gun"},
			GrenadeLauncher: ItemDesc{"GL", "Grenade Launcher"},
			RocketLauncher:  ItemDesc{"RL", "Rocket Launcher"},
			Railgun:         ItemDesc{"RG", "Railgun"},
			LightningGun:    ItemDesc{"LG", "Lightning Gun"},
			BFG:             ItemDesc{"BFG", "BFG"},

			// items
			RedArmor:     ItemDesc{"RA", "Red Armor"},
			YellowArmor:  ItemDesc{"YA", "Yellow Armor"},
			GreenArmor:   ItemDesc{"GA", "Green Armor"},
			MegaHealth:   ItemDesc{"MH", "Mega Health"},
			QuadDamage:   ItemDesc{"Q", "Quad Damage"},
			BattleSuit:   ItemDesc{"BS", "Battle Suit"},
			Regeneration: ItemDesc{"REG", "Regeneration"},
			Haste:        ItemDesc{"HS", "Haste"},
			Flight:       ItemDesc{"FL", "Flight"},
			Invisibility: ItemDesc{"INV", "Invisibility"},
		},
	}
)

// obtain pointer to item descriptor
func GetDesc() *DescData {
	return descData
}

// obtain description of given item/weapon
func GetItemDesc(name string) ItemDesc {
	return descData.ItemMap[name]
}

// returns true if effects of given item are time bounded
func ItemHasDuration(itype string) bool {
	switch itype {
	case QuadDamage, BattleSuit, Regeneration, Haste, Flight:
		return true
	}
	return false
}
