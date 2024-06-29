package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()

	ctx := context.WithValue(context.Background(), "username", map[string]string{
		"name": "benjamin",
	})

	userID, err := fetchUserId(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("the response took %v -> %+v\n", time.Since(start), userID)
}

func fetchUserId(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	val := ctx.Value("username")

	fmt.Println("the value -> ", val)

	type result struct {
		userID string
		err    error
	}

	resultch := make(chan result, 1)

	go func() {
		res, err := thirdPartyHTTPCall()
		resultch <- result{
			userID: res,
			err:    err,
		}
	}()

	select {
	// Done()
	// 1. -> the context timeout is exceeded
	// 2. -> the context has been manually canceled -> Cancel()
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-resultch:
		return res.userID, res.err
	}

}

func thirdPartyHTTPCall() (string, error) {
	time.Sleep(time.Millisecond * 100)
	return "user id 1", nil
}
