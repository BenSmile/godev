package main

import (
	"fmt"
)

type NumberStorer interface {
	GetAll() ([]int, error)
	Put(int) error
}

type MongoDBNumberStore struct {
	// some values
}

type PostgreDBNumberStore struct {
	// some values
}

func (mongoDBNumberStore MongoDBNumberStore) GetAll() ([]int, error) {
	return []int{1, 2, 3, 4, 5, 6, 7}, nil
}

func (mongoDBNumberStore MongoDBNumberStore) Put(value int) error {
	fmt.Printf("Store the number %d into the mongoDB storage\n", value)
	return nil
}

func (s PostgreDBNumberStore) GetAll() ([]int, error) {
	return []int{1, 2, 3, 4}, nil
}

func (s PostgreDBNumberStore) Put(value int) error {
	fmt.Printf("Store the number %d into the postgres storage\n", value)
	return nil
}

type ApiServer struct {
	numberStore NumberStorer
}

func main() {
	apiServer := ApiServer{
		numberStore: MongoDBNumberStore{},
	}

	if err := apiServer.numberStore.Put(1); err != nil {
		apiServer.numberStore = PostgreDBNumberStore{}
	}

	numbers, err := apiServer.numberStore.GetAll()
	if err != nil {
		panic(err)
	}

	fmt.Println(numbers)

	err = apiServer.numberStore.Put(2)
	if err != nil {
		panic(err)
	}
}
