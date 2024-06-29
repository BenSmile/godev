package main

import "fmt"

type Entity struct {
	name    string
	id      string
	version string
	Position
}

type Position struct {
	x int
	y int
}

type SpecialEntity struct {
	Entity

	specialField float64
}

type Player struct {
	Position
}

type Color int

const (
	ColorBlue Color = iota
	ColorBlack
	ColorYellow
	ColorPink
)

func (p *Position) Move(val int) {
	fmt.Println("Position is moving by :", val)
}

func (c Color) String() string {
	switch c {
	case ColorBlack:
		return "BLACK"
	case ColorBlue:
		return "BLUE"
	case ColorYellow:
		return "YELLOW"
	case ColorPink:
		return "PINK"
	default:
		panic("no color")
	}
}

func main() {

	fmt.Println(ColorBlack)

	p := Player{}
	p.Move(2)

}

func foo() {
	e := Entity{
		name:    "my entity",
		version: "version 1.1",
		Position: Position{
			x: 100,
			y: 200,
		},
	}

	se := SpecialEntity{
		specialField: 80.90,
		Entity:       e,
	}

	se.id = "id 1"
	fmt.Printf("%+v", se.x)
}
