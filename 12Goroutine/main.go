package main

import (
	"fmt"
	"time"
)

func main() {
	msgch := make(chan string, 128)

	msgch <- "A"
	msgch <- "B"
	msgch <- "C"
	close(msgch)

	for {
		msg, ok := <-msgch
		if !ok {
			break
		}
		fmt.Println("the message is:", msg)
	}

	// ranging over channel
	// this piece of code is our consumer
	for msg := range msgch {
		fmt.Println("the message is:", msg)
	}
}

func main3() {

	// 1 unbuffered channel -> synchronous
	// 2 buffered channel -> asynch
	// always blocks if it's full
	resultch := make(chan string) // -> unbuffered channel

	go func() {
		result := <-resultch
		fmt.Println(result)
	}()

	resultch <- "benjamin" // -> now it's full -> it will block -> BLOCK HERE

}

func fetchResource(n int) string {
	time.Sleep(time.Second * 2)
	return fmt.Sprintf("result %d", n)
}

func main2() {

	// 1 unbuffered channel -> synchronous
	// 2 buffered channel
	// always blocks if it's full
	resultch := make(chan string) // -> unbuffered channel
	// resultch2 := make(chan string, 10) // -> buffered channel
	resultch <- "foo" // -> now it's full -> it will block -> BLOCK HERE

	// the code below will never execute
	result1 := <-resultch

	fmt.Printf("the result from ch :", result1)

	result := fetchResource(1)
	fmt.Printf("the result :", result)
	// go fetchResource()
	// go fetchResource()
	// go fetchResource()
	// go fetchResource()

	go func() {
		fetchResource(2)
	}() // anonymous func

}
