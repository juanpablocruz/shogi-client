package shogi

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
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
	parts := strings.Split(strings.TrimSpace(sfen), " ")
	if len(parts) < 3 {
		return Board{}, fmt.Errorf("shogi: invalid sfen string received, expecting at least 3 parts, received: %s", sfen)
	}

	b := Board{
		CurrentMove: 0,
	}

	bParts := strings.Split(parts[0], "/")
	if len(bParts) < 9 {
		return Board{}, fmt.Errorf("shogi: invalid sfen string received, expecting board with 9 /, received: %s", parts[0])
	}

	bb := make([]string, 81)
	for rank, rankStr := range bParts {
		fileIndx := 0
		for _, c := range rankStr {
			if unicode.IsDigit(c) {
				num := int(c) - 48
				fileIndx += num
				continue
			}
			bb[(rank*numOfSquaresInRow)+fileIndx] = string(c)
			fileIndx++
		}
	}

	b.BitBoard = bb

	turn := parts[1]
	if turn != "b" && turn != "w" {
		return Board{}, fmt.Errorf("shogi: invalid sfen string received, expecting turn to be b or w, but received: %s", turn)
	}

	if turn == "b" {
		b.Turn = Black
	} else {
		b.Turn = White
	}

	var h Hand

	b.Hand = h

	if len(parts) == 4 {
		currMovements, err := strconv.Atoi(parts[3])
		if err != nil {
			return Board{}, err
		}
		b.CurrentMove = currMovements
	}

	return b, nil
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

func (n Notation) DecodeHodgesMove(move string) (Move, error) {
	if len(move) < 4 {
		return Move{}, fmt.Errorf("shogi: Couldn't decode hodges movement, expecting length 4, received: %s", move)
	}

	if len(move) == 5 {
		move = move[1:]
	}

	org := move[:2]
	dest := move[2:4]

	orgParts := strings.Split(org, "")
	destParts := strings.Split(dest, "")

	fromFile, err := strconv.ParseInt(orgParts[0], 10, 8)
	if err != nil {
		return Move{}, err
	}

	destFile, err := strconv.ParseInt(destParts[0], 10, 8)
	if err != nil {
		return Move{}, err
	}
	orgSquare := NewSquare(File(fromFile-1), Rank(rankAsNum[orgParts[1]]-1))
	destSquare := NewSquare(File(destFile-1), Rank(rankAsNum[destParts[1]]-1))

	m := Move{
		Origin:      orgSquare,
		Destination: destSquare,
		Type:        SimpleMovement,
	}

	return m, nil
}
