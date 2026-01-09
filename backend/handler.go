package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)
	client := NewClient(hub, conn)
	hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
