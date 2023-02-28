package main

import (
	"net"
)

// External players and API stuff

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

type RemotePlayer struct {
	Conn  net.Conn
	Color int
}

func (p *RemotePlayer) GetMove(b *Board) Move {
	// If we dont have the player
	// try to connect to the player
	if p.Conn == nil {
		conn, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
		if err != nil {
			panic(err)
		}
		p.Conn = conn
	}

	// Send the board to the player
	// and get the move back

	return Move{}
}
