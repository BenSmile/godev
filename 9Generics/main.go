package main

import "fmt"

type CustomMap[K comparable, V any] struct {
	data map[K]V
}

func (m *CustomMap[K, V]) Insert(k K, v V) error {
	m.data[k] = v
	return nil
}

func NewCustomMap[K comparable, V any]() *CustomMap[K, V] {
	return &CustomMap[K, V]{
		data: make(map[K]V),
	}
}

func foo[T comparable](val T) {
	fmt.Println(val)
}

func foo2[T any, S any](val T, val2 S) {
	fmt.Println(val)
	fmt.Println(val2)
}

func main() {
	m1 := NewCustomMap[string, int]()
	m1.Insert("benjamin", 20)
	m1.Insert("sophie", 25)
	m1.Insert("melissa", 50)

	fmt.Println(m1)

	m2 := NewCustomMap[int, float64]()
	m2.Insert(1, 20)
	m2.Insert(2, 25)
	m2.Insert(3, 50)

	fmt.Println(m2)

	foo(1)
	foo2("benjamin", "samy")
}
