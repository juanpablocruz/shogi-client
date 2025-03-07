# Shogo client

Shogo Client is a terminal-based Shogi game written in Go. It implements full Shogi game logic—from board setup and move validation (using SFEN notation) to piece-specific movement rules and captured-piece (hand) management. The project also features a terminal user interface built with [tcell](https://github.com/gdamore/tcell), support for networked play, and integration with AI agents and engine commands using a USI-inspired protocol.

## Features

- **Complete Shogi Game Logic:**  
  - Board representation and SFEN encoding/decoding for positions.
  - Movement rules for all piece types, including promotions, captures, and drops.
  - Management of captured pieces (hand) and notation handling.  
  _See:_ `internal/shogi/board.go`, `piece.go`, `notation.go`, etc.

- **Terminal-Based GUI:**  
  - Interactive rendering of the Shogi board, moves, logs, and prompts using the tcell library.
  - Dynamic UI updates, including hint display and move logs.  
  _See:_ `internal/gui/gui.go`, `render.go`

- **Command Processing & Input Handling:**  
  - Processes user commands for moves, hints, saving game state, resetting, and quitting.
  - Keyboard input management and command parsing.  
  _See:_ `internal/cmd/cmd.go`, `internal/input/input.go`

- **Engine & AI Integration:**  
  - Supports integration with AI agents (e.g., OpenAI and Claude) to suggest moves or respond with hints.
  - Implements USI-like commands (e.g., `usi`, `id`, `position`, `go`, etc.) for engine communication.  
  _See:_ `internal/agent/`, `internal/engine/engine.go`, `gui_engine.go`

- **Modular and Extensible Architecture:**  
  - Clearly separated packages for game logic, UI, engine communication, configuration, and agent integration.
  - Includes tests for critical modules such as board logic, notation, and engine functionality.

## Requirements

- **Go:** Version 1.16 or later.
- **Environment Variables:** Configuration is managed via a `.env` file (loaded using [godotenv](https://github.com/joho/godotenv)).
- **API Keys:** For AI integration (if you choose to use an AI agent), supply your keys in the configuration.

## Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/juanpablocruz/shogi-client.git
   cd shogi-client
   ```
2. Download Dependencies:

```bash
go mod tidy
```

3. Build the Application:

```bash
go build -o shogo ./cmd/shogo
```


## Configuration

Create a .env file in the repository root with settings similar to the example below. Adjust values as needed:

```dotenv
# AI Agent selection: "openai" or "claude"
AGENT=openai

# API keys for AI integration
OPENAI_API_KEY=your_openai_api_key_here
# or if using Claude:
CLAUDE_API_KEY=your_claude_api_key_here

# Network and game settings
PORT=8080
SENTE_PLAYER=Player1
GOTE_PLAYER=Player2
```


The application loads these values during initialization.

## Usage 

1. Start the Application:

```bash
./shogo
```

2. Gameplay Instructions:
- Starting Position: The game begins with a standard SFEN starting position.
- Input Moves: Enter moves using Shogi notation. The UI supports move entry, AI hints (by typing `hint`), saving the board state 
via (`save`), and resetting the game (`reset`).
- Exit: Use __Escape__ or __Ctrl+C__ to quit.
- AI Integration: When you enter `hint`, the board's SFEN string is sent to the configured AI agent which returns a suggested move in Hodges notation.
- Engine Commands: The client supports USI-style commands (e.g., position, go, stop) to facilitate network play and engine integration.

3. GUI & Logs:
The terminal UI displays the board, current moves, logs, and hints dynamically, updating after each command.

## Repository Structure

```graphql

.
├── cmd
│   └── shogo
│       └── main.go         # Entry point; loads configuration, initializes GUI and game
├── go.mod
├── go.sum
└── internal
    ├── agent               # AI agent implementations
    │   ├── agent.go         # Agent interface definitions
    │   ├── claude_agent.go  # Claude agent implementation
    │   └── openai_agent.go  # OpenAI agent implementation
    ├── cmd                 # Command processing logic
    │   └── cmd.go         # Handles move commands, hints, game resets, etc.
    ├── config              # Application configuration
    │   └── config.go      # Environment configuration loader
    ├── engine              # Engine communication and USI-like command processing
    │   ├── engine.go           # Core engine command handling
    │   ├── engine_api.go       # Engine messaging interface
    │   ├── engine_test.go      # Engine module tests
    │   ├── gui_engine.go       # GUI integration with engine commands
    │   └── gui_engine_test.go  # GUI engine tests
    ├── gui                 # Terminal user interface components
    │   ├── gui.go         # Initialization and event processing for the GUI
    │   └── render.go      # Board and UI rendering functions
    ├── input               # User input handling
    │   └── input.go       # Input buffering and processing
    ├── shogi               # Core Shogi game logic and types
    │   ├── board.go           # Board representation and SFEN parsing
    │   ├── board_test.go      # Board tests
    │   ├── color.go           # Definitions for piece colors
    │   ├── command.go         # Command constants for engine/GUI communication
    │   ├── game.go            # Game state and move processing
    │   ├── hand.go            # Management of captured pieces ("hand")
    │   ├── hand_test.go       # Tests for hand functionality
    │   ├── move.go            # Move structure and move type definitions
    │   ├── notation.go        # SFEN notation and movement encoding/decoding
    │   ├── notation_test.go   # Notation module tests
    │   ├── piece.go           # Piece types, movement rules, and rendering
    │   ├── piece_test.go      # Piece module tests
    │   ├── player.go          # Player structure (for future expansion)
    │   ├── square.go          # Square type and algebraic notation utilities
    │   └── utils.go           # Utility functions (e.g., path checking, sign/abs)
    └── theme               # Theme and styling definitions for the GUI
        └── theme.go       # GUI color and style settings
```
## Testing

Several modules include test to verify game logic, board handling, and engine integration. To run the tests:

```bash
go test ./internal/...
```

## License

This project is licensed under the MIT License. See [License](./LICENSE.md) for details.
