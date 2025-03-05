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
		h, err := ai.AskHint(game.Board().String())
		if err != nil {
			fmt.Printf("shogo: ask hint error %v", err)
			return ""
		}

		gui.SetHint(h)
		return fmt.Sprintf("%s%s", h, strings.Repeat(" ", 80-len(h)))

	} else {
		fmt.Printf("No ai client found \n")
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
		if err := game.MoveStr(cmd); err != nil {
			return "\u26A0 Illegal. Try again.", game
		}
	}
	return strings.Repeat(" ", 80), game
}
