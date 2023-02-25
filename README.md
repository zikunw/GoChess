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

## My current progress
- [x] basic game logic (check detection, legal moves etc)
- [ ] overall game logic (starting a game, turns, etc)
- [ ] player control from commandline
- [ ] simple chess playing bots
- [ ] A javascript interface (maybe)

