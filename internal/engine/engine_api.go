package engine

import "github.com/juanpablocruz/shogo/clientr/internal/shogi"

type EngineOption struct {
	Name        string
	Description string
	Type        string
	Default     string
}

type EngineAPI interface {
	SendMessage(string) error
	ReceiveMessage() (string, error)
}

type LocalEngine struct {
	Game     *shogi.Game
	engineCh chan string
	guiCh    chan string
	engine   *Engine
}

type ServerLocalEngine struct {
	engineCh chan string
	guiCh    chan string
}

func (e ServerLocalEngine) SendMessage(s string) error {
	e.guiCh <- s
	return nil
}

func (e ServerLocalEngine) ReceiveMessage() (string, error) {
	msg := <-e.engineCh
	return msg, nil
}

func NewLocalEngine(g *shogi.Game) LocalEngine {
	engineCh := make(chan string, 2)
	guiCh := make(chan string, 2)

	sle := ServerLocalEngine{
		engineCh: engineCh,
		guiCh:    guiCh,
	}
	engineOptions := make(map[string]EngineOption)

	engine := NewEngine(sle, g, engineOptions)

	return LocalEngine{
		Game:     g,
		engineCh: engineCh,
		guiCh:    guiCh,
		engine:   engine,
	}
}

func (e LocalEngine) SendMessage(s string) error {
	e.engineCh <- s
	return nil
}

func (e LocalEngine) ReceiveMessage() (string, error) {
	msg := <-e.guiCh
	return msg, nil
}
