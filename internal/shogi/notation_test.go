package shogi_test

import (
	"testing"

	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

func TestNotation_EncodeMovement(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		m    shogi.Move
		want string
	}{
		{
			name: "SimpleMove",
			m: shogi.Move{
				Type:        shogi.SimpleMovement,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(1), shogi.Rank(2)),
				Origin:      shogi.NewSquare(shogi.File(1), shogi.Rank(1)),
			},
			want: "P-2c",
		},
		{
			name: "SimpleMoveWhite",
			m: shogi.Move{
				Type:        shogi.SimpleMovement,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.White},
				Destination: shogi.NewSquare(shogi.File(1), shogi.Rank(2)),
				Origin:      shogi.NewSquare(shogi.File(1), shogi.Rank(1)),
			},
			want: "p-2c",
		},
		{
			name: "SimpleMoveFull",
			m: shogi.Move{
				Type:        shogi.SimpleMovement,
				Piece:       shogi.Piece{Type: shogi.Gold, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(6), shogi.Rank(7)),
				Origin:      shogi.NewSquare(shogi.File(6), shogi.Rank(6)),
				IsPromoting: true,
			},
			want: "G7g-7h+",
		},
		{
			name: "SimpleMovePromoting",
			m: shogi.Move{
				Type:        shogi.SimpleMovement,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(1), shogi.Rank(2)),
				Origin:      shogi.NewSquare(shogi.File(1), shogi.Rank(1)),
				IsPromoting: true,
			},
			want: "P-2c+",
		},
		{
			name: "Capture",
			m: shogi.Move{
				Type:        shogi.Capture,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(1), shogi.Rank(2)),
				Origin:      shogi.NewSquare(shogi.File(1), shogi.Rank(1)),
			},
			want: "Px2c",
		},
		{
			name: "Drop",
			m: shogi.Move{
				Type:        shogi.Drop,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(1), shogi.Rank(2)),
				Origin:      shogi.NewSquare(shogi.File(1), shogi.Rank(1)),
			},
			want: "P*2c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var n shogi.Notation
			n.Board = shogi.NewBoard()

			allPieces := []shogi.Piece{
				{Type: shogi.Pawn, Color: shogi.Black, Square: shogi.NewSquare(shogi.File(2), shogi.Rank(2))},
				{Type: shogi.Gold, Color: shogi.Black, Square: shogi.NewSquare(shogi.File(6), shogi.Rank(6))},
				{Type: shogi.Gold, Color: shogi.Black, Square: shogi.NewSquare(shogi.File(7), shogi.Rank(7))},
				{Type: shogi.Gold, Color: shogi.Black, Square: shogi.NewSquare(shogi.File(6), shogi.Rank(8))},
			}

			n.Board.Pieces[shogi.Black] = allPieces
			for _, p := range allPieces {
				n.Board.BitBoard[p.Square] = p.String()
			}

			got := n.EncodeMovement(tt.m)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				n.Board.Debug()
				t.Errorf("EncodeMovement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotation_DecodeMovement(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		sfen    string
		want    shogi.Move
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "SimpleMove",
			want: shogi.Move{
				Type:        shogi.SimpleMovement,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(2), shogi.Rank(3)),
			},
			sfen: "P-2c",
		},
		{
			name: "SimpleMoveWhite",
			want: shogi.Move{
				Type:        shogi.SimpleMovement,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.White},
				Destination: shogi.NewSquare(shogi.File(2), shogi.Rank(3)),
			},
			sfen: "p-2c",
		},
		{
			name: "SimpleMoveFull",
			want: shogi.Move{
				Type:        shogi.SimpleMovement,
				Piece:       shogi.Piece{Type: shogi.Gold, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(7), shogi.Rank(8)),
				Origin:      shogi.NewSquare(shogi.File(7), shogi.Rank(7)),
				IsPromoting: true,
			},
			sfen: "G7g-7h+",
		},
		{
			name: "SimpleMovePromoting",
			want: shogi.Move{
				Type:        shogi.SimpleMovement,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(2), shogi.Rank(3)),
				IsPromoting: true,
			},
			sfen: "P-2c+",
		},
		{
			name: "Capture",
			want: shogi.Move{
				Type:        shogi.Capture,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(2), shogi.Rank(3)),
			},
			sfen: "Px2c",
		},
		{
			name: "Drop",
			want: shogi.Move{
				Type:        shogi.Drop,
				Piece:       shogi.Piece{Type: shogi.Pawn, Color: shogi.Black},
				Destination: shogi.NewSquare(shogi.File(2), shogi.Rank(3)),
			},
			sfen: "P*2c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var n shogi.Notation

			n.Board = shogi.NewBoard()

			n.Board.Pieces[shogi.Black] = []shogi.Piece{
				{Type: shogi.Pawn, Color: shogi.Black, Square: shogi.NewSquare(shogi.File(2), shogi.Rank(2))},
				{Type: shogi.Gold, Color: shogi.Black, Square: shogi.NewSquare(shogi.File(7), shogi.Rank(7))},
				{Type: shogi.Gold, Color: shogi.Black, Square: shogi.NewSquare(shogi.File(9), shogi.Rank(7))},
				{Type: shogi.Gold, Color: shogi.Black, Square: shogi.NewSquare(shogi.File(8), shogi.Rank(6))},
			}

			got, gotErr := n.DecodeMovement(tt.sfen)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DecodeMovement() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DecodeMovement() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.

			if got.Type != tt.want.Type || got.Piece.Type != tt.want.Piece.Type ||
				got.IsPromoting != tt.want.IsPromoting || got.Destination.String() != tt.want.Destination.String() {
				t.Errorf("DecodeMovement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotation_DecodeHodgesMove(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		move    string
		want    shogi.Move
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Hodges decode",
			move: "7g7f",
			want: shogi.Move{
				Destination: shogi.NewSquare(shogi.File(6), shogi.Rank(5)),
				Origin:      shogi.NewSquare(shogi.File(6), shogi.Rank(6)),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var n shogi.Notation
			got, gotErr := n.DecodeHodgesMove(tt.move)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DecodeHodgesMove() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DecodeHodgesMove() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if tt.want.Destination != got.Destination || tt.want.Origin != got.Origin {
				t.Errorf("DecodeHodgesMove() = %v, want %v", got, tt.want)
			}
		})
	}
}
