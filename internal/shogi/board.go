package shogi

import (
	"fmt"
	"strconv"
	"strings"
)

// L N S G K G S N L
//
//	R           B
//
// # P P P P P P P P P
//
// P P P P P P P P P
//
//	B           R
//
// L N S G K G S N L
type Board struct {
	BitBoard    []string
	Pieces      map[Color][]Piece
	Turn        Color
	CurrentMove int
	Hand        Hand
}

func (b Board) getAllPiecesOfType(p Piece) []Piece {
	allPieces := []Piece{}
	for _, bp := range b.Pieces[p.Color] {
		if bp.Type == p.Type {
			allPieces = append(allPieces, bp)
		}
	}
	return allPieces
}

func (b Board) PieceCanMove(p Piece, m Move) bool {
	return p.CanMove(m.Destination, b)
}

func (b Board) GetPiecesThatCanMove(m Move) []Piece {
	// 1. Get all Pieces of type m.Type
	allPiecesOfType := b.getAllPiecesOfType(m.Piece)
	// 2. Get all of those pieces that can move to m.Destination
	piecesThatCanMove := []Piece{}
	for _, p := range allPiecesOfType {
		if b.PieceCanMove(p, m) {
			piecesThatCanMove = append(piecesThatCanMove, p)
		}
	}
	return piecesThatCanMove
}

// In cases where the moving piece is ambiguous, the starting square is added after the letter for the piece but before the movement
func (b Board) DisambiguityNeeded(m Move) bool {
	piecesThatCanMove := b.GetPiecesThatCanMove(m)
	// 3. If there is more than 1 possibility return true else return false
	return len(piecesThatCanMove) > 1
}

func NewBoard() Board {
	pieces := make(map[Color][]Piece)
	pieces[Black] = make([]Piece, 0)
	pieces[White] = make([]Piece, 0)

	bitBoard := make([]string, 81)

	return Board{
		BitBoard:    bitBoard,
		Pieces:      pieces,
		Hand:        Hand{},
		Turn:        Black,
		CurrentMove: 1,
	}
}

func (b Board) Debug() {
	fmt.Print("  1 2 3 4 5 6 7 8 9")
	for rank, r := range b.BitBoard {
		if rank%numOfSquaresInRow == 0 {
			if rank != 0 {
				fmt.Printf(" %s", string(numAsRank[(rank/numOfSquaresInRow)-1]))
			}
			fmt.Printf("\n%d ", (rank/numOfSquaresInRow)+1)
		}
		if r == "" {
			fmt.Print("  ")
		} else {
			fmt.Printf("%s ", r)
		}
	}

	fmt.Printf(" %s", string(numAsRank[len(numAsRank)-1]))
	fmt.Println("")
}

func (b *Board) LoadSfen(sfen string) error {
	parts := strings.Split(sfen, " ")
	if len(parts) < 3 {
		return fmt.Errorf("shogi: error parsing sfen, expected at least 3 parts, got: %d (%v)", len(parts), parts)
	}
	if len(parts) > 4 {
		return fmt.Errorf("shogi: error parsing sfen, expected at maximum 4 parts, got: %d (%v)", len(parts), parts)
	}

	placements := parts[0]

	rankPieces := strings.Split(placements, "/")
	if len(rankPieces) != 9 {
		return fmt.Errorf("shogi: error parsing sfen, expected 9 ranks, got: %d (%v)", len(rankPieces), rankPieces)
	}
	for rank, r := range rankPieces {
		files := strings.Split(r, "")
		fileIdx := 0
		for _, p := range files {
			ws, err := strconv.Atoi(p)
			if err != nil {
				b.BitBoard[(rank*numOfSquaresInRow)+fileIdx] = p
			} else {
				fileIdx += ws - 1
			}
			fileIdx++
		}
	}

	turn := parts[1]
	switch turn {
	case "b":
		b.Turn = Black
	case "w":
		b.Turn = White
	default:
		return fmt.Errorf("shogi: error parsing sfen, expected turn to be b or w, got: %s", turn)
	}
	// handPieces := parts[2]

	currentMovePart := parts[3]
	if cm, err := strconv.Atoi(currentMovePart); err == nil {
		b.CurrentMove = cm
	}

	return nil
}

func (b Board) Drop(pt PieceType, c Color, s Square) bool {
	return false
}

// Return SFEN encoding of the board:
// This encode contains four fields separated by a space:
// - Piece placement on the board from Black's perspective: lnsgk2nl/1r4gs1/p1pppp1pp/1p4p2/7P1/2P6/PP1PPPP1P/1SG4R1/LN2KGSNL
// - Who has the next move: b|w
// - Pieces in hand: Bb
// - Move count (optional): 1
func (b Board) String() string {
	// each piece is represented with a single letter. Gote's pieces are lowercase while Sente's are uppercase.
	// the set of letters are the ones in western notation
	// each rank is separated by '/'. The listing of ranks is from top (rank1) to bottom (rank9), and
	// the order to pieces is from file 9 to 1 (left to right as viewed on typical shogi diagram with gote as top player).
	// empty squares are indicated with numeral corresponding to the number of adjacent empty squares on the same rank.
	// for example in lnsgk2nl the rank is as follows: |l|n|s|g|k| | |n|l|
	piecePlacement := ""
	wSpree := 0
	for file, p := range b.BitBoard {
		if file%9 == 0 && file != 0 && file != 81 {
			if wSpree > 0 {
				piecePlacement = fmt.Sprintf("%s%d", piecePlacement, wSpree)
				wSpree = 0
			}
			piecePlacement += "/"

		}
		if p == "" {
			wSpree++
			continue
		}
		if wSpree > 0 {
			piecePlacement = fmt.Sprintf("%s%d", piecePlacement, wSpree)
			wSpree = 0
		}
		piecePlacement += p
	}

	// b for Black's turn or w for White's
	turn := b.Turn.String()

	// all the pieces in hand held by each player. Black's pieces in hand use capital letters while White's use lowercase.
	// Bb indicates that Black has one bishop in hand and White also has one bishop in hand.
	// If there are more than one of a type in hand it is preceded by the piece count, e.g. 3P
	handPieces := b.Hand.String()

	// moves are counted under shogi convention, meaning move count is incremented with a single player's action.
	// this field is optional, but most programs will include the move count field when exporting positions.
	currentMove := b.CurrentMove

	return fmt.Sprintf("%s %s %s %d", piecePlacement, turn, handPieces, currentMove)
}

func (b Board) GetPieceAtSquare(p Piece, sq Square) (Piece, error) {
	if b.BitBoard[sq] == p.String() {
		p.Square = sq
		return p, nil
	}
	return Piece{}, fmt.Errorf("shogi: no piece %s found at square (%s,%s)", p.String(), sq.File().String(), sq.Rank().String())
}

func (b *Board) ProcessMove(m *Move) error {
	var p Piece
	var err error
	if m.Origin == NewSquare(0, 0) {
		candidates := b.GetPiecesThatCanMove(*m)
		if len(candidates) != 1 {
			b.Debug()
			return fmt.Errorf("shogi: no valid candidates %s to move to (%s,%s)", m.Piece.String(), m.Destination.File().String(), m.Destination.Rank().String())
		}
		p = candidates[0]
	} else {
		p, err = b.GetPieceAtSquare(m.Piece, m.Origin)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s (%s,%s)\n", p.String(), p.Square.File().String(), p.Square.Rank().String())
	return nil
}
