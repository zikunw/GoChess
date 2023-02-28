package main

import (
	"bufio"
	"fmt"
	"math/rand"
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
func (p *HumanPlayer) GetMove(b *Board) Move {

	reader := bufio.NewReader(os.Stdin)
	var move Move
	var from Location
	var to Location
	var isValid bool

	// Prompt user for which piece to move
	// If the piece is not owned by the player, prompt again
	for {
		fmt.Print("Which piece do you want to move? (e.g. a1): ")
		text, _ := reader.ReadString('\n')
		if len(text) < 2 {
			fmt.Print("Invalid piece. ")
			continue
		}
		if text[0] < 'a' || text[0] > 'h' || text[1] < '1' || text[1] > '8' {
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
		fmt.Print("Where do you want to move the piece? (e.g. a1): ")
		text, _ := reader.ReadString('\n')
		if len(text) < 2 {
			fmt.Print("Invalid piece. ")
			continue
		}
		if text[0] < 'a' || text[0] > 'h' || text[1] < '1' || text[1] > '8' {
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

type RandomComputerPlayer struct {
	Color int // 1 - white, 2 - black
}

// This player will randomly pick a legal move from the list
func (p *RandomComputerPlayer) GetMove(b *Board) Move {
	moves := b.GetPlayerLegalMoves(p.Color)
	// Choose a random move from the list
	move := moves[rand.Intn(len(moves))]
	return move
}

type MinimaxComputerPlayer struct {
	Color int // 1 - white, 2 - black
	Depth int
}

// This player will use the minimax algorithm to find the best move
func (p *MinimaxComputerPlayer) GetMove(b *Board) Move {
	var maximizingPlayer bool
	if p.Color == 1 {
		maximizingPlayer = true
	} else {
		maximizingPlayer = false
	}
	_, move := AlphaBetaAlg(b, p.Depth, maximizingPlayer, -100000, 100000)
	return move
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinMaxHeuristic(b *Board, playerColor int) int {
	// if white is checkmated, return a very small number
	if b.CheckPlayerInCheckmate(1) {
		return -10000
	}
	// if black is checkmate, return a very large number
	if b.CheckPlayerInCheckmate(2) {
		return 10000
	}
	return b.GetPieceValue(1) - b.GetPieceValue(2)
}

// MinMax algorithm (No alpha-beta pruning, so pretty slow)
// Input parameters:
// - b: the board
// - depth: the depth of the search tree
// - maximizingPlayer: true if the current player is the maximizing player
// - playerColor: the color of the current player
// Output:
// - the best move & heuristic value
func MinMaxAlg(b *Board, depth int, maximizingPlayer bool) (int, Move) {
	var playerColor int
	if maximizingPlayer {
		playerColor = 1
	} else {
		playerColor = 2
	}

	if depth == 0 || b.IsTerminal() {
		return MinMaxHeuristic(b, playerColor), Move{}
	}

	var bestMove Move
	if maximizingPlayer {
		value := -10000000
		moves := b.GetPlayerLegalMoves(playerColor)
		moves = ShuffleMoves(moves) // To make it more fun
		for _, move := range moves {
			newBoard := b.Copy()
			newBoard.MakeMove(move)
			newValue, _ := MinMaxAlg(&newBoard, depth-1, false)
			value = Max(value, newValue)
			if value == newValue {
				bestMove = move
			}
		}
		return value, bestMove
	} else {
		value := 10000000
		moves := b.GetPlayerLegalMoves(playerColor)
		moves = ShuffleMoves(moves) // To make it more fun
		for _, move := range moves {
			newBoard := b.Copy()
			newBoard.MakeMove(move)
			newValue, _ := MinMaxAlg(&newBoard, depth-1, true)
			value = Min(value, newValue)
			if value == newValue {
				bestMove = move
			}
		}
		return value, bestMove
	}
}

// Same as MinMax Alg but with alpha-beta pruning
func AlphaBetaAlg(b *Board, depth int, maximizingPlayer bool, alpha int, beta int) (int, Move) {
	var playerColor int
	if maximizingPlayer {
		playerColor = 1
	} else {
		playerColor = 2
	}

	if depth == 0 || b.IsTerminal() {
		return MinMaxHeuristic(b, playerColor), Move{}
	}

	var bestMove Move
	if maximizingPlayer {
		value := -10000000
		moves := b.GetPlayerLegalMoves(playerColor)
		moves = ShuffleMoves(moves) // To make it more fun
		for _, move := range moves {
			newBoard := b.Copy()
			newBoard.MakeMove(move)
			newValue, _ := AlphaBetaAlg(&newBoard, depth-1, false, alpha, beta)
			if newValue > beta {
				return newValue, move
			}
			value = Max(value, newValue)
			alpha = Max(alpha, value)
			if value == newValue {
				bestMove = move
			}
		}
		return value, bestMove
	} else {
		value := 10000000
		moves := b.GetPlayerLegalMoves(playerColor)
		moves = ShuffleMoves(moves) // To make it more fun
		for _, move := range moves {
			newBoard := b.Copy()
			newBoard.MakeMove(move)
			newValue, _ := AlphaBetaAlg(&newBoard, depth-1, true, alpha, beta)
			if newValue < alpha {
				return newValue, move
			}
			beta = Min(beta, value)
			value = Min(value, newValue)
			if value == newValue {
				bestMove = move
			}
		}
		return value, bestMove
	}
}
