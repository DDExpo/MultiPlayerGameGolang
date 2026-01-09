package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan WSMessage
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan WSMessage, 1024),
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

		var wsMesage WSMessage
		if err := json.Unmarshal(message, &wsMesage); err != nil {
			log.Printf("error marshalling even :%v", err)
		}

		c.hub.broadcast <- wsMesage
	}
}
func (c *Client) WritePump() {
	defer c.conn.Close()

	for msg := range c.send {
		data, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			return
		}
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("write error:", err)
			}
			return
		}
	}
}
