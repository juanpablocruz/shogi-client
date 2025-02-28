package engine_test

import (
	"context"
	"testing"
	"time"

	"github.com/juanpablocruz/shogo/clientr/internal/engine"
	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

func TestGUIEngine_ProcessUSI(t *testing.T) {
	localApi := engine.ServerLocalEngine{
		EngineCh: make(chan string, 2),
		GUICh:    make(chan string, 2),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	e := engine.NewGUIEngine(localApi)

	go func() {
		msg, err := receiveMessage(ctx, localApi.GUICh)
		if err != nil {
			t.Errorf("ProcessCMD() %s: %v", "usi", err)
		}

		if msg != "usi" {
			t.Errorf("ProcessCMD() %s failed want %v received %v", "usi", "usi", msg)
		}
		localApi.EngineCh <- "id name engine"
		localApi.EngineCh <- "option name USI_ShowCurrLine type check"
		localApi.EngineCh <- "usiok"
	}()

	gotErr := e.ProcessCMD(shogi.USI)
	if gotErr != nil {
		t.Errorf("ProcessCMD() failed: %v", gotErr)
		return
	}

	if e.EngineID != "engine" {
		t.Errorf("ProcessCMD() %s: %s, expecting 'engine' got '%s'", "usi", "id not set", e.EngineID)
	}

	opt, exists := e.EngineOptions["USI_ShowCurrLine"]
	if !exists {
		t.Errorf("ProcessCMD() %s: options not set, expecting 'USI_ShowCurrLine received: %v", "usi", e.EngineOptions)
	}
	if opt.Name != "USI_ShowCurrLine" || opt.Type != "check" {
		t.Errorf("ProcessCMD() %s: option not configured properly: %v", "usi", opt)
	}
}

func TestGUIEngine_ProcessCMD(t *testing.T) {
	localApi := engine.ServerLocalEngine{
		EngineCh: make(chan string, 2),
		GUICh:    make(chan string, 2),
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		e engine.EngineAPI
		// Named input parameters for target function.
		cmd     shogi.GUICommand
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "isready",
			e:       localApi,
			cmd:     shogi.IsReady,
			want:    "isready",
			wantErr: false,
		},
		{
			name:    "usinewgame",
			e:       localApi,
			cmd:     shogi.USINewGame,
			want:    "usinewgame",
			wantErr: false,
		},
		{
			name:    "register",
			e:       localApi,
			cmd:     shogi.Register,
			args:    []string{"name", "Stefan", "MK", "code", "4359"},
			want:    "register name Stefan MK code 4359",
			wantErr: false,
		},
		{
			name:    "register fails if name has no args",
			e:       localApi,
			cmd:     shogi.Register,
			args:    []string{"name", "code", "4359"},
			wantErr: true,
		},
		{
			name:    "register fails if code has no args",
			e:       localApi,
			cmd:     shogi.Register,
			args:    []string{"name", "Stefan", "MK", "code"},
			wantErr: true,
		},
		{
			name:    "register fails if no args",
			e:       localApi,
			cmd:     shogi.Register,
			wantErr: true,
		},
		{
			name:    "debug default",
			e:       localApi,
			cmd:     shogi.Debug,
			want:    "debug off",
			wantErr: false,
		},
		{
			name:    "debug on",
			e:       localApi,
			cmd:     shogi.Debug,
			want:    "debug on",
			args:    []string{"on"},
			wantErr: false,
		},
		{
			name:    "debug off",
			e:       localApi,
			cmd:     shogi.Debug,
			want:    "debug off",
			args:    []string{"off"},
			wantErr: false,
		},
		{
			name:    "stop",
			e:       localApi,
			cmd:     shogi.Stop,
			want:    "stop",
			wantErr: false,
		},
		{
			name:    "ponderhit",
			e:       localApi,
			cmd:     shogi.Ponderhit,
			want:    "ponderhit",
			wantErr: false,
		},
		{
			name:    "setoption",
			e:       localApi,
			cmd:     shogi.SetOption,
			want:    "setoption name test value 1",
			args:    []string{"test", "1"},
			wantErr: false,
		},
		{
			name:    "setoption fails if no args",
			e:       localApi,
			cmd:     shogi.SetOption,
			wantErr: true,
		},
		{
			name:    "setoption fails if no value",
			e:       localApi,
			cmd:     shogi.SetOption,
			args:    []string{"test"},
			wantErr: true,
		},
		{
			name:    "position",
			e:       localApi,
			cmd:     shogi.Position,
			want:    "position sfen moves move1 move2 move3",
			args:    []string{"sfen", "move1", "move2", "move3"},
			wantErr: false,
		},
		{
			name:    "position fails if no args",
			e:       localApi,
			cmd:     shogi.Position,
			wantErr: true,
		},
		{
			name:    "position fails if no moves",
			e:       localApi,
			cmd:     shogi.Position,
			args:    []string{"sfen"},
			wantErr: true,
		},
		{
			name: "gameover",
			e:    localApi,
			cmd:  shogi.Gameover,
			want: "gameover win",
			args: []string{"win"},
		},
		{
			name:    "gameover bad outcome fails",
			e:       localApi,
			cmd:     shogi.Gameover,
			args:    []string{"fail"},
			wantErr: true,
		},
		{
			name:    "quit",
			e:       localApi,
			cmd:     shogi.Quit,
			want:    "quit",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			e := engine.NewGUIEngine(tt.e)
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
