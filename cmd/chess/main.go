package main

import (
	"fmt"
	"os"

	"github.com/wlbr/chess"

	"github.com/nsf/termbox-go"
)

var game *chess.Game

func main() {
	chess.Configure()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	drawMenu()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Ch {
			case '1':
				playVSPlayer()
				return
			case '2':
				playAgainstAI()
				return
			case '3':
				exportPGN()
				return
			case 'q':
				return
			}
		}
	}
}

func drawMenu() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	msg1 := "Choose game mode:"
	msg2 := "1. Player vs Player"
	msg3 := "2. Player vs AI"
	msg4 := "3. Export PGN"
	msg5 := "q. Quit"
	for i, r := range msg1 {
		termbox.SetCell(i, 0, r, termbox.ColorWhite, termbox.ColorDefault)
	}
	for i, r := range msg2 {
		termbox.SetCell(i, 2, r, termbox.ColorWhite, termbox.ColorDefault)
	}
	for i, r := range msg3 {
		termbox.SetCell(i, 3, r, termbox.ColorWhite, termbox.ColorDefault)
	}
	for i, r := range msg4 {
		termbox.SetCell(i, 4, r, termbox.ColorWhite, termbox.ColorDefault)
	}
	for i, r := range msg5 {
		termbox.SetCell(i, 5, r, termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.Flush()
}

func playVSPlayer() {
	game = chess.NewGame()
	drawBoard()

	for {
		if chess.IsCheckmate(game.Board(), game.Turn()) {
			drawMessage("Checkmate!")
			break
		}
		if chess.IsStalemate(game.Board(), game.Turn()) {
			drawMessage("Stalemate!")
			break
		}
		if chess.IsThreefoldRepetition(game.BoardHistory()) {
			drawMessage("Threefold repetition!")
			break
		}

		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeyArrowUp:
				moveCursor(0, -1)
			case termbox.KeyArrowDown:
				moveCursor(0, 1)
			case termbox.KeyArrowLeft:
				moveCursor(-1, 0)
			case termbox.KeyArrowRight:
				moveCursor(1, 0)
			case termbox.KeyEnter:
				selectPiece()
			}
		}
		drawBoard()
	}

	// Wait for user to press Esc to exit
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			break
		}
	}
}

func playAgainstAI() {
	game = chess.NewGame()
	drawBoard()

	for {
		if chess.IsCheckmate(game.Board(), game.Turn()) {
			drawMessage("Checkmate!")
			break
		}
		if chess.IsStalemate(game.Board(), game.Turn()) {
			drawMessage("Stalemate!")
			break
		}
		if chess.IsThreefoldRepetition(game.BoardHistory()) {
			drawMessage("Threefold repetition!")
			break
		}

		if game.Turn() == chess.White {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					return
				case termbox.KeyArrowUp:
					moveCursor(0, -1)
				case termbox.KeyArrowDown:
					moveCursor(0, 1)
				case termbox.KeyArrowLeft:
					moveCursor(-1, 0)
				case termbox.KeyArrowRight:
					moveCursor(1, 0)
				case termbox.KeyEnter:
					selectPiece()
				}
			}
		} else {
			from, to := chess.FindBestMove(game.Board(), chess.Black, 3)
			piece := game.Board().PieceAt(from.Row, from.Col)
			game.MoveLog().AddMove(from, to, *piece)
			game.Board().MovePiece(from, to)
			game.SetTurn(chess.White)
		}

		drawBoard()
	}

	// Wait for user to press Esc to exit
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			break
		}
	}
}

func exportPGN() {
	if game == nil || len(game.MoveLog().Moves()) == 0 {
		drawMessage("No game to export.")
		return
	}

	// TODO: prompt for filename
	filename := "game.pgn"
	pgn := chess.ExportToPGN(game.MoveLog())
	file, err := os.Create(filename)
	if err != nil {
		drawMessage(fmt.Sprintf("Error creating file: %s", err.Error()))
		return
	}
	defer file.Close()

	_, err = file.WriteString(pgn)
	if err != nil {
		drawMessage(fmt.Sprintf("Error writing to file: %s", err.Error()))
		return
	}

	drawMessage(fmt.Sprintf("Game exported to %s", filename))
}

func drawMessage(msg string) {
	msg2 := "Press Esc to exit."
	for i, r := range msg {
		termbox.SetCell(i, 9, r, termbox.ColorWhite, termbox.ColorDefault)
	}
	for i, r := range msg2 {
		termbox.SetCell(i, 10, r, termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.Flush()
}

func moveCursor(dx, dy int) {
	game.Cursor().Col += dx
	game.Cursor().Row += dy

	if game.Cursor().Col < 0 {
		game.Cursor().Col = 0
	}
	if game.Cursor().Col > 7 {
		game.Cursor().Col = 7
	}
	if game.Cursor().Row < 0 {
		game.Cursor().Row = 0
	}
	if game.Cursor().Row > 7 {
		game.Cursor().Row = 7
	}
}

func selectPiece() {
	if game.Selected() == nil {
		piece := game.Board().PieceAt(game.Cursor().Row, game.Cursor().Col)
		if piece != nil && piece.Color() == game.Turn() {
			game.SetSelected(&chess.Position{Row: game.Cursor().Row, Col: game.Cursor().Col})
		}
	} else {
		to := chess.Position{Row: game.Cursor().Row, Col: game.Cursor().Col}
		if game.Selected().Row == to.Row && game.Selected().Col == to.Col {
			game.SetSelected(nil)
			return
		}

		if chess.IsValidMove(game.Board(), *game.Selected(), to) {
			// Make the move on a temporary board
			tempBoard := game.Board().Clone()
			piece := tempBoard.PieceAt(game.Selected().Row, game.Selected().Col)
			tempBoard.MovePiece(*game.Selected(), to)
			piece.SetHasMoved(true)

			// Handle castling
			if piece.Type() == chess.King && chess.Abs(to.Col-game.Selected().Col) == 2 {
				if to.Col > game.Selected().Col { // Kingside
					tempBoard.MovePiece(chess.Position{Row: to.Row, Col: 7}, chess.Position{Row: to.Row, Col: to.Col - 1})
				} else { // Queenside
					tempBoard.MovePiece(chess.Position{Row: to.Row, Col: 0}, chess.Position{Row: to.Row, Col: to.Col + 1})
				}
			}

			// Check if the king is in check
			if !chess.IsCheck(tempBoard, game.Turn()) {
				game.MoveLog().AddMove(*game.Selected(), to, *piece)
				game.SetBoard(tempBoard)
				game.AddToBoardHistory(game.Board().Clone())
				if game.Turn() == chess.White {
					game.SetTurn(chess.Black)
				} else {
					game.SetTurn(chess.White)
				}
			}
		}
		game.SetSelected(nil)
	}
}

func drawBoard() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw column labels
	for i := 0; i < 8; i++ {
		termbox.SetCell(i*2+2, 0, rune('a'+i), termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw row labels
	for i := 0; i < 8; i++ {
		termbox.SetCell(0, i+1, rune('8'-i), termbox.ColorWhite, termbox.ColorDefault)
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			bg := termbox.ColorDarkGray
			if (i+j)%2 != 0 {
				bg = termbox.ColorBlack
			}

			if game.Selected() != nil && game.Selected().Row == i && game.Selected().Col == j {
				bg = termbox.ColorGreen
			} else if game.Cursor().Row == i && game.Cursor().Col == j {
				bg = termbox.ColorYellow
			}

			piece := game.Board().PieceAt(i, j)
			var r rune
			if piece != nil {
				r = piece.Rune()
			}

			fg := termbox.ColorDefault
			if piece != nil {
				if piece.Color() == chess.White {
					fg = termbox.ColorWhite
				} else {
					fg = termbox.ColorBlue
				}
			}

			termbox.SetCell(j*2+2, i+1, r, fg, bg)
			termbox.SetCell(j*2+3, i+1, ' ', fg, bg)
		}
	}

	// Draw move log
	for i := 0; i < len(game.MoveLog().Moves()); i += 2 {
		row := (i / 2) + 1
		moveNum := (i / 2) + 1
		moveNumStr := fmt.Sprintf("%d. ", moveNum)
		for j, r := range moveNumStr {
			termbox.SetCell(20+j, row, r, termbox.ColorWhite, termbox.ColorDefault)
		}

		whiteMove := game.MoveLog().Moves()[i]
		for j, r := range whiteMove.Notation() {
			termbox.SetCell(20+j+len(moveNumStr), row, r, termbox.ColorWhite, termbox.ColorDefault)
		}

		if i+1 < len(game.MoveLog().Moves()) {
			blackMove := game.MoveLog().Moves()[i+1]
			for j, r := range blackMove.Notation() {
				termbox.SetCell(35+j, row, r, termbox.ColorWhite, termbox.ColorDefault)
			}
		}
	}

	termbox.Flush()
}
