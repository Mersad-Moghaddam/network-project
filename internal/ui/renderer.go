package ui

import (
	"fmt"
	"image/color"
	"network-pong-battle/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// Renderer handles the game graphics rendering
type Renderer struct {
	gameState    game.GameState
	fieldSize    int
	scale        float64
	playerID     int
	gameStarted  bool
	gameOver     bool
	winner       int
	
	// UI state
	showMenu     bool
	menuOption   int
	menuOptions  []string
	
	// Assets
	font         font.Face
	colors       map[string]color.Color
}

// NewRenderer creates a new game renderer
func NewRenderer(fieldSize int) *Renderer {
	return &Renderer{
		fieldSize:   fieldSize,
		scale:       1.0,
		showMenu:    true,
		menuOption:  0,
		menuOptions: []string{"Connect to Server", "Host Server", "Settings", "Exit"},
		font:        basicfont.Face7x13,
		colors: map[string]color.Color{
			"background": color.RGBA{20, 20, 40, 255},
			"field":      color.RGBA{30, 30, 60, 255},
			"paddle1":    color.RGBA{100, 200, 100, 255},
			"paddle2":    color.RGBA{200, 100, 100, 255},
			"ball":       color.RGBA{255, 255, 255, 255},
			"text":       color.RGBA{255, 255, 255, 255},
			"score":      color.RGBA{255, 255, 0, 255},
			"menu":       color.RGBA{100, 150, 255, 255},
		},
	}
}

// Update updates the renderer state
func (r *Renderer) Update() error {
	return nil
}

// Draw draws the game graphics
func (r *Renderer) Draw(screen *ebiten.Image) {
	if r.showMenu {
		r.drawMenu(screen)
	} else if r.gameOver {
		r.drawGameOver(screen)
	} else if r.gameStarted {
		r.drawGame(screen)
	} else {
		r.drawWaiting(screen)
	}
}

// Layout returns the logical screen size
func (r *Renderer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return r.fieldSize, r.fieldSize
}

// SetGameState updates the game state for rendering
func (r *Renderer) SetGameState(state game.GameState) {
	r.gameState = state
}

// SetPlayerID sets the current player ID
func (r *Renderer) SetPlayerID(playerID int) {
	r.playerID = playerID
}

// SetGameStarted sets whether the game has started
func (r *Renderer) SetGameStarted(started bool) {
	r.gameStarted = started
}

// SetGameOver sets whether the game is over
func (r *Renderer) SetGameOver(over bool, winner int) {
	r.gameOver = over
	r.winner = winner
}

// SetShowMenu sets whether to show the menu
func (r *Renderer) SetShowMenu(show bool) {
	r.showMenu = show
}

// SetMenuOption sets the selected menu option
func (r *Renderer) SetMenuOption(option int) {
	if option >= 0 && option < len(r.menuOptions) {
		r.menuOption = option
	}
}

// GetMenuOption returns the current menu option
func (r *Renderer) GetMenuOption() int {
	return r.menuOption
}

// drawMenu draws the main menu
func (r *Renderer) drawMenu(screen *ebiten.Image) {
	// Draw background
	screen.Fill(r.colors["background"])
	
	// Draw title
	title := "Network Pong Battle"
	titleBounds := text.BoundString(r.font, title)
	titleX := (r.fieldSize - titleBounds.Dx()) / 2
	titleY := r.fieldSize / 3
	text.Draw(screen, title, r.font, titleX, titleY, r.colors["text"])
	
	// Draw menu options
	optionY := r.fieldSize / 2
	for i, option := range r.menuOptions {
		color := r.colors["text"]
		if i == r.menuOption {
			color = r.colors["menu"]
		}
		
		bounds := text.BoundString(r.font, option)
		x := (r.fieldSize - bounds.Dx()) / 2
		y := optionY + i*30
		text.Draw(screen, option, r.font, x, y, color)
	}
	
	// Draw instructions
	instructions := "Use ↑↓ to navigate, Enter to select"
	instBounds := text.BoundString(r.font, instructions)
	instX := (r.fieldSize - instBounds.Dx()) / 2
	instY := r.fieldSize - 50
	text.Draw(screen, instructions, r.font, instX, instY, r.colors["text"])
}

// drawGame draws the actual game
func (r *Renderer) drawGame(screen *ebiten.Image) {
	// Draw background
	screen.Fill(r.colors["background"])
	
	// Draw field border
	fieldRect := ebiten.NewImage(r.fieldSize, r.fieldSize)
	fieldRect.Fill(r.colors["field"])
	
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(fieldRect, op)
	
	// Draw paddles
	for _, paddle := range r.gameState.Paddles {
		r.drawPaddle(screen, paddle)
	}
	
	// Draw balls
	for _, ball := range r.gameState.Balls {
		r.drawBall(screen, ball)
	}
	
	// Draw scores
	r.drawScores(screen)
	
	// Draw player info
	r.drawPlayerInfo(screen)
}

// drawPaddle draws a paddle
func (r *Renderer) drawPaddle(screen *ebiten.Image, paddle game.Paddle) {
	paddleImg := ebiten.NewImage(int(paddle.Width), int(paddle.Height))
	
	// Choose color based on player
	var paddleColor color.Color
	if paddle.PlayerID == 1 {
		paddleColor = r.colors["paddle1"]
	} else {
		paddleColor = r.colors["paddle2"]
	}
	
	paddleImg.Fill(paddleColor)
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(paddle.X, paddle.Y)
	screen.DrawImage(paddleImg, op)
}

// drawBall draws a ball
func (r *Renderer) drawBall(screen *ebiten.Image, ball game.Ball) {
	ballImg := ebiten.NewImage(int(ball.Radius*2), int(ball.Radius*2))
	ballImg.Fill(r.colors["ball"])
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ball.X-ball.Radius, ball.Y-ball.Radius)
	screen.DrawImage(ballImg, op)
}

// drawScores draws the score display
func (r *Renderer) drawScores(screen *ebiten.Image) {
	scoreText := fmt.Sprintf("P1: %d  P2: %d", r.gameState.Scores.Player1, r.gameState.Scores.Player2)
	text.Draw(screen, scoreText, r.font, 10, 30, r.colors["text"])
}

// drawPlayerInfo draws player information
func (r *Renderer) drawPlayerInfo(screen *ebiten.Image) {
	playerText := fmt.Sprintf("You are Player %d", r.playerID)
	text.Draw(screen, playerText, r.font, 10, r.fieldSize-20, r.colors["text"])
}

// drawWaiting draws the waiting screen
func (r *Renderer) drawWaiting(screen *ebiten.Image) {
	screen.Fill(r.colors["background"])
	
	waitingText := "Waiting for players..."
	waitingBounds := text.BoundString(r.font, waitingText)
	waitingX := (r.fieldSize - waitingBounds.Dx()) / 2
	waitingY := r.fieldSize / 2
	text.Draw(screen, waitingText, r.font, waitingX, waitingY, r.colors["text"])
	
	playerText := fmt.Sprintf("Connected as Player %d", r.playerID)
	playerBounds := text.BoundString(r.font, playerText)
	playerX := (r.fieldSize - playerBounds.Dx()) / 2
	playerY := waitingY + 30
	text.Draw(screen, playerText, r.font, playerX, playerY, r.colors["text"])
}

// drawGameOver draws the game over screen
func (r *Renderer) drawGameOver(screen *ebiten.Image) {
	screen.Fill(r.colors["background"])
	
	// Draw game over text
	gameOverText := "Game Over!"
	gameOverBounds := text.BoundString(r.font, gameOverText)
	gameOverX := (r.fieldSize - gameOverBounds.Dx()) / 2
	gameOverY := r.fieldSize / 3
	text.Draw(screen, gameOverText, r.font, gameOverX, gameOverY, r.colors["text"])
	
	// Draw winner
	var winnerText string
	if r.winner == 0 {
		winnerText = "It's a tie!"
	} else {
		winnerText = fmt.Sprintf("Player %d wins!", r.winner)
	}
	winnerBounds := text.BoundString(r.font, winnerText)
	winnerX := (r.fieldSize - winnerBounds.Dx()) / 2
	winnerY := gameOverY + 40
	text.Draw(screen, winnerText, r.font, winnerX, winnerY, r.colors["score"])
	
	// Draw final scores
	finalScoreText := fmt.Sprintf("Final Score - P1: %d, P2: %d", r.gameState.Scores.Player1, r.gameState.Scores.Player2)
	scoreBounds := text.BoundString(r.font, finalScoreText)
	scoreX := (r.fieldSize - scoreBounds.Dx()) / 2
	scoreY := winnerY + 40
	text.Draw(screen, finalScoreText, r.font, scoreX, scoreY, r.colors["text"])
	
	// Draw return to menu instruction
	menuText := "Press ESC to return to menu"
	menuBounds := text.BoundString(r.font, menuText)
	menuX := (r.fieldSize - menuBounds.Dx()) / 2
	menuY := r.fieldSize - 50
	text.Draw(screen, menuText, r.font, menuX, menuY, r.colors["text"])
}
