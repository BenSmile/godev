package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type UserProfile struct {
	ID       int
	Comments []string
	Likes    int
	Friends  []int
}

type Response struct {
	data any
	err  error
}

func main() {
	start := time.Now()
	userProfile, err := handleGetUserProfile(10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(userProfile)
	fmt.Println("fetching the user profile took", time.Since(start))

}

func handleGetUserProfile(id int) (*UserProfile, error) {
	var (
		respch = make(chan Response, 3)
		wg     = &sync.WaitGroup{}
	)

	// we are making 3 request inside thier own goroutine

	go getComments(id, respch, wg)
	go getLikes(id, respch, wg)
	go getFriends(id, respch, wg)
	// adding 3 to the waitgroup
	wg.Add(3)
	wg.Wait() // block until the wg counter == 0
	close(respch)

	userProfile := &UserProfile{}
	for resp := range respch {
		if resp.err != nil {
			return nil, resp.err
		}
		switch msg := resp.data.(type) {
		case int:
			userProfile.Likes = msg
		case []int:
			userProfile.Friends = msg
		case []string:
			userProfile.Comments = msg
		default:
		}
	}

	return userProfile, nil
}

func getComments(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)

	comments := []string{
		"hey, that was great",
		"yeah Buddy",
		"Ow",
	}
	respch <- Response{
		data: comments,
		err:  nil,
	}
	//work is done
	wg.Done()
}

func getLikes(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 100)
	respch <- Response{
		data: 34,
		err:  nil,
	}
	//work is done
	wg.Done()
}

func getFriends(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	friendIds := []int{11, 34, 43, 5, 6}
	respch <- Response{
		data: friendIds,
		err:  nil,
	}
	//work is done

	wg.Done()
}
