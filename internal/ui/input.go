package ui

import (
	"network-pong-battle/internal/net"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputHandler handles keyboard input for the game
type InputHandler struct {
	renderer *Renderer
	client   interface{} // Will be *net.GameClient when connected

	// Input state
	keysPressed map[ebiten.Key]bool
	lastInput   map[ebiten.Key]bool

	// Paddle movement
	paddle1Y float64
	paddle2X float64

	// Menu navigation
	menuOption int
}

// NewInputHandler creates a new input handler
func NewInputHandler(renderer *Renderer) *InputHandler {
	return &InputHandler{
		renderer:    renderer,
		keysPressed: make(map[ebiten.Key]bool),
		lastInput:   make(map[ebiten.Key]bool),
		paddle1Y:    300,
		paddle2X:    300,
		menuOption:  0,
	}
}

// SetClient sets the network client for sending input
func (ih *InputHandler) SetClient(client interface{}) {
	ih.client = client
}

// Update updates the input handler state
func (ih *InputHandler) Update() {
	// Update key states
	for key := range ih.keysPressed {
		ih.lastInput[key] = ih.keysPressed[key]
		ih.keysPressed[key] = ebiten.IsKeyPressed(key)
	}

	// Handle menu input
	if ih.renderer.showMenu {
		ih.handleMenuInput()
	} else if ih.renderer.gameStarted && !ih.renderer.gameOver {
		ih.handleGameInput()
	} else if ih.renderer.gameOver {
		ih.handleGameOverInput()
	}
}

// handleMenuInput handles input when in the menu
func (ih *InputHandler) handleMenuInput() {
	// Menu navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		ih.menuOption = (ih.menuOption - 1 + len(ih.renderer.menuOptions)) % len(ih.renderer.menuOptions)
		ih.renderer.SetMenuOption(ih.menuOption)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		ih.menuOption = (ih.menuOption + 1) % len(ih.renderer.menuOptions)
		ih.renderer.SetMenuOption(ih.menuOption)
	}

	// Menu selection
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ih.handleMenuSelection()
	}
}

// handleGameInput handles input during gameplay
func (ih *InputHandler) handleGameInput() {
	var moved bool

	// Paddle 1 movement (vertical - left side)
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		ih.paddle1Y -= 5
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		ih.paddle1Y += 5
		moved = true
	}

	// Paddle 2 movement (horizontal - top side)
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		ih.paddle2X -= 5
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		ih.paddle2X += 5
		moved = true
	}

	// Constrain paddle positions
	if ih.paddle1Y < 0 {
		ih.paddle1Y = 0
	}
	if ih.paddle1Y > 500 { // 600 - 100 (paddle height)
		ih.paddle1Y = 500
	}
	if ih.paddle2X < 0 {
		ih.paddle2X = 0
	}
	if ih.paddle2X > 500 { // 600 - 100 (paddle width)
		ih.paddle2X = 500
	}

	// Send input to server if connected and movement occurred
	if moved && ih.client != nil {
		if client, ok := ih.client.(*net.GameClient); ok {
			client.SendInput(ih.paddle1Y, ih.paddle2X)
		}
	}

	// Return to menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ih.renderer.SetShowMenu(true)
	}
}

// handleGameOverInput handles input when the game is over
func (ih *InputHandler) handleGameOverInput() {
	// Return to menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ih.renderer.SetShowMenu(true)
	}
}

// handleMenuSelection handles menu option selection
func (ih *InputHandler) handleMenuSelection() {
	switch ih.menuOption {
	case 0: // Connect to Server
		ih.handleConnectToServer()
	case 1: // Host Server
		ih.handleHostServer()
	case 2: // Settings
		ih.handleSettings()
	case 3: // Exit
		ih.handleExit()
	}
}

// handleConnectToServer handles connecting to a server
func (ih *InputHandler) handleConnectToServer() {
	// This will be implemented when we create the connection UI
	ih.renderer.SetShowMenu(false)
	// TODO: Show connection dialog
}

// handleHostServer handles hosting a server
func (ih *InputHandler) handleHostServer() {
	// This will be implemented when we create the server hosting UI
	ih.renderer.SetShowMenu(false)
	// TODO: Start server and show waiting screen
}

// handleSettings handles the settings menu
func (ih *InputHandler) handleSettings() {
	// This will be implemented when we create the settings UI
	// TODO: Show settings dialog
}

// handleExit handles exiting the game
func (ih *InputHandler) handleExit() {
	// This will be implemented when we create the main game loop
	// TODO: Exit game
}

// GetPaddlePositions returns the current paddle positions
func (ih *InputHandler) GetPaddlePositions() (float64, float64) {
	return ih.paddle1Y, ih.paddle2X
}

// SetPaddlePositions sets the paddle positions
func (ih *InputHandler) SetPaddlePositions(paddle1Y, paddle2X float64) {
	ih.paddle1Y = paddle1Y
	ih.paddle2X = paddle2X
}

// IsKeyJustPressed checks if a key was just pressed
func (ih *InputHandler) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

// IsKeyPressed checks if a key is currently pressed
func (ih *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}
