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

// Post function
// From: https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
async function postData(url = "", data = {}) {
  // Default options are marked with *
  const response = await fetch(url, {
    method: "POST", // *GET, POST, PUT, DELETE, etc.
    mode: "cors", // no-cors, *cors, same-origin
    cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
    credentials: "same-origin", // include, *same-origin, omit
    headers: {
      "Content-Type": "application/json",
      // 'Content-Type': 'application/x-www-form-urlencoded',
    },
    redirect: "follow", // manual, *follow, error
    referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: JSON.stringify(data), // body data type must match "Content-Type" header
  });
  return response.json(); // parses JSON response into native JavaScript objects
}

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
    // Keep track of whose turn is it
	  // 1 - white player's turn
	  // 2 - black player's turn
	  // 3 - white player won
	  // 4 - black player won
	  // 5 - draw
    this.state = 1
    this.readFEN('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR')
  }

  readFEN(fen) {
    // Clear the squares
    this.squares = []

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

  readBoardState(boardState) {
    const data = boardState.split(' ')
    this.readFEN(data[0])
    this.state = parseInt(data[1])
  }

  movePiece(from, to) {
    this.squares[to].piece = this.squares[from].piece
    this.squares[from].piece = ""
  }

  // Ideally I dont want to use this function
  // Instead TODO: request server for legal moves
  async pieceLegalMoves(square) {
    console.log(square)
    // []string
    let moves = await requestLegalMoves(square)
    console.log(moves)
    // []int
    let moveIndices = moves.map(location => parseIndexSquare(location))
    console.log(moveIndices)
    return moveIndices
  }
}

// parse board state string from server
// and return a board object
function parseBoardState(boardState) {
  let board = new Board()
  board.readBoardState(boardState)
  return board
}

// parse the square index in chess notation
// Input: 0-63
// return a string e.g. 'a1'
function parseSquareIndex(squareIndex) {
  let file = squareIndex % 8
  let rank = 8 - Math.floor(squareIndex / 8)
  let fileString = ''
  let rankString = rank.toString()
  switch (file) {
    case 0:
      fileString = 'a'
      break
    case 1:
      fileString = 'b'
      break
    case 2:
      fileString = 'c'
      break
    case 3:
      fileString = 'd'
      break
    case 4:
      fileString = 'e'
      break
    case 5:
      fileString = 'f'
      break
    case 6:
      fileString = 'g'
      break
    case 7:
      fileString = 'h'
      break
  }
  return fileString + rankString
}

// parse square back to index
function parseIndexSquare(square) {
  let file = 0
  let row = 8 - parseInt(square[1])

  switch (square[0]) {
    case "a":
      file = 0
      break
    case "b":
      file = 1
      break
    case "c":
      file = 2
      break
    case "d":
      file = 3
      break
    case "e":
      file = 4
      break
    case "f":
      file = 5
      break
    case "g":
      file = 6
      break
    case "h":
      file = 7
  }

  return file + row * 8
}

// Request a list of legal moves from the server
// TODO: Finish this part
async function requestLegalMoves(index) {
  let location = parseSquareIndex(index)
  let possibleMoves = await postData("http://localhost:9988/valid_moves", {"piece": location})
  console.log(possibleMoves)
  return possibleMoves.validsquares.split(" ")
}

async function updateBoardState() {
  let response = await postData("http://localhost:9988/state", {})
  console.log(response)
  return response.BoardState
}

function App() {

  const [board, setBoard] = useState(new Board())
  const [selectedSquare, setSelectedSquare] = useState(-1)

  const [legalMoves, setLegalMoves] = useState([])

  const [player, setPlayer] = useState('')

  const handlePieceMove = async (from, to) => {
    let move = parseSquareIndex(from) + parseSquareIndex(to)
    let response = await postData("http://localhost:9988/move", {"move": move})
    console.log(response)
    updateBoardState().then((boardState) => {
      setBoard(parseBoardState(boardState))
    })
    setSelectedSquare(-1)
  }

  // Register to the server
  // and get the board state
  useEffect(() => {
    // Initialization
    fetch('http://localhost:9988/init', {// Adding method type
      method: "POST",
      body: JSON.stringify({uid: '1234'}),
      headers: {"Content-type": "application/json; charset=UTF-8"}
    }).
    then((response) => {
      return response.json()
    })
    .then((data) => {
      console.log(data)
      if (data.color === 1) {
        setPlayer('white')
      } else {
        setPlayer('black')
      }
    }).catch((error) => {
      console.log(error)
      setPlayer('')
    })

    // get the board state every 1 second
    const interval = setInterval(() => {
      // If not registered, try to register again
      if (player === '') {
        fetch('http://localhost:9988/init', {// Adding method type
        method: "POST",
        body: JSON.stringify({
            uid: '1234'
        }),
        headers: {"Content-type": "application/json; charset=UTF-8"}
      })
        .then((response) => {
          return response.json()
        })
        .then((data) => {
          //console.log(data)
          if (data.color === 1) {
            setPlayer('white')
          } else {
            setPlayer('black')
          }
        }).catch((error) => {
          console.log(error)
          setPlayer('')
        })
      }

      fetch('http://localhost:9988/state', {uid: '1234'})
      .then((response) => {
        return response.json()
      })
      .then((data) => {
        //console.log(data)
        setBoard(parseBoardState(data.BoardState))
      }).catch((error) => {
        console.log(error)
        setPlayer('')
      })
    }, 1000)

    // Clear the interval when the component unmounts
    return () => clearInterval(interval)

  },[])


  return (
    <main className='w-auto h-screen bg-stone-800'>

      {player==="" && <PopUp content="Connecting to the server..." />}

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

function PopUp ({content}) {
  return (
    <div className="absolute flex flex-col items-center justify-center z-10 w-screen h-screen bg-stone-900/50">
      <div className="w-60 h-28 bg-slate-50 shadow-md rounded-md flex flex-col items-center justify-start">
        <p className='my-auto'>{content}</p>
      </div>
    </div>
  )
}

function SquareDisplay ({index, square, board, selectedSquare,legalMoves, setSelectedSquare, handlePieceMove, setLegalMoves}) {

  const handleOnClick = async () => {
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
      let moves = await board.pieceLegalMoves(index)
      setLegalMoves(moves)
      setSelectedSquare(index)
    }
    
  }

  return (
    <div 
      onClick={handleOnClick} 
      className={`aspect-square static ${(index===selectedSquare || (legalMoves.includes(index)) && <div className="aspect-square bg-green-500 opacity-50 static"></div>) ? 'bg-red-400' :square.color ? 'bg-stone-300' : 'bg-stone-500'}`}>
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
