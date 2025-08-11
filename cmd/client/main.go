package main

import (
	"flag"
	"log"
	"network-pong-battle/internal/game"
	"network-pong-battle/internal/net"
	"network-pong-battle/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Parse command line flags
	serverAddr := flag.String("server", "localhost:8080", "Server address to connect to")
	playerName := flag.String("name", "Player", "Player name")
	flag.Parse()

	log.Println("Starting Network Pong Battle Client...")
	log.Printf("Connecting to server: %s", *serverAddr)

	// Create renderer
	renderer := ui.NewRenderer(600)

	// Create input handler
	inputHandler := ui.NewInputHandler(renderer)

	// Create network client
	client := net.NewClient(*serverAddr, *playerName)

	// Set up callbacks
	client.SetCallbacks(
		// onStateUpdate
		func(state game.GameState) {
			renderer.SetGameState(state)
		},
		// onGameStart
		func(settings game.GameSettings) {
			renderer.SetGameStarted(true)
			renderer.SetShowMenu(false)
		},
		// onGameEnd
		func(winner int, scores game.Scores, gameTime int64) {
			renderer.SetGameOver(true, winner)
		},
		// onJoin
		func(playerID int, playerName string) {
			renderer.SetPlayerID(playerID)
			renderer.SetShowMenu(false)
			// Don't set gameStarted yet - wait for start message
		},
	)

	// Connect to server
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	// Set client in input handler
	inputHandler.SetClient(client)

	// Set up Ebiten game
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("Network Pong Battle - Client")
	ebiten.SetWindowResizable(true)

	// Create game struct that implements ebiten.Game
	game := &Game{
		renderer:     renderer,
		inputHandler: inputHandler,
		client:       client,
	}

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game error: %v", err)
	}
}

// Game implements ebiten.Game interface
type Game struct {
	renderer     *ui.Renderer
	inputHandler *ui.InputHandler
	client       *net.GameClient
}

// Update updates the game state
func (g *Game) Update() error {
	g.inputHandler.Update()
	return nil
}

// Draw draws the game
func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen)
}

// Layout returns the logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.renderer.Layout(outsideWidth, outsideHeight)
}
