package chess

import (
	"math"
)

func FindBestMove(board *Board, color Color, depth int) (Position, Position) {
	bestScore := math.Inf(-1)
	var bestMoveFrom Position
	var bestMoveTo Position

	for r1 := 0; r1 < 8; r1++ {
		for c1 := 0; c1 < 8; c1++ {
			piece := board.PieceAt(r1, c1)
			if piece != nil && piece.Color() == color {
				from := Position{Row: r1, Col: c1}
				for r2 := 0; r2 < 8; r2++ {
					for c2 := 0; c2 < 8; c2++ {
						to := Position{Row: r2, Col: c2}
						if IsValidMove(board, from, to) {
							tempBoard := board.Clone()
							tempBoard.MovePiece(from, to)

							score := minimax(tempBoard, depth-1, false, color)
							if score > bestScore {
								bestScore = score
								bestMoveFrom = from
								bestMoveTo = to
							}
						}
					}
				}
			}
		}
	}
	return bestMoveFrom, bestMoveTo
}

func minimax(board *Board, depth int, isMaximizingPlayer bool, color Color) float64 {
	if depth == 0 {
		return evaluate(board, color)
	}

	if isMaximizingPlayer {
		bestScore := math.Inf(-1)
		for r1 := 0; r1 < 8; r1++ {
			for c1 := 0; c1 < 8; c1++ {
				piece := board.PieceAt(r1, c1)
				if piece != nil && piece.Color() == color {
					from := Position{Row: r1, Col: c1}
					for r2 := 0; r2 < 8; r2++ {
						for c2 := 0; c2 < 8; c2++ {
							to := Position{Row: r2, Col: c2}
							if IsValidMove(board, from, to) {
								tempBoard := board.Clone()
								tempBoard.MovePiece(from, to)
								score := minimax(tempBoard, depth-1, false, color)
								bestScore = math.Max(bestScore, score)
							}
						}
					}
				}
			}
		}
		return bestScore
	} else {
		bestScore := math.Inf(1)
		for r1 := 0; r1 < 8; r1++ {
			for c1 := 0; c1 < 8; c1++ {
				piece := board.PieceAt(r1, c1)
				if piece != nil && piece.Color() != color {
					from := Position{Row: r1, Col: c1}
					for r2 := 0; r2 < 8; r2++ {
						for c2 := 0; c2 < 8; c2++ {
							to := Position{Row: r2, Col: c2}
							if IsValidMove(board, from, to) {
								tempBoard := board.Clone()
								tempBoard.MovePiece(from, to)
								score := minimax(tempBoard, depth-1, true, color)
								bestScore = math.Min(bestScore, score)
							}
						}
					}
				}
			}
		}
		return bestScore
	}
}

func evaluate(board *Board, color Color) float64 {
	var score float64
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			piece := board.PieceAt(r, c)
			if piece != nil {
				value := getPieceValue(piece.Type())
				if piece.Color() == color {
					score += value
				} else {
					score -= value
				}
			}
		}
	}
	return score
}

func getPieceValue(pieceType PieceType) float64 {
	switch pieceType {
	case Pawn:
		return 1
	case Knight:
		return 3
	case Bishop:
		return 3
	case Rook:
		return 5
	case Queen:
		return 9
	case King:
		return 900
	default:
		return 0
	}
}
