package shogi

import "strconv"

type Hand struct {
	BlackPieces map[PieceType]int
	WhitePieces map[PieceType]int
}

var pieceOrder = []PieceType{King, Rook, Bishop, Gold, Silver, Knight, Lance, Pawn}

func (h Hand) String() string {
	handStr := ""
	if len(h.BlackPieces) == 0 && len(h.WhitePieces) == 0 {
		return "-"
	}
	for _, pType := range pieceOrder {
		if val, exists := h.BlackPieces[pType]; exists && val > 0 {
			p := Piece{Type: pType, Color: Black}

			if val > 1 {
				handStr += strconv.Itoa(val)
			}
			handStr += p.String()
		}
	}

	for _, pType := range pieceOrder {
		if val, exists := h.WhitePieces[pType]; exists && val > 0 {
			p := Piece{Type: pType, Color: White}

			if val > 1 {
				handStr += strconv.Itoa(val)
			}
			handStr += p.String()
		}
	}

	return handStr
}
