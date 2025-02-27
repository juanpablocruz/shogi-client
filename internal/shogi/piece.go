package shogi

import (
	"strings"
)

type PieceType int8

const (
	NoPiece PieceType = iota
	King
	Rook
	Bishop
	Gold
	Silver
	Knight
	Lance
	Pawn
)

type Piece struct {
	Type       PieceType
	Color      Color
	IsPromoted bool
	Square     Square
}

func (p Piece) String() string {
	typeString := p.Type.String()
	promotedFlag := ""
	if p.IsPromoted {
		promotedFlag = "+"
	}

	if p.Color == Black {
		return promotedFlag + typeString
	}
	return promotedFlag + strings.ToLower(typeString)
}

func (p Piece) Render() (rune, Color) {
	pStr := strings.ToUpper(p.String())
	switch pStr {
	case "K":
		return '王', p.Color
	case "R":
		return '飛', p.Color
	case "+R":
		return '龍', p.Color
	case "B":
		return '角', p.Color
	case "+B":
		return '馬', p.Color
	case "G":
		return '金', p.Color
	case "S":
		return '銀', p.Color
	case "+S":
		return '全', p.Color
	case "N":
		return '桂', p.Color
	case "+N":
		return '圭', p.Color
	case "L":
		return '香', p.Color
	case "+L":
		return '杏', p.Color
	case "P":
		return '歩', p.Color
	case "+P":
		return 'と', p.Color
	}
	return 'a', p.Color
}

func (pt PieceType) String() string {
	switch pt {
	case King:
		return "K"
	case Rook:
		return "R"
	case Bishop:
		return "B"
	case Gold:
		return "G"
	case Silver:
		return "S"
	case Knight:
		return "N"
	case Lance:
		return "L"
	case Pawn:
		return "P"
	}
	return ""
}

func PieceTypeFromCode(code string) PieceType {
	switch code {
	case "K":
		return King
	case "R":
		return Rook
	case "B":
		return Bishop
	case "G":
		return Gold
	case "S":
		return Silver
	case "N":
		return Knight
	case "L":
		return Lance
	case "P":
		return Pawn
	}
	return NoPiece
}

var allPiecesCodes []string = []string{"K", "R", "B", "G", "S", "N", "L", "P"}

func NewPiece(code string, isPromoted bool) Piece {
	var color Color
	if strings.ToLower(code) == code {
		color = White
	} else {
		color = Black
	}
	pieceType := PieceTypeFromCode(strings.ToUpper(code))
	return Piece{
		Type:       pieceType,
		Color:      color,
		IsPromoted: isPromoted,
	}
}

func (p Piece) CanMove(s Square, board Board) bool {
	switch p.Type {
	case King:
		return PieceKingCanMove(p.Square, s, board)
	case Rook:
		if p.IsPromoted {
			return PiecePromotedRookCanMove(p.Square, s, board)
		}
		return PieceRookCanMove(p.Square, s, board)
	case Bishop:
		if p.IsPromoted {
			return PiecePromotedBishopCanMove(p.Square, s, board)
		}
		return PieceBishopCanMove(p.Square, s, board)
	case Gold:
		return PieceGoldCanMove(p.Square, s, board)
	case Silver:
		if p.IsPromoted {
			return PieceGoldCanMove(p.Square, s, board)
		}
		return PieceSilverCanMove(p.Square, s, board)
	case Lance:
		if p.IsPromoted {
			return PieceGoldCanMove(p.Square, s, board)
		}
		return PieceLanceCanMove(p.Square, s, board)
	case Knight:
		if p.IsPromoted {
			return PieceGoldCanMove(p.Square, s, board)
		}
		return PieceKnightCanMove(p.Square, s, board)
	case Pawn:
		if p.IsPromoted {
			return PieceGoldCanMove(p.Square, s, board)
		}
		return PiecePawnCanMove(p.Square, s, board)
	}
	return false
}

// Rook moves any number of squares in an orthogonal direction
//
//	| | ||| | |
//	| | ||| | |
//	|-|-|R|-|-|
//	| | ||| | |
//	| | ||| | |
func PieceRookCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()
	if df == 0 && dr == 0 {
		return false // no move
	}
	if df != 0 && dr != 0 {
		return false // not a straight line
	}

	// Determine step direction
	stepFile := sign(int(df))
	stepRank := sign(int(dr))

	return isPathClear(o, s, board, stepFile, stepRank)
}

// Promoted Rook moves as a Rook and as a King
//
//	| | ||| | |
//	| |*|||*| |
//	|-|-|龍|-|-|
//	| |*|||*| |
//	| | ||| | |
func PiecePromotedRookCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()

	if abs(int(df)) <= 1 && abs(int(dr)) <= 1 {
		return true
	}
	return PieceRookCanMove(o, s, board)
}

// King moves one square in any direction, orthogonal or diagonal
// | | | | | |
// | |*|*|*| |
// | |*|K|*| |
// | |*|*|*| |
// | | | | | |
func PieceKingCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()
	if df == 0 && dr == 0 {
		return false // no move
	}
	if abs(int(df)) <= 1 && abs(int(dr)) <= 1 {
		return true
	}

	return false
}

// Bishop moves any number of squares in diagonal direction
// |\| | | |/|
// | |\| |/| |
// | | |B| | |
// | |/| |\| |
// |/| | | |\|
func PieceBishopCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()
	if abs(int(df)) != abs(int(dr)) || df == 0 {
		return false // not a diagonal move
	}
	stepFile := sign(int(df))
	stepRank := sign(int(dr))
	return isPathClear(o, s, board, stepFile, stepRank)
}

// Promoted Bishop moves as a Bishop and as a King
// |\| | | |/|
// | |\|*|/| |
// | |*|馬|*| |
// | |/|*|\| |
// |/| | | |\|
func PiecePromotedBishopCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()
	// One-square move in any direction
	if abs(int(df)) <= 1 && abs(int(dr)) <= 1 && (df != 0 || dr != 0) {
		return true
	}
	return PieceBishopCanMove(o, s, board)
}

// Gold general moves one square orthogonally, or one square diagonally forward
// Promoted Silver general moves like a Gold General
// Promoted Knight general moves like a Gold General
// Promoted Lance general moves like a Gold General
// Promoted Pawn general moves like a Gold General
// | | | | | |
// | |*|*|*| |
// | |*|G|*| |
// | | |*| | |
// | | | | | |
func PieceGoldCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()
	fwd := forwardDirection(board, o)

	var allowed [][2]int
	if fwd == -1 {
		// For pieces moving upward:
		allowed = [][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {0, 1}}
	} else {
		// For pieces moving downward:
		allowed = [][2]int{{-1, 1}, {0, 1}, {1, 1}, {-1, 0}, {1, 0}, {0, -1}}
	}

	for _, delta := range allowed {
		if int(df) == delta[0] && int(dr) == delta[1] {
			return true
		}
	}

	return false
}

// Silver general moves one square diagonally, or one square straight forward
// | | | | | |
// | |*|*|*| |
// | | |S| | |
// | |*| |*| |
// | | | | | |
func PieceSilverCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()
	fwd := forwardDirection(board, o)

	// Allowed moves for silver: forward and all four diagonals.
	allowed := [][2]int{
		{0, fwd},   // straight forward
		{-1, fwd},  // diagonal forward-left
		{1, fwd},   // diagonal forward-right
		{-1, -fwd}, // diagonal backward-left
		{1, -fwd},  // diagonal backward-right
	}
	for _, delta := range allowed {
		if int(df) == delta[0] && int(dr) == delta[1] {
			return true
		}
	}
	return false
}

// Knight jumps at an angle intermediate to orthogonal and diagonal, amounting to one
// square straight forward plus one square diagonally forward, in a single move.
// | |*| |*| |
// | | | | | |
// | | |K| | |
// | | | | | |
// | | | | | |
func PieceKnightCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()
	fwd := forwardDirection(board, o)
	// For knight, allowed moves are:
	//   for fwd == -1 (moving up): (-1, -2) and (1, -2)
	//   for fwd ==  1 (moving down): (-1, 2) and (1, 2)
	allowed := [][2]int{
		{-1, 2 * fwd},
		{1, 2 * fwd},
	}
	for _, delta := range allowed {
		if int(df) == delta[0] && int(dr) == delta[1] {
			return true
		}
	}
	return false
}

// Lance moves just like the rook except it cannot move backwards or to the sides.
// | | ||| | |
// | | ||| | |
// | | |L| | |
// | | | | | |
// | | | | | |
func PieceLanceCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()
	fwd := forwardDirection(board, o)
	// Lance moves only in the forward direction and along the same file.
	if df != 0 {
		return false
	}
	// For a valid move, the rank change must be in the forward direction.
	if (fwd == -1 && dr >= 0) || (fwd == 1 && dr <= 0) {
		return false
	}
	// Check that every square in front is clear.
	return isPathClear(o, s, board, 0, fwd)
}

// Pawn moves one square straight forward
// | | | | | |
// | | |*| | |
// | | |P| | |
// | | | | | |
// | | | | | |
func PiecePawnCanMove(o Square, s Square, board Board) bool {
	df := s.File() - o.File()
	dr := s.Rank() - o.Rank()
	fwd := forwardDirection(board, o)
	return df == 0 && int(dr) == fwd
}
