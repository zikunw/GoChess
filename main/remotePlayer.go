package main

import (
	"fmt"
)

// External players and API stuff

type RemotePlayer struct {
	Color int

	playerMove chan Move

	board *Board
}

// Initialize remote player
func (p *RemotePlayer) Init(color int, b *Board) {
	p.Color = color
	p.playerMove = make(chan Move, 1)
	p.board = b
}

// TODO: implement this
func (p *RemotePlayer) GetMove(b *Board) Move {
	// ==== Get the move ==== //
	fmt.Println("Waiting for move from client")
	// Wait for the client to send the move
	move := <-p.playerMove
	fmt.Println("Move received from client")

	return move
}

func (p *RemotePlayer) SendMove(move Move) {
	p.playerMove <- move
}
