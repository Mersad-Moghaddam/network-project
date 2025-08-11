package game

import (
	"math/rand"
	"time"
)

// Game represents the main game controller
type Game struct {
	state     *GameState
	tickRate  time.Duration
	lastTick  time.Time
	running   bool
}

// NewGame creates a new game instance
func NewGame() *Game {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	return &Game{
		state:    NewGameState(),
		tickRate: time.Second / 60, // 60 FPS
		running:  false,
	}
}

// Start starts the game
func (g *Game) Start() {
	g.state.InitializeGame()
	g.running = true
	g.lastTick = time.Now()
}

// Stop stops the game
func (g *Game) Stop() {
	g.running = false
}

// IsRunning returns whether the game is currently running
func (g *Game) IsRunning() bool {
	return g.running
}

// GetState returns the current game state
func (g *Game) GetState() GameState {
	return g.state.GetState()
}

// UpdatePaddle updates a paddle position
func (g *Game) UpdatePaddle(playerID, paddleID int, x, y float64) {
	g.state.UpdatePaddle(playerID, paddleID, x, y)
}

// Update performs one game tick update
func (g *Game) Update() {
	if !g.running {
		return
	}

	now := time.Now()
	if now.Sub(g.lastTick) < g.tickRate {
		return
	}

	g.lastTick = now

	// Update ball positions
	g.updateBalls()

	// Check collisions
	g.checkCollisions()

	// Check if game should end
	if g.state.CheckGameEnd() {
		g.running = false
	}
}

// updateBalls updates all ball positions and checks wall collisions
func (g *Game) updateBalls() {
	state := g.state.GetState()
	
	for i := range g.state.Balls {
		ball := &g.state.Balls[i]
		ball.Update(state.Settings.FieldSize)
		
		// Check wall collisions and handle scoring
		if ball.CheckWallCollision(state.Settings.FieldSize) {
			// Determine which player scores based on which wall was hit
			if ball.X <= ball.Radius || ball.Y <= ball.Radius {
				// Left or top wall - Player 2 scores
				g.state.AddScore(2)
			} else {
				// Right or bottom wall - Player 1 scores
				g.state.AddScore(1)
			}
			
			// Reset ball to center
			ball.Reset(state.Settings.FieldSize)
		}
	}
}

// checkCollisions checks for ball-paddle collisions
func (g *Game) checkCollisions() {
	state := g.state.GetState()
	
	for i := range g.state.Balls {
		ball := &g.state.Balls[i]
		
		for _, paddle := range state.Paddles {
			if ball.CheckPaddleCollision(paddle) {
				break // Ball can only hit one paddle at a time
			}
		}
	}
}

// GetScore returns the current scores
func (g *Game) GetScore() Scores {
	state := g.state.GetState()
	return state.Scores
}

// IsGameOver returns whether the game is over
func (g *Game) IsGameOver() bool {
	state := g.state.GetState()
	return state.GameOver
}

// GetWinner returns the winner (0 for tie, 1 or 2 for players)
func (g *Game) GetWinner() int {
	state := g.state.GetState()
	return state.Winner
}

// GetGameTime returns the elapsed game time
func (g *Game) GetGameTime() time.Duration {
	state := g.state.GetState()
	if state.GameOver {
		return state.EndTime.Sub(state.StartTime)
	}
	return time.Since(state.StartTime)
}

// GetRemainingTime returns the remaining time if there's a time limit
func (g *Game) GetRemainingTime() time.Duration {
	state := g.state.GetState()
	elapsed := time.Since(state.StartTime)
	remaining := state.Settings.TimeLimit - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}
