package main

import (
	"log"

	"github.com/gorilla/websocket"
)

const (
	MsgTypeUser          byte = 1
	MsgTypeChat          byte = 2
	MsgTypeUserState     byte = 3
	MsgTypeResumeSession byte = 4
)

type Client struct {
	sessionID string
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn, sessionid string) *Client {
	return &Client{
		sessionID: sessionid,
		hub:       hub,
		conn:      conn,
		send:      make(chan []byte, 1024),
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
		case MsgTypeUserState:
			c.handleUserStateMessage(message[1:])
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
	LocalData.localUsername = string(data)
}
func (c *Client) handleChatMessage(data []byte) {
	msg := append([]byte{MsgTypeChat}, data...)
	c.hub.broadcast <- msg
}
func (c *Client) handleUserStateMessage(data []byte) {
	if len(data) < 13 {
		log.Println("user state message too short")
		return
	}
	msg := append([]byte{MsgTypeUserState}, data...)
	c.hub.broadcast <- msg
}
func (c *Client) handleUserResumeSession(data []byte) {
	msg := append([]byte{MsgTypeResumeSession}, data...)
	c.hub.broadcast <- msg
}
