package main

import (
	"encoding/binary"
	"log"
	"math"
	"multiplayerGame/game"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	player *game.Player
}

func NewClient(hub *Hub, conn *websocket.Conn, player *game.Player) *Client {

	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 1024),
		player: player,
	}
}
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}
			break
		}
		switch message[0] {
		case MsgTypeUserChat:
			c.handleChatMessage(message[1:])
		case StateTypeUserReg:
			c.handleUserRegistration()
		case StateTypeUserResume:
			c.handleUserResumeSession()
		case StateTypeUserInput:
			c.handleUserInput(message[1:])
		case StateTypeUserPressedShoot:
			c.handleUserPressedShoot()
		case MsgTypeUserResumedDeath:
			c.handleUserResumedDeath()
		default:
			log.Printf("unknown message type: %d", message[0])
		}
	}
}
func (c *Client) WritePump() {
	defer c.conn.Close()

	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("write error:", err)
			}
			return
		}
	}
}

func (c *Client) handleUserInput(data []byte) {

	offset := 0
	moveX := math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	moveY := math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	angle := math.Float32frombits(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	dash := data[offset] != 0

	c.hub.gameCmd <- UserInputCmd{c.player, game.PlayerInput{MoveX: moveX, MoveY: moveY, Angle: angle, Dash: dash}}
}
func (c *Client) handleUserRegistration() {
	c.hub.gameCmd <- UserRegistrationCmd{c.player}
}
func (c *Client) handleUserPressedShoot() {
	c.hub.gameCmd <- SpawnProjectileCmd{c.player}
}
func (c *Client) handleUserResumedDeath() {
	c.hub.gameCmd <- UserResumedDeathCmd{c.player}
}
func (c *Client) handleUserResumeSession() {
	c.hub.gameCmd <- UserResumeSession{c.player}
}
func (c *Client) handleChatMessage(data []byte) {
	text, color, err := DeserializeUserMsg(data)
	if err != nil {
		log.Printf("Problem with deserializing user text: %s color: %s err: %v", text, color, err)
	}
	c.hub.broadcast <- SerializeUserChat(c.player.Meta.Username, text, time.Now().UTC().Format("2006-01-02 15:04"), color)
}
