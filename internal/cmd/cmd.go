package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/juanpablocruz/shogo/clientr/internal/gui"
	"github.com/juanpablocruz/shogo/clientr/internal/input"
	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

func resetGame(game *shogi.Game) *shogi.Game {
	newGame := shogi.NewGame(game.SentePlayer(), game.GotePlayer())
	return newGame
}

func saveGame(game *shogi.Game) string {
	return game.Board().String()
}

func hint(game *shogi.Game, gui *gui.GUI, in *input.Input) string {
	gui.DrawMsgLabel("Thinking...", gui.Theme)
	gui.Render(game, in)

	if game.GetAIClient() != nil {
		ai := game.GetAIClient()
		boardStr := game.Board().String()
		gui.AppendLog(fmt.Sprintf("Sending AI client: %s", boardStr))
		h, err := ai.AskHint(boardStr)
		if err != nil {
			gui.AppendLog(fmt.Sprintf("shogo error: ask hint error %v", err))
			return ""
		}

		gui.AppendLog(fmt.Sprintf("AI responds: %s", h))
		gui.SetHint(h)
		return h
	} else {
		gui.AppendLog("No ai client found")
	}

	return strings.Repeat(" ", 80)
}

func ProcessCmd(cmd string, game *shogi.Game, gui *gui.GUI, in *input.Input) (string, *shogi.Game) {
	cmd = strings.TrimSpace(cmd)

	if len(cmd) == 0 {
		return strings.Repeat(" ", 80), game
	}

	switch cmd {
	case "quit":
		gui.Quit()
		os.Exit(0)
		return "", game
	case "save":
		return saveGame(game), game
	case "reset":
		return strings.Repeat(" ", 80), resetGame(game)
	case "hint":
		return hint(game, gui, in), game
	case "y":
		if gui.Hint != "" {
			m, err := game.Notation().DecodeHodgesMove(gui.Hint)
			gui.Hint = ""
			if err != nil {
				return strings.Repeat(" ", 80), game
			}
			_ = game.Move(m)
			return strings.Repeat(" ", 80), game

		}
		gui.Hint = ""
	case "n":

		gui.Hint = ""
	default:

		m, err := game.Notation().DecodeHodgesMove(cmd)
		if err == nil {
			game.Move(m)

			gui.AppendLog(fmt.Sprintf("%s -> %s", cmd, game.Board().String()))

			return strings.Repeat(" ", 80), game
		}
		if err := game.MoveStr(cmd); err != nil {
			return "\u26A0 Illegal. Try again.", game
		}

		gui.AppendLog(fmt.Sprintf("%s -> %s", cmd, game.Board().String()))
	}
	return strings.Repeat(" ", 80), game
}
