package main

import "fmt"

var (
	name             = "Foo"
	firstName string = "Foo"
	lastName  string
)

const (
	version = 1
	Version = 2 // exportable
)

func main() {

	version := 2 //infer int

	fmt.Print(version)

}
