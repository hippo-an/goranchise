package main

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)

type Player struct {
	Id string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

func NewPlayer(id string) *Player {
	return &Player{
		Id: id,
	}
}

type GameState struct {
	mu      sync.Mutex
	Players map[string]*Player
}

func NewGameState() *GameState {
	return &GameState{
		Players: make(map[string]*Player),
	}
}

func (g *GameState) UpdatePlayer(id string, x, y int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if player, exist := g.Players[id]; exist {
		player.X = x
		player.Y = y
	} else {
		g.Players[id] = &Player{
			Id: id,
			X:  x,
			Y:  y,
		}
	}
}

func (g *GameState) RemovePlayer(id string) {
	g.mu.Lock()
	delete(g.Players, id)
	g.mu.Unlock()
}

func (g *GameState) Broadcast(connMap map[string]net.Conn) {
	g.mu.Lock()
	defer g.mu.Unlock()

	state, err := json.Marshal(g.Players)
	if err != nil {
		log.Println("ERror serializing game state;", err)
		return
	}

	for _, conn := range connMap {
		conn.Write(state)
	}
}
