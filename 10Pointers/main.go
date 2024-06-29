package main

import "fmt"

type Database struct {
	user string
}

type Server struct {
	db *Database // uintptr -> 8 bytes long
}

func (s *Server) GetUserFromDB() (string, error) {
	// golang is going to derefence the db pointer
	// it's going to lookup the memory address of the pointer
	if s.db == nil {
		// panic("database is not initialized")
		return "", fmt.Errorf("database is not initialized")
	}
	return s.db.user, nil
}

type Player struct {
	HP int
}

// function receiver -> function bind to a struct
func (p *Player) TakeDamage(amount int) {
	p.HP -= amount
	fmt.Printf("player is taking damage ->%d. New HP -> %d\n", amount, p.HP)
}

// if the player is not a pointer, we are adjusting the copy of the player
// not the actual player
func TakeDamageAgain(p *Player, amount int) {
	p.HP -= amount
	fmt.Printf("player is taking damage ->%d. New HP -> %d\n", amount, p.HP)
}

func main() {
	player := &Player{
		HP: 100,
	}

	player.TakeDamage(2)
	player.TakeDamage(2)
	player.TakeDamage(2)
	player.TakeDamage(2)

	TakeDamageAgain(player, 10)

	fmt.Println("==========")

	s := &Server{}
	user, err := s.GetUserFromDB()
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(user)
}
