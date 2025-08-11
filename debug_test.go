package main

import (
	"fmt"
	"network-pong-battle/internal/game"
	"time"
)

func main() {
	fmt.Println("Testing Network Pong Battle Game Logic...")

	// Create a new game
	g := game.NewGame()
	fmt.Printf("Game created: running=%v\n", g.IsRunning())

	// Start the game
	g.Start()
	fmt.Printf("Game started: running=%v\n", g.IsRunning())

	// Get initial state
	state := g.GetState()
	fmt.Printf("Initial state: balls=%d, paddles=%d, scores=%+v\n",
		len(state.Balls), len(state.Paddles), state.Scores)

	// Update game a few times
	for i := 0; i < 5; i++ {
		g.Update()
		time.Sleep(100 * time.Millisecond)
		state = g.GetState()
		fmt.Printf("Update %d: ball1 pos=(%.1f, %.1f), scores=%+v\n",
			i+1, state.Balls[0].X, state.Balls[0].Y, state.Scores)
	}

	fmt.Println("Test completed successfully!")
}
