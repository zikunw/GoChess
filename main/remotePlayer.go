package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// External players and API stuff

// Here im using a rest http api in order for
// the javascript client to communicate with the server
// I will replace this once I have better method.

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

type RemotePlayer struct {
	Color int
	uid   string

	playerMove chan Move

	board *Board
}

// Initialize remote player
func (p *RemotePlayer) Init(color int, b *Board) {
	p.Color = color
	p.uid = ""
	p.playerMove = make(chan Move)
	p.board = b
}

// TODO: implement this
func (p *RemotePlayer) GetMove(b *Board) Move {
	// ==== Initialization ==== //
	// If the uid is not set, then the client has not registered yet
	if p.uid == "" {
		p.MakeServer()
	}
	// Wait until the client register
	for p.uid == "" {
		time.Sleep(100 * time.Millisecond)
	}

	// ==== Get the move ==== //
	fmt.Println("Waiting for move from client")
	// Wait for the client to send the move
	move := <-p.playerMove
	fmt.Println("Move received from client")

	return move
}

// TODO: implement this
func (p *RemotePlayer) MakeServer() {
	// set all the handlers
	http.HandleFunc("/init", p.initHandler)
	http.HandleFunc("/state", p.boardHandler)
	http.HandleFunc("/move", p.moveHandler)
	http.HandleFunc("/valid_moves", p.validMovesHandler)
	// start the server
	fmt.Println("Server started at " + SERVER_HOST + ":" + SERVER_PORT)
	go http.ListenAndServe(SERVER_HOST+":"+SERVER_PORT, nil)
}

type InitRequest struct {
	uid int
}

type InitResponse struct {
	Color int
}

// TODO: implement this
// The client will send a request with uid to this handler
// and the server will register the client
func (p *RemotePlayer) initHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	//fmt.Println("initHandler called")

	if p.uid != "" {
		fmt.Println("Client already registered")
		fmt.Fprintf(w, "client already registered")
		return
	}

	// Recieve the uid from the client
	// register the client
	reqBody, _ := ioutil.ReadAll(r.Body)
	var initRequest InitRequest
	json.Unmarshal(reqBody, &initRequest)

	fmt.Println(string(reqBody))

	response := InitResponse{Color: p.Color}
	jsonResponse, _ := json.Marshal(response)

	// Send the color to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// The client send the move to the server
// TODO: Need to return a json
func (p *RemotePlayer) moveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("moveHandler called")

	// Recieve the move from the client
	isValid, move := DeserializeMove(r.FormValue("move"), p.Color, p.board)

	if isValid == false {
		fmt.Println("Invalid move received from client")
		fmt.Fprintf(w, "invalid move")
		return
	}

	fmt.Fprintf(w, "move received")

	// Send the move to the player
	p.playerMove <- move
}

type validMoveRequest struct {
	Piece string `json:"piece"`
}

type validMoveResponse struct {
	Err          string `json:"err"`
	ValidSquares string `json:"validsquares"`
}

// The client will request valid moves from the server
// TODO: Need to return a json
func (p *RemotePlayer) validMovesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("validMovesHandler called")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var request validMoveRequest
	json.Unmarshal(reqBody, &request)

	// Check the location
	fmt.Println(request.Piece)
	fmt.Println(GridToLocation(request.Piece))
	isValid, location := GridToLocation(request.Piece)
	if isValid == false {
		fmt.Println("Invalid location received from client")
		response := validMoveResponse{"Invalid location", ""}
		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)
		return
	}

	// Check the piece
	piece := p.board.GetPieceAtLocation(location)
	if piece.IsEmpty() == true {
		fmt.Println("Empty piece received from client")
		response := validMoveResponse{"Empty piece", ""}
		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)
		return
	}
	if piece.GetPlayer() != p.Color {
		fmt.Println("Invalid piece received from client")
		response := validMoveResponse{"Invalid piece", ""}
		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)
		return
	}

	// Get the valid moves
	moves := []Move{}
	for _, move := range p.board.GetPlayerLegalMoves(p.Color) {
		if move.From == location {
			moves = append(moves, move)
		}
	}

	//fmt.Println(moves)

	validSquare := []string{}
	for _, move := range moves {
		validSquare = append(validSquare, LocationToGrid(move.To))
	}

	//fmt.Println(validSquare)

	// Send the valid moves to the client
	response := validMoveResponse{"", strings.Join(validSquare, " ")}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Not able to convert to json")
	}
	fmt.Println(string(jsonResponse))
	w.Write(jsonResponse)
}

type BoardStateResponse struct {
	BoardState string
}

// return the current board state to the client
func (p *RemotePlayer) boardHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	//fmt.Println("boardHandler called")

	boardState := p.board.Serialize()
	boardStateResponse := BoardStateResponse{BoardState: boardState}
	jsonResponse, _ := json.Marshal(boardStateResponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}
