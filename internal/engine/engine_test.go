package engine_test

import (
	"context"
	"fmt"
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
