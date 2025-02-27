package shogi

import "unicode"

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	return 1
}

// isPathClear checks that every square between origin (o) and target
// in the direction specified by (stepFile, stepRank) is empty.
// Note: s itself is not checked so that capturing logic can be handled separately.
func isPathClear(o, s Square, board Board, stepFile, stepRank int) bool {
	curFile := o.File() + File(stepFile)
	curRank := o.Rank() + Rank(stepRank)
	targetFile := s.File()
	targetRank := s.Rank()

	// loop until we reach the target Square
	for curFile != targetFile || curRank != targetRank {
		index := curRank*numOfSquaresInRow + Rank(curFile)
		if board.BitBoard[index] != "" {
			return false
		}
		curFile += File(stepFile)
		curRank += Rank(stepRank)
	}
	return true
}

// forwardDirection returns the "forward" direction for a piece based on its letter color
// Here we assume that uppercase pieces move "up" and lowercase pieces move "down"
func forwardDirection(board Board, o Square) int {
	pieceLetter := board.BitBoard[o]
	if pieceLetter != "" && unicode.IsUpper(rune(pieceLetter[0])) {
		return -1 // e.g. Sente: forward means upward (decreasing rank)
	}
	return 1 // e.g. Gote: forward means downward (increasing rank)
}
