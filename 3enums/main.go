package main

import (
	"fmt"
)

// weapon type
// axe
// sword
// wood stick
// knife

type WeaponType int

// const (
// 	Axe         WeaponType = 1
// 	Sword       WeaponType = 2
// 	WoodenStick WeaponType = 3
// 	Knife       WeaponType = 4
// )

const (
	Axe WeaponType = iota // will increment all below values
	Sword
	WoodenStick
	Knife
)

func (w WeaponType) String() string {
	switch w {
	case Axe:
		return "Axe"
	case Sword:
		return "Sword"
	case WoodenStick:
		return "WoodenStick"
	case Knife:
		return "Knife"
	default:
		panic(fmt.Errorf("weapon not exists"))
	}
}

func getDammage(weaponType WeaponType) int {
	switch weaponType {
	case Axe:
		return 100
	case Sword:
		return 90
	case WoodenStick:
		return 1
	case Knife:
		return 40
	default:
		panic(fmt.Errorf("weapon not exists"))
	}
}

func main() {
	weapon := Axe

	fmt.Printf("dammage of weapon (%s) is (%d)\n", weapon, getDammage(weapon))
	fmt.Printf("dammage of weapon (%s) is (%d)\n", Sword, getDammage(Sword))
	fmt.Printf("dammage of weapon (%s) is (%d)\n", WoodenStick, getDammage(WoodenStick))
	fmt.Printf("dammage of weapon (%s) is (%d)\n", Knife, getDammage(Knife))

	fmt.Printf("Tostring : %s", weapon)
}
