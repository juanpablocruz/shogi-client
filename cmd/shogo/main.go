package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/joho/godotenv"
	"github.com/juanpablocruz/shogo/clientr/internal/agent"
	"github.com/juanpablocruz/shogo/clientr/internal/cmd"
	"github.com/juanpablocruz/shogo/clientr/internal/config"
	"github.com/juanpablocruz/shogo/clientr/internal/gui"
	"github.com/juanpablocruz/shogo/clientr/internal/input"
	"github.com/juanpablocruz/shogo/clientr/internal/shogi"
	"github.com/juanpablocruz/shogo/clientr/internal/theme"
)

type Message struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

func ConnectClient(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error establishing connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	joinMsg := Message{
		Action: "join",
		Data:   make(map[string]interface{}),
	}
	data, err := json.Marshal(joinMsg)
	if err != nil {
		fmt.Println("Error sending join: ", err)
		os.Exit(1)
	}

	conn.Write(data)

	decoder := json.NewDecoder(conn)

	for {
		var msg Message
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			os.Exit(1)
		}

		fmt.Printf("Received: %s: %v", msg.Action, msg.Data)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	var aiClient agent.Agent
	agentType := os.Getenv("AGENT")
	if agentType == "openai" {
		openai_key := os.Getenv("OPENAI_API_KEY")

		if openai_key == "" {
			log.Fatal("No OPENAI_API_KEY found in env")
		}
		aiClient = agent.NewOpenAIAgent(openai_key)
	} else {

		claude_key := os.Getenv("CLAUDE_API_KEY")

		if claude_key == "" {
			log.Fatal("No CLAUDE_API_KEY found in env")
		}
		aiClient = agent.NewClaudeAgent(claude_key)
	}

	config := config.Init()

	address := fmt.Sprintf("127.0.0.1:%d", config.Port)
	fmt.Println("Connecting to ", address)

	gs := *shogi.NewGame(config.SentePlayer, config.GotePlayer,
		func(g *shogi.Game) {
			g.SetAIClient(aiClient)
		})

	gui := gui.NewGUI()
	gui.Theme = theme.ThemeBasic

	gui.AppendLog("Initializing....")

	in := input.NewInput()

	defer gui.Quit()

	gui.Render(&gs, in)

	gui.AppendLog("Initialized.")

	for {
		_ = Interact(gui, in, &gs)

		gui.Render(&gs, in)
	}
}

func Interact(gui *gui.GUI, in *input.Input, gs *shogi.Game) bool {
	rescore := true
	ev := (*gui.Screen).PollEvent()
	quit := func() {
		gui.Quit()
		os.Exit(0)
	}
	msg := ""

	// gui.Update()

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			quit()
		case tcell.KeyEnter:
			msg, gs = cmd.ProcessCmd(in.Current(), gs, gui, in)
			gui.DrawMsgLabel(msg, gui.Theme)
			in.Clear()
			gui.Render(gs, in)

		case tcell.KeyBackspace, tcell.KeyBackspace2:
			rescore = false
			in.Backspace()
		default:
			in.Append(ev.Rune())
		}

	case *tcell.EventResize:
		(*gui.Screen).Sync()
	}
	return rescore
}
