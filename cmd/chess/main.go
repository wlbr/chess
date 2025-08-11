package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/wlbr/chess"

	"github.com/nsf/termbox-go"
)

const (
	boardHeight   = 9
	messageRow    = boardHeight + 1
	messageHeight = 4
)

var game *chess.Game

func main() {
	chess.Configure()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	for {
		drawMenu()
		mode := waitForModeChoice()
		switch mode {
		case '1':
			playVSPlayer()
		case '2':
			playAgainstAI()
		case 'q':
			return
		}
	}
}

func waitForModeChoice() rune {
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			return ev.Ch
		}
	}
}

func drawMenu() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	msg1 := "Choose game mode:"
	msg2 := "1. Player vs Player"
	msg3 := "2. Player vs AI"
	msg4 := "q. Quit"
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
	termbox.Flush()
}

func playVSPlayer() {
	game = chess.NewGame()
	gameLoop(false)
}

func playAgainstAI() {
	game = chess.NewGame()
	gameLoop(true)
}

func gameLoop(ai bool) {
	for {
		drawEverything()

		if handleEndGameConditions() {
			return
		}

		if ai && game.Turn() == chess.Black {
			from, to := chess.FindBestMove(game, 3)
			piece := game.Board().PieceAt(from.Row, from.Col)
			finalizeMove(from, to, *piece, nil)
			continue
		}

		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Ch {
			case 'e':
				handleExport()
			}
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
	}
}

func drawEverything() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawBoard()
	drawMoveLog()
	drawMessages(game.Status())
	termbox.Flush()
}

func handleEndGameConditions() bool {
	var msg string
	if chess.IsCheckmate(game.Board(), game.MoveLog(), game.Turn()) {
		msg = "Checkmate!"
	} else if chess.IsStalemate(game.Board(), game.MoveLog(), game.Turn()) {
		msg = "Stalemate!"
	} else if chess.IsThreefoldRepetition(game.BoardHistory()) {
		msg = "Threefold repetition!"
	}

	if msg != "" {
		game.SetStatus(msg)
		drawEverything()
		drawMessages(game.Status(), "Export to PGN? (y/n)")
		termbox.Flush()
		waitForPGNChoice()
		return true
	}
	return false
}

func waitForPGNChoice() {
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Ch {
			case 'y':
				msg := doExportPGN()
				drawMessages(game.Status(), msg, "Press Esc to continue.")
				termbox.Flush()
				for {
					ev := termbox.PollEvent()
					if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
						return
					}
				}
			case 'n':
				return
			}
		}
	}
}

func handleExport() {
	msg := doExportPGN()
	drawMessages(msg, "Press any key to continue.")
	termbox.Flush()
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			return
		}
	}
}

func doExportPGN() string {
	if game == nil || len(game.MoveLog().Moves()) == 0 {
		return "No game to export."
	}

	filename := "game.pgn"
	pgn := chess.ExportToPGN(game.MoveLog())
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Sprintf("Error creating file: %s", err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(pgn)
	if err != nil {
		return fmt.Sprintf("Error writing to file: %s", err.Error())
	}

	return fmt.Sprintf("Game exported to %s", filename)
}

func drawMessages(messages ...string) {
	for i, msg := range messages {
		for j, r := range msg {
			termbox.SetCell(j, messageRow+i, r, termbox.ColorWhite, termbox.ColorDefault)
		}
	}
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
		from := *game.Selected()
		to := chess.Position{Row: game.Cursor().Row, Col: game.Cursor().Col}
		game.SetSelected(nil)

		if from.Row == to.Row && from.Col == to.Col {
			return
		}

		if chess.IsValidMove(game.Board(), game.MoveLog(), from, to) {
			piece := game.Board().PieceAt(from.Row, from.Col)
			var promotion *chess.PieceType

			if piece.Type() == chess.Pawn && (to.Row == 0 || to.Row == 7) {
				promoType := promptForPromotion()
				promotion = &promoType
			}
			finalizeMove(from, to, *piece, promotion)
		}
	}
}

func finalizeMove(from, to chess.Position, piece chess.Piece, promotion *chess.PieceType) {
	// Clone board to check for check
	tempBoard := game.Board().Clone()
	pieceAtTo := tempBoard.PieceAt(to.Row, to.Col)
	isCapture := pieceAtTo != nil

	// --- Make all moves on the temporary board ---
	tempBoard.MovePiece(from, to)
	movedPiece := tempBoard.PieceAt(to.Row, to.Col)
	movedPiece.SetHasMoved(true)

	// Handle en passant capture
	if piece.Type() == chess.Pawn && to.Col != from.Col && pieceAtTo == nil {
		isCapture = true
		if piece.Color() == chess.White {
			tempBoard.SetPieceAt(to.Row+1, to.Col, nil)
		} else {
			tempBoard.SetPieceAt(to.Row-1, to.Col, nil)
		}
	}

	// Handle promotion
	if promotion != nil {
		movedPiece.SetType(*promotion)
	}

	// Handle castling
	if piece.Type() == chess.King && chess.Abs(to.Col-from.Col) == 2 {
		var rookFrom, rookTo chess.Position
		if to.Col > from.Col { // Kingside
			rookFrom = chess.Position{Row: to.Row, Col: 7}
			rookTo = chess.Position{Row: to.Row, Col: to.Col - 1}
		} else { // Queenside
			rookFrom = chess.Position{Row: to.Row, Col: 0}
			rookTo = chess.Position{Row: to.Row, Col: to.Col + 1}
		}
		tempBoard.MovePiece(rookFrom, rookTo)
		rook := tempBoard.PieceAt(rookTo.Row, rookTo.Col)
		rook.SetHasMoved(true)
	}

	// --- Validate and commit the move ---
	if !chess.IsCheck(tempBoard, game.MoveLog(), game.Turn()) {
		// Determine check/checkmate status AFTER the move
		var checkStatus string
		opponentColor := chess.White
		if game.Turn() == chess.White {
			opponentColor = chess.Black
		}
		if chess.IsCheckmate(tempBoard, game.MoveLog(), opponentColor) {
			checkStatus = "++"
		} else if chess.IsCheck(tempBoard, game.MoveLog(), opponentColor) {
			checkStatus = "+"
		}

		game.MoveLog().AddMove(from, to, piece, isCapture, checkStatus)
		game.SetBoard(tempBoard)
		game.AddToBoardHistory(game.Board().Clone())
		// Switch turn and update status message
		if game.Turn() == chess.White {
			game.SetTurn(chess.Black)
		} else {
			game.SetTurn(chess.White)
		}
		if checkStatus == "+" || checkStatus == "++" {
			game.SetStatus("Check!")
		} else {
			game.SetStatus("")
		}
	}
}

func promptForPromotion() chess.PieceType {
	drawEverything()
	drawMessages("Promote to: [Q]ueen, [R]ook, [B]ishop, [K]night")
	termbox.Flush()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch strings.ToUpper(string(ev.Ch)) {
			case "Q":
				return chess.Queen
			case "R":
				return chess.Rook
			case "B":
				return chess.Bishop
			case "K":
				return chess.Knight
			default:
				continue
			}
		}
	}
}

func drawBoard() {
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
}

func drawMoveLog() {
	moves := game.MoveLog().Moves()
	xOffset := 20
	yOffsetStart := 1
	_, termHeight := termbox.Size()
	movesPerColumn := termHeight - yOffsetStart - 1
	if movesPerColumn < 1 {
		movesPerColumn = 1
	}

	for i, move := range moves {
		column := i / movesPerColumn
		row := i % movesPerColumn

		finalX := xOffset + (column * 15)
		finalY := yOffsetStart + row

		if finalY >= termHeight {
			continue
		}

		moveNumStr := fmt.Sprintf("%d. ", i+1)
		line := moveNumStr + move.Notation()

		for j, r := range line {
			termbox.SetCell(finalX+j, finalY, r, termbox.ColorWhite, termbox.ColorDefault)
		}
	}
}
