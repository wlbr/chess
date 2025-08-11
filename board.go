package chess

// Piece represents a chess piece
type Piece struct {
	pieceType PieceType
	color     Color
	hasMoved  bool
}

func (p *Piece) Type() PieceType {
	return p.pieceType
}

func (p *Piece) Color() Color {
	return p.color
}

func (p *Piece) HasMoved() bool {
	return p.hasMoved
}

func (p *Piece) SetHasMoved(hasMoved bool) {
	p.hasMoved = hasMoved
}

func (p *Piece) SetType(pt PieceType) {
	p.pieceType = pt
}

func (p *Piece) Rune() rune {
	switch p.pieceType {
	case King:
		if p.color == White {
			return '♔'
		} else {
			return '♚'
		}
	case Queen:
		if p.color == White {
			return '♕'
		} else {
			return '♛'
		}
	case Rook:
		if p.color == White {
			return '♖'
		} else {
			return '♜'
		}
	case Bishop:
		if p.color == White {
			return '♗'
		} else {
			return '♝'
		}
	case Knight:
		if p.color == White {
			return '♘'
		} else {
			return '♞'
		}
	case Pawn:
		if p.color == White {
			return '♙'
		} else {
			return '♟'
		}
	}
	return ' '
}

// PieceType represents the type of a chess piece
type PieceType int

const (
	King PieceType = iota
	Queen
	Rook
	Bishop
	Knight
	Pawn
)

// Color represents the color of a chess piece
type Color int

const (
	White Color = iota
	Black
)

// Board represents the chess board
type Board [8][8]*Piece

func (b *Board) PieceAt(row, col int) *Piece {
	return b[row][col]
}

func (b *Board) MovePiece(from, to Position) {
	b[to.Row][to.Col] = b[from.Row][from.Col]
	b[from.Row][from.Col] = nil
}

func (b *Board) SetPieceAt(row, col int, p *Piece) {
	b[row][col] = p
}

func NewPiece(pieceType PieceType, color Color) *Piece {
	return &Piece{pieceType: pieceType, color: color, hasMoved: false}
}

func (b *Board) Clone() *Board {
	newBoard := &Board{}
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			if b[r][c] != nil {
				newBoard[r][c] = &Piece{pieceType: b[r][c].pieceType, color: b[r][c].color, hasMoved: b[r][c].hasMoved}
			}
		}
	}
	return newBoard
}

// NewBoard creates a new board with pieces in their starting positions
func NewBoard() *Board {
	board := &Board{}

	// Place black pieces
	board[0][0] = &Piece{pieceType: Rook, color: Black, hasMoved: false}
	board[0][1] = &Piece{pieceType: Knight, color: Black, hasMoved: false}
	board[0][2] = &Piece{pieceType: Bishop, color: Black, hasMoved: false}
	board[0][3] = &Piece{pieceType: Queen, color: Black, hasMoved: false}
	board[0][4] = &Piece{pieceType: King, color: Black, hasMoved: false}
	board[0][5] = &Piece{pieceType: Bishop, color: Black, hasMoved: false}
	board[0][6] = &Piece{pieceType: Knight, color: Black, hasMoved: false}
	board[0][7] = &Piece{pieceType: Rook, color: Black, hasMoved: false}
	for i := 0; i < 8; i++ {
		board[1][i] = &Piece{pieceType: Pawn, color: Black, hasMoved: false}
	}

	// Place white pieces
	board[7][0] = &Piece{pieceType: Rook, color: White, hasMoved: false}
	board[7][1] = &Piece{pieceType: Knight, color: White, hasMoved: false}
	board[7][2] = &Piece{pieceType: Bishop, color: White, hasMoved: false}
	board[7][3] = &Piece{pieceType: Queen, color: White, hasMoved: false}
	board[7][4] = &Piece{pieceType: King, color: White, hasMoved: false}
	board[7][5] = &Piece{pieceType: Bishop, color: White, hasMoved: false}
	board[7][6] = &Piece{pieceType: Knight, color: White, hasMoved: false}
	board[7][7] = &Piece{pieceType: Rook, color: White, hasMoved: false}
	for i := 0; i < 8; i++ {
		board[6][i] = &Piece{pieceType: Pawn, color: White, hasMoved: false}
	}

	return board
}

func pieceToRune(p *Piece) rune {
	switch p.Type() {
	case King:
		if p.Color() == White {
			return '♔'
		} else {
			return '♚'
		}
	case Queen:
		if p.Color() == White {
			return '♕'
		} else {
			return '♛'
		}
	case Rook:
		if p.Color() == White {
			return '♖'
		} else {
			return '♜'
		}
	case Bishop:
		if p.Color() == White {
			return '♗'
		} else {
			return '♝'
		}
	case Knight:
		if p.Color() == White {
			return '♘'
		} else {
			return '♞'
		}
	case Pawn:
		if p.Color() == White {
			return '♙'
		} else {
			return '♟'
		}
	}
	return ' '
}
