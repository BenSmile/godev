package main

import (
	"sync"
	"sync/atomic"
)

// state -> set OR update OR delete

type State struct {
	mu    sync.Mutex
	count int
}

func (state *State) setState(i int) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.count = i
}

type State2 struct {
	count int32
}

func (state *State2) setState(i int) {
	atomic.AddInt32(&state.count, int32(i))
}

func main() {

}
