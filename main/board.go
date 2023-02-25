package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Board struct {
	Width  int
	Height int

	// The board is represented as a 2D slice of Pieces.
	Locations [][]Piece

	// Keep track of the last move.
	LastMove Move

	// Keep track of whose turn is it
	// 1 - white player's turn
	// 2 - black player's turn
	// 3 - white player won
	// 4 - black player won
	// 5 - draw
	State int

	// Castling rights
	WhiteQueenSideCastle bool
	WhiteKingSideCastle  bool
	BlackQueenSideCastle bool
	BlackKingSideCastle  bool

	// En passant target square
	EnPassantTargetSquare Location

	// Halfmove clock
	HalfmoveClock int

	// Fullmove number
	FullmoveNumber int
}

type Location struct {
	X int
	Y int
}

func (b *Board) Init() {
	b.Height = 8
	b.Width = 8

	// Init with empty move
	b.LastMove = Move{
		Type: ' ',
	}

	// Initialize the board with empty pieces.
	b.Locations = make([][]Piece, b.Height)
	for i := 0; i < b.Height; i++ {
		b.Locations[i] = make([]Piece, b.Width)
		for j := 0; j < b.Width; j++ {
			b.Locations[i][j] = EmptyPiece{}
		}
	}
	// Initialize the pieces.
	// White pieces.
	b.Locations[0][0] = PlayerPiece{1, 'R', Location{0, 0}}
	b.Locations[0][1] = PlayerPiece{1, 'N', Location{0, 1}}
	b.Locations[0][2] = PlayerPiece{1, 'B', Location{0, 2}}
	b.Locations[0][3] = PlayerPiece{1, 'Q', Location{0, 3}}
	b.Locations[0][4] = PlayerPiece{1, 'K', Location{0, 4}}
	b.Locations[0][5] = PlayerPiece{1, 'B', Location{0, 5}}
	b.Locations[0][6] = PlayerPiece{1, 'N', Location{0, 6}}
	b.Locations[0][7] = PlayerPiece{1, 'R', Location{0, 7}}
	for i := 0; i < b.Width; i++ {
		b.Locations[1][i] = PlayerPiece{1, 'P', Location{1, i}}
	}
	// Black pieces.
	b.Locations[7][0] = PlayerPiece{2, 'R', Location{7, 0}}
	b.Locations[7][1] = PlayerPiece{2, 'N', Location{7, 1}}
	b.Locations[7][2] = PlayerPiece{2, 'B', Location{7, 2}}
	b.Locations[7][3] = PlayerPiece{2, 'Q', Location{7, 3}}
	b.Locations[7][4] = PlayerPiece{2, 'K', Location{7, 4}}
	b.Locations[7][5] = PlayerPiece{2, 'B', Location{7, 5}}
	b.Locations[7][6] = PlayerPiece{2, 'N', Location{7, 6}}
	b.Locations[7][7] = PlayerPiece{2, 'R', Location{7, 7}}
	for i := 0; i < b.Width; i++ {
		b.Locations[6][i] = PlayerPiece{2, 'P', Location{6, i}}
	}
}

// decode FEN string and return a board based on the string
func InitFEN(fen string) Board {
	var b Board
	b.Height = 8
	b.Width = 8

	// Init with empty move
	b.LastMove = Move{
		Type: ' ',
	}

	// Initialize the board with empty pieces.
	b.Locations = make([][]Piece, b.Height)
	for i := 0; i < b.Height; i++ {
		b.Locations[i] = make([]Piece, b.Width)
		for j := 0; j < b.Width; j++ {
			b.Locations[i][j] = EmptyPiece{}
		}
	}

	// Decode FEN string
	// FEN string is in the format:
	// <piece placement> <active color> <castling availability> <en passant target square> <halfmove clock> <fullmove number>

	// Split the string into 6 parts
	fenParts := strings.Split(fen, " ")

	// Piece placement
	piecePlacement := fenParts[0]
	ranks := strings.Split(piecePlacement, "/")
	for i := 0; i < b.Height; i++ {
		rank := ranks[7-i]
		file := 0
		for _, char := range rank {
			if char >= '1' && char <= '8' {
				file += int(char - '0')
			} else {
				// Check player
				if char >= 'A' && char <= 'Z' {
					b.Locations[i][file] = PlayerPiece{1, char, Location{i, file}}
				} else {
					b.Locations[i][file] = PlayerPiece{2, char & '_', Location{i, file}}
				}
				file++
			}
		}
	}

	// Active color
	activeColor := fenParts[1]
	if activeColor == "w" {
		b.State = 1
	} else {
		b.State = 2
	}

	// Castling availability
	castlingAvailability := fenParts[2]
	b.WhiteKingSideCastle = strings.Contains(castlingAvailability, "K")
	b.WhiteQueenSideCastle = strings.Contains(castlingAvailability, "Q")
	b.BlackKingSideCastle = strings.Contains(castlingAvailability, "k")
	b.BlackQueenSideCastle = strings.Contains(castlingAvailability, "q")

	// En passant target square
	enPassantTargetSquare := fenParts[3]
	if enPassantTargetSquare != "-" {
		b.EnPassantTargetSquare = Location{int(enPassantTargetSquare[1] - '1'), int(enPassantTargetSquare[0] - 'a')}
	}

	// Halfmove clock
	halfmoveClock := fenParts[4]
	b.HalfmoveClock, _ = strconv.Atoi(halfmoveClock)

	// Fullmove number
	fullmoveNumber := fenParts[5]
	b.FullmoveNumber, _ = strconv.Atoi(fullmoveNumber)

	return b
}

// Print the board.
func (b *Board) Print() {
	for i := b.Height - 1; i >= 0; i-- {
		for j := 0; j < b.Width; j++ {
			fmt.Printf("%c ", b.Locations[i][j].GetChar())
		}
		fmt.Println()
	}
}

// Print the board with border.
func (b *Board) PrintWithBorder() {
	fmt.Print("  ")
	for i := 0; i < b.Width; i++ {
		fmt.Printf("%c ", 'a'+i)
	}
	fmt.Println()
	for i := b.Height - 1; i >= 0; i-- {
		fmt.Printf("%d ", i+1)
		for j := 0; j < b.Width; j++ {
			fmt.Printf("%c ", b.Locations[i][j].GetChar())
		}
		fmt.Printf("%d", i+1)
		fmt.Println()
	}
	fmt.Print("  ")
	for i := 0; i < b.Width; i++ {
		fmt.Printf("%c ", 'a'+i)
	}
	fmt.Println()
}

// Deep copy the board.
func (b *Board) Copy() Board {
	newBoard := Board{
		Width:  b.Width,
		Height: b.Height,
	}
	newBoard.Locations = make([][]Piece, b.Height)
	for i := 0; i < b.Height; i++ {
		newBoard.Locations[i] = make([]Piece, b.Width)
		for j := 0; j < b.Width; j++ {
			newBoard.Locations[i][j] = b.Locations[i][j]
		}
	}
	return newBoard
}

// Get the piece at the given location.
func (b *Board) GetPiece(x, y int) Piece {
	// Check if the location is out of bound.
	// If so, return an empty piece.
	if x < 0 || x >= b.Height || y < 0 || y >= b.Width {
		return EmptyPiece{}
	}
	return b.Locations[x][y]
}
func (b *Board) GetPieceAtLocation(location Location) Piece {
	return b.Locations[location.X][location.Y]
}

// Get all pieces of the given player.
func (b *Board) GetPlayerPieces(player int) []Piece {
	pieces := make([]Piece, 0)
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			if b.Locations[i][j].GetPlayer() == player {
				pieces = append(pieces, b.Locations[i][j])
			}
		}
	}
	return pieces
}

// Get all moves for the given player.
// This function does not check if the player is in check.
func (b *Board) GetPlayerMoves(player int) []Move {
	pieces := b.GetPlayerPieces(player)
	moves := make([]Move, 0)
	for _, piece := range pieces {
		moves = append(moves, piece.GetMoves(b)...)
	}
	return moves
}

// Get all moves from the given player except the king
// since the king cannot attack another king
// while this avoid to have an infinite loop for detecting moves
// HACK: This is a hack, we should find a better way to do this.
func (b *Board) GetPlayerMovesExceptKing(player int) []Move {
	pieces := b.GetPlayerPieces(player)
	moves := make([]Move, 0)
	for _, piece := range pieces {
		if piece.GetType() != 'K' {
			moves = append(moves, piece.GetMoves(b)...)
		}
	}
	return moves
}

// Get all the legal moves for the given player.
// This function uses GetPlayerMoves() and filters out moves that put the player in check.
func (b *Board) GetPlayerLegalMoves(player int) []Move {
	moves := b.GetPlayerMoves(player)
	legalMoves := make([]Move, 0)
	for _, move := range moves {
		newBoard := b.Copy()
		newBoard.MakeMove(move)
		if !newBoard.CheckPlayerInCheck(player) {
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}

// Disambiguate given a set of moves.
// When multiple same type of pieces can move to the same location
// TODO: Implement this function.
func DisambiguateMoves(moves []Move) {

}

// Check if the given player is in check.
func (b *Board) CheckPlayerInCheck(player int) bool {
	// Get the king's location.
	kingLoc := Location{}
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			if b.Locations[i][j].GetType() == 'K' && b.Locations[i][j].GetPlayer() == player {
				kingLoc = Location{i, j}
			}
		}
	}
	// Check if any of the opponent's pieces can move to the king's location.
	opponentMoves := b.GetPlayerMovesExceptKing(3 - player)
	for _, move := range opponentMoves {
		if move.To == kingLoc {
			return true
		}
	}
	return false
}

// Move a piece on the board.
func (b *Board) MakeMove(m Move) {

	player := b.Locations[m.From.X][m.From.Y].GetPlayer()

	// Update the board state
	b.LastMove = m

	// Update the castling rights
	if m.Piece == 'K' {
		if player == 1 {
			b.WhiteQueenSideCastle = false
			b.WhiteKingSideCastle = false
		} else {
			b.BlackQueenSideCastle = false
			b.BlackKingSideCastle = false
		}
	}
	if m.Piece == 'R' {
		if player == 1 {
			if m.From.Y == 0 && m.From.X == 0 {
				b.WhiteQueenSideCastle = false
			}
			if m.From.Y == 7 && m.From.X == 0 {
				b.WhiteKingSideCastle = false
			}
		} else {
			if m.From.Y == 0 && m.From.X == 7 {
				b.BlackQueenSideCastle = false
			}
			if m.From.Y == 7 && m.From.X == 7 {
				b.BlackKingSideCastle = false
			}
		}
	}

	// Special case for castling.
	if m.Type == 'K' && m.From.Y-m.To.Y == 2 {
		// Long castling.
		b.Locations[m.From.X][m.From.Y] = EmptyPiece{}
		b.Locations[m.From.X][0] = EmptyPiece{}
		b.Locations[m.To.X][m.To.Y] = PlayerPiece{player, 'K', m.To}
		b.Locations[m.To.X][m.To.Y-1] = PlayerPiece{player, 'R', Location{m.To.X, m.To.Y - 1}}
		return
	}
	if m.Type == 'K' && m.From.Y-m.To.Y == -2 {
		// Short castling.
		b.Locations[m.From.X][m.From.Y] = EmptyPiece{}
		b.Locations[m.From.X][7] = EmptyPiece{}
		b.Locations[m.To.X][m.To.Y] = PlayerPiece{player, 'K', m.To}
		b.Locations[m.To.X][m.To.Y+1] = PlayerPiece{player, 'R', Location{m.To.X, m.To.Y + 1}}
		return
	}

	// Special case for en passant.
	if m.Type == 'E' {
		b.Locations[m.From.X][m.From.Y] = EmptyPiece{}
		b.Locations[m.From.X][m.To.Y] = EmptyPiece{}
		b.Locations[m.To.X][m.To.Y] = PlayerPiece{player, 'P', m.To}
		return
	}

	// Default case of moving / promoting / capturing.
	if m.Type == 'M' || m.Type == 'P' || m.Type == 'C' {
		b.Locations[m.From.X][m.From.Y] = EmptyPiece{}
		b.Locations[m.To.X][m.To.Y] = PlayerPiece{player, m.Piece, m.To}
		return
	}
}

// Chekc if the given player is in checkmate.
