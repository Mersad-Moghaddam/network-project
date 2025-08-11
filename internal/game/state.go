package game

import (
	"sync"
	"time"
)

// GameState represents the complete state of the game
type GameState struct {
	mu       sync.RWMutex
	Balls    []Ball
	Paddles  []Paddle
	Scores   Scores
	GameOver bool
	Winner   int
	StartTime time.Time
	EndTime   time.Time
	Settings  GameSettings
}

// GameSettings holds configurable game parameters
type GameSettings struct {
	FieldSize     int
	BallCount     int
	TargetScore   int
	TimeLimit     time.Duration
	PaddleSpeed  float64
	BallSpeed    float64
}

// Scores holds player scores
type Scores struct {
	Player1 int
	Player2 int
}

// NewGameState creates a new game state with default settings
func NewGameState() *GameState {
	return &GameState{
		Balls:    make([]Ball, 0),
		Paddles:  make([]Paddle, 0),
		Scores:   Scores{Player1: 0, Player2: 0},
		GameOver: false,
		Winner:   0,
		StartTime: time.Now(),
		Settings: GameSettings{
			FieldSize:    600,
			BallCount:    2,
			TargetScore:  10,
			TimeLimit:    5 * time.Minute,
			PaddleSpeed: 5.0,
			BallSpeed:   3.0,
		},
	}
}

// InitializeGame sets up the initial game state
func (gs *GameState) InitializeGame() {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	// Create paddles for both players
	gs.Paddles = []Paddle{
		// Player 1: Left and Top edges
		NewPaddle(1, 1, 0, 300, 20, 100),      // Left paddle
		NewPaddle(1, 2, 300, 0, 100, 20),      // Top paddle
		// Player 2: Right and Bottom edges
		NewPaddle(2, 1, 580, 300, 20, 100),    // Right paddle
		NewPaddle(2, 2, 300, 580, 100, 20),    // Bottom paddle
	}

	// Create balls
	gs.Balls = make([]Ball, gs.Settings.BallCount)
	for i := 0; i < gs.Settings.BallCount; i++ {
		gs.Balls[i] = NewBall(300, 300, gs.Settings.BallSpeed)
	}

	gs.Scores = Scores{Player1: 0, Player2: 0}
	gs.GameOver = false
	gs.Winner = 0
	gs.StartTime = time.Now()
}

// GetState returns a copy of the current game state for safe reading
func (gs *GameState) GetState() GameState {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	return GameState{
		Balls:    append([]Ball{}, gs.Balls...),
		Paddles:  append([]Paddle{}, gs.Paddles...),
		Scores:   gs.Scores,
		GameOver: gs.GameOver,
		Winner:   gs.Winner,
		StartTime: gs.StartTime,
		EndTime:   gs.EndTime,
		Settings:  gs.Settings,
	}
}

// UpdatePaddle updates a specific paddle position
func (gs *GameState) UpdatePaddle(playerID, paddleID int, x, y float64) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	for i := range gs.Paddles {
		if gs.Paddles[i].PlayerID == playerID && gs.Paddles[i].PaddleID == paddleID {
			gs.Paddles[i].X = x
			gs.Paddles[i].Y = y
			break
		}
	}
}

// CheckGameEnd checks if the game should end and updates the winner
func (gs *GameState) CheckGameEnd() bool {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	// Check target score
	if gs.Scores.Player1 >= gs.Settings.TargetScore {
		gs.GameOver = true
		gs.Winner = 1
		gs.EndTime = time.Now()
		return true
	}
	if gs.Scores.Player2 >= gs.Settings.TargetScore {
		gs.GameOver = true
		gs.Winner = 2
		gs.EndTime = time.Now()
		return true
	}

	// Check time limit
	if time.Since(gs.StartTime) >= gs.Settings.TimeLimit {
		gs.GameOver = true
		if gs.Scores.Player1 > gs.Scores.Player2 {
			gs.Winner = 1
		} else if gs.Scores.Player2 > gs.Scores.Player1 {
			gs.Winner = 2
		} else {
			gs.Winner = 0 // Tie
		}
		gs.EndTime = time.Now()
		return true
	}

	return false
}

// AddScore increments the score for a player
func (gs *GameState) AddScore(playerID int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if playerID == 1 {
		gs.Scores.Player1++
	} else if playerID == 2 {
		gs.Scores.Player2++
	}
}
