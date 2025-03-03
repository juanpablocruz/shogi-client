package engine

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

type Engine struct {
	IsInitialized bool
	EngineID      string
	EngineOptions map[string]EngineOption
	EngineAPI     EngineAPI
	Game          *shogi.Game
}

func NewEngine(id string, api EngineAPI, game *shogi.Game, options map[string]EngineOption) *Engine {
	return &Engine{
		IsInitialized: false,
		EngineID:      id,
		EngineOptions: options,
		EngineAPI:     api,
		Game:          game,
	}
}

func (e *Engine) ListenCMD() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	m, _ := e.EngineAPI.ReceiveMessage(ctx)
	cmd, args, err := e.parseGUICommand(m)
	if err != nil {
		return err
	}

	return e.ProcessGUICMD(cmd, args)
}

func (e *Engine) parseGUICommand(str string) (shogi.GUICommand, []string, error) {
	parts := strings.Split(str, " ")
	if len(parts) < 1 {
		return shogi.NoGUICommand, []string{}, fmt.Errorf("unable to parse gui command %s", str)
	}

	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}
	switch parts[0] {
	case "usi":
		return shogi.USI, args, nil
	case "debug":
		return shogi.Debug, args, nil
	case "isready":
		return shogi.IsReady, args, nil
	case "setoption":
		return shogi.SetOption, args, nil
	case "register":
		return shogi.Register, args, nil
	case "usinewgame":
		return shogi.USINewGame, args, nil
	case "position":
		return shogi.Position, args, nil
	case "go":
		return shogi.Go, args, nil
	case "stop":
		return shogi.Stop, args, nil
	case "ponderhit":
		return shogi.Ponderhit, args, nil
	case "gameover":
		return shogi.Gameover, args, nil
	case "quit":
		return shogi.Quit, args, nil
	default:
		return shogi.NoGUICommand, []string{}, fmt.Errorf("unable to parse gui command %s", str)
	}
}

func (e *Engine) ProcessGUICMD(cmd shogi.GUICommand, args []string) error {
	// TODO: Implement
	switch cmd {
	case shogi.USI:
		return e.ProcessCMD(shogi.Id)
	case shogi.Position:
		return e.ProcessPosition(args)
	}
	return nil
}

func (e *Engine) ProcessCMD(cmd shogi.EngineCommand, args ...string) error {
	switch cmd {
	case shogi.Id:
		return e.sendId()
	case shogi.BestMove:
		return e.sendBestMove()
	case shogi.Checkmate:
		return e.sendCheckMate(args)
	case shogi.Info:
		return e.sendInfo()
	case shogi.Option:
		return e.sendOptions()
	case shogi.ReadyOk:
		return e.sendReady()
	case shogi.USIOk:
	default:
		panic(fmt.Sprintf("unexpected shogi.EngineCommand: %#v", cmd))
	}
	return nil
}

// id
// `name <x>` - This must be sent after receiving the `usi` command to identify the engine, e.g. id name Shredder X.Y\n
// `author <x>` - This must be sent after receiving the `usi` command to identify the engine, eg.g. id author Stefan MK\n
func (e *Engine) sendId() error {
	err := e.EngineAPI.SendMessage(fmt.Sprintf("id name %s", e.EngineID))
	if err != nil {
		return err
	}

	err = e.sendOptions()
	if err != nil {
		return err
	}
	// usiok
	// Must be sent after the `id` and optional options to tell the GUI that the engine has sent all infos and is ready in usi mode.
	return e.EngineAPI.SendMessage("usiok")
}

// readyok
// This must be sent when the engine has received an `isready` command and has processed all input and is ready to accept new commands now.
// It is usually sent after a command that can take some time to be able to wait for the engine, but it can be used anytime,
// event when the engine is searching, and must always be answered with `readyok`.
func (e *Engine) sendReady() error {
	return e.EngineAPI.SendMessage("readyok")
}

// bestmove <move1> [ponder <move2>]
// bestmove [resign | win]
// The engine has stopped searching and found the move <move> best in this position. The engine can send the move it likes to ponder on.
// The engine must not start pondering automatically. This command must always be sent if the engine stops searching,
// also in pondering mode if there is a `stop` command, so for every `go` command a `bestmove` command is needed!
func (e *Engine) sendBestMove() error {
	// TODO: implement
	return fmt.Errorf("notimplemented")
}

// checkmate [<move1> ... <movei> | nomate | timeout | notimplemented]
// As `go mate` is not supported we always reply with `checkmate notimplemented`
func (e *Engine) sendCheckMate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("checkmate expecting arguments, none received")
	}
	return e.EngineAPI.SendMessage(fmt.Sprintf("checkmate %s", strings.Join(args, " ")))
}

// info
// The engine wants to send information to the GUI. This should be done whenever one of the info has changed.
// The engine can send only selected infos or multiple infos with one info command, e.g. info currmove 2g2f currmovenumber 1
// Also all infos belonging to the pv should be sent together, e.g.
// info depth 2 score cp 214 time 1242 nodes 2124 nps 34928 pv 2g2f 8c8d 2f2e
// I suggest to start sending `currmove`,`currmovenumber`,`currline` and `refutation` only after one second in order to avoid too much traffic.
//
// `depth <x>` - search depth in plies.
// `seldepth <x>` - selective search depth in plies. If the engine sends seldepth there must also be a depth present in the same string.
// `time <x>` - The time searched in ms. This should be sent together with the pv.
// `nodes <x>` - x nodes searched. The engine should send this info regularly.
// `pv <move1> ... <movei>` - The best line found.
// `multipv <num>` - This for the multi pv mode. For the best move/pv add `multipv 1` in the string when you send the pv.
//
//	In k-best mode, always send all k variants in k strings together.
//
// `score`
//
//	`cp <x>` - The score from the engine's point of view, in centipawns.
//	`mate <y>` - Mate in y plies. If the engine is getting mated, use negative values for y.
//	`lowerbound` - The score is just a lower bound.
//	`upperbound` - The score is just an upper bound.
//
// `currmove <move>` - Currently searching this move
// `currmovenumber <x>` - Currently searching move number x, for the first move x should be 1, not 0.
// `hashfull <x>` - The hash is x permill full. The engine should send this info regularly.
// `npx <x>` - x nodes per second searched. The engine should send this info regularly.
// `cpuload <x>` - The cpu usage of the engine is x permill.
// `string <str>` - Any string str which will be displayed by the engine. If there is a string command the rest of the line will be interpreted as <str>.
// `refutation <move1> <move2> ... <movei>` - Move <move1> is refuted by the line <move2> ... <movei>, where i can be any number >= 1.
//
//	example: after move 8h2b+ is searched, the engine can send `info refutation 8h2b+ 1c2b` if 1c2b is the best answer after 8h2b+
//	          or if 1c2b refutes the move 8h2b+. If there is no refutation for 8h2b+ found, the engine should just send `info refutation 8h2b+`.
//	The engine should only send this if the option `USI_ShowRefutations` is set to true.
//
// `currline <cpunr> <move1> ... <movei>` - This is the current line the engine is calculating. <cpunr> is the number of the cpu if the engine
// is running on more than one cpu. <cpunr> = 1,2,3,... If the engine is just using one cpu, <cpunr> can be omitted.
// If <cpunr> is greater than 1, always send all k lines in k strings together. The engine should only send this if the option USI_ShowCurrLine is set to true.
func (e *Engine) sendInfo() error {
	// TODO: implement
	return fmt.Errorf("notimplemented")
}

// option
// This command tells the GUI which parameters can be changed in the engine. This should be sent once at engine startup after the `usi` and the `id` commands
// if any parameter can be changed in the engine. The GUI should parse this and build a dialog for the user to change the settings.
// Note that not every option should appear in this dialog, as some options like USI_Ponder, USI_AnalyseMode, etc. are better handled elsewhere or are set automatically.
// If the user wants to change some settings, the GUI will send a `setoption` command to the engine.
// Note that the GUI need not send the setoption command when starting the engine for every option if it doesn't want to change the default value.
// For all allowed combinations see the examples below, as some combinations of this tokens don't make sense.
// One string will be sent for each parameter.
//
// `name <id>` - The option has the name <id>. Whitespace is not allowed in an option name. Note that the name should normally not be displayed directly in the GUI:
// The GUI should look up the option name in the translation file, and present the translation into the users preferred language in the engine's option dialog.
// Certain options have a fixed value for <id>, which means that the semantics of this option is fixed. Usually those options should not be displayed in the normal engine
// options window of the GUI but get a special treatment. USI_Pondering for example should be set automatically when pondering is enabled or disabled in the GUI options.
// The same for USI_AnalyseMode which should also be set automatically by the GUI. All those certain options have the prefix USI_.
// If the GUI gets an unknown option with the prefix USI_, it should just ignore it and not display it in the engine's options dialog.
//
// The options with fixed semantics are:
// <id> = USI_Hash, type spin
// The value in MB for memory for hash tables can be changed, this should be answered with the first `setoptions` command at program boot if the engine has sent
// the appropiate option name `Hash` command, wchich should be supported by all engines! So the engine should use a very small hash first as default.
// <id> = USI_Ponder, type check
// This means that the engine is able to ponder (i.e. think during the opponent's time). The GUI will send this whenever pondering is possible or not. Note: The engine
// should not start pondering on its own if this is enabled, this option is only needed because the engine might change its time management algorithm when pondering is allowed.
// <id> = USI_OwnBook, type check
// This means that the engine has its own opening book which is accessed by the engine itself. If this is set, the engine takes care of the opening book and the GUI will
// never execute a move out of this book for the engine. If this is set to false by the GUI, the engine should not access its own book.
// <id> = USI_MultiPV, type spin
// The engine supports multi best line or k-best mode. The default value is 1.
// <id> = USI_ShowCurrLine, type check
// The engine can show the current line it is calculating. See `info currline` above. This option should be false by default.
// <id> = USI_ShowRefutations, type check
// The engine can show a move and its refutation in a line. See `info refutations` above. This option should be false by default.
// <id> = USI_LimitStrength, type check
// The engine is able to limit its strength to a specific dan/kyu number. This should always be implemented together with USI_Strength. This option should be false by default.
// <id> = USI_Strength, type spin
// The engine can limit its strength within the given interval. Negative numbers are kyu levels, while positive numbers are amateur dan levels.
// If USI_LimitStrength is set to false, this value should be ignored. If USI_LimitStrength is set to true, the engine should play with this specific strength.
// This option should always be implemented together with USI_LimitStrength
// <id> = USI_AnalyseMode, type check
// The engine wants to behave differently when analysing or playing a game. For example when playing it can use some kind of learning, or an asymetric evaluation function.
// The GUI should set this option to false if the engine is playing a game, and to true if the engine is analysing.
//
// `type <t>` - The option has type t. There are 5 different types of options the engine can send:
// check - A checkbox that can either be true or false
// spin - A spin wheel or slider that can be an integer in a certain range.
// combo - A combo box that can have different predefined strings as a value.
// button - A button that can be pressed to send a command to the engine.
// string - A text field that has a string as a value, an empty string has the value <empty>.
// filename - Similar to string, but is presented as a file browser instead of a text field in the GUI.
//
// `default <x>` - The default value of this parameter is x.
// `min <x>` - The minimum value of this parameter is x.
// `max <x>` - The maximum value of this parameter is x.
//
// Examples:
// "option name Nullmove type check default true\n"
// "option name Selectivity type spin default 2 min 0 max 4\n"
// "option name Style type combo default Normal var Solid var Normal var Risky\n"
// "option name LearningFile type filename default /shogi/my-shogi-engine/learn.bin"
// "option name ResetLearning type button\n"
func (e *Engine) sendOptions() error {
	for key, opt := range e.EngineOptions {
		defaultVal := ""
		if opt.Default != "" {
			defaultVal = fmt.Sprintf(" default %s", opt.Default)
		}
		optStr := fmt.Sprintf("option name %s type %s%s\n", key, opt.Type, defaultVal)
		_ = e.EngineAPI.SendMessage(optStr)
	}
	return nil
}

/*
		Process the position command from the gui

	  position [sfen <sfenstring> | startpos] moves <move1> ... <movei>

	  Set up the position described in sfenstring on the internal board and play the moves on the internal board.
	  If the game was played from the start position, the string `startpos` will be sent.
*/
func (e *Engine) ProcessPosition(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("position requires at least 2 arguments: moves and the movement string, received: %v", args)
	}

	startPost := ""
	isMove := false
	currStr := ""
	movements := []string{}
	for _, arg := range args {
		switch arg {
		case "moves":
			isMove = true
			startPost = currStr
			currStr = ""
		default:
			if isMove {
				movements = append(movements, arg)
			} else {
				currStr = fmt.Sprintf("%s %s", currStr, arg)
			}
		}
	}

	if len(movements) < 1 {
		return fmt.Errorf("invalid position command, expecting `moves`, received: %v", args)
	}

	if startPost != "" {
		b, err := e.Game.Notation().DecodeBoard(startPost)
		if err != nil {
			return err
		}

		e.Game.SetBoard(&b)
	}

	for _, m := range movements {
		mo, err := e.Game.Notation().DecodeMovement(m)
		if err != nil {
			continue
		}

		err = e.Game.Move(mo)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}

	return nil
}
