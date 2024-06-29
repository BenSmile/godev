package main

import (
	"fmt"
	"time"
)

type Server struct {
	quitch chan struct{} // 0 byte
	msgch  chan string
}

func NewServer() *Server {
	return &Server{
		quitch: make(chan struct{}),
		msgch:  make(chan string, 128),
	}
}

func (s *Server) start() {
	fmt.Println("server starting...")
	s.loop() // block
}

func (s *Server) loop() {
mainloop: // labelling loop
	for {
		select {
		case <-s.quitch:
			fmt.Println("quitting server")
			break mainloop
		case msg := <-s.msgch:
			// to some stuff when we have a message
			s.handleMessage(msg)
		default:
		}
	}
	fmt.Println("server is shutting down gracefully...")
}

func (s *Server) sendMessage(msg string) {
	s.msgch <- msg
}

func (s *Server) handleMessage(msg string) {
	fmt.Println("we received a message:", msg)
}

func (s *Server) quit() {
	// close(s.quitch)
	s.quitch <- struct{}{}
}

func main() {
	server := NewServer()
	go func() {
		time.Sleep(time.Second * 5)
		server.quit()
	}()
	server.start()
}

func main1() {
	s := NewServer()
	go s.start()

	for i := 0; i < 100; i++ {
		s.sendMessage(fmt.Sprintf("handle this number %d", i))
	}

	time.Sleep(time.Second * 5)
}
