import { useEffect, useState } from 'react'
import useWebSocket from 'react-use-websocket'

const WS_URL = 'wss://gochess.app:8000/ws';
//const WS_URL = 'ws://localhost:8000/ws'
//https://gochess-aber2fx4bq-ue.a.run.app


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

// help icon
import helpIcon from './assets/help-icon.svg'

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

  getBoardState() {
    if (this.state == 1) {
      return "white's turn"
    } else if (this.state == 2) {
      return "black's turn"
    } else if (this.state == 3) {
      return "white won"
    } else if (this.state == 4) {
      return "black won"
    } else if (this.state == 5) {
      return "draw"
    }
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

function App() {

  const [board, setBoard] = useState(new Board())
  const [selectedSquare, setSelectedSquare] = useState(-1)
  const [legalMoves, setLegalMoves] = useState([])

  const [isWhite, setIsWhite] = useState(true)
  const [username, setUsername] = useState('')
  const [opponent, setOpponent] = useState('')

  const [room, setRoom] = useState('')
  const [roomStatus, setRoomStatus] = useState('')
	// 1 - waiting for opponent
	// 2 - game in progress
	// 3 - game over

  const [errorMsg, setErrorMsg] = useState('')

  const [isConnected, setIsConnected] = useState(false)

  // about popup
  const [showAbout, setShowAbout] = useState(false)

  // reset state
  const resetState = () => {
    setBoard(new Board())
    setSelectedSquare(-1)
    setLegalMoves([])
    setIsWhite(true)
    setOpponent('')
    setRoom('')
    setRoomStatus('')
    setErrorMsg('')
  }

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
      //console.log(event)
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

    let boardState;

    switch (data.type) {
      case "roomCreated":
        setRoom(data.data)
        break
      case "roomJoined":
        //setRoom(data.data)
        break
      case "roomStatus":

        // if the room status is complete
        // do nothing
        if (roomStatus === 3) {
          return
        }

        let roomInfo = JSON.parse(data.data)
        // if it is empty string, reinitialize the all the states
        if (roomInfo.name === "") {
          setRoom("")
          setRoomStatus("")
          setIsWhite(true)
          setOpponent("")
          setBoard(new Board())
          setSelectedSquare(-1)
          setLegalMoves([])
          return
        }

        // set room status
        setRoom(roomInfo.name)
        setRoomStatus(roomInfo.status)
        // set player color
        let playerColor = roomInfo.white === username ? "white" : "black"
        setIsWhite(playerColor === "white")
        let opp = roomInfo.white === username ? roomInfo.black : roomInfo.white
        setOpponent(opp)
        break
      case "legalMoves":
        let moves = data.data === "" ? [] : data.data.split(',')
        let moveIndices = moves.map(move => parseIndexSquare(move))
        setLegalMoves(moveIndices)
        break

      case "gameState":
        boardState = parseBoardState(data.data)
        setBoard(boardState)
        break

      case "gameResult":
        // This means the game is over
        setRoomStatus(3)
        boardState = parseBoardState(data.data)
        setBoard(boardState)
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
    console.log("send join room request")
    sendJsonMessage({type: "joinRoom", data: room})
  }

  // request list of rooms
  const handleRequestRooms = async () => {
    sendJsonMessage({type: "requestRooms", data: ""})
  }

  // request legal moves for a piece
  const handleRequestLegalMoves = async (index) => {
    // if the player is not in a game, do nothing
    if (roomStatus !== 2) {
      return
    }
    let location = parseSquareIndex(index)
    sendJsonMessage({type: "requestLegalMoves", data: location})
  }

  const handlePieceMove = async (from, to) => {
    let fromSquare = parseSquareIndex(from)
    let toSquare = parseSquareIndex(to)
    sendJsonMessage({type: "movePiece", data: fromSquare + "," + toSquare})
  }

  // ping the server to stay connected
  useEffect(() => {
    const interval = setInterval(() => {
      sendJsonMessage({type: "ping", data: ""})
    }, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <main className='w-auto h-screen bg-stone-800'>

      <button className='absolute top-2 left-2 z-50 font-bold bg-stone-600 text-stone-200 p-2 text-sm shadow-lg rounded' onClick={() => setShowAbout(true)}>
        About this project
      </button>

      {showAbout && <AboutPopup setShowAbout={setShowAbout}/>}

      {!isConnected && <PopUp content="Connecting to the server..." />}
      {isConnected && !username && <UsernamePopUp handleUsernameChange={handleUsernameChange} />}
      {isConnected && username && !room && <RoomPopUp handleCreateRoom={handleCreateRoom} handleJoinRoom={handleJoinRoom} />}
      {isConnected && username && room && roomStatus === 1 && <PopUp content={`Waiting for opponent...\n Room code: ${room}`} />}

      {isConnected && username && room && roomStatus === 3 && <GameOverPopUp board={board} resetState={resetState} />}

      <div className='flex flex-row items-center justify-center h-full '>

        {!!errorMsg && <ErrorMsg content={errorMsg} />}

        <GameBoard 
          isWhite={isWhite}
          roomStatus = {roomStatus}
          board={board}
          selectedSquare={selectedSquare} 
          legalMoves={legalMoves}
          setSelectedSquare={setSelectedSquare}
          handlePieceMove={handlePieceMove}
          setLegalMoves={setLegalMoves}
          handleRequestLegalMoves={handleRequestLegalMoves}
        />

        <div className="bg-stone-700 w-36 h-96 p-2 text-sm">
              <p className='text-white'>{username} v.s. {opponent}</p>
              <p className='text-white'>Piece: {isWhite ? "White" : "Black"}</p>
              <p className='text-white'>Room: {room}</p>
              <p className='text-white'>Status: {board.getBoardState()}</p>
        </div>

      </div>
    </main>
  )
}

function GameBoard ({isWhite, roomStatus, board, selectedSquare, legalMoves, setSelectedSquare, handlePieceMove, setLegalMoves, handleRequestLegalMoves}) {
  return (
    <div className='flex flex-col items-center justify-center'>
          <HorizontalLabel />
          <div className='flex flex-row'>
            <VerticalLabel />
            <div className="w-96 aspect-square grid grid-cols-8 grid-rows-8 shadow-xl">
              {
                isWhite && board.squares.map((square, index) => (
              <SquareDisplay 
                isWhite={isWhite}
                key={index} 
                index={index} 
                square={square}
                board={board}
                selectedSquare={selectedSquare} 
                legalMoves={legalMoves}
                setSelectedSquare={setSelectedSquare}
                handlePieceMove={handlePieceMove}
                setLegalMoves={setLegalMoves}
                handleRequestLegalMoves={handleRequestLegalMoves}
                  /> ))
              }
              {
                !isWhite && board.squares.slice().reverse().map((square, index) => (
              <SquareDisplay
                isWhite={isWhite}
                key={blackIndex(index)}
                index={blackIndex(index)}
                square={square}
                board={board}
                selectedSquare={selectedSquare}
                legalMoves={legalMoves}
                setSelectedSquare={setSelectedSquare}
                handlePieceMove={handlePieceMove}
                setLegalMoves={setLegalMoves}
                handleRequestLegalMoves={handleRequestLegalMoves}
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
        <p className='my-auto mx-2'>{content}</p>
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

function GameOverPopUp ({board, resetState}) {
  const result = board.state === 2 ? 'White Wins!' : board.state === 3 ? 'Black Wins!' : 'Draw!'
  return (
    <div className="absolute flex flex-col items-center justify-center z-10 w-screen h-screen backdrop-blur-sm bg-stone-900/50">
      <div className="w-60 p-5 bg-slate-50 shadow-md rounded-md flex flex-col items-center justify-start">
        <p className='my-4'>{result}</p>
        <button className="w-full h-8 rounded-md bg-stone-800 my-4 text-white" onClick={resetState}>Play Again</button>
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

function SquareDisplay ({isWhite, index, square, board, selectedSquare,legalMoves, setSelectedSquare, handlePieceMove, setLegalMoves, handleRequestLegalMoves}) {

  const handleOnClick = async () => {
    //console.log("clicked", index)
    setLegalMoves([])

    // if the square is already selected
    // cancel selection
    if (selectedSquare === index) {
      setSelectedSquare(-1)
      return
    }

    // if there is a piece on the square
    // try to make move
    if (selectedSquare !== -1) {
      // check if legal move
      if (legalMoves.includes(index)){
        handlePieceMove(selectedSquare, index)
      }
      setSelectedSquare(-1)
      return
    }
    
    // if there is a piece on the square
    // set legal moves
    if (square.piece){
      // if the piece is not the current player's piece
      // do nothing
      console.log("piece", square.piece)
      if (square.piece.color === "white" && !isWhite) return
      // if it is not the current player's turn
      // do nothing
      if (board.state === 1 && !isWhite) return
      await handleRequestLegalMoves(index)
      setSelectedSquare(index)
    }
  }

  return (
    <div 
      onClick={handleOnClick} 
      className={`aspect-square relative ${(index===selectedSquare || (legalMoves.includes(index)) && <div className="aspect-square bg-green-500 opacity-50 static"></div>) ? 'bg-red-400' :square.color ? 'bg-stone-300' : 'bg-stone-500'}`}>
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

function AboutPopup({setShowAbout}) {
  return(
    <div  className="absolute flex flex-col items-center justify-center z-50 w-screen h-screen backdrop-blur-sm bg-stone-900/50 text-white text-sm">
      <div className="max-w-xl p-5 bg-stone-600 shadow-md rounded-md flex flex-col">
        <p className='my-2'>About GoChess (v0.1):</p>
        <p className='my-2'>GoChess is a chess game built with Go and React. I'm still testing the server! Feel free to raise any issues on Github: <a className="underline"href="https://github.com/zikunw/GoChess">https://github.com/zikunw/GoChess</a></p>
        <p className='my-2'><a className="underline underline-offset-2 font-bold shadow-none hover:shadow-2xl" href="https://www.zikunw.com">Learn more about what I do here!</a></p>
        <p className='my-2'>© 2023 Zikun Wang. All rights reserved.</p>
        <button className='my-2 underline border py-2 border-white rounded-lg' onClick={()=>setShowAbout(false)}>Back</button>
      </div>
    </div>
  )
}

export default App
