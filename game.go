package chess

// Game represents the state of the chess game
type Game struct {
	board        *Board
	turn         Color
	selected     *Position
	cursor       *Position
	moveLog      *MoveLog
	boardHistory []*Board
	status       string
}

func (g *Game) Board() *Board {
	return g.board
}

func (g *Game) SetBoard(b *Board) {
	g.board = b
}

func (g *Game) Turn() Color {
	return g.turn
}

func (g *Game) SetTurn(c Color) {
	g.turn = c
}

func (g *Game) Selected() *Position {
	return g.selected
}

func (g *Game) SetSelected(p *Position) {
	g.selected = p
}

func (g *Game) Cursor() *Position {
	return g.cursor
}

func (g *Game) MoveLog() *MoveLog {
	return g.moveLog
}

func (g *Game) BoardHistory() []*Board {
	return g.boardHistory
}

func (g *Game) AddToBoardHistory(history *Board) {
	g.boardHistory = append(g.boardHistory, history)
}

func (g *Game) Status() string {
	return g.status
}

func (g *Game) SetStatus(s string) {
	g.status = s
}

// Position represents a position on the board
type Position struct {
	Row int
	Col int
}

// NewGame creates a new game
func NewGame() *Game {
	return &Game{
		board:        NewBoard(),
		turn:         White,
		cursor:       &Position{Row: 0, Col: 0},
		moveLog:      NewMoveLog(),
		boardHistory: []*Board{NewBoard()},
	}
}
