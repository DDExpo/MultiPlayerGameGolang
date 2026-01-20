package main

import (
	"encoding/binary"
	"log"
	"math"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	player *Player
}

func NewClient(hub *Hub, conn *websocket.Conn, player *Player) *Client {

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
		case MsgTypeUser:
			c.handleUserRegistration()
		case MsgTypeChat:
			c.handleChatMessage(message[1:])
		case MsgTypeInput:
			c.handleUserInput(message[1:])
		case MsgTypePressedShoot:
			c.handleUserPressedShoot(message[1:])
		case MsgTypeShootStatus:
			c.handleUserShootStatus()
		case MsgTypeResumeSession:
			c.handleUserResumeSession()
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

func (c *Client) handleUserRegistration() {
	msg := SerializeUserStateDelta(MsgTypeUserState, c.player, UserStateDeltaPOS|UserStateDeltaSTATS|UserStateDeltaWEAPON)
	c.hub.broadcast <- msg
}

func (c *Client) handleChatMessage(data []byte) {
	text, color, err := DeserializeUserMsg(data)
	if err != nil {
		log.Printf("Problem with deserializing user text: %s color: %s err: %v", text, color, err)
	}

	msg := SerializeUserMsg(c.player.Meta.Username, text, time.Now().UTC().Format("2006-01-02 15:04"), color)
	c.hub.broadcast <- msg
}

func (c *Client) handleUserInput(data []byte) {

	offset := 0
	moveX := int8(data[offset])
	offset++
	moveY := int8(data[offset])
	offset++
	angle := math.Float32frombits(binary.LittleEndian.Uint32(data[offset:]))
	offset += 4

	dash := data[offset] != 0
	c.player.Input = PlayerInput{
		MoveX: moveX,
		MoveY: moveY,
		Dash:  dash,
		Angle: angle,
	}
}

func (c *Client) handleUserPressedShoot(data []byte) {
	mu.Lock()
	Projectiles = append(Projectiles, CreateProjectile(c.player, binary.BigEndian.Uint32(data)))
	mu.Unlock()

	msg := SerializeUserStateDelta(MsgTypePressedShoot, c.player, UserStateDeltaPOS|UserStateDeltaSTATS|UserStateDeltaWEAPON)
	c.hub.broadcast <- msg
}

func (c *Client) handleUserShootStatus() {
}

func (c *Client) handleUserResumeSession() {
	msg := SerializeUserStateDelta(MsgTypeUserState, c.player, UserStateDeltaPOS|UserStateDeltaSTATS|UserStateDeltaWEAPON)
	c.hub.broadcast <- msg
}
