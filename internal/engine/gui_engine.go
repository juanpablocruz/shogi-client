package engine

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
)

type GUIEngine struct {
	IsInitialized bool
	EngineID      string
	EngineOptions map[string]EngineOption
	EngineAPI     EngineAPI
}

func NewGUIEngine(e EngineAPI) *GUIEngine {
	return &GUIEngine{
		IsInitialized: false,
		EngineID:      "",
		EngineOptions: make(map[string]EngineOption),
		EngineAPI:     e,
	}
}

func (e GUIEngine) ListenCMD() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m, _ := e.EngineAPI.ReceiveMessage()
				e.ProcessEngineCMD(m)
			}
		}
	}()
}

func (e GUIEngine) ProcessEngineCMD(str string) error {
	// TODO: Implement this
	return nil
}

func (e GUIEngine) ProcessCMD(cmd shogi.GUICommand, args ...string) error {
	switch cmd {
	case shogi.USI:
		return e.initializeUSI()
	case shogi.Debug:
		isOn := false
		if len(args) > 0 && args[0] == "on" {
			isOn = true
		}
		return e.toggleDebug(isOn)
	case shogi.IsReady:
		return e.synchReady()
	case shogi.Register:
		return e.register(args)
	case shogi.USINewGame:
		return e.newGame()
	case shogi.Position:
		if len(args) < 2 {
			return fmt.Errorf("invalid command position arguments, %v", args)
		}
		sfen := args[0]
		return e.position(sfen, args[1:])
	case shogi.SetOption:

		if len(args) < 2 {
			return fmt.Errorf("invalid command setoption arguments, %v", args)
		}
		return e.setOption(args[0], args[1])
	case shogi.Go:
		return e.sendGo()
	case shogi.Stop:
		return e.stop()
	case shogi.Ponderhit:
		return e.ponderHit()
	case shogi.Gameover:

		if len(args) < 1 {
			return fmt.Errorf("invalid command gameover arguments, %v", args)
		}
		return e.gameOver(args[1])
	case shogi.Quit:
		return e.quit()
	}
	return nil
}

func (e GUIEngine) sendCommand(cmd string) error {
	return e.EngineAPI.SendMessage(cmd)
}

func (e GUIEngine) readMsg() (string, error) {
	return e.EngineAPI.ReceiveMessage()
}

func (e GUIEngine) receiveMessage(ctx context.Context, msg chan string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m, _ := e.readMsg()
				msg <- m
			}
		}
	}()
}

func maxTimeoutWait(ms time.Duration, cancel context.CancelFunc) {
	time.AfterFunc(ms, cancel)
}

func (e GUIEngine) receiveOptions(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("invalid response received, expecting option <args>, but received %v", args)
	}
	var opt EngineOption
	key := ""

	for _, ar := range args {
		switch ar {
		case "name":
			key = "Name"
		case "type":
			key = "Type"
		case "default":
			key = "Default"
		default:
			if key == "" {
				return fmt.Errorf("invalid option, expecting name|type|default, received %s", ar)
			}
			switch key {
			case "Name":
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
				opt.Name = ar
			case "Type":
				// `type <t>` - The option has type t. There are 5 different types of options the engine can send:
				// check - A checkbox that can either be true or false
				// spin - A spin wheel or slider that can be an integer in a certain range.
				// combo - A combo box that can have different predefined strings as a value.
				// button - A button that can be pressed to send a command to the engine.
				// string - A text field that has a string as a value, an empty string has the value <empty>.
				// filename - Similar to string, but is presented as a file browser instead of a text field in the GUI.
				opt.Type = ar
			case "Default":
				// `default <x>` - The default value of this parameter is x.
				// `min <x>` - The minimum value of this parameter is x.
				// `max <x>` - The maximum value of this parameter is x.
				opt.Default = fmt.Sprintf("%s %s", opt.Default, ar)
			default:
			}
		}
	}
	e.EngineOptions[opt.Name] = opt
	return nil
}

// usi
// Tell engine to use the USI.
// This will be sent once as a first command after program boot to tell the engine to switch to USI mode.
// After receiving the usi command the engine must identify itself with the `id` command and send the `option` command
// to tell the GUI which engine settings the engine support.
// After that, the engine should send `usiok` to acknowledge the USI mode.
// If no `usiok` is sent within a certain time period, the engine task will be killed by the GUI.
func (e GUIEngine) initializeUSI() error {
	if err := e.sendCommand("usi"); err != nil {
		return err
	}
	// wait for usiok
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxTimeoutWait(3*time.Second, cancel)

	ch := make(chan string)
	defer close(ch)

	e.receiveMessage(ctx, ch)

	state := 0

	for {
		select {
		case msg := <-ch:
			parts := strings.Split(msg, " ")
			if len(parts) < 1 {
				return fmt.Errorf("invalid msg received %s", msg)
			}
			cmd := parts[0]
			args := parts[1:]
			switch state {
			case 0:

				if cmd == "id" {
					state = 1
					if len(args) < 1 {
						return fmt.Errorf("invalid response received, expecting id <id>, received id %v", args)
					}
					e.EngineID = args[0]
				} else {
					return fmt.Errorf("invalid response received, expecting id <id>, received %v", msg)
				}
			case 1:
				if cmd == "option" {
					if err := e.receiveOptions(args); err != nil {
						return err
					}
				}
				if cmd == "usiok" {
					return nil
				}
			}
		case <-ctx.Done():
			return fmt.Errorf("expecting usiok within a timeframe, but no usiok received")
		}
	}
}

// debug [on|off]
func (e GUIEngine) toggleDebug(isOn bool) error {
	cmd := "debug"
	if isOn {
		cmd = fmt.Sprintf("%s on", cmd)
	} else {
		cmd = fmt.Sprintf("%s off", cmd)
	}
	return e.sendCommand(cmd)
}

func (e GUIEngine) synchReady() error {
	return e.sendCommand("isready")
}

func (e GUIEngine) setOption(option string, value string) error {
	return e.sendCommand(fmt.Sprintf("setoption name %s value %s", option, value))
}

// register later
// register name Stefan MK code 4359874324
func (e GUIEngine) register(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("invalid arguments, received: %s", args)
	}

	type RegisterCMD struct {
		cmd      string
		argument string
	}

	allCMDs := []RegisterCMD{}
	currCMD := RegisterCMD{
		cmd: "",
	}
	for _, s := range args {
		switch s {
		case "later":
			if currCMD.cmd != "" {
				allCMDs = append(allCMDs, currCMD)
				currCMD.cmd = ""
				currCMD.argument = ""
			}
			allCMDs = append(allCMDs, RegisterCMD{cmd: "later"})
		case "name":
			if currCMD.cmd != "" {
				allCMDs = append(allCMDs, currCMD)
			}
			currCMD = RegisterCMD{
				cmd:      "name",
				argument: "",
			}
		case "code":
			if currCMD.cmd != "" {
				allCMDs = append(allCMDs, currCMD)
			}
			currCMD = RegisterCMD{
				cmd:      "code",
				argument: "",
			}
		default:
			currCMD.argument = fmt.Sprintf("%s %s", currCMD.argument, s)
		}
	}
	if currCMD.cmd != "" {
		allCMDs = append(allCMDs, currCMD)
	}

	allCMDscmdStr := "register"
	for _, cmd := range allCMDs {
		switch cmd.cmd {
		case "later":
			allCMDscmdStr = fmt.Sprintf("%s later", allCMDscmdStr)
		case "name":
			if len(cmd.argument) != 1 {
				return fmt.Errorf("invalid command, register name <x> requires an argument x, none received")
			}
			allCMDscmdStr = fmt.Sprintf("%s name %s", allCMDscmdStr, cmd.cmd)
		case "code":
			if len(cmd.argument) != 1 {
				return fmt.Errorf("invalid command, register code <x> requires an argument x, none received")
			}
			allCMDscmdStr = fmt.Sprintf("%s code %s", allCMDscmdStr, cmd.cmd)
		}
	}

	if allCMDscmdStr == "register" {
		return fmt.Errorf("invalid command register without arguments")
	}
	return e.sendCommand(allCMDscmdStr)
}

func (e GUIEngine) newGame() error {
	return e.sendCommand("usinewgame")
}

func (e GUIEngine) position(sfen string, moves []string) error {
	movesString := "moves "
	for _, m := range moves {
		movesString = fmt.Sprintf("%s %s", movesString, m)
	}
	return e.sendCommand(fmt.Sprintf("position %s %s", sfen, movesString))
}

func (e GUIEngine) sendGo() error {
	// TODO: implement this?
	return nil
}

func (e GUIEngine) stop() error {
	return e.sendCommand("stop")
}

func (e GUIEngine) ponderHit() error {
	return e.sendCommand("ponderhit")
}

func (e GUIEngine) gameOver(outcome string) error {
	if outcome != "win" && outcome != "lose" && outcome != "draw" {
		return fmt.Errorf("invalid command, expecting gameover [win|lose|draw], received argument %s", outcome)
	}
	return e.sendCommand(fmt.Sprintf("gameover %s", outcome))
}

func (e GUIEngine) quit() error {
	return e.sendCommand("quit")
}
