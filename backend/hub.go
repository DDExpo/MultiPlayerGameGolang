package main

import (
	"log"
	"multiplayerGame/game"
	"sync"
)

type Hub struct {
	clients         map[*Client]bool
	register        chan *Client
	unregister      chan *Client
	broadcast       chan []byte
	gameCmd         chan GameCommand
	players         map[string]*game.Player
	activeUsernames map[string]bool
	mu              sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		broadcast:       make(chan []byte, 1024),
		gameCmd:         make(chan GameCommand, 256),
		players:         make(map[string]*game.Player),
		activeUsernames: make(map[string]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				h.mu.Lock()
				delete(h.players, client.player.Meta.SessionID)
				h.activeUsernames[client.player.Meta.Username] = false
				h.mu.Unlock()

				game.FastProjectileCheck.Remove(client.player.Meta.Username)

				log.Printf("Client unregistered: %s", client.player.Meta.Username)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
