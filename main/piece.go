package main

import "fmt"

type PlayerPiece struct {
	Player   int  // 1 - white player, 2 - black player
	Type     rune // 'P' - pawn, 'R' - rook, 'N' - knight, 'B' - bishop, 'Q' - queen, 'K' - king
	Location Location
}

func (p PlayerPiece) IsEmpty() bool  { return false }
func (p PlayerPiece) GetPlayer() int { return p.Player }
func (p PlayerPiece) GetType() rune  { return p.Type }
func (p PlayerPiece) GetValue() int {
	switch p.Type {
	case 'P':
		return 1
	case 'R':
		return 5
	case 'N':
		return 3
	case 'B':
		return 3
	case 'Q':
		return 9
	case 'K':
		return 0
	}
	return 0
}

// GetMoves returns a list of all possible moves for the piece.
func (p PlayerPiece) GetMoves(b *Board) []Move {
	switch p.Type {
	case 'P':
		return p.GetPawnMoves(b)
	case 'R':
		return p.GetRookMoves(b)
		//case 'N':
		//	return p.GetKnightMoves()
		//case 'B':
		//	return p.GetBishopMoves()
		//case 'Q':
		//	return p.GetQueenMoves()
		//case 'K':
		//	return p.GetKingMoves()
	}
	return []Move{}
}

// Get all the move from a pawn piece
func (p PlayerPiece) GetPawnMoves(b *Board) []Move {

	moves := make([]Move, 0)

	if p.Player == 1 { // White pawn
		// A pawn can move two spaces from a starting position.
		if p.Location.X == 1 {
			if b.GetPiece(p.Location.X+1, p.Location.Y).IsEmpty() {
				moves = append(moves, Move{
					Type:             'M',
					Piece:            'P',
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X + 1, p.Location.Y},
				})
				if b.GetPiece(p.Location.X+2, p.Location.Y).IsEmpty() {
					moves = append(moves, Move{
						Type:             'M',
						Piece:            'P',
						IsDisambiguation: false,
						From:             p.Location,
						To:               Location{p.Location.X + 2, p.Location.Y},
					})
				}
			}
		}
		// A pawn can move one space forward.
		if p.Location.X < 6 && p.Location.X != 1 {
			if b.GetPiece(p.Location.X+1, p.Location.Y).IsEmpty() {
				moves = append(moves, Move{
					Type:             'M',
					Piece:            'P',
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X + 1, p.Location.Y},
				})
			}
		}
		// A pawn can promote to any piece except a king.
		if p.Location.X == 6 && b.GetPiece(p.Location.X+1, p.Location.Y).IsEmpty() {
			const promotionPieces = "QRBN"
			for _, piece := range promotionPieces {
				// add the move to the list of moves
				moves = append(moves, Move{
					Type:             'P',
					Piece:            piece,
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X + 1, p.Location.Y},
				})
			}
		}

		// A pawn can move one space diagonally forward to capture an enemy piece.
		if p.Location.X < 6 {
			// left
			if p.Location.Y > 0 {
				leftPiece := b.GetPiece(p.Location.X+1, p.Location.Y-1)
				if !leftPiece.IsEmpty() && leftPiece.GetPlayer() != p.Player {
					moves = append(moves, Move{
						Type:             'C',
						Piece:            'P',
						IsDisambiguation: false,
						From:             p.Location,
						To:               Location{p.Location.X + 1, p.Location.Y - 1},
					})
				}
			}
			// right
			if p.Location.Y < 7 {
				rightPiece := b.GetPiece(p.Location.X+1, p.Location.Y+1)
				if !rightPiece.IsEmpty() && rightPiece.GetPlayer() != p.Player {
					moves = append(moves, Move{
						Type:             'C',
						Piece:            'P',
						IsDisambiguation: false,
						From:             p.Location,
						To:               Location{p.Location.X + 1, p.Location.Y + 1},
					})
				}
			}
		}

		// En passant (dont you dare forget this)
		if b.LastMove.Type == 'M' && b.LastMove.Piece == 'P' {
			// left en passant
			if p.Location.X == 5 && b.LastMove.To.Y == p.Location.Y-1 && b.LastMove.To.X == p.Location.X {
				moves = append(moves, Move{
					Type:             'E',
					Piece:            'P',
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X + 1, p.Location.Y - 1},
				})
			}
			// right en passant
			if p.Location.X == 5 && b.LastMove.To.Y == p.Location.Y+1 && b.LastMove.To.X == p.Location.X {
				moves = append(moves, Move{
					Type:             'E',
					Piece:            'P',
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X + 1, p.Location.Y + 1},
				})
			}
		}

	} else { // Black pawn
		// TODO: The X and Y is flipped for black pawns. Fix this.
		Debug(2, "Pawn is black.")
		// A pawn can move two spaces from a starting position.
		if p.Location.X == 7 {
			if b.GetPiece(p.Location.X-1, p.Location.Y).IsEmpty() {
				moves = append(moves, Move{
					Type:             'M',
					Piece:            'P',
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X - 1, p.Location.Y},
				})
			}
			if b.GetPiece(p.Location.X-2, p.Location.Y).IsEmpty() {
				moves = append(moves, Move{
					Type:             'M',
					Piece:            'P',
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X - 2, p.Location.Y},
				})
			}
		}
		// A pawn can move one space forward.
		if p.Location.X > 1 && p.Location.X != 7 {
			if b.GetPiece(p.Location.X-1, p.Location.Y).IsEmpty() {
				moves = append(moves, Move{
					Type:             'M',
					Piece:            'P',
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X - 1, p.Location.Y},
				})
			}
		}
		// A pawn can promote to any piece except a king.
		if p.Location.X == 1 && b.GetPiece(p.Location.X-1, p.Location.Y).IsEmpty() {
			const promotionPieces = "QRBN"
			for _, piece := range promotionPieces {
				// add the move to the list of moves
				moves = append(moves, Move{
					Type:             'P',
					Piece:            piece,
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X - 1, p.Location.Y},
				})
			}
		}

		// A pawn can move one space diagonally forward to capture an enemy piece.
		if p.Location.X > 0 {
			// left
			if p.Location.Y > 0 {
				leftPiece := b.GetPiece(p.Location.X-1, p.Location.Y-1)
				if !leftPiece.IsEmpty() && leftPiece.GetPlayer() != p.Player {
					// check if the pawn can be promoted
					if p.Location.X == 1 {
						const promotionPieces = "QRBN"
						for _, piece := range promotionPieces {
							// add the move to the list of moves
							moves = append(moves, Move{
								Type:             'P',
								Piece:            piece,
								IsDisambiguation: false,
								From:             p.Location,
								To:               Location{p.Location.X - 1, p.Location.Y - 1},
							})
						}
					} else {
						moves = append(moves, Move{
							Type:             'C',
							Piece:            'P',
							IsDisambiguation: false,
							From:             p.Location,
							To:               Location{p.Location.X - 1, p.Location.Y - 1},
						})
					}
				}
			}
			// right
			if p.Location.X < 7 {
				rightPiece := b.GetPiece(p.Location.X-1, p.Location.Y+1)
				if !rightPiece.IsEmpty() && rightPiece.GetPlayer() != p.Player {
					// check if the pawn can be promoted
					if p.Location.X == 1 {
						const promotionPieces = "QRBN"
						for _, piece := range promotionPieces {
							// add the move to the list of moves
							moves = append(moves, Move{
								Type:             'P',
								Piece:            piece,
								IsDisambiguation: false,
								From:             p.Location,
								To:               Location{p.Location.X - 1, p.Location.Y + 1},
							})
						}
					} else {
						moves = append(moves, Move{
							Type:             'C',
							Piece:            'P',
							IsDisambiguation: false,
							From:             p.Location,
							To:               Location{p.Location.X - 1, p.Location.Y + 1},
						})
					}
				}
			}
		}

		// En passant (dont you dare forget this)
		if b.LastMove.Type == 'M' && b.LastMove.Piece == 'P' {
			// left en passant
			if p.Location.X == 4 && b.LastMove.To.Y == p.Location.Y-1 && b.LastMove.To.X == p.Location.X {
				moves = append(moves, Move{
					Type:             'E',
					Piece:            'P',
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X - 1, p.Location.Y - 1},
				})
			}
			// right en passant
			if p.Location.X == 4 && b.LastMove.To.Y == p.Location.Y+1 && b.LastMove.To.X == p.Location.X {
				moves = append(moves, Move{
					Type:             'E',
					Piece:            'P',
					IsDisambiguation: false,
					From:             p.Location,
					To:               Location{p.Location.X - 1, p.Location.Y + 1},
				})
			}
		}
	}

	return moves
}

// Get all the possible moves for a rook.
func (p PlayerPiece) GetRookMoves(b *Board) []Move {
	moves := []Move{}
	rowNum := p.Location.X
	colNum := p.Location.Y

	// Check all the sqaures in the same row
	// left
	for i := p.Location.Y - 1; i >= 0; i-- {
		piece := b.GetPiece(rowNum, i)
		if piece.IsEmpty() {
			moves = append(moves, Move{'M', 'R', false, p.Location, Location{rowNum, i}})
		} else {
			if piece.GetPlayer() != p.Player {
				moves = append(moves, Move{'C', 'R', false, p.Location, Location{rowNum, i}})
			}
			break // stop checking this direction
		}
	}
	// right
	for i := p.Location.Y + 1; i < 8; i++ {
		piece := b.GetPiece(rowNum, i)
		if piece.IsEmpty() {
			moves = append(moves, Move{'M', 'R', false, p.Location, Location{rowNum, i}})
		} else {
			if piece.GetPlayer() != p.Player {
				moves = append(moves, Move{'C', 'R', false, p.Location, Location{rowNum, i}})
			}
			break // stop checking this direction
		}
	}

	// Check all the squares in the same column
	// up
	for i := p.Location.X - 1; i >= 0; i-- {
		piece := b.GetPiece(i, colNum)
		if piece.IsEmpty() {
			moves = append(moves, Move{'M', 'R', false, p.Location, Location{i, colNum}})
		} else {
			if piece.GetPlayer() != p.Player {
				moves = append(moves, Move{'C', 'R', false, p.Location, Location{i, colNum}})
			}
			break // stop checking this direction
		}
	}
	// down
	for i := p.Location.X + 1; i < 8; i++ {
		piece := b.GetPiece(i, colNum)
		if piece.IsEmpty() {
			moves = append(moves, Move{'M', 'R', false, p.Location, Location{i, colNum}})
		} else {
			if piece.GetPlayer() != p.Player {
				moves = append(moves, Move{'C', 'R', false, p.Location, Location{i, colNum}})
			}
			break // stop checking this direction
		}
	}

	return moves
}

func (p PlayerPiece) GetChar() rune {
	switch p.Type {
	case 'P':
		if p.Player == 1 {
			return '♟'
		} else {
			return '♙'
		}
	case 'R':
		if p.Player == 1 {
			return '♜'
		} else {
			return '♖'
		}
	case 'N':
		if p.Player == 1 {
			return '♞'
		} else {
			return '♘'
		}
	case 'B':
		if p.Player == 1 {
			return '♝'
		} else {
			return '♗'
		}
	case 'Q':
		if p.Player == 1 {
			return '♛'
		} else {
			return '♕'
		}
	case 'K':
		if p.Player == 1 {
			return '♚'
		} else {
			return '♔'
		}
	}
	return ' '
}

func (p PlayerPiece) String() string {
	return fmt.Sprintf("[Player Piece] Player: %d, Type: %c, X: %d, Y: %d", p.Player, p.Type, p.Location.X, p.Location.Y)
}

type EmptyPiece struct {
}

func (p EmptyPiece) IsEmpty() bool            { return true }
func (p EmptyPiece) GetPlayer() int           { return 0 }
func (p EmptyPiece) GetType() rune            { return ' ' }
func (p EmptyPiece) GetValue() int            { return 0 }
func (p EmptyPiece) GetMoves(b *Board) []Move { return []Move{} }
func (p EmptyPiece) GetChar() rune            { return ' ' }
func (p EmptyPiece) String() string           { return "[Empty Piece]" }

type Piece interface {
	// Check if empty.
	IsEmpty() bool
	// Returns the player that owns this piece.
	GetPlayer() int
	// Returns the type of this piece.
	GetType() rune
	// Returns the value of this piece.
	GetValue() int
	// Returns the possible moves for this piece.
	GetMoves(*Board) []Move
	// Return unicode character for this piece.
	GetChar() rune
	// to string
	String() string
}

type Move struct {
	Type             rune // 'M' - move, 'C' - capture, 'E' - en passant, 'P' - pawn promotion, 'K' - castle, 'I' - initial position
	Piece            rune // 'P' - pawn, 'R' - rook, 'N' - knight, 'B' - bishop, 'Q' - queen, 'K' - king
	IsDisambiguation bool // If true, the move is disambiguated by the FromX and FromY fields.
	From             Location
	To               Location
}

// Translate the move to algebraic notation
func (m Move) ToString() string {
	// Panic if the move is invalid.
	if m.Type == ' ' {
		panic("Invalid move")
	}
	// Panic if out of bounds.
	if m.From.X < 0 || m.From.X > 7 || m.From.Y < 0 || m.From.Y > 7 || m.To.X < 0 || m.To.X > 7 || m.To.Y < 0 || m.To.Y > 7 {
		panic("Move out of bounds")
	}
	// Translate the move to a string.
	moveString := string(m.Piece)
	switch m.Type {
	case 'M':
		if m.IsDisambiguation {
			moveString = moveString + LocationToGrid(m.From) + LocationToGrid(m.To)
		} else {
			moveString = moveString + LocationToGrid(m.To)
		}
	case 'C':
		if m.IsDisambiguation {
			moveString = moveString + LocationToGrid(m.From) + "x" + LocationToGrid(m.To)
		} else {
			moveString = moveString + "x" + LocationToGrid(m.To)
		}
	case 'E':
		moveString = moveString + "x" + LocationToGrid(m.To)
	case 'P':
		moveString = moveString + LocationToGrid(m.To) + "=" + string(m.Piece)
	case 'K':
		if m.To.X == 6 {
			moveString = "O-O"
		} else {
			moveString = "O-O-O"
		}
	}
	// if m.IsCheck() {
	// 	moveString = moveString + "+"
	// }
	return moveString
}

// x,y to grid location
func LocationToGrid(l Location) string {
	return string('a'+l.Y) + string('1'+l.X)
}

// grid location to x,y
func GridToLocation(grid string) Location {
	//return int(grid[0] - 'a'), int(grid[1] - '1')
	return Location{int(grid[0] - 'a'), int(grid[1] - '1')}
}

// Check if a given move is check.
// TODO
func (m Move) IsCheck(b *Board) bool {
	// // Make a copy of the board.
	// newBoard := b.Copy()
	// // Make the move.
	// newBoard.MakeMove(m)
	// // Check if the king is in check.
	// return newBoard.IsInCheck(b.Turn)

	return false
}
