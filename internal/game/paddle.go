package game

// Paddle represents a paddle in the game
type Paddle struct {
	PlayerID int     // Which player owns this paddle (1 or 2)
	PaddleID int     // Which paddle for this player (1 or 2)
	X, Y     float64 // Position
	Width    float64 // Width of the paddle
	Height   float64 // Height of the paddle
	Speed    float64 // Movement speed
}

// NewPaddle creates a new paddle
func NewPaddle(playerID, paddleID int, x, y, width, height float64) Paddle {
	return Paddle{
		PlayerID: playerID,
		PaddleID: paddleID,
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		Speed:    5.0,
	}
}

// Move moves the paddle in the specified direction
func (p *Paddle) Move(dx, dy float64, fieldSize int) {
	newX := p.X + dx*p.Speed
	newY := p.Y + dy*p.Speed
	
	// Constrain paddle movement based on its position
	switch p.PaddleID {
	case 1: // Left paddle (vertical movement only)
		if newY >= 0 && newY+p.Height <= float64(fieldSize) {
			p.Y = newY
		}
	case 2: // Top paddle (horizontal movement only)
		if newX >= 0 && newX+p.Width <= float64(fieldSize) {
			p.X = newX
		}
	case 3: // Right paddle (vertical movement only)
		if newY >= 0 && newY+p.Height <= float64(fieldSize) {
			p.Y = newY
		}
	case 4: // Bottom paddle (horizontal movement only)
		if newX >= 0 && newX+p.Width <= float64(fieldSize) {
			p.X = newX
		}
	}
}

// GetCenter returns the center point of the paddle
func (p *Paddle) GetCenter() (float64, float64) {
	return p.X + p.Width/2, p.Y + p.Height/2
}

// Contains checks if a point is within the paddle bounds
func (p *Paddle) Contains(x, y float64) bool {
	return x >= p.X && x <= p.X+p.Width && y >= p.Y && y <= p.Y+p.Height
}
