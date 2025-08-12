package net

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"network-pong-battle/internal/game"
	"sync"
	"time"
)

// Client represents a connected client
type Client struct {
	conn       net.Conn
	playerID   int
	playerName string
	mu         sync.Mutex
}

// Server represents the game server
type Server struct {
	listener    net.Listener
	clients     map[int]*Client
	game        *game.Game
	mu          sync.RWMutex
	port        string
	running     bool
	gameStarted bool
	clientCount int
}

// NewServer creates a new game server
func NewServer(port string) *Server {
	return &Server{
		clients:     make(map[int]*Client),
		game:        game.NewGame(),
		port:        port,
		running:     false,
		gameStarted: false,
		clientCount: 0,
	}
}

// Start starts the server
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	s.running = true
	log.Printf("Server started on port %s", s.port)

	// Start accepting clients
	go s.acceptClients()

	// Start game loop
	go s.gameLoop()

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}

	// Close all client connections
	s.mu.Lock()
	for _, client := range s.clients {
		client.conn.Close()
	}
	s.clients = make(map[int]*Client)
	s.mu.Unlock()

	log.Println("Server stopped")
}

// acceptClients accepts incoming client connections
func (s *Server) acceptClients() {
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.running {
				log.Printf("Error accepting connection: %v", err)
			}
			continue
		}

		// Check if we can accept more clients
		s.mu.Lock()
		if len(s.clients) >= 2 {
			conn.Close()
			log.Println("Game is full, rejecting connection")
			s.mu.Unlock()
			continue
		}
		s.mu.Unlock()

		// Handle new client
		go s.handleClient(conn)
	}
}

// handleClient handles a single client connection
func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	// Assign player ID
	s.mu.Lock()
	s.clientCount++
	playerID := s.clientCount
	client := &Client{
		conn:       conn,
		playerID:   playerID,
		playerName: fmt.Sprintf("Player %d", playerID),
	}
	s.clients[playerID] = client

	// Check if we have 2 players and can start the game
	shouldStartGame := len(s.clients) == 2 && !s.gameStarted
	s.mu.Unlock()

	log.Printf("Client %d connected from %s", playerID, conn.RemoteAddr())

	// Send join confirmation
	joinMsg := CreateJoinMessage(playerID, client.playerName)
	if data, err := EncodeMessage(joinMsg); err == nil {
		conn.Write(data)
		conn.Write([]byte("\n"))
	}

	// Start game if we have both players
	if shouldStartGame {
		log.Printf("Two players connected, starting game...")
		s.startGame()
	}

	// Handle client messages
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() && s.running {
		data := scanner.Bytes()
		if len(data) == 0 {
			continue
		}

		s.handleMessage(playerID, data)
	}

	// Client disconnected
	s.mu.Lock()
	delete(s.clients, playerID)
	log.Printf("Client %d disconnected", playerID)

	// Stop game if not enough players
	if len(s.clients) < 2 && s.gameStarted {
		s.stopGame()
	}
	s.mu.Unlock()
}

// handleMessage processes a message from a client
func (s *Server) handleMessage(playerID int, data []byte) {
	var msg InputMessage
	if err := DecodeMessage(data, &msg); err != nil {
		log.Printf("Error decoding message from client %d: %v", playerID, err)
		return
	}

	if msg.Type == MessageTypeInput {
		// Update paddle positions
		s.game.UpdatePaddle(playerID, 1, 0, msg.Paddle1Y)
		s.game.UpdatePaddle(playerID, 2, msg.Paddle2X, 0)
	}
}

// startGame starts the game
func (s *Server) startGame() {
	s.mu.Lock()
	if s.gameStarted {
		s.mu.Unlock()
		return // Game already started
	}

	log.Println("Starting game with 2 players...")
	s.gameStarted = true
	s.game.Start()
	s.mu.Unlock()

	// Send start message to all clients
	startMsg := CreateStartMessage(s.game.GetState().Settings)
	log.Printf("Broadcasting start message to %d clients", len(s.clients))
	s.broadcastMessage(startMsg)

	log.Println("Game started!")
}

// stopGame stops the game
func (s *Server) stopGame() {
	s.gameStarted = false
	s.game.Stop()

	// Send end message to all clients
	endMsg := CreateEndMessage(0, game.Scores{}, 0)
	s.broadcastMessage(endMsg)

	log.Println("Game stopped")
}

// gameLoop runs the main game loop
func (s *Server) gameLoop() {
	ticker := time.NewTicker(time.Second / 60) // 60 FPS
	defer ticker.Stop()

	for s.running {
		select {
		case <-ticker.C:
			if s.gameStarted {
				s.game.Update()

				// Broadcast game state to all clients
				stateMsg := CreateStateMessage(s.game.GetState())
				s.broadcastMessage(stateMsg)

				// Check if game ended
				if s.game.IsGameOver() {
					s.gameStarted = false
					endMsg := CreateEndMessage(
						s.game.GetWinner(),
						s.game.GetScore(),
						int64(s.game.GetGameTime().Milliseconds()),
					)
					s.broadcastMessage(endMsg)
					log.Printf("Game ended! Winner: Player %d", s.game.GetWinner())
				}
			}
		}
	}
}

// broadcastMessage sends a message to all connected clients
func (s *Server) broadcastMessage(msg interface{}) {
	data, err := EncodeMessage(msg)
	if err != nil {
		log.Printf("Error encoding message: %v", err)
		return
	}
	data = append(data, '\n')

	s.mu.RLock()
	for _, client := range s.clients {
		client.mu.Lock()
		client.conn.Write(data)
		client.mu.Unlock()
	}
	s.mu.RUnlock()
}

// GetClientCount returns the number of connected clients
func (s *Server) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

// IsGameStarted returns whether the game has started
func (s *Server) IsGameStarted() bool {
	return s.gameStarted
}
