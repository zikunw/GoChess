package main

import (
	"bufio"
	"fmt"
	"os"
)

// player is the struct that represents a player in the game.
// The game sturct is going ask this struct for the next move.

// Therefore we can have different types of players:
// - human player
// - computer player (and more variations of alg)

type PlayerController interface {
	GetMove(*Board) Move
}

type HumanPlayer struct {
	Color int // 1 - white, 2 - black
}

// Prompts the user to move a piece.
// TODO: Check if this function work
func (p *HumanPlayer) GetMove(b *Board) Move {

	reader := bufio.NewReader(os.Stdin)
	var move Move
	var from Location
	var to Location
	var isValid bool

	// Prompt user for which piece to move
	// If the piece is not owned by the player, prompt again
	for {
		fmt.Print("Which piece do you want to move? (e.g. a2): ")
		text, _ := reader.ReadString('\n')
		if len(text) < 2 {
			fmt.Print("Invalid piece. ")
			continue
		}
		// parse the text into a move
		from = GridToLocation(text)
		piece := b.GetPieceAtLocation(from)
		if piece.IsEmpty() != true && piece.GetPlayer() == p.Color {
			break
		}
		// prompt again
		fmt.Print("Invalid piece. ")
	}

	// Prompt user for where to move the piece
	// If the move is not legal, prompt again
	for {
		fmt.Print("Where do you want to move the piece? (e.g. a3): ")
		text, _ := reader.ReadString('\n')
		if len(text) < 2 {
			fmt.Print("Invalid piece. ")
			continue
		}
		// parse the text into a move
		to = GridToLocation(text)
		// check if the move is legal
		isValid, move = ValidMove(from, to, p.Color, b)
		if isValid {
			break
		}
		// prompt again
		fmt.Print("Invalid move. ")
	}

	return move
}
