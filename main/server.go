package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// This is the main file for the server
// It spawns a http server handles all the multiplayer logic

type Server struct {
	players map[string]*Player
	rooms   map[string]*Room

	mu sync.Mutex
}

type Player struct {
	conn *websocket.Conn
	name string
	room *Room
}

type Room struct {
	name   string
	status int
	// 1 - waiting for opponent
	// 2 - game in progress
	// 3 - game over
	game  *Game
	white *Player
	black *Player
}

// Random string generator for room names
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// pongWait is the maximum time in seconds to wait for a pong
//const pongWait = 1 * time.Second

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// type Game is in game.go

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// reader reads messages from the websocket connection
func (s *Server) reader(conn *websocket.Conn) {
	// register user
	// create a player
	player := &Player{
		conn: conn,
		name: "",
	}
	s.addPlayer(player)

	//conn.SetReadDeadline(time.Now().Add(pongWait))
	//conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// reinitialize the player board status
	s.sendEmptyStatus(player)

	for {
		// read in a message
		// we assume all mesaages are json
		_, p, err := conn.ReadMessage()
		if err != nil {
			// if the connection is closed, remove the player
			s.disconnect(player)
			log.Println(err)
			return
		}

		go s.respond(player, p)
	}
}

// Client message
type ClientMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// Server message
type ServerMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// Server Error message
type ServerError struct {
	Error string `json:"error"`
}

// server responds to messages from the client
func (s *Server) respond(player *Player, message []byte) {
	//fmt.Println(string(message))
	// parse json message
	clientMessage := &ClientMessage{}
	err := json.Unmarshal(message, clientMessage)
	if err != nil {
		log.Println(err)
		return
	}

	//fmt.Println(clientMessage.Type, clientMessage.Data)

	switch clientMessage.Type {
	// ping
	case "ping":
		// do nothing

	case "registerUsername":
		player.name = clientMessage.Data
		fmt.Println("Registered username: ", player.name)

	// create a room
	case "createRoom":
		roomName := s.createRoom(player)
		fmt.Println("Created room: ", roomName)
		// send room name back to client
		serverMessage := &ServerMessage{
			Type: "roomCreated",
			Data: roomName,
		}
		serverMessageJSON, err := json.Marshal(serverMessage)
		if err != nil {
			log.Println(err)
			return
		}
		err = player.conn.WriteMessage(websocket.TextMessage, serverMessageJSON)
		if err != nil {
			log.Println(err)
			return
		}

	// join a room
	case "joinRoom":
		fmt.Println("Joining room: ", clientMessage.Data)
		roomName := clientMessage.Data
		room, ok := s.rooms[roomName]
		// check if the room exists
		if !ok {
			// send error message back to client
			serverError := &ServerError{
				Error: "Room does not exist",
			}
			serverErrorJSON, err := json.Marshal(serverError)
			if err != nil {
				log.Println(err)
				return
			}
			err = player.conn.WriteMessage(websocket.TextMessage, serverErrorJSON)
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
		// check if the room is full
		if room.black != nil && room.white != nil {
			// send error message back to client
			serverError := &ServerError{
				Error: "Room is full",
			}
			serverErrorJSON, err := json.Marshal(serverError)
			if err != nil {
				log.Println(err)
				return
			}
			err = player.conn.WriteMessage(websocket.TextMessage, serverErrorJSON)
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
		// join the room
		s.joinRoom(room, player)
		// send room name back to client
		serverMessage := &ServerMessage{
			Type: "roomJoined",
			Data: roomName,
		}
		serverMessageJSON, err := json.Marshal(serverMessage)
		if err != nil {
			log.Println(err)
			return
		}
		err = player.conn.WriteMessage(websocket.TextMessage, serverMessageJSON)
		if err != nil {
			log.Println(err)
			return
		}

	// request room names
	// for debugging
	case "requestRooms":
		rooms := []string{}
		for _, room := range s.rooms {
			rooms = append(rooms, room.name)
		}
		// send room names back to client
		serverMessage := &ServerMessage{
			Type: "rooms",
			Data: strings.Join(rooms, ","),
		}
		serverMessageJSON, err := json.Marshal(serverMessage)
		if err != nil {
			log.Println(err)
			return
		}
		err = player.conn.WriteMessage(websocket.TextMessage, serverMessageJSON)
		if err != nil {
			log.Println(err)
			return
		}

	// request legal moves
	case "requestLegalMoves":
		// check if player is in a room
		if player.room == nil {
			fmt.Println("Player is not in a room")
			return
		}
		// check if player is in a game
		if player.room.game == nil {
			fmt.Println("There is no game in this room")
			return
		}
		// check if it is the player's turn
		playerColor := -1
		if player == player.room.white {
			playerColor = 1
		} else if player == player.room.black {
			playerColor = 2
		} else {
			return
		}
		if playerColor != player.room.game.Board.State {
			return
		}

		// get legal moves
		// Check the location
		isValid, location := GridToLocation(clientMessage.Data)
		if isValid == false {
			fmt.Println("Invalid location received from client")
			return
		}

		// Check the piece
		piece := player.room.game.Board.GetPieceAtLocation(location)
		if piece.IsEmpty() == true {
			fmt.Println("Empty piece received from client")
			return
		}
		if piece.GetPlayer() != playerColor {
			fmt.Println("Invalid piece received from client")
			return
		}

		// Get the valid moves
		moves := []Move{}
		for _, move := range player.room.game.Board.GetPlayerLegalMoves(playerColor) {
			if move.From == location {
				moves = append(moves, move)
			}
		}

		validSquare := []string{}
		for _, move := range moves {
			validSquare = append(validSquare, LocationToGrid(move.To))
		}

		// send legal moves back to client
		serverMessage := &ServerMessage{
			Type: "legalMoves",
			Data: strings.Join(validSquare, ","),
		}
		serverMessageJSON, err := json.Marshal(serverMessage)
		if err != nil {
			log.Println(err)
			return
		}
		err = player.conn.WriteMessage(websocket.TextMessage, serverMessageJSON)
		if err != nil {
			log.Println(err)
			return
		}

	// make a move
	case "movePiece":
		// check if player is in a room
		if player.room == nil {
			fmt.Println("Player is not in a room")
			return
		}
		// check if player is in a game
		if player.room.game == nil {
			fmt.Println("There is no game in this room")
			return
		}
		// check if it is the player's turn
		playerColor := -1
		if player == player.room.white {
			playerColor = 1
		} else if player == player.room.black {
			playerColor = 2
		} else {
			return
		}
		if playerColor != player.room.game.Board.State {
			return
		}

		// check locations
		fromLocation := strings.Split(clientMessage.Data, ",")[0]
		toLocation := strings.Split(clientMessage.Data, ",")[1]
		isValid, from := GridToLocation(fromLocation)
		if isValid == false {
			fmt.Println("Invalid location received from client")
			return
		}
		isValid, to := GridToLocation(toLocation)
		if isValid == false {
			fmt.Println("Invalid location received from client")
			return
		}

		// check if the move is legal
		var move Move
		isValid, move = ValidMove(from, to, playerColor, &player.room.game.Board)
		if isValid == false {
			fmt.Println("Invalid move received from client")
			break
		}

		// make the move
		if playerColor == 1 {
			player.room.game.WhitePlayer.SendMove(move)
		} else if playerColor == 2 {
			player.room.game.BlackPlayer.SendMove(move)
		}

	default:
		fmt.Println("Unknown message type: ", clientMessage.Type)
	}
}

// createRoom creates a room
func (s *Server) createRoom(player *Player) string {
	// generate a random room name
	rand.Seed(time.Now().UnixNano())
	roomName := randSeq(5)
	_, ok := s.rooms[roomName]
	// Check if the room name is already taken
	for ok {
		rand.Seed(time.Now().UnixNano())
		roomName = randSeq(5)
		_, ok = s.rooms[roomName]
	}

	// create a room
	room := &Room{
		name:   roomName,
		status: 1, // 1 = waiting for players
		game:   nil,
		white:  player,
		black:  nil,
	}
	player.room = room

	// add room to the server
	s.addRoom(room)

	// broadcast room creation to all players
	s.sendRoomStatus(room)

	// return
	return roomName
}

// addRoom adds a room to the server
func (s *Server) addRoom(room *Room) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rooms[room.name] = room
}

// add player to the server
func (s *Server) addPlayer(player *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.players[player.name] = player
}

// createServer creates a socket.io server
func createServer() {

	server := &Server{
		players: make(map[string]*Player),
		rooms:   make(map[string]*Room),
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("Client connected")

		go server.reader(ws)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))

}

// Disconnect removes the player from the server
// and removes the player from the room if they are in one
func (s *Server) disconnect(player *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// remove player from the server
	delete(s.players, player.name)

	// remove player from the room
	for _, room := range s.rooms {
		if room.white == player {
			s.leaveRoom(room, player)
		}
		if room.black == player {
			s.leaveRoom(room, player)
		}
	}
}

// remove player from the room
func (s *Server) leaveRoom(room *Room, player *Player) {
	player.room = nil

	// remove player from the room
	if room.white == player {
		room.white = nil
	}
	if room.black == player {
		room.black = nil
	}

	// if the room is empty, delete it
	if room.white == nil && room.black == nil {
		delete(s.rooms, room.name)
		return
	}

	// if the room is waiting for a player, change the status
	if room.white == nil || room.black == nil {
		room.status = 1 // 1 = waiting for players
		fmt.Println("Player left room, waiting for players")
	}

	// send room status to the players
	s.sendRoomStatus(room)
}

// join player to the room
func (s *Server) joinRoom(room *Room, player *Player) {
	// join player to the room
	if room.white == nil {
		room.white = player
		player.room = room
	}
	if room.black == nil {
		room.black = player
		player.room = room
	}

	// if the room is full, change the status
	if room.white != nil && room.black != nil {
		// start the game
		s.startGame(room)
	}
}

// StartGame starts the game
func (s *Server) startGame(room *Room) {

	room.status = 2 // 2 = game in progress
	// send room status to the players
	s.sendRoomStatus(room)

	// create a game
	whitePlayerController := &RemotePlayer{}
	blackPlayerController := &RemotePlayer{}
	game := &Game{}
	game.Init(whitePlayerController, blackPlayerController)
	// init players
	whitePlayerController.Init(1, &game.Board)
	blackPlayerController.Init(2, &game.Board)

	// add game to the room
	room.game = game

	// game loop
	// TODO:
	for {
		fmt.Println("")
		fmt.Println(game.Board.FullmoveNumber)
		game.Print()

		// send game state to the players
		gameState := game.Board.Serialize()
		room.white.conn.WriteJSON(&ClientMessage{
			Type: "gameState",
			Data: gameState,
		})
		room.black.conn.WriteJSON(&ClientMessage{
			Type: "gameState",
			Data: gameState,
		})

		isEnd := game.Play()
		if isEnd {
			break
		}
	}

	// send game result to the players
}

type RoomStatus struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
	White  string `json:"white"`
	Black  string `json:"black"`
}

// send room status to the players
func (s *Server) sendRoomStatus(room *Room) {
	// send room status to the players
	if room.black == nil && room.white == nil {
		return
	}

	var blackPlayerName string
	var whitePlayerName string

	if room.black == nil {
		blackPlayerName = ""
	} else {
		blackPlayerName = room.black.name
	}
	if room.white == nil {
		whitePlayerName = ""
	} else {
		whitePlayerName = room.white.name
	}

	roomStatus := &RoomStatus{
		Name:   room.name,
		Status: room.status,
		White:  whitePlayerName,
		Black:  blackPlayerName,
	}

	roomStatusJSON, err := json.Marshal(roomStatus)
	if err != nil {
		log.Println(err)
		return
	}
	serverMessage := &ServerMessage{
		Type: "roomStatus",
		Data: string(roomStatusJSON),
	}
	serverMessageJSON, err := json.Marshal(serverMessage)
	if err != nil {
		log.Println(err)
		return
	}
	if room.white != nil {
		err = room.white.conn.WriteMessage(websocket.TextMessage, serverMessageJSON)
		if err != nil {
			log.Println(err)
			return
		}
	}
	if room.black != nil {
		err = room.black.conn.WriteMessage(websocket.TextMessage, serverMessageJSON)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// send empty status to reinitialize the client
func (s *Server) sendEmptyStatus(player *Player) {
	roomStatus := &RoomStatus{
		Name:   "",
		Status: 0,
		White:  "",
		Black:  "",
	}

	roomStatusJSON, err := json.Marshal(roomStatus)
	if err != nil {
		log.Println(err)
		return
	}
	serverMessage := &ServerMessage{
		Type: "roomStatus",
		Data: string(roomStatusJSON),
	}
	serverMessageJSON, err := json.Marshal(serverMessage)
	if err != nil {
		log.Println(err)
		return
	}
	err = player.conn.WriteMessage(websocket.TextMessage, serverMessageJSON)
	if err != nil {
		log.Println(err)
		return
	}
}
