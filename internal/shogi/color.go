package shogi

type Color int8

const (
	Black Color = iota // Sente
	White              // Gote
)

func (c Color) String() string {
	if c == Black {
		return "b"
	}
	return "w"
}
