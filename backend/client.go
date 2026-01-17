package main

import (
	"encoding/binary"
	"log"
	"math"

	"github.com/gorilla/websocket"
)

const (
	MsgTypeUser          byte = 1
	MsgTypeChat          byte = 2
	MsgTypeUserState     byte = 3
	MsgTypeResumeSession byte = 4
	MsgTypeInput         byte = 5
)

type Client struct {
	sessionID string
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	player    *UserState
}

func NewClient(hub *Hub, conn *websocket.Conn, sessionid string) *Client {
	player, exists := hub.players[sessionid]
	if !exists {
		newPlayer := LocalUserData
		newPlayer.CurrentSessionID = sessionid
		player = &newPlayer
		hub.players[sessionid] = player
	}

	return &Client{
		sessionID: sessionid,
		hub:       hub,
		conn:      conn,
		send:      make(chan []byte, 1024),
		player:    player,
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
			c.handleUserRegistration(message[1:])
		case MsgTypeChat:
			c.handleChatMessage(message[1:])
		case MsgTypeInput:
			c.handleUserInput(message[1:])
		case MsgTypeResumeSession:
			c.handleUserResumeSession(message[1:])
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

func (c *Client) handleUserRegistration(data []byte) {
	LocalUserData.Username = string(data)
	msg := SerializeUserStateDelta(&LocalUserData, UserStateDeltaPOS|UserStateDeltaSTATS|UserStateDeltaWEAPON)
	c.hub.broadcast <- msg
}
func (c *Client) handleChatMessage(data []byte) {
	msg := append([]byte{MsgTypeChat}, data...)
	c.hub.broadcast <- msg
}

func (c *Client) handleUserInput(data []byte) {

	offset := 0
	seq := binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	moveX := int8(data[offset])
	offset++
	moveY := int8(data[offset])
	offset++
	angle := math.Float32frombits(binary.LittleEndian.Uint32(data[offset:]))
	offset += 4

	dash := data[offset] != 0

	c.player.Input = PlayerInput{
		Seq:   seq,
		MoveX: moveX,
		MoveY: moveY,
		Dash:  dash,
		Angle: angle,
	}
}
func (c *Client) handleUserResumeSession(data []byte) {
	msg := SerializeUserStateDelta(&LocalUserData, UserStateDeltaPOS|UserStateDeltaSTATS|UserStateDeltaWEAPON)
	c.hub.broadcast <- msg
}
