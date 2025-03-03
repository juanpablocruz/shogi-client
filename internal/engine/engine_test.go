package engine_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/juanpablocruz/shogo/clientr/internal/engine"
	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

func TestEngine_Process_ID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	localApi := engine.ServerLocalEngine{
		EngineCh: make(chan string, 2),
		GUICh:    make(chan string, 2),
	}
	e := engine.NewEngine(
		"id",
		localApi,
		shogi.NewGame("sente", "gote"),
		make(map[string]engine.EngineOption),
	)
	err := e.ProcessCMD(shogi.Id)
	if err != nil {
		t.Errorf("Process ID failed: %v", err)
		return
	}

	msg, err := receiveMessage(ctx, localApi.GUICh)
	if err != nil {
		t.Fatalf("Process ID failed: %v", err)
	}
	if msg != "id name id" {
		t.Errorf("Process ID failed: id command expecting id name engine.ID")
	}

	msg, err = receiveMessage(ctx, localApi.GUICh)
	if err != nil {
		t.Fatalf("Process ID failed: %v", err)
	}
	if msg != "usiok" {
		t.Errorf("Process ID failed: id command expecting usiok after id identification")
	}
}

func TestEngine_Process_ID_With_Opts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	localApi := engine.ServerLocalEngine{
		EngineCh: make(chan string, 5),
		GUICh:    make(chan string, 5),
	}
	opts := make(map[string]engine.EngineOption)
	opts["opt1"] = engine.EngineOption{
		Name: "opt1",
		Type: "opt1type",
	}
	opts["optwd"] = engine.EngineOption{
		Name:    "optwd",
		Type:    "opt2type",
		Default: "32",
	}

	e := engine.NewEngine(
		"id",
		localApi,
		shogi.NewGame("sente", "gote"),
		opts,
	)
	err := e.ProcessCMD(shogi.Id)
	if err != nil {
		t.Errorf("Process ID failed: %v", err)
		return
	}

	msg, err := receiveMessage(ctx, localApi.GUICh)
	if err != nil {
		t.Fatalf("Process ID failed: %v", err)
	}
	if msg != "id name id" {
		t.Errorf("Process ID failed: id command expecting id name engine.ID")
	}

	msg, err = receiveMessage(ctx, localApi.GUICh)
	if err != nil {
		t.Fatalf("Process ID failed: %v", err)
	}
	if msg != "option name opt1 type opt1type\n" {
		t.Errorf("Process ID failed: id command expecting option sent, received: %s", msg)
	}

	msg, err = receiveMessage(ctx, localApi.GUICh)
	if err != nil {
		t.Fatalf("Process ID failed: %v", err)
	}
	if msg != "option name optwd type opt2type default 32\n" {
		t.Errorf("Process ID failed: id command expecting option sent, received: %s", msg)
	}

	msg, err = receiveMessage(ctx, localApi.GUICh)
	if err != nil {
		t.Fatalf("Process ID failed: %v", err)
	}
	if msg != "usiok" {
		t.Errorf("Process ID failed: id command expecting usiok after id identification")
	}
}

func receiveMessage(ctx context.Context, msg chan string) (string, error) {
	select {
	case m := <-msg:
		return m, nil
	case <-ctx.Done():
		return "", fmt.Errorf("timeout waiting for message")
	}
}

func TestEngine_ProcessCMD(t *testing.T) {
	localApi := engine.ServerLocalEngine{
		EngineCh: make(chan string, 2),
		GUICh:    make(chan string, 2),
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		api     engine.EngineAPI
		id      string
		game    *shogi.Game
		options map[string]engine.EngineOption
		// Named input parameters for target function.
		cmd     shogi.EngineCommand
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "readyok",
			api:     localApi,
			id:      "id",
			game:    shogi.NewGame("sente", "gote"),
			options: make(map[string]engine.EngineOption),
			cmd:     shogi.ReadyOk,
			args:    []string{},
			want:    "readyok",
			wantErr: false,
		},
		{
			name:    "checkmate fails if no arguments",
			api:     localApi,
			id:      "id",
			game:    shogi.NewGame("sente", "gote"),
			options: make(map[string]engine.EngineOption),
			cmd:     shogi.Checkmate,
			wantErr: true,
		},
		{
			name:    "bestmove sends moves",
			api:     localApi,
			id:      "id",
			game:    shogi.NewGame("sente", "gote"),
			options: make(map[string]engine.EngineOption),
			cmd:     shogi.Checkmate,
			args:    []string{"move1", "move2", "move3"},
			want:    "checkmate move1 move2 move3",
		},
		{
			name:    "checkmate sends outcome",
			api:     localApi,
			id:      "id",
			game:    shogi.NewGame("sente", "gote"),
			options: make(map[string]engine.EngineOption),
			cmd:     shogi.Checkmate,
			args:    []string{"nomate"},
			want:    "checkmate nomate",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			e := engine.NewEngine(tt.id, tt.api, tt.game, tt.options)
			gotErr := e.ProcessCMD(tt.cmd, tt.args...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ProcessCMD() failed: %v", gotErr)
				}
				return
			}

			msg, err := receiveMessage(ctx, localApi.GUICh)
			if err != nil {
				t.Fatalf("ProcessCMD() %s: %v", tt.name, err)
			}
			if tt.want != msg {
				t.Fatalf("ProcessCMD() %s failed want %v received %v", tt.name, tt.want, msg)
			}
			if tt.wantErr {
				t.Fatal("ProcessCMD() succeeded unexpectedly")
			}
		})
	}
}

func TestEngine_ProcessPosition_full(t *testing.T) {
	localApi := engine.ServerLocalEngine{
		EngineCh: make(chan string, 2),
		GUICh:    make(chan string, 2),
	}

	expectedBb := []string{
		"l", "n", "s", "g", "k", "g", "s", "n", "l", "", "r", "", "", "", "", "", "b", "", "", "p", "p", "p", "p", "p", "p", "p", "p", "p", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "P", "P", "P", "P", "P", "P", "P", "P", "P", "", "B", "", "", "", "", "", "R", "", "L", "N", "S", "G", "K", "G", "S", "N", "L",
	}
	e := engine.NewEngine("id", localApi, shogi.NewGame("sente", "gote"), make(map[string]engine.EngineOption))
	gotErr := e.ProcessPosition([]string{
		"lnsgkgsnl/1r5b1/1pppppppp/p8/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL",
		"b",
		"-",
		"1",
		"moves",
		"p-4d",
		"P-6f",
	})
	if gotErr != nil {
		t.Errorf("ProcessCMD() failed: %v", gotErr)
	}

	if !reflect.DeepEqual(expectedBb, e.Game.Board().BitBoard) {

		e.Game.Board().Debug()
		t.Errorf("ProcessCMD() failed: bitboard differs: want %v got: %v", expectedBb, e.Game.Board().BitBoard)
	}

	if e.Game.Board().Turn != shogi.Black {
		t.Errorf("ProcessCMD() failed: turn differs: want %s got: %s", shogi.Black, e.Game.Board().Turn)
	}

	if len(e.Game.Moves()) != 2 {
		t.Errorf("ProcessCMD() failed: moves not loaded, want %d got: %d", 2, len(e.Game.Moves()))
	}
}

func TestEngine_ProcessPosition_moves(t *testing.T) {
	localApi := engine.ServerLocalEngine{
		EngineCh: make(chan string, 2),
		GUICh:    make(chan string, 2),
	}

	g := shogi.NewGame("sente", "gote")
	g.Move(shogi.Move{})

	e := engine.NewEngine("id", localApi, g, make(map[string]engine.EngineOption))
	gotErr := e.ProcessPosition([]string{
		"moves",
		"p-4d",
		"P-6f",
	})
	if gotErr != nil {
		t.Errorf("ProcessCMD() failed: %v", gotErr)
	}

	if len(e.Game.Moves()) != 3 {
		t.Errorf("ProcessCMD() failed: moves not loaded, want %d got: %d", 3, len(e.Game.Moves()))
	}
}

func TestEngine_ProcessPosition_no_moves(t *testing.T) {
	localApi := engine.ServerLocalEngine{
		EngineCh: make(chan string, 2),
		GUICh:    make(chan string, 2),
	}

	e := engine.NewEngine("id", localApi, shogi.NewGame("sente", "gote"), make(map[string]engine.EngineOption))
	gotErr := e.ProcessPosition([]string{"move1", "move2"})
	if gotErr == nil {
		t.Errorf("ProcessCMD() failed: %v", gotErr)
	}

	if len(e.Game.Moves()) != 0 {
		t.Errorf("ProcessCMD() failed: moves not loaded, want %d got: %d", 0, len(e.Game.Moves()))
	}
}
