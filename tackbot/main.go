package main

import (
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"tackbot/state/board"
	"tackbot/tackgo/uct"
	"time"
)

const (
	Address string = "0.0.0.0:12345"
)

type Server int

type MoveArg struct {
	Algorithm string
	Board     board.Board
	Timeout   uint
}

type MoveResult struct {
	H, V int32
}

func (s *Server) MakeMove(arg *MoveArg, result *MoveResult) (err error) {
	var agent board.IAlgorithm

	switch arg.Algorithm {
	case "uct":
		agent = new(uct.UCT)
	default:
		agent = new(uct.UCT)
	}

	b := board.NewBoard(arg.Board.H, arg.Board.V, arg.Board.S[0], arg.Board.S[1], arg.Board.Now, arg.Board.Turn)
	result.H, result.V, err = agent.MakeMove(b, arg.Timeout, true)

	return
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	server := rpc.NewServer()
	server.Register(new(Server))
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	l, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	log.Println("Server is running at", Address)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Accept new connection", conn.RemoteAddr())
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
