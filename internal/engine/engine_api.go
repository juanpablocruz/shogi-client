package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

type EngineOption struct {
	Name        string
	Description string
	Type        string
	Default     string
}

type EngineAPI interface {
	SendMessage(string) error
	ReceiveMessage(context.Context) (string, error)
}

type LocalEngine struct {
	Game     *shogi.Game
	engineCh chan string
	guiCh    chan string
	engine   *Engine
}

type ServerLocalEngine struct {
	EngineCh chan string
	GUICh    chan string
}

func (e ServerLocalEngine) SendMessage(s string) error {
	e.GUICh <- s
	return nil
}

func (e ServerLocalEngine) ReceiveMessage(ctx context.Context) (string, error) {
	select {
	case m := <-e.EngineCh:
		return strings.ReplaceAll(m, "\n", ""), nil
	case <-ctx.Done():
		return "", fmt.Errorf("timeout waiting for message")
	}
}

func NewLocalEngine(g *shogi.Game) LocalEngine {
	engineCh := make(chan string, 2)
	guiCh := make(chan string, 2)

	sle := ServerLocalEngine{
		EngineCh: engineCh,
		GUICh:    guiCh,
	}
	engineOptions := make(map[string]EngineOption)
	id := uuid.New().String()
	engine := NewEngine(id, sle, g, engineOptions)

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

func (e LocalEngine) ReceiveMessage(ctx context.Context) (string, error) {
	select {
	case m := <-e.guiCh:
		return strings.ReplaceAll(m, "\n", ""), nil
	case <-ctx.Done():
		return "", fmt.Errorf("timeout waiting for message")
	}
}
