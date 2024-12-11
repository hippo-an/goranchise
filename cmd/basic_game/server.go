package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	clients map[string]net.Conn
	game    *GameState
	mu      sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]net.Conn),
		game:    NewGameState(),
	}
}

func (s *Server) Start(port uint64) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	log.Printf("Game server started on port %d\n", port)

	go s.gameLoop()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error: ", err)
			continue
		}

		go s.handleConnection(conn)

	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	s.mu.Lock()
	s.clients[clientAddr] = conn
	s.mu.Unlock()

	log.Println("Player connected:", clientAddr)

	defer func() {
		log.Println("Player disconnected:", clientAddr)
		s.mu.Lock()
		delete(s.clients, clientAddr)
		s.mu.Unlock()
		s.game.RemovePlayer(clientAddr)
	}()

	p := NewPlayer(clientAddr)
	s.game.UpdatePlayer(p.Id, p.X, p.Y)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		msg := scanner.Bytes()

		var playerUpdate Player
		if err := json.Unmarshal(msg, &playerUpdate); err != nil {
			log.Println("Invalid player update;", err)
			continue
		}

		s.game.UpdatePlayer(clientAddr, playerUpdate.X, playerUpdate.Y)
	}

}

func (s *Server) gameLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		s.game.Broadcast(s.clients)
		s.mu.Unlock()
	}
}
