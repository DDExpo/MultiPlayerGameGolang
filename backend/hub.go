package main

import (
	"log"
	"sync"
	"time"
)

type Hub struct {
	clients         map[*Client]bool
	register        chan *Client
	unregister      chan *Client
	broadcast       chan []byte
	players         map[string]*Player
	activeUsernames map[string]bool
	mu              sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		broadcast:       make(chan []byte, 1024),
		players:         make(map[string]*Player),
		activeUsernames: make(map[string]bool),
	}
}

func (h *Hub) RunGameLoop() {
	ticker := time.NewTicker(TickDuration)
	defer ticker.Stop()

	for range ticker.C {
		for _, p := range h.players {

			mu.Lock()
			aliveProjectiles := Projectiles[:0]
			now := time.Now()
			for _, prj := range Projectiles {
				alive, hit := simulateProjectile(prj, now)
				if !alive {
					if len(hit) > 0 {
						log.Println(hit)
						htarget, ok1 := h.players[hit]
						attacker, ok2 := h.players[prj.OwnerId]
						if ok1 && ok2 {
							ApplyDamage(htarget, attacker)
						}
					}
					msg := SerializeUserShootStatus(alive, prj.ProjectileId)
					h.broadcast <- msg
					continue
				}
				aliveProjectiles = append(aliveProjectiles, prj)
			}
			Projectiles = aliveProjectiles
			mu.Unlock()

			applyInput(p)
			deltaMask := computeDeltaMask(p)
			if deltaMask != 0 {
				if p.Combat.HP <= 0 {
					b := make([]byte, 0, 1)
					b = append(b, MsgTypeUserDead)
					h.broadcast <- b
				}
				msg := SerializeUserStateDelta(MsgTypeUserState, p, deltaMask)
				h.broadcast <- msg
				updateLastSent(p, deltaMask)
			}
		}
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
				delete(h.activeUsernames, client.player.Meta.Username)
				h.mu.Unlock()

				FastProjectileCheck.mu.Lock()
				FastProjectileCheck.Remove(client.player.Meta.Username)
				FastProjectileCheck.mu.Unlock()

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
