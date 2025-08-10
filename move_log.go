package chess

import (
	"fmt"
	"strings"
)

// Move represents a single move in the game
type Move struct {
	from     Position
	to       Position
	piece    Piece
	notation string
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
func (ml *MoveLog) AddMove(from, to Position, piece Piece) {
	move := &Move{
		from:  from,
		to:    to,
		piece: piece,
	}
	move.notation = moveToAlgebraic(move)
	ml.moves = append(ml.moves, move)
}

func moveToAlgebraic(move *Move) string {
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
	sb.WriteString("-")
	sb.WriteString(fmt.Sprintf("%c%d", 'a'+move.to.Col, 8-move.to.Row))

	return sb.String()
}

func NewMove(from, to Position, piece Piece) *Move {
	move := &Move{
		from:  from,
		to:    to,
		piece: piece,
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
	toCol := int(moveStr[2] - 'a')
	toRow := 8 - int(moveStr[3]-'0')

	from := Position{Row: fromRow, Col: fromCol}
	to := Position{Row: toRow, Col: toCol}
	piece := board.PieceAt(fromRow, fromCol)

	// TODO: handle promotion

	return &Move{
		from:  from,
		to:    to,
		piece: *piece,
	}
}
