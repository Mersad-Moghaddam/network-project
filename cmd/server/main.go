package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"network-pong-battle/internal/net"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	log.Println("Starting Network Pong Battle Server...")
	log.Printf("Server will listen on port %s", *port)

	// Create and start server
	server := net.NewServer(*port)
	
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	log.Println("Server is running. Press Ctrl+C to stop.")
	
	// Wait for signal
	<-sigChan
	
	log.Println("Shutting down server...")
	server.Stop()
	log.Println("Server stopped.")
}
