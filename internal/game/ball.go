package game

import (
	"math"
	"math/rand"
)

// Ball represents a ball in the game
type Ball struct {
	X, Y    float64 // Position
	DX, DY  float64 // Direction (velocity)
	Radius  float64
	Speed   float64
}

// NewBall creates a new ball with random direction
func NewBall(x, y, speed float64) Ball {
	// Random angle between 0 and 2π
	angle := rand.Float64() * 2 * math.Pi
	
	return Ball{
		X:      x,
		Y:      y,
		DX:     math.Cos(angle) * speed,
		DY:     math.Sin(angle) * speed,
		Radius: 8,
		Speed:  speed,
	}
}

// Update moves the ball based on its current velocity
func (b *Ball) Update(fieldSize int) {
	b.X += b.DX
	b.Y += b.DY
}

// CheckWallCollision checks and handles wall collisions
func (b *Ball) CheckWallCollision(fieldSize int) bool {
	collision := false
	
	// Left wall (Player 2 scores)
	if b.X-b.Radius <= 0 {
		b.X = b.Radius
		b.DX = -b.DX
		collision = true
		return true
	}
	
	// Right wall (Player 1 scores)
	if b.X+b.Radius >= float64(fieldSize) {
		b.X = float64(fieldSize) - b.Radius
		b.DX = -b.DX
		collision = true
		return true
	}
	
	// Top wall (Player 2 scores)
	if b.Y-b.Radius <= 0 {
		b.Y = b.Radius
		b.DY = -b.DY
		collision = true
		return true
	}
	
	// Bottom wall (Player 1 scores)
	if b.Y+b.Radius >= float64(fieldSize) {
		b.Y = float64(fieldSize) - b.Radius
		b.DY = -b.DY
		collision = true
		return true
	}
	
	return collision
}

// CheckPaddleCollision checks collision with a paddle
func (b *Ball) CheckPaddleCollision(paddle Paddle) bool {
	// Check if ball is within paddle bounds
	if b.X+b.Radius >= paddle.X && 
	   b.X-b.Radius <= paddle.X+paddle.Width &&
	   b.Y+b.Radius >= paddle.Y && 
	   b.Y-b.Radius <= paddle.Y+paddle.Height {
		
		// Determine which side of the paddle was hit
		// and adjust ball direction accordingly
		if b.X < paddle.X+paddle.Width/2 {
			// Hit left side of paddle
			b.DX = -math.Abs(b.DX)
		} else {
			// Hit right side of paddle
			b.DX = math.Abs(b.DX)
		}
		
		if b.Y < paddle.Y+paddle.Height/2 {
			// Hit top side of paddle
			b.DY = -math.Abs(b.DY)
		} else {
			// Hit bottom side of paddle
			b.DY = math.Abs(b.DY)
		}
		
		// Add some randomness to prevent infinite loops
		b.DX += (rand.Float64() - 0.5) * 0.5
		b.DY += (rand.Float64() - 0.5) * 0.5
		
		// Normalize speed
		speed := math.Sqrt(b.DX*b.DX + b.DY*b.DY)
		b.DX = (b.DX / speed) * b.Speed
		b.DY = (b.DY / speed) * b.Speed
		
		return true
	}
	
	return false
}

// Reset resets the ball to center with random direction
func (b *Ball) Reset(fieldSize int) {
	b.X = float64(fieldSize) / 2
	b.Y = float64(fieldSize) / 2
	
	// Random angle between 0 and 2π
	angle := rand.Float64() * 2 * math.Pi
	b.DX = math.Cos(angle) * b.Speed
	b.DY = math.Sin(angle) * b.Speed
}
