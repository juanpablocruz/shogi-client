package shogi

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Notation struct {
	// placement of the pieces of the board written rank by rank, with each rank separated by a slash '/'.
	// Ranks are ordered from the white side, beginning with rank 'a' and ending with rank 'i'
	// For each rank squares are specified from file 9 to file 1.
	Board Board

	Turn rune // w or b
	// Pieces in hand encoded in uppercase for black, and lowercase for white
	// and a digit before the letter for multiple pieces of type.
	// The pieces are ordered in: rook > bishop > gold > silver > knight > lance > pawn
	// and all black pieces before white pieces.
	// If neither player has pieces in hand a single '-' character is used
	Hand Hand

	MoveCount int32 // (optional) number of current move in the game.
}

const (
	StartingPosition = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
)

func (n Notation) DecodeBoard(sfen string) (Board, error) {
	return Board{}, nil
}

func (n Notation) EncodeMovement(m Move) string {
	encoded := ""
	encoded += m.Piece.String()

	if n.Board.DisambiguityNeeded(m) {
		encoded += m.Origin.String()
	}

	encoded += m.Type.String()

	encoded += m.Destination.String()

	if m.IsPromoting {
		encoded += "+"
	}
	return encoded
}

func GetNextPart(parts []string) (string, []string) {
	if len(parts) < 1 {
		return "", []string{}
	}
	return parts[0], parts[1:]
}

func (n Notation) ParsePieceWithPromotion(parts []string) (Piece, []string, error) {
	isPromoted := false
	part, parts := GetNextPart(parts)
	if part == "+" {
		isPromoted = true
	} else if p, err := n.ParsePiece(part, false); err == nil {
		return p, parts, nil
	} else {
		return Piece{}, parts, fmt.Errorf("shogi: Couldn't parse piece, expecting (+)PieceCode, received: %s", part)
	}

	part, parts = GetNextPart(parts)
	p, err := n.ParsePiece(part, isPromoted)
	return p, parts, err
}

func (n Notation) ParseSquare(parts []string) (Square, []string, error) {
	part, newParts := GetNextPart(parts)
	i, err := strconv.ParseInt(part, 10, 8)
	if err != nil {
		return Square(-1), parts, err
	}
	part, newParts = GetNextPart(newParts)

	i2, exists := rankAsNum[part]
	if !exists {
		return Square(-1), parts, fmt.Errorf("shogi: Couldn't parse square rank, expecting a-i, received %s", part)
	}

	return NewSquare(File(i), Rank(i2)), newParts, nil
}

func (n Notation) ParseMovement(parts []string) (MoveType, []string, error) {
	part, parts := GetNextPart(parts)

	if part == "-" {
		return SimpleMovement, parts, nil
	}
	if part == "*" {
		return Drop, parts, nil
	}
	if part == "x" {
		return Capture, parts, nil
	}

	return SimpleMovement, parts, fmt.Errorf("shogi: Couldn't parse movement type, expecting -,x,* but received: %s", part)
}

func (n Notation) DecodeMovement(sfen string) (Move, error) {
	var m Move

	parts := strings.Split(sfen, "")

	// Parse Piece
	p, parts, err := n.ParsePieceWithPromotion(parts)
	if err != nil {
		return m, err
	}
	m.Piece = p

	// Parse optional origin
	origin, parts, err := n.ParseSquare(parts)
	if err == nil {
		m.Origin = origin
	}
	// Parse movement
	mType, parts, err := n.ParseMovement(parts)
	if err != nil {
		return m, err
	}
	m.Type = mType

	// Parse Destination
	dest, parts, err := n.ParseSquare(parts)
	if err != nil {
		return m, err
	}
	m.Destination = dest

	// Parse optional promotion
	if len(parts) > 0 {
		part, _ := GetNextPart(parts)
		if part == "+" {
			m.IsPromoting = true
		}
	}

	return m, nil
}

func (n Notation) ParsePiece(sfen string, isPromoted bool) (Piece, error) {
	if !slices.Contains(allPiecesCodes, strings.ToUpper(sfen)) {
		return Piece{}, fmt.Errorf("shogi: Couldn't parse PieceCode, received: %s", sfen)
	}
	return NewPiece(sfen, isPromoted), nil
}
