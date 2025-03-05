package gui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/juanpablocruz/shogo/clientr/internal/input"
	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
	"github.com/juanpablocruz/shogo/clientr/internal/theme"
)

const (
	leftMargin          = 4
	topMargin           = 4
	numOfSquaresInBoard = 81
	numOfSquaresInRow   = 9
)

func (gui GUI) drawRune(x, y int, style tcell.Style, r rune) {
	(*gui.Screen).SetContent(x, y, r, nil, style)
}

func squareColor(sq shogi.Square) shogi.Color {
	if sq&1 == 0 {
		return shogi.Black
	}
	return shogi.White
}

func stylePiece(p shogi.Piece, sqBg tcell.Color, t theme.Theme) tcell.Style {
	pieceStyle := tcell.StyleDefault.Background(sqBg)

	if p.Color == shogi.White {
		return pieceStyle.Foreground(t.Gote).Bold(true)
	}
	return pieceStyle.Foreground(t.Sente).Bold(true)
}

func squareBg(sq shogi.Square, t theme.Theme) tcell.Color {
	squareColor := squareColor(sq)
	if squareColor == shogi.Black {
		return t.SquareDark
	}
	return t.SquareLight
}

func (gui GUI) drawLabel(x, y int, style tcell.Style, text string) {
	for _, r := range text {
		(*gui.Screen).SetContent(x, y, r, nil, style)
		x++
	}
}

func (gui GUI) drawSquare(col, row int, p shogi.Piece, sqBg tcell.Color, t theme.Theme) {
	if p.Type == shogi.NoPiece {
		(*gui.Screen).SetContent(col, row, ' ', nil, tcell.StyleDefault.Background((sqBg)))
		(*gui.Screen).SetContent(col+1, row, ' ', nil, tcell.StyleDefault.Background((sqBg)))
	} else {
		code, _ := p.Render()
		piece := code
		pieceStyle := stylePiece(p, sqBg, t)

		(*gui.Screen).SetContent(col+1, row, ' ', nil, tcell.StyleDefault.Background((sqBg)))
		(*gui.Screen).SetContent(col, row, piece, nil, pieceStyle)
	}
}

func (gui GUI) drawRank(col, row int, r shogi.Rank, t theme.Theme) {
	rank := r.Rune()
	rankStyle := tcell.StyleDefault.Foreground(t.Rank)
	gui.drawRune(col, row, rankStyle, rank)
}

func (gui GUI) DrawMsgLabel(msg string, t theme.Theme) {
	topMargin := topMargin + 10
	labelStyle := tcell.StyleDefault.Foreground(t.Msg)
	gui.drawLabel(leftMargin, topMargin, labelStyle, msg)
}

func (gui GUI) drawMoveLabel(gs *shogi.Game) {
	labelStyle := tcell.StyleDefault.Background(gui.Theme.MoveLabelBg).Foreground(gui.Theme.MoveLabelFg)
	var nextPlayer string
	if gs.Board().Turn == shogi.White {
		nextPlayer = " ☖ White to Move "
	} else {
		nextPlayer = " ☗ Black to Move "
	}
	gui.drawLabel(leftMargin+2, topMargin-2, labelStyle, nextPlayer)
}

func (gui GUI) drawMoves(gs *shogi.Game) {
	leftMargin := leftMargin + 22
	boxStyle := tcell.StyleDefault.Foreground(gui.Theme.MoveBox)
	gui.drawLabel(leftMargin, topMargin, boxStyle, "┏━━━━━━━━━━━━━━━━━━━━━┓")
	moves := gs.Moves()
	for i := 0; i < 5; i++ {
		if len(moves)-1 < i {

			row := fmt.Sprintf("┃ %-3v %-7v %-7v ┃", i+1, "", "")
			gui.drawLabel(leftMargin, topMargin+i+1, boxStyle, row)
			continue
		}
		move := moves[len(moves)-1-i]
		moveStr := gs.Notation().EncodeMovement(*move)
		row := fmt.Sprintf("┃ %-3v %-7v %-7v ┃", i+1, moveStr, "")
		gui.drawLabel(leftMargin, topMargin+i+1, boxStyle, row)
	}
	gui.drawLabel(leftMargin, topMargin+6, boxStyle, "┗━━━━━━━━━━━━━━━━━━━━━┛")
}

func (gui GUI) drawPlayers(game *shogi.Game) {
	leftMargin := leftMargin + 22
	emojiStyle := tcell.StyleDefault.Foreground(gui.Theme.Emoji)
	black := fmt.Sprintf("%v %v (Sente)", "☗", game.SentePlayer())
	gui.drawLabel(leftMargin, topMargin+8, emojiStyle, black)
	white := fmt.Sprintf("%v %v (Gonte)", "☖", game.GotePlayer())
	gui.drawLabel(leftMargin, topMargin-2, emojiStyle, white)
}

func (gui GUI) Render(gs *shogi.Game, i *input.Input) {
	gui.drawMoveLabel(gs)
	gui.drawBoard(gs, gui.Theme)
	gui.drawPrompt(i, gui.Theme)
	gui.drawPlayers(gs)
	gui.drawMoves(gs)
	gui.drawHint(gs)

	(*gui.Screen).Show()
}

func (gui GUI) drawPrompt(i *input.Input, t theme.Theme) {
	topMargin := topMargin + 11
	promptStyle := tcell.StyleDefault.Foreground(t.Prompt)
	gui.drawRune(leftMargin, topMargin, promptStyle, '>')
	inputStyle := tcell.StyleDefault.Foreground(t.Input)
	gui.drawLabel(leftMargin+2, topMargin, inputStyle, i.Current())
	(*gui.Screen).ShowCursor(leftMargin+2+i.Length(), topMargin)
}

// idxToRank converts an index to its corresponding rank string
func idxToRank(idx shogi.Rank) string {
	ranks := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	return ranks[idx]
}

// idxToFile converts an index to its corresponding file string
func idxToFile(idx int) string {
	files := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	return files[idx]
}

// idxToSquare returns a string representing the algebraic notation
// for a square given a rank index and a file index
func idxToSquare(rIdx shogi.Rank, fIdx int) string {
	return fmt.Sprintf("%v%v", idxToFile(fIdx), idxToRank(rIdx))
}

// getSquare returns a chess square given a file and a rank
func getSquare(f shogi.File, r shogi.Rank) shogi.Square {
	return shogi.Square((int(r) * 8) + int(f))
}

func (gui GUI) drawBoard(g *shogi.Game, t theme.Theme) {
	row := topMargin

	var r shogi.Rank
	for r = 0; r < numOfSquaresInRow; r++ {
		col := leftMargin
		gui.drawRank(col, row, r, t)
		col += 2
		for f := 0; f < numOfSquaresInRow; f++ {
			sq := shogi.NewSquare(shogi.File(f), shogi.Rank(r))
			sqBg := squareBg(sq, t)
			boardPiece := g.Board().BitBoard[sq]
			var p shogi.Piece
			if boardPiece == "" {
				p = shogi.Piece{
					Type: shogi.NoPiece,
				}
			} else {
				isPromoted := false
				if strings.Contains(boardPiece, "+") {
					isPromoted = true
					boardPiece = boardPiece[1:]
				}
				p = shogi.NewPiece(boardPiece, isPromoted)
			}
			gui.drawSquare(col, row, p, sqBg, t)
			col += 2
		}
		row++
	}

	fileStyle := tcell.StyleDefault.Foreground(t.File)
	gui.drawLabel(leftMargin+2, row, fileStyle, "1 2 3 4 5 6 7 8 9")
}

func (gui *GUI) SetHint(movement string) {
	gui.Hint = movement
}

func (gui *GUI) drawHint(g *shogi.Game) {
	m, err := g.Notation().DecodeHodgesMove(gui.Hint)
	if err != nil {
		return
	}

	srcFile := int(m.Origin.File())
	srcRank := int(m.Origin.Rank())
	dstFile := int(m.Destination.File())
	dstRank := int(m.Destination.Rank())

	srcX := leftMargin + 2 + 2*srcFile
	srcY := topMargin + srcRank
	dstX := leftMargin + 2 + 2*dstFile
	dstY := topMargin + dstRank

	boardPiece := g.Board().BitBoard[m.Origin]
	var piece shogi.Piece
	if boardPiece == "" {
		piece = shogi.Piece{Type: shogi.NoPiece}
	} else {
		isPromoted := false
		if strings.Contains(boardPiece, "+") {
			isPromoted = true
			boardPiece = boardPiece[1:]
		}
		piece = shogi.NewPiece(boardPiece, isPromoted)
	}

	srcBg := squareBg(m.Origin, gui.Theme)
	// dstBg := squareBg(m.Destination, gui.Theme)

	highlightFg := gui.Theme.PieceHint
	highlightBg := gui.Theme.SquareHint

	srcHighlightStyle := tcell.StyleDefault.Background(srcBg).Foreground(highlightFg).Bold(true)

	// Redraw the origin square with the highlighted piece.
	pieceRune, _ := piece.Render()
	gui.drawRune(srcX, srcY, srcHighlightStyle, pieceRune)

	// For the destination square, get any piece present.
	// g.Board().BitBoard[shogi.NewSquare(shogi.File(srcFile), shogi.Rank(srcRank))]

	destCode := g.Board().BitBoard[m.Destination]
	var destPiece shogi.Piece
	if destCode == "" {
		destPiece = shogi.Piece{Type: shogi.NoPiece}
	} else {
		isPromoted := false
		if strings.Contains(destCode, "+") {
			isPromoted = true
			destCode = destCode[1:]
		}
		destPiece = shogi.NewPiece(destCode, isPromoted)
	}

	// Redraw the destination square with the highlight background.
	gui.drawSquare(dstX, dstY, destPiece, highlightBg, gui.Theme)

	gui.DrawMsgLabel(fmt.Sprintf("(%s) Accept hint? y/n", gui.Hint), gui.Theme)
	(*gui.Screen).Show()
}
