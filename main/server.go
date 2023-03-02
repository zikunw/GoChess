package main

import (
	"fmt"
	"net/http"
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

type RemotePlayer struct {
	Color int
	uid   string
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

	return Move{}
}

// TODO: implement this
func (p *RemotePlayer) MakeServer() {
	// set all the handlers
	http.HandleFunc("/init", p.initHandler)
	// start the server
	fmt.Println("Server started at " + SERVER_HOST + ":" + SERVER_PORT)
	go http.ListenAndServe(SERVER_HOST+":"+SERVER_PORT, nil)
}

type InitResponse struct {
	Color int
}

// TODO: implement this
// The client will send a request with uid to this handler
// and the server will register the client
func (p *RemotePlayer) initHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("initHandler called")

	// Recieve the uid from the client
	// register the client
	p.uid = r.FormValue("uid")

	// Send the color to the client
	fmt.Fprintf(w, "Hello World!")
}
