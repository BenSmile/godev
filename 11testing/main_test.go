package main

import (
	"reflect"
	"testing"
)

// to test manually -> go test ./ -run TestCalculateValues
// -count=1 -> no cached
func TestCalculateValues(t *testing.T) {
	var (
		expected = 8
		a        = 3
		b        = 5
	)
	have := calculateValues(a, b)

	if have != expected {
		t.Errorf("expected %d but have %d", expected, have)
	}
}

func TestEqualPlayers(t *testing.T) {
	expected := Player{
		name: "Mark",
		hp:   100,
	}

	have := Player{
		name: "Mark",
		hp:   100,
	}

	if !reflect.DeepEqual(expected, have) {
		t.Errorf("expected : %+v but got %+v", expected, have)
	}
}
