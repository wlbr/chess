package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wlbr/chess"
)

var game *chess.Game

func main() {
	chess.Configure()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "uci":
			fmt.Println("id name MyChessEngine")
			fmt.Println("id author Your Name")
			fmt.Println("uciok")
		case "isready":
			fmt.Println("readyok")
		case "ucinewgame":
			game = chess.NewGame()
		case "position":
			if len(parts) > 1 {
				switch parts[1] {
				case "startpos":
					game = chess.NewGame()
					if len(parts) > 2 && parts[2] == "moves" {
						for _, moveStr := range parts[3:] {
							move := chess.AlgebraicToMove(game.Board(), moveStr)
							game.Board().MovePiece(move.From(), move.To())
						}
					}
				case "fen":
					// TODO: implement fen
				}
			}
		case "go":
			from, to := chess.FindBestMove(game.Board(), game.Turn(), 5) // Using a depth of 5 for the UCI
			piece := game.Board().PieceAt(from.Row, from.Col)
			move := chess.NewMove(from, to, *piece)
			fmt.Printf("bestmove %s\n", move.Notation())
		case "stop":
			// TODO: implement stop command
		case "quit":
			return
		}
	}
}
