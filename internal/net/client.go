package net

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"network-pong-battle/internal/game"
	"sync"
	"time"
)

// GameClient represents a game client
type GameClient struct {
	conn       net.Conn
	serverAddr string
	playerID   int
	playerName string
	connected  bool
	mu         sync.RWMutex
	
	// Callbacks for handling server messages
	onStateUpdate func(game.GameState)
	onGameStart   func(game.GameSettings)
	onGameEnd     func(int, game.Scores, int64)
	onJoin        func(int, string)
	
	// Input channel
	inputChan chan *InputMessage
	stopChan  chan bool
}

// NewClient creates a new game client
func NewClient(serverAddr, playerName string) *GameClient {
	return &GameClient{
		serverAddr: serverAddr,
		playerName: playerName,
		connected:  false,
		inputChan:  make(chan *InputMessage, 100),
		stopChan:   make(chan bool),
	}
}

// Connect connects to the game server
func (c *GameClient) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	c.connected = true
	log.Printf("Connected to server at %s", c.serverAddr)

	// Start message handling
	go c.handleServerMessages()
	go c.inputHandler()

	return nil
}

// Disconnect disconnects from the server
func (c *GameClient) Disconnect() {
	c.mu.Lock()
	c.connected = false
	c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
	}

	// Signal stop
	select {
	case c.stopChan <- true:
	default:
	}
}

// IsConnected returns whether the client is connected
func (c *GameClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// GetPlayerID returns the player ID assigned by the server
func (c *GameClient) GetPlayerID() int {
	return c.playerID
}

// SetCallbacks sets the callback functions for handling server messages
func (c *GameClient) SetCallbacks(
	onStateUpdate func(game.GameState),
	onGameStart func(game.GameSettings),
	onGameEnd func(int, game.Scores, int64),
	onJoin func(int, string),
) {
	c.onStateUpdate = onStateUpdate
	c.onGameStart = onGameStart
	c.onGameEnd = onGameEnd
	c.onJoin = onJoin
}

// SendInput sends player input to the server
func (c *GameClient) SendInput(paddle1Y, paddle2X float64) {
	if !c.IsConnected() {
		return
	}

	msg := CreateInputMessage(c.playerID, paddle1Y, paddle2X)
	select {
	case c.inputChan <- msg:
	default:
		log.Println("Input channel full, dropping input")
	}
}

// handleServerMessages handles incoming messages from the server
func (c *GameClient) handleServerMessages() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() && c.IsConnected() {
		data := scanner.Bytes()
		if len(data) == 0 {
			continue
		}

		c.processMessage(data)
	}

	log.Println("Server connection closed")
}

// processMessage processes a single message from the server
func (c *GameClient) processMessage(data []byte) {
	// Try to determine message type first
	var baseMsg struct {
		Type MessageType `json:"type"`
	}
	
	if err := json.Unmarshal(data, &baseMsg); err != nil {
		log.Printf("Error parsing message type: %v", err)
		return
	}

	switch baseMsg.Type {
	case MessageTypeJoin:
		var msg JoinMessage
		if err := DecodeMessage(data, &msg); err != nil {
			log.Printf("Error decoding join message: %v", err)
			return
		}
		c.playerID = msg.PlayerID
		if c.onJoin != nil {
			c.onJoin(msg.PlayerID, msg.PlayerName)
		}
		log.Printf("Joined game as Player %d", c.playerID)

	case MessageTypeStart:
		var msg StartMessage
		if err := DecodeMessage(data, &msg); err != nil {
			log.Printf("Error decoding start message: %v", err)
			return
		}
		if c.onGameStart != nil {
			c.onGameStart(msg.Settings)
		}
		log.Println("Game started!")

	case MessageTypeState:
		var msg StateMessage
		if err := DecodeMessage(data, &msg); err != nil {
			log.Printf("Error decoding state message: %v", err)
			return
		}
		
		// Convert to game state
		state := game.GameState{
			Balls:    msg.Balls,
			Paddles:  msg.Paddles,
			GameOver: msg.GameOver,
			Winner:   msg.Winner,
		}
		
		if c.onStateUpdate != nil {
			c.onStateUpdate(state)
		}

	case MessageTypeEnd:
		var msg EndMessage
		if err := DecodeMessage(data, &msg); err != nil {
			log.Printf("Error decoding end message: %v", err)
			return
		}
		if c.onGameEnd != nil {
			c.onGameEnd(msg.Winner, msg.FinalScores, msg.GameTime)
		}
		log.Printf("Game ended! Winner: Player %d", msg.Winner)

	default:
		log.Printf("Unknown message type: %s", baseMsg.Type)
	}
}

// inputHandler sends input messages to the server
func (c *GameClient) inputHandler() {
	ticker := time.NewTicker(time.Second / 60) // 60 FPS
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			select {
			case msg := <-c.inputChan:
				if c.IsConnected() {
					data, err := EncodeMessage(msg)
					if err != nil {
						log.Printf("Error encoding input message: %v", err)
						continue
					}
					data = append(data, '\n')
					
					c.mu.Lock()
					if c.conn != nil {
						c.conn.Write(data)
					}
					c.mu.Unlock()
				}
			default:
				// No input to send
			}
		}
	}
}
