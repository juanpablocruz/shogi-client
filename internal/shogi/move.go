package shogi

// - Simple movement
// x capture (opponent's piece)
// * drop (your own piece)
type MoveType int8

const (
	SimpleMovement MoveType = iota
	Capture
	Drop
)

func (mt MoveType) String() string {
	switch mt {
	case SimpleMovement:
		return "-"
	case Capture:
		return "x"
	case Drop:
		return "*"
	}
	return ""
}

//	1       2         3          4             5
//
// piece  (origin)  movement  destination   (promotion)
type Move struct {
	Type        MoveType
	Piece       Piece
	IsPromoting bool
	Destination Square
	Origin      Square
}
