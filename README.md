# GoChess

GoChess is currently work-in-progress.

## Usage

Fork the project and run the following command in the git directory:

```
go run main/*.go
```

You can generate a board state using FEN under `main()` in `main/init.go`:

```go
board := InitFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
board.PrintWithBorder()
```

Output:

```
  a b c d e f g h 
8 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ 8
7 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙ 7
6                 6
5                 5
4                 4
3                 3
2 ♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟ 2
1 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ 1
  a b c d e f g h 
```

## Agents

Rather than directly prompting user for commandline input, I wrote two simple AI implementation.

1. The first one is `RandomComputerPlayer` which plays its legal moves randomly.
2. The second one is `MinimaxComputerPlayer` which uses a min-max algorithm with piece value as its
heuristic to cauculate the next best move.

Run the following command in `main()` will pitch two minmax AI against each other:

```go
	game := Game{}
	game.Init(&MinimaxComputerPlayer{1, 3}, &MinimaxComputerPlayer{2, 2})
	for {
		fmt.Println("")
		fmt.Println(game.Board.FullmoveNumber)
		game.Print()
		isEnd := game.Play()
		if isEnd {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	game.Print()
```

## My current progress
- [x] basic game logic (check detection, legal moves etc)
- [x] overall game logic (starting a game, turns, etc)
- [x] player control from commandline
- [x] simple chess playing bots
- [ ] A javascript interface (maybe)

