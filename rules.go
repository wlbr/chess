package chess

func IsValidMove(board *Board, moveLog *MoveLog, from, to Position) bool {
	piece := board[from.Row][from.Col]
	if piece == nil {
		return false
	}

	switch piece.Type() {
	case Pawn:
		return isValidPawnMove(board, moveLog, from, to)
	case Rook:
		return isValidRookMove(board, from, to)
	case Knight:
		return isValidKnightMove(board, from, to)
	case Bishop:
		return isValidBishopMove(board, from, to)
	case Queen:
		return isValidQueenMove(board, from, to)
	case King:
		return isValidKingMove(board, moveLog, from, to)
	}

	return false
}

func isValidPawnMove(board *Board, moveLog *MoveLog, from, to Position) bool {
	piece := board.PieceAt(from.Row, from.Col)
	dx := to.Col - from.Col
	dy := to.Row - from.Row

	if piece.Color() == White {
		// Standard 1-step forward
		if dx == 0 && dy == -1 && board.PieceAt(to.Row, to.Col) == nil {
			return true
		}
		// Standard 2-step forward
		if from.Row == 6 && dx == 0 && dy == -2 && board.PieceAt(to.Row, to.Col) == nil && board.PieceAt(from.Row-1, from.Col) == nil {
			return true
		}
		// Standard capture
		if Abs(dx) == 1 && dy == -1 && board.PieceAt(to.Row, to.Col) != nil && board.PieceAt(to.Row, to.Col).Color() == Black {
			return true
		}
		// En passant capture
		if from.Row == 3 && Abs(dx) == 1 && dy == -1 && board.PieceAt(to.Row, to.Col) == nil {
			lastMove := moveLog.LastMove()
			if lastMove != nil && lastMove.piece.Type() == Pawn && lastMove.to.Row == 3 && lastMove.to.Col == to.Col && Abs(lastMove.from.Row-lastMove.to.Row) == 2 {
				return true
			}
		}
	} else { // Black
		// Standard 1-step forward
		if dx == 0 && dy == 1 && board.PieceAt(to.Row, to.Col) == nil {
			return true
		}
		// Standard 2-step forward
		if from.Row == 1 && dx == 0 && dy == 2 && board.PieceAt(to.Row, to.Col) == nil && board.PieceAt(from.Row+1, from.Col) == nil {
			return true
		}
		// Standard capture
		if Abs(dx) == 1 && dy == 1 && board.PieceAt(to.Row, to.Col) != nil && board.PieceAt(to.Row, to.Col).Color() == White {
			return true
		}
		// En passant capture
		if from.Row == 4 && Abs(dx) == 1 && dy == 1 && board.PieceAt(to.Row, to.Col) == nil {
			lastMove := moveLog.LastMove()
			if lastMove != nil && lastMove.piece.Type() == Pawn && lastMove.to.Row == 4 && lastMove.to.Col == to.Col && Abs(lastMove.from.Row-lastMove.to.Row) == 2 {
				return true
			}
		}
	}
	return false
}

func isValidRookMove(board *Board, from, to Position) bool {
	dx := to.Col - from.Col
	dy := to.Row - from.Row

	if dx != 0 && dy != 0 {
		return false
	}

	// Check for pieces in the way
	if dx > 0 {
		for i := 1; i < dx; i++ {
			if board.PieceAt(from.Row, from.Col+i) != nil {
				return false
			}
		}
	} else if dx < 0 {
		for i := -1; i > dx; i-- {
			if board.PieceAt(from.Row, from.Col+i) != nil {
				return false
			}
		}
	} else if dy > 0 {
		for i := 1; i < dy; i++ {
			if board.PieceAt(from.Row+i, from.Col) != nil {
				return false
			}
		}
	} else if dy < 0 {
		for i := -1; i > dy; i-- {
			if board.PieceAt(from.Row+i, from.Col) != nil {
				return false
			}
		}
	}

	return board.PieceAt(to.Row, to.Col) == nil || board.PieceAt(to.Row, to.Col).Color() != board.PieceAt(from.Row, from.Col).Color()
}

func isValidKnightMove(board *Board, from, to Position) bool {
	dx := Abs(to.Col - from.Col)
	dy := Abs(to.Row - from.Row)

	if (dx == 1 && dy == 2) || (dx == 2 && dy == 1) {
		return board.PieceAt(to.Row, to.Col) == nil || board.PieceAt(to.Row, to.Col).Color() != board.PieceAt(from.Row, from.Col).Color()
	}

	return false
}

func isValidBishopMove(board *Board, from, to Position) bool {
	dx := to.Col - from.Col
	dy := to.Row - from.Row

	if Abs(dx) != Abs(dy) {
		return false
	}

	// Check for pieces in the way
	xDir := 1
	if dx < 0 {
		xDir = -1
	}
	yDir := 1
	if dy < 0 {
		yDir = -1
	}

	for i := 1; i < Abs(dx); i++ {
		if board.PieceAt(from.Row+i*yDir, from.Col+i*xDir) != nil {
			return false
		}
	}

	return board.PieceAt(to.Row, to.Col) == nil || board.PieceAt(to.Row, to.Col).Color() != board.PieceAt(from.Row, from.Col).Color()
}

func isValidQueenMove(board *Board, from, to Position) bool {
	return isValidRookMove(board, from, to) || isValidBishopMove(board, from, to)
}

func isValidKingMove(board *Board, moveLog *MoveLog, from, to Position) bool {
	if isValidCastling(board, moveLog, from, to) {
		return true
	}
	dx := Abs(to.Col - from.Col)
	dy := Abs(to.Row - from.Row)

	if dx > 1 || dy > 1 {
		return false
	}

	return board.PieceAt(to.Row, to.Col) == nil || board.PieceAt(to.Row, to.Col).Color() != board.PieceAt(from.Row, from.Col).Color()
}

func IsCheck(board *Board, moveLog *MoveLog, color Color) bool {
	kingPos := findKing(board, color)
	if kingPos == nil {
		return false
	}

	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			piece := board.PieceAt(r, c)
			if piece != nil && piece.Color() != color {
				if IsValidMove(board, moveLog, Position{Row: r, Col: c}, *kingPos) {
					return true
				}
			}
		}
	}

	return false
}

func findKing(board *Board, color Color) *Position {
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			piece := board.PieceAt(r, c)
			if piece != nil && piece.Type() == King && piece.Color() == color {
				return &Position{Row: r, Col: c}
			}
		}
	}
	return nil
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func IsCheckmate(board *Board, moveLog *MoveLog, color Color) bool {
	if !IsCheck(board, moveLog, color) {
		return false
	}

	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			piece := board.PieceAt(r, c)
			if piece != nil && piece.Color() == color {
				for r2 := 0; r2 < 8; r2++ {
					for c2 := 0; c2 < 8; c2++ {
						to := Position{Row: r2, Col: c2}
						if IsValidMove(board, moveLog, Position{Row: r, Col: c}, to) {
							// Make the move on a temporary board
							tempBoard := board.Clone()
							tempBoard.MovePiece(Position{Row: r, Col: c}, to)

							// Check if the king is in check
							if !IsCheck(tempBoard, moveLog, color) {
								return false
							}
						}
					}
				}
			}
		}
	}

	return true
}

func IsStalemate(board *Board, moveLog *MoveLog, color Color) bool {
	if IsCheck(board, moveLog, color) {
		return false
	}

	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			piece := board.PieceAt(r, c)
			if piece != nil && piece.Color() == color {
				for r2 := 0; r2 < 8; r2++ {
					for c2 := 0; c2 < 8; c2++ {
						to := Position{Row: r2, Col: c2}
						if IsValidMove(board, moveLog, Position{Row: r, Col: c}, to) {
							// Make the move on a temporary board
							tempBoard := board.Clone()
							tempBoard.MovePiece(Position{Row: r, Col: c}, to)

							// Check if the king is in check
							if !IsCheck(tempBoard, moveLog, color) {
								return false
							}
						}
					}
				}
			}
		}
	}

	return true
}

func IsThreefoldRepetition(history []*Board) bool {
	if len(history) < 9 {
		return false
	}

	lastBoard := history[len(history)-1]
	count := 0
	for _, board := range history {
		if boardsEqual(lastBoard, board) {
			count++
		}
	}

	return count >= 3
}

func boardsEqual(b1, b2 *Board) bool {
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			p1 := b1[r][c]
			p2 := b2[r][c]
			if p1 == nil && p2 == nil {
				continue
			}
			if p1 == nil || p2 == nil {
				return false
			}
			if p1.Type() != p2.Type() || p1.Color() != p2.Color() {
				return false
			}
		}
	}
	return true
}

func isValidCastling(board *Board, moveLog *MoveLog, from, to Position) bool {
	piece := board.PieceAt(from.Row, from.Col)
	if piece == nil || piece.Type() != King || piece.HasMoved() {
		return false
	}

	dx := to.Col - from.Col
	dy := to.Row - from.Row

	if Abs(dx) != 2 || dy != 0 {
		return false
	}

	// Check if king is in check
	if IsCheck(board, moveLog, piece.Color()) {
		return false
	}

	// Check for pieces between king and rook
	if dx > 0 { // Kingside castling
		rook := board.PieceAt(from.Row, 7)
		if rook == nil || rook.Type() != Rook || rook.HasMoved() {
			return false
		}
		if board.PieceAt(from.Row, from.Col+1) != nil || board.PieceAt(from.Row, from.Col+2) != nil {
			return false
		}
		// Check if squares king moves through are under attack
		if isSquareAttacked(board, moveLog, from.Row, from.Col+1, piece.Color()) || isSquareAttacked(board, moveLog, from.Row, from.Col+2, piece.Color()) {
			return false
		}
	} else { // Queenside castling
		rook := board.PieceAt(from.Row, 0)
		if rook == nil || rook.Type() != Rook || rook.HasMoved() {
			return false
		}
		if board.PieceAt(from.Row, from.Col-1) != nil || board.PieceAt(from.Row, from.Col-2) != nil || board.PieceAt(from.Row, from.Col-3) != nil {
			return false
		}
		// Check if squares king moves through are under attack
		if isSquareAttacked(board, moveLog, from.Row, from.Col-1, piece.Color()) || isSquareAttacked(board, moveLog, from.Row, from.Col-2, piece.Color()) {
			return false
		}
	}

	return true
}

func isSquareAttacked(board *Board, moveLog *MoveLog, row, col int, color Color) bool {
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			piece := board.PieceAt(r, c)
			if piece != nil && piece.Color() != color {
				if IsValidMove(board, moveLog, Position{Row: r, Col: c}, Position{Row: row, Col: col}) {
					return true
				}
			}
		}
	}
	return false
}
