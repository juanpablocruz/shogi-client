package shogi_test

import (
	"testing"

	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

func TestHand_String(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		h    shogi.Hand
		want string
	}{
		{
			name: "one of each black and white", h: shogi.Hand{
				BlackPieces: map[shogi.PieceType]int{
					shogi.Pawn:   1,
					shogi.Lance:  1,
					shogi.Knight: 1,
					shogi.Silver: 1,
					shogi.Gold:   1,
					shogi.Bishop: 1,
					shogi.Rook:   1,
					shogi.King:   1,
				},
				WhitePieces: map[shogi.PieceType]int{
					shogi.Pawn:   1,
					shogi.Lance:  1,
					shogi.Knight: 1,
					shogi.Silver: 1,
					shogi.Gold:   1,
					shogi.Bishop: 1,
					shogi.Rook:   1,
					shogi.King:   1,
				},
			},
			want: "KRBGSNLPkrbgsnlp",
		},
		{
			name: "several with black and white", h: shogi.Hand{
				BlackPieces: map[shogi.PieceType]int{
					shogi.Pawn: 6,
					shogi.Rook: 1,
				},
				WhitePieces: map[shogi.PieceType]int{
					shogi.Pawn: 2,
				},
			},
			want: "R6P2p",
		},
		{
			name: "empty hand",
			h:    shogi.Hand{},
			want: "-",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			got := tt.h.String()
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
