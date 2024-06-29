package main

import "fmt"

var (
	floatVar  float32 = 0.1
	floatVar2 float64 = 0.1
	name      string  = "foo"
	age       int32   = 30
	intVar64  int64   = 48484
	uint8Var  uint8   = 0x11
	int8Var   int8    = -128
	runVar    rune    = 'a'
)

type Player struct {
	name        string
	health      int
	attachPower float64
}

type Weapon string

type version int

func getWeapon(weapon Weapon) string {
	return string(weapon)
}

func (p *Player) getHealth() int {
	return p.health
}

func main() {
	player := Player{
		name:        "Captain jack",
		health:      100,
		attachPower: 45.1,
	}

	fmt.Printf("players : %+v", player)

	users := map[string]int{
		"smile": 10,
	} // empty map

	users["foo"] = 20
	users["smilee"] = 20

	fmt.Printf("users : %+v\n", users)

	userss := make(map[string]int)

	fmt.Printf("userss : %+v\n", userss)

	age := users["foo"]

	fmt.Printf("age : %d\n", age)

	age, ok := users["smilee"]

	if !ok {
		fmt.Println("user not exists")
	} else {
		fmt.Printf("user age : %d\n", age)
	}

	delete(users, "foo")

	age = users["foo"]

	fmt.Printf("foo age : %d\n", age)

	for k, v := range users {
		fmt.Printf("the key : %s | value : %d\n", k, v)
	}

	for k := range users {
		fmt.Printf("the key : %s\n", k)
	}

	for _, v := range users {
		fmt.Printf("value : %d\n", v)
	}

	numbers := []int{} // empty slice
	fmt.Printf("numbers : %v\n", numbers)

	otherNumbers := make([]int, 5)
	otherNumbers[1] = 9
	otherNumbers = append(otherNumbers, 2)
	fmt.Printf("otherNumbers : %v\n", otherNumbers)

	numbers2 := [2]int{} // arrays

	fmt.Printf("numbers2 : %v\n", numbers2)

	numbers2[1] = 2

}
