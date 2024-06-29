package main

// 1. advanced interface mechanics
// 2. typed function

import (
	"fmt"
	"strings"
)

// 1. advanced interface mechanics

type Putter interface {
	Put(id int, value any) error
}

type SimplePutter struct {
}

func (sp *SimplePutter) Put(id int, value any) error {
	return nil
}

type Storage interface {
	fmt.Stringer
	Putter
	Get(id int) (any, error)
}

type FooStorage struct{}

func (f *FooStorage) Get(id int) (any, error) {
	return 1, nil
}

func (f *FooStorage) String() string {
	return ""
}

func (f *FooStorage) Put(id int, value any) error {
	return nil
}

type Server struct {
	Store Storage
}

func updateValue(id int, val any, p Putter) error {
	return p.Put(id, val)
}

func foo() {
	s := &Server{
		Store: &FooStorage{},
	}

	s.Store.Get(1)

	updateValue(1, "value", &SimplePutter{})
}

// 2. typed function

type TransformFunc func(s string) string

func Uppercase(s string) string {
	return strings.ToUpper(s)
}

func Prefixer(prefix string) TransformFunc {
	return func(s string) string {
		return prefix + s
	}
}

func transformString(s string, fn TransformFunc) string {
	return fn(s)
}

func main() {
	fmt.Println(transformString("benjamin", Uppercase))
	fmt.Println(transformString("benjamin", Prefixer("FOO_")))
	fmt.Println(transformString("Smile", Prefixer("BAR_")))
}
