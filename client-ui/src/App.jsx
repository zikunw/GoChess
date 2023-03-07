import { useEffect, useState } from 'react'
import useWebSocket from 'react-use-websocket'

const WS_URL = 'ws://127.0.0.1:8000/ws';

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

// cauculate the correct index when play black
function blackIndex(index) {
  return 63 - index
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

  const [isWhite, setIsWhite] = useState(true)
  const [username, setUsername] = useState('')
  const [opponent, setOpponent] = useState('')
  const [room, setRoom] = useState('')

  const [errorMsg, setErrorMsg] = useState('')

  const [isConnected, setIsConnected] = useState(false)

  const {
    sendMessage,
    sendJsonMessage,
    lastMessage,
    lastJsonMessage,
    readyState,
    getWebSocket
  } = useWebSocket(WS_URL, {
    onOpen: () => {
      setIsConnected(true) 
      console.log('WebSocket connection established.');
    },
    onMessage: (event) => {
      console.log(event)
      handleIncomingMessage(event.data)
    },
    onError: (event) => {
      console.log(event)
    },
    onClose: () => {
      setIsConnected(false)
      console.log('WebSocket connection closed.');
    },
  });

  // handle incoming messages
  const handleIncomingMessage = (message) => {
    console.log(message)
    let data = JSON.parse(message)

    // Check error message
    if (data.error) {
      setErrorMsg(data.error)
      return
    }

    switch (data.type) {
      case "roomCreated":
        setRoom(data.data)
        break
      case "roomJoined":
        setRoom(data.data)
        break
      case "roomStatus":
        let playerColor = JSON.parse(data.data).white === username ? "white" : "black"
        setIsWhite(playerColor === "white")
        let opp = JSON.parse(data.data).white === username ? JSON.parse(data.data).black : JSON.parse(data.data).white
        setOpponent(opp)
        break

    }
  }

  // handle username change and register to the server
  const handleUsernameChange = async (username) => {
    setUsername(username)
    sendJsonMessage({type: "registerUsername", data: username})
  }

  // create a new room
  const handleCreateRoom = async () => {
    console.log("create room")
    sendJsonMessage({type: "createRoom", data: ""})
  }

  // join a room
  const handleJoinRoom = async (room) => {
    sendJsonMessage({type: "joinRoom", data: room})
  }

  // request list of rooms
  const handleRequestRooms = async () => {
    sendJsonMessage({type: "requestRooms", data: ""})
  }

  const handlePieceMove = async (from, to) => {

  }


  return (
    <main className='w-auto h-screen bg-stone-800'>

      {!isConnected && <PopUp content="Connecting to the server..." />}
      {isConnected && !username && <UsernamePopUp handleUsernameChange={handleUsernameChange} />}
      {isConnected && username && !room && <RoomPopUp handleCreateRoom={handleCreateRoom} handleJoinRoom={handleJoinRoom} />}


      <div className='flex flex-row items-center justify-center h-full '>

        {!!errorMsg && <ErrorMsg content={errorMsg} />}

        <GameBoard 
          isWhite={isWhite}
          board={board}
          selectedSquare={selectedSquare} 
          legalMoves={legalMoves}
          setSelectedSquare={setSelectedSquare}
          handlePieceMove={handlePieceMove}
          setLegalMoves={setLegalMoves}
        />

        <div className="bg-stone-700 w-36 h-96 p-2">
              <p className='text-white'>{username} v.s. {opponent}</p>
        </div>

      </div>
    </main>
  )
}

function GameBoard ({isWhite, board, selectedSquare, legalMoves, setSelectedSquare, handlePieceMove, setLegalMoves}) {
  return (
    <div className='flex flex-col items-center justify-center'>
          <HorizontalLabel />
          <div className='flex flex-row'>
            <VerticalLabel />
            <div className="w-96 aspect-square grid grid-cols-8 grid-rows-8 shadow-xl">
              {
                isWhite && board.squares.map((square, index) => (
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
              {
                !isWhite && board.squares.slice().reverse().map((square, index) => (
              <SquareDisplay
                key={blackIndex(index)}
                index={blackIndex(index)}
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
            <VerticalLabel />
          </div>
          <HorizontalLabel />
        </div>
  )
}

function ErrorMsg (props) {
  return (
    <div className="bg-red-500 z-30 absolute left-2 top-2 w-96 h-8 flex flex-row items-center justify-center">
      <p className="text-white font-mono font-thin text-sm">Error: {props.content}</p>
    </div>
  )
}

function VerticalLabel () {
  return (
    <div className="w-12 h-96 flex flex-col items-center justify-around">
      <p className="text-stone-400 text-md font-mono font-thin">8</p>
      <p className="text-stone-400 text-md font-mono font-thin">7</p>
      <p className="text-stone-400 text-md font-mono font-thin">6</p>
      <p className="text-stone-400 text-md font-mono font-thin">5</p>
      <p className="text-stone-400 text-md font-mono font-thin">4</p>
      <p className="text-stone-400 text-md font-mono font-thin">3</p>
      <p className="text-stone-400 text-md font-mono font-thin">2</p>
      <p className="text-stone-400 text-md font-mono font-thin">1</p>
    </div>
  )
}

function HorizontalLabel () {
  return (
    <div className="w-96 h-12 flex flex-row items-center justify-around">
      <p className="text-stone-400 text-md font-mono font-thin">a</p>
      <p className="text-stone-400 text-md font-mono font-thin">b</p>
      <p className="text-stone-400 text-md font-mono font-thin">c</p>
      <p className="text-stone-400 text-md font-mono font-thin">d</p>
      <p className="text-stone-400 text-md font-mono font-thin">e</p>
      <p className="text-stone-400 text-md font-mono font-thin">f</p>
      <p className="text-stone-400 text-md font-mono font-thin">g</p>
      <p className="text-stone-400 text-md font-mono font-thin">h</p>
    </div>
  )
}

function PopUp ({content}) {
  return (
    <div className="absolute flex flex-col items-center justify-center z-10 w-screen h-screen backdrop-blur-sm bg-stone-900/50">
      <div className="w-60 h-28 bg-slate-50 shadow-md rounded-md flex flex-col items-center justify-start">
        <p className='my-auto'>{content}</p>
      </div>
    </div>
  )
}

function UsernamePopUp ({handleUsernameChange}) {
  const [usernameInput, setUsernameInput] = useState('')
  

  return (
    <div className="absolute flex flex-col items-center justify-center z-10 w-screen h-screen backdrop-blur-sm bg-stone-900/50">
      <div className="w-60 p-5 bg-slate-50 shadow-md rounded-md flex flex-col items-center justify-start">
        <p className='my-4'>Enter your username:</p>
        <input className="w-full px-2 h-8 mb-2 rounded-md border-2" type="text" onChange={(e) => setUsernameInput(e.target.value)} />
        <button className="w-full h-8 rounded-md bg-stone-800 my-4 text-white" onClick={() => handleUsernameChange(usernameInput)}>Submit</button>
      </div>
    </div>
  )
}

function RoomPopUp ({handleCreateRoom, handleJoinRoom}) {
  const [roomInput, setRoomInput] = useState('')

  return (
    <div className="absolute flex flex-col items-center justify-center z-10 w-screen h-screen backdrop-blur-sm bg-stone-900/50">
      <div className="w-60 p-2 bg-slate-50 shadow-md rounded-md flex flex-col items-center justify-start divide-y-2">
        <div className='w-full'>
          <p className='mt-2'>Create a room?</p>
          <button className="w-full p-2 rounded-md bg-stone-800 my-4 text-white" onClick={handleCreateRoom}>Create</button>
        </div>
        
        <div>
          <p className='my-4'>Or join a room:</p>
          <div className="w-full mb-2 flex flex-row items-center justify-center">
            <input className="w-full p-2  rounded-md border-2" type="text" onChange={(e) => setRoomInput(e.target.value)}/>
            <button className="p-2 rounded-md bg-stone-800  text-white" onClick={() => handleJoinRoom(roomInput)}>Join</button>
          </div>
        </div>
      </div>
    </div>
  )
}

function SquareDisplay ({index, square, board, selectedSquare,legalMoves, setSelectedSquare, handlePieceMove, setLegalMoves}) {

  const handleOnClick = async () => {
    console.log("clicked", index)
    setLegalMoves([])

    // if the square is already selected
    // cancel selection
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
      className={`aspect-square relative ${(index===selectedSquare || (legalMoves.includes(index)) && <div className="aspect-square bg-green-500 opacity-50 static"></div>) ? 'bg-red-400' :square.color ? 'bg-stone-300' : 'bg-stone-500'}`}>
        <p className="absolute left-1 top-0 text-sm">{index}</p>
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
