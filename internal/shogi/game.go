package shogi

import "github.com/juanpablocruz/shogo/clientr/internal/agent"

type Outcome string

const (
	NoOutcome Outcome = "*"
	WhiteWon  Outcome = "1-0"
	BlackWon  Outcome = "0-1"
	Draw      Outcome = "1/2-1/2"
)

func (o Outcome) String() string {
	return string(o)
}

type Game struct {
	sentePlayer string
	gotePlayer  string
	notation    Notation
	moves       []*Move
	board       *Board
	ai          agent.Agent
}

func NewGame(sentePlayer, GotePlayer string, options ...func(*Game)) *Game {
	board := NewBoard()
	board.LoadSfen(StartingPosition)

	game := &Game{
		sentePlayer: sentePlayer,
		gotePlayer:  GotePlayer,
		notation: Notation{
			Board:     board,
			Turn:      rune(board.Turn.String()[0]),
			Hand:      Hand{},
			MoveCount: int32(board.CurrentMove),
		},
		board: &board,
		moves: []*Move{},
	}

	for _, f := range options {
		f(game)
	}

	return game
}

func (g *Game) SetAIClient(c agent.Agent) {
	g.ai = c
}

func (g *Game) GetAIClient() agent.Agent {
	return g.ai
}

func (g Game) GotePlayer() string {
	return g.gotePlayer
}

func (g Game) SentePlayer() string {
	return g.sentePlayer
}

func (g Game) Notation() Notation {
	return g.notation
}

func (g Game) Moves() []*Move {
	return g.moves
}

func (g Game) Board() *Board {
	return g.board
}

func (g *Game) SetBoard(b *Board) {
	g.board = b
}

func (g Game) MoveStr(cmd string) error {
	m, err := g.notation.DecodeMovement(cmd)
	if err != nil {
		return err
	}
	return g.Move(m)
}

func (g *Game) Move(m Move) error {
	g.moves = append(g.moves, &m)
	return g.board.ProcessMove(&m)
}
