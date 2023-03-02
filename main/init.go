package main

import "fmt"

func main() {

	remotePlayer := &RemotePlayer{}

	game := Game{}
	game.Init(remotePlayer, &MinimaxComputerPlayer{2, 3})

	remotePlayer.Init(1, &game.Board)

	game.Print()
	for {
		fmt.Println("")
		fmt.Println(game.Board.FullmoveNumber)
		game.Print()
		isEnd := game.Play()
		if isEnd {
			break
		}
	}
	game.Print()

	//b := Board{}
	//b.Init()
	////b.Print()
	//fmt.Println(b.Serialize())

	//b2 := Board{}
	//b2.Init()
	//b2.Deserialize("RNBQKBNRPPPPEEPPEEEEEEEEEEEEEQqEEEEEEEEEEEEEEEEEpppppppprnbqkbnr 0 1 true true true true")
	//b2.Print()
	//fmt.Print(b2.Serialize())
}
