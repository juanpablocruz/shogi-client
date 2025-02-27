package shogi_test

import (
	"testing"

	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

func TestPiece_Render(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		code       string
		isPromoted bool
		want       rune
		want2      shogi.Color
	}{
		{
			name:       "Returns white pawn kanji",
			code:       "p",
			isPromoted: false,
			want:       '歩',
			want2:      shogi.White,
		},
		{
			name:       "Returns black pawn kanji",
			code:       "P",
			isPromoted: false,
			want:       '歩',
			want2:      shogi.Black,
		},
		{
			name:       "Returns promoted black pawn kanji",
			code:       "P",
			isPromoted: true,
			want:       'と',
			want2:      shogi.Black,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := shogi.NewPiece(tt.code, tt.isPromoted)
			got, got2 := p.Render()
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("Render() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("Render() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
