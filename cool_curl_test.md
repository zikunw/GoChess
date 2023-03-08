``` shell
curl --data "uid=42" localhost:9988/init
curl localhost:9988/state
curl --data "move=d2d4" localhost:9988/move
curl --data "location=e2" localhost:9988/valid_moves
```

Json
``` shell
curl -X POST localhost:9988/init -H 'Content-Type: application/json' -d '{"uid":"1234"}'
curl -X POST localhost:9988/move -H 'Content-Type: application/json' -d '{"move":"d2d4"}'
```

fen/state/halfmove/fullmove/wqcastle/wkcastle/bqcastle/bkcastle

{
    BoardState: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR 1 0 1 true true true true"
}


&{GET /valid_moves HTTP/1.1 1 1 map[Accept:[*/*] Accept-Encoding:[gzip, deflate] Accept-Language:[en-US,en;q=0.9] Connection:[keep-alive] Origin:[http://localhost:5173] Referer:[http://localhost:5173/] User-Agent:[Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Safari/605.1.15]] {} <nil> 0 [] false localhost:9988 map[] map[] <nil> map[] 127.0.0.1:62540 /valid_moves <nil> <nil> <nil> 0x14000290000}

&{POST /valid_moves HTTP/1.1 1 1 map[Accept:[*/*] Content-Length:[11] Content-Type:[application/x-www-form-urlencoded] User-Agent:[curl/7.79.1]] 0x1400002a4c0 <nil> 11 [] false localhost:9988 map[] map[] <nil> map[] 127.0.0.1:62582 /valid_moves <nil> <nil> <nil> 0x1400002a500}

Docker:
1. env GOOS=linux GOARCH=amd64 go build -o ./app/main ./main
2. docker build -t zikunw/gochess:1.0 .
