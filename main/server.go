package main

// This is the main file for the server
// It spawns a http server handles all the multiplayer logic

type Server struct {
	players map[string]*Player
	rooms   map[string]*Room
}

type Player struct {
	uid  string
	name string
}

type Room struct {
	uid     string
	game    *Game
	players []*Player
}

// type Game is in game.go

// TODO: implement this
func createServer() {

}

// TODO: implement this
func (s *Server) registerPlayer() {

}

// TODO: implement this
func (s *Server) spawnGame() {

}
