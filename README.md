# Network Pong Battle

A multiplayer real-time paddle game implemented in Go, where each player controls two paddles on adjacent edges of a square playfield. Players must prevent balls from hitting their edges, earning points when the opponent fails to block.

## Features

- **Two-player multiplayer** over network (LAN or Internet)
- **Square playfield** (600×600 pixels)
- **Dual paddles per player** on adjacent edges
- **Multiple balls** with configurable count
- **Real-time synchronization** of paddle positions, ball movements, and scores
- **TCP networking** for reliable gameplay
- **Modern UI** built with Ebiten graphics library
- **Configurable settings** for game parameters

## Game Mechanics

- **Player 1**: Controls left (vertical) and top (horizontal) paddles
- **Player 2**: Controls right (vertical) and bottom (horizontal) paddles
- **Scoring**: Points are awarded when balls hit the opponent's walls
- **Game End**: First player to reach target score or when time limit expires
- **Controls**: 
  - Player 1: W/S for left paddle, A/D for top paddle
  - Player 2: W/S for right paddle, A/D for bottom paddle

## Architecture

The project follows a clean architecture with separate concerns:

- **Game Logic** (`internal/game/`): Core game mechanics, collision detection, scoring
- **Networking** (`internal/net/`): TCP client-server communication, message protocol
- **User Interface** (`internal/ui/`): Graphics rendering, input handling, menu system
- **Executables** (`cmd/`): Server and client entry points

## Prerequisites

- Go 1.24 or later
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd network-pong-battle
```

2. Install dependencies:
```bash
go mod tidy
```

## Usage

### Starting the Server

1. Start the game server:
```bash
go run cmd/server/main.go
```

By default, the server listens on port 8080. You can specify a different port:

```bash
go run cmd/server/main.go -port 9000
```

### Starting the Client

1. In a new terminal, start the client:
```bash
go run cmd/client/main.go
```

By default, the client connects to `localhost:8080`. You can specify a different server:

```bash
go run cmd/client/main.go -server 192.168.1.100:8080
```

You can also set a custom player name:

```bash
go run cmd/client/main.go -name "Player1"
```

### Game Controls

- **Menu Navigation**: ↑/↓ arrows, Enter to select
- **Paddle Movement**: W/A/S/D keys
- **Return to Menu**: ESC key

## Building

### Build Server
```bash
go build -o pong-server cmd/server/main.go
```

### Build Client
```bash
go build -o pong-client cmd/client/main.go
```

## Network Protocol

The game uses JSON messages over TCP for communication:

### Client → Server
```json
{
  "type": "input",
  "playerId": 1,
  "paddle1Y": 120,
  "paddle2X": 480
}
```

### Server → Client
```json
{
  "type": "state",
  "balls": [
    {"x": 200, "y": 300, "dx": 3, "dy": -2}
  ],
  "paddles": [
    {"playerId": 1, "paddle1Y": 120, "paddle2X": 480}
  ],
  "scores": {"Player1": 5, "Player2": 3},
  "gameOver": false
}
```

## Configuration

Game settings can be modified in `internal/game/state.go`:

- Field size: 600×600 pixels
- Ball count: 2 balls
- Target score: 10 points
- Time limit: 5 minutes
- Paddle speed: 5.0
- Ball speed: 3.0

## Project Structure

```
network-pong-battle/
├── cmd/
│   ├── server/
│   │   └── main.go         # Server entry point
│   └── client/
│       └── main.go         # Client entry point
├── internal/
│   ├── game/
│   │   ├── game.go         # Core game logic
│   │   ├── paddle.go       # Paddle implementation
│   │   ├── ball.go         # Ball implementation
│   │   └── state.go        # Game state management
│   ├── net/
│   │   ├── server.go       # Server networking
│   │   ├── client.go       # Client networking
│   │   └── protocol.go     # Message protocol
│   └── ui/
│       ├── renderer.go     # Graphics rendering
│       └── input.go        # Input handling
├── assets/                 # Game assets (images, sounds)
├── go.mod                  # Go module file
└── README.md               # This file
```

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Linting
```bash
go vet ./...
```

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the server port using the `-port` flag
2. **Connection refused**: Ensure the server is running before starting clients
3. **Game not starting**: Wait for two players to connect

### Debug Mode

Enable debug logging by setting the log level in the server and client code.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- Built with [Ebiten](https://ebiten.org/) - A dead simple 2D game library for Go
- Inspired by classic Pong games
- Designed for learning Go networking and game development concepts
