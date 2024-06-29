package main

import "fmt"

func main() {

	for i := 0; i < 10; i++ {
		fmt.Println("it :", i)
	}

	numbers := []int{1, 2, 3, 4, 5, 6, 7}

	for i := 0; i < len(numbers); i++ {
		fmt.Println("num :", i)
	}

	names := []string{"a", "b", "c"}

	for index, name := range names {
		if name == "b" {
			break
		}
		fmt.Println("names :", name)
		fmt.Println("index :", index)
	}

	// continue -> skip that element
	// breack -> exit the loop , stop looping
	// return -> exit the function

	users := map[string]int{
		"foo":   1,
		"ben":   2,
		"alice": 3,
	}

	for k, v := range users {
		fmt.Printf("%s : %d\n", k, v)
	}

	name := "Alice"
	switch name {
	case "Alice":
		fmt.Println("The name Alice")
	case "Bob":
		fmt.Println("The name Alice")
	default:
		fmt.Println("The name Default")
	}

}
