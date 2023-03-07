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
}

type Room struct {
	name   string
	status string
	game   *Game
	white  *Player
	black  *Player
}

// Random string generator for room names
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
		// print out that message for clarity
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

	fmt.Println(clientMessage.Type, clientMessage.Data)

	switch clientMessage.Type {
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
		status: "waiting",
		game:   nil,
		white:  player,
		black:  nil,
	}

	// add room to the server
	s.addRoom(room)

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

		server.reader(ws)
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
	}

	// if the room is waiting for a player, change the status
	if room.white == nil || room.black == nil {
		room.status = "waiting"
	}
}

// join player to the room
func (s *Server) joinRoom(room *Room, player *Player) {
	// join player to the room
	if room.white == nil {
		room.white = player
	}
	if room.black == nil {
		room.black = player
	}

	// if the room is full, change the status
	if room.white != nil && room.black != nil {
		// send room status to the players
		s.sendRoomStatus(room)
		room.status = "playing"
	}
}

type RoomStatus struct {
	White string `json:"white"`
	Black string `json:"black"`
}

// send room status to the players
func (s *Server) sendRoomStatus(room *Room) {
	// send room status to the players
	roomStatus := &RoomStatus{
		White: room.white.name,
		Black: room.black.name,
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
