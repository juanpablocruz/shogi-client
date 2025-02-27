package shogi_test

import (
	"reflect"
	"testing"

	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

func TestBoard_LoadSfen(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		sfen    string
		want    shogi.Board
		wantErr bool
	}{
		{
			name: "test",
			sfen: shogi.StartingPosition,
			want: shogi.Board{
				BitBoard:    []string{"l", "n", "s", "g", "k", "g", "s", "n", "l", "", "r", "", "", "", "", "", "b", "", "p", "p", "p", "p", "p", "p", "p", "p", "p", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "P", "P", "P", "P", "P", "P", "P", "P", "P", "", "B", "", "", "", "", "", "R", "", "L", "N", "S", "G", "K", "G", "S", "N", "L"},
				Turn:        shogi.Black,
				CurrentMove: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := shogi.NewBoard()
			gotErr := b.LoadSfen(tt.sfen)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("LoadSfen() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("LoadSfen() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(tt.want.BitBoard, b.BitBoard) {
				t.Errorf("LoadSfen() failed: bitboard differs: want %v got: %v", tt.want.BitBoard, b.BitBoard)
			}
			if tt.want.Turn != b.Turn {
				t.Errorf("LoadSfen() failed: turn differs: want %s got: %s", tt.want.Turn, b.Turn)
			}
			if tt.want.CurrentMove != b.CurrentMove {
				t.Errorf("LoadSfen() failed: current move differs: want %d got: %d", tt.want.CurrentMove, b.CurrentMove)
			}
		})
	}
}

func TestBoard_String(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		sfen string
		want string
	}{
		{
			name: "test",
			sfen: shogi.StartingPosition,
			want: shogi.StartingPosition,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := shogi.NewBoard()
			b.LoadSfen(tt.sfen)
			got := b.String()
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
