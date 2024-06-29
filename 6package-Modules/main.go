package main

import (
	"fmt"
	"mygoproject/types"
	"mygoproject/utils"
)

func main() {
	number := utils.GetNumber()
	fmt.Println("Number is :", number)
	fmt.Println("Hello World!")

	user := types.User{
		Username:  utils.GetUsername(),
		FirstName: "admin",
		Age:       utils.GetAge(),
	}

	fmt.Printf("The user is :%+v\n", user)

}
