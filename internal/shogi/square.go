package shogi

import "fmt"

type (
	Square int8
	File   int8
	Rank   int8
)

const (
	numOfSquaresInBoard = 81
	numOfSquaresInRow   = 9
)

type SquareInfo struct {
	File  File
	Rank  Rank
	Piece Piece
}

func (sq Square) File() File {
	return File(int(sq) % numOfSquaresInRow)
}

func (sq Square) Rank() Rank {
	return Rank(int(sq) / numOfSquaresInRow)
}

func (sq Square) String() string {
	return sq.File().String() + sq.Rank().String()
}

func NewSquare(f File, r Rank) Square {
	return Square(int8(r)*numOfSquaresInRow + int8(f))
}

func (f File) String() string {
	return fmt.Sprintf("%d", f+1)
}

func (r Rank) String() string {
	if r < 1 || int(r) > len(numAsRank) {
		return ""
	}
	return string(numAsRank[r])
}

func (r Rank) Rune() rune {
	if r < 0 || int(r) > len(numAsRank) {
		return '-'
	}
	return numAsRank[r]
}

var rankAsNum = map[string]int{
	"a": 1, "b": 2, "c": 3,
	"d": 4, "e": 5, "f": 6,
	"g": 7, "h": 8, "i": 9,
}
var numAsRank = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i'}
