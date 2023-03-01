import { useEffect, useState } from 'react'

// piece png
import blackPawn from './assets/black_pawn.png'
import blackRook from './assets/black_rook.png'
import blackKnight from './assets/black_knight.png'
import blackBishop from './assets/black_bishop.png'
import blackQueen from './assets/black_queen.png'
import blackKing from './assets/black_king.png'
import whitePawn from './assets/white_pawn.png'
import whiteRook from './assets/white_rook.png'
import whiteKnight from './assets/white_knight.png'
import whiteBishop from './assets/white_bishop.png'
import whiteQueen from './assets/white_queen.png'
import whiteKing from './assets/white_king.png'

// Square direction constants
const UP = -8
const DOWN = 8
const LEFT = -1
const RIGHT = 1
const UP_LEFT = UP + LEFT
const UP_RIGHT = UP + RIGHT
const DOWN_LEFT = DOWN + LEFT
const DOWN_RIGHT = DOWN + RIGHT

// information about a board square
class Square {
  constructor(color, piece) {
    this.color = color
    this.piece = piece
  }
}

// information about a piece
class Piece {
  constructor(color, type) {
    this.color = color
    this.type = type
  }
}

// information about a board
class Board {

  constructor() {
    this.squares = []
    this.readFEN('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR')
  }

  readFEN(fen) {
    let pieceLocation = fen.split(' ')[0]
    let rows = pieceLocation.split('/')
    let squareCount = 0
    let squareColorIsLight = true
    for (let row of rows) {
      for (let square of row) {
        if (parseInt(square)) {
          for (let i = 0; i < square; i++) {
            squareCount++
            if ((squareCount - 1) % 8 == 0) {
              squareColorIsLight = !squareColorIsLight
            }
            squareColorIsLight = !squareColorIsLight
            this.squares.push(new Square(squareColorIsLight, ""))
          }
        } else {
          squareCount++
            if ((squareCount - 1) % 8 == 0) {
              squareColorIsLight = !squareColorIsLight
            }
            squareColorIsLight = !squareColorIsLight
          if (square === square.toUpperCase()) {
            this.squares.push(new Square(squareColorIsLight, new Piece('white', square)))
          } else {
            this.squares.push(new Square(squareColorIsLight, new Piece('black', square)))
          }
        }
      }
    }
  }

  movePiece(from, to) {
    this.squares[to].piece = this.squares[from].piece
    this.squares[from].piece = ""
  }

  pieceLegalMoves(piece, square) {
    console.log(piece, square)
    let moves = []

    if (piece.type === 'P') { // white pawn
      // move forward
      if (square > 8 && this.squares[square + UP].piece === "") {
        moves.push(square + UP)
      }
      // move forward 2
      if (square <= 55 && square >= 48 && this.squares[square + UP].piece === "" && this.squares[square + UP + UP].piece === "") {
        console.log("Move forward 2")
        moves.push(square + UP + UP)
      }
      // capture left
      if (square > 8 && square % 8 !== 0 && this.squares[square + UP_LEFT].piece !== "" && this.squares[square + UP_LEFT].piece.color !== piece.color) {
        moves.push(square + UP_LEFT)
      }
      // capture right
      if (square > 8 && square % 8 !== 7 && this.squares[square + UP_RIGHT].piece !== "" && this.squares[square + UP_RIGHT].piece.color !== piece.color) {
        console.log("capture right")
        moves.push(square + UP_RIGHT)
      }
      // en passant left TODO: check if last move was a double pawn move
      // en passant right

      // promotion TODO: deal with this later
    }

    if (piece.type === 'p') { // black pawn
      // move forward
      if (square < 55 && this.squares[square + DOWN].piece === "") {
        moves.push(square + DOWN)
      }
      // move forward 2
      if (square <= 15 && square >= 8 && this.squares[square + DOWN].piece === "" && this.squares[square + DOWN + DOWN].piece === "") {
        moves.push(square + DOWN + DOWN)
      }
      // capture left
      if (square < 55 && square % 8 !== 0 && this.squares[square + DOWN_LEFT].piece !== "" && this.squares[square + DOWN_LEFT].piece.color !== piece.color) {
        moves.push(square + DOWN_LEFT)
      }
      // capture right
      if (square < 55 && square % 8 !== 7 && this.squares[square + DOWN_RIGHT].piece !== "" && this.squares[square + DOWN_RIGHT].piece.color !== piece.color) {
        moves.push(square + DOWN_RIGHT)
      }
      // en passant left TODO: check if last move was a double pawn move
      // en passant right

      // promotion TODO: deal with this later
    }
    return moves
  }
}

function App() {

  const [board, setBoard] = useState(new Board())
  const [selectedSquare, setSelectedSquare] = useState(-1)

  const [legalMoves, setLegalMoves] = useState([])

  //const [player, setPlayer] = useState('white')

  const handlePieceMove = (from, to) => {
    board.movePiece(from, to)
    setBoard(board)
    setSelectedSquare(-1)
  }

  return (
    <main className='w-auto h-screen bg-stone-800'>
      <div className='flex flex-col items-center justify-center h-full '>

        <div className="w-96 aspect-square grid grid-cols-8 grid-rows-8 shadow-xl">
          {
            board.squares.map((square, index) => (
            <SquareDisplay 
              key={index} 
              index={index} 
              square={square}
              board={board}
              selectedSquare={selectedSquare} 
              legalMoves={legalMoves}
              setSelectedSquare={setSelectedSquare}
              handlePieceMove={handlePieceMove}
              setLegalMoves={setLegalMoves}
              /> ))
          }
        </div>

      </div>
    </main>
  )
}

function SquareDisplay ({index, square, board, selectedSquare,legalMoves, setSelectedSquare, handlePieceMove, setLegalMoves}) {

  const handleOnClick = () => {
    console.log("clicked", index)
    setLegalMoves([])

    if (selectedSquare === index) {
      setSelectedSquare(-1)
      return
    }

    if (selectedSquare !== -1) {
      // TODO: handle request from the server
      // check if legal move
      if (legalMoves.includes(index)){
        handlePieceMove(selectedSquare, index)
      } else {
        setSelectedSquare(-1)
      }
      
      return
    }
    
    // if there is a piece on the square
    // set legal moves
    if (square.piece){
      setLegalMoves(board.pieceLegalMoves(square.piece, index))
      setSelectedSquare(index)
    }
    
  }

  return (
    <div 
      onClick={handleOnClick} 
      className={`aspect-square ${(index===selectedSquare || legalMoves.includes(index))? "bg-red-400" :square.color ? 'bg-stone-300' : 'bg-stone-500'}`}>
      {square.piece && <PieceDisplay piece={square.piece} />}
    </div>
  )
}

function PieceDisplay ({piece}) {
  return (
    <div className="aspect-square duration-200 hover:scale-110">
      <img src={
        piece.type === 'P' ? whitePawn :
        piece.type === 'R' ? whiteRook :
        piece.type === 'N' ? whiteKnight :
        piece.type === 'B' ? whiteBishop :
        piece.type === 'Q' ? whiteQueen :
        piece.type === 'K' ? whiteKing :
        piece.type === 'p' ? blackPawn :
        piece.type === 'r' ? blackRook :
        piece.type === 'n' ? blackKnight :
        piece.type === 'b' ? blackBishop :
        piece.type === 'q' ? blackQueen :
        piece.type === 'k' ? blackKing : null} 
        style={{
          MozWindowDragging: 'none',
          WebkitAppRegion: 'no-drag',
          WebkitUserSelect: 'none',
          WebkitTouchCallout: 'none',
          WebkitUserDrag: 'none',

        }}
        />
        
    </div>
  )
}

export default App
