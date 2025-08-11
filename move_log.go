package chess

import (
	"fmt"
	"strings"
)

// Move represents a single move in the game
type Move struct {
	from        Position
	to          Position
	piece       Piece
	notation    string
	isCapture   bool
	checkStatus string
}

func (m *Move) Notation() string {
	return m.notation
}

// MoveLog stores the history of moves in the game
type MoveLog struct {
	moves []*Move
}

func (ml *MoveLog) Moves() []*Move {
	return ml.moves
}

// NewMoveLog creates a new move log
func NewMoveLog() *MoveLog {
	return &MoveLog{
		moves: []*Move{},
	}
}

// AddMove adds a move to the log
func (ml *MoveLog) AddMove(from, to Position, piece Piece, isCapture bool, checkStatus string) {
	move := &Move{
		from:        from,
		to:          to,
		piece:       piece,
		isCapture:   isCapture,
		checkStatus: checkStatus,
	}
	move.notation = moveToAlgebraic(move)
	ml.moves = append(ml.moves, move)
}

func (ml *MoveLog) LastMove() *Move {
	if len(ml.moves) == 0 {
		return nil
	}
	return ml.moves[len(ml.moves)-1]
}

func moveToAlgebraic(move *Move) string {
	// Handle castling first
	if move.piece.Type() == King {
		dx := move.to.Col - move.from.Col
		if dx == 2 {
			return "O-O" + move.checkStatus
		}
		if dx == -2 {
			return "O-O-O" + move.checkStatus
		}
	}

	var sb strings.Builder

	switch move.piece.Type() {
	case King:
		sb.WriteString("K")
	case Queen:
		sb.WriteString("Q")
	case Rook:
		sb.WriteString("R")
	case Bishop:
		sb.WriteString("B")
	case Knight:
		sb.WriteString("N")
	}

	sb.WriteString(fmt.Sprintf("%c%d", 'a'+move.from.Col, 8-move.from.Row))
	if move.isCapture {
		sb.WriteString("x")
	} else {
		sb.WriteString("-")
	}
	sb.WriteString(fmt.Sprintf("%c%d", 'a'+move.to.Col, 8-move.to.Row))
	sb.WriteString(move.checkStatus)

	return sb.String()
}

func NewMove(from, to Position, piece Piece, isCapture bool, checkStatus string) *Move {
	move := &Move{
		from:        from,
		to:          to,
		piece:       piece,
		isCapture:   isCapture,
		checkStatus: checkStatus,
	}
	move.notation = moveToAlgebraic(move)
	return move
}

func (m *Move) From() Position {
	return m.from
}

func (m *Move) To() Position {
	return m.to
}

func AlgebraicToMove(board *Board, moveStr string) *Move {
	fromCol := int(moveStr[0] - 'a')
	fromRow := 8 - int(moveStr[1]-'0')
	isCapture := moveStr[2] == 'x'
	var toCol, toRow int
	if isCapture {
		toCol = int(moveStr[3] - 'a')
		toRow = 8 - int(moveStr[4]-'0')
	} else {
		toCol = int(moveStr[3] - 'a')
		toRow = 8 - int(moveStr[4]-'0')
	}

	from := Position{Row: fromRow, Col: fromCol}
	to := Position{Row: toRow, Col: toCol}
	piece := board.PieceAt(fromRow, fromCol)

	// TODO: handle promotion

	return &Move{
		from:      from,
		to:        to,
		piece:     *piece,
		isCapture: isCapture,
	}
}
