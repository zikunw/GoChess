import { useState } from 'react'

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

}

function App() {

  const [board, setBoard] = useState(new Board())

  return (
    <main className='w-auto h-screen bg-stone-700'>
      <div className='flex flex-col items-center justify-center h-full '>

        <div className="border-2 w-96 aspect-square grid grid-cols-8 grid-rows-8 shadow-xl">
          {board.squares.map((square, index) => (
            <div key={index} className={`aspect-square ${square.color ? 'bg-stone-300' : 'bg-stone-500'}`}>
              {square.piece && <PieceDisplay piece={square.piece} />}
            </div>
          ))}
        </div>
      </div>
    </main>
  )
}

function PieceDisplay ({piece}) {
  return (
    <div className="aspect-square">
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
        piece.type === 'k' ? blackKing : null} />
        
    </div>
  )
}

export default App
