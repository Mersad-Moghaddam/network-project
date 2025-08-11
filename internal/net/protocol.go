package net

import (
	"encoding/json"
	"network-pong-battle/internal/game"
)

// MessageType represents the type of network message
type MessageType string

const (
	MessageTypeInput MessageType = "input"
	MessageTypeState MessageType = "state"
	MessageTypeJoin  MessageType = "join"
	MessageTypeStart MessageType = "start"
	MessageTypeEnd   MessageType = "end"
)

// InputMessage represents player input sent from client to server
type InputMessage struct {
	Type      MessageType `json:"type"`
	PlayerID  int         `json:"playerId"`
	Paddle1Y  float64     `json:"paddle1Y,omitempty"`
	Paddle2X  float64     `json:"paddle2X,omitempty"`
}

// StateMessage represents the complete game state sent from server to clients
type StateMessage struct {
	Type      MessageType       `json:"type"`
	Balls     []game.Ball       `json:"balls"`
	Paddles   []game.Paddle     `json:"paddles"`
	GameOver  bool              `json:"gameOver"`
	Winner    int               `json:"winner"`
	GameTime  int64             `json:"gameTime"`  // in milliseconds
	Remaining int64              `json:"remaining"` // remaining time in milliseconds
}

// JoinMessage represents a player joining the game
type JoinMessage struct {
	Type      MessageType `json:"type"`
	PlayerID  int         `json:"playerId"`
	PlayerName string     `json:"playerName"`
}

// StartMessage represents the game starting
type StartMessage struct {
	Type      MessageType `json:"type"`
	Settings  game.GameSettings `json:"settings"`
}

// EndMessage represents the game ending
type EndMessage struct {
	Type      MessageType `json:"type"`
	Winner    int         `json:"winner"`
	FinalScores game.Scores `json:"finalScores"`
	GameTime  int64       `json:"gameTime"`
}

// EncodeMessage encodes a message to JSON bytes
func EncodeMessage(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

// DecodeMessage decodes JSON bytes to a message
func DecodeMessage(data []byte, msg interface{}) error {
	return json.Unmarshal(data, msg)
}

// CreateInputMessage creates an input message for a player
func CreateInputMessage(playerID int, paddle1Y, paddle2X float64) *InputMessage {
	return &InputMessage{
		Type:     MessageTypeInput,
		PlayerID: playerID,
		Paddle1Y: paddle1Y,
		Paddle2X: paddle2X,
	}
}

// CreateStateMessage creates a state message from game state
func CreateStateMessage(state game.GameState) *StateMessage {
	return &StateMessage{
		Type:     MessageTypeState,
		Balls:    state.Balls,
		Paddles:  state.Paddles,
		GameOver: state.GameOver,
		Winner:   state.Winner,
		GameTime: int64(state.StartTime.UnixMilli()),
		Remaining: int64(state.Settings.TimeLimit.Milliseconds()),
	}
}

// CreateJoinMessage creates a join message
func CreateJoinMessage(playerID int, playerName string) *JoinMessage {
	return &JoinMessage{
		Type:      MessageTypeJoin,
		PlayerID:  playerID,
		PlayerName: playerName,
	}
}

// CreateStartMessage creates a start message
func CreateStartMessage(settings game.GameSettings) *StartMessage {
	return &StartMessage{
		Type:     MessageTypeStart,
		Settings: settings,
	}
}

// CreateEndMessage creates an end message
func CreateEndMessage(winner int, scores game.Scores, gameTime int64) *EndMessage {
	return &EndMessage{
		Type:        MessageTypeEnd,
		Winner:      winner,
		FinalScores: scores,
		GameTime:    gameTime,
	}
}
