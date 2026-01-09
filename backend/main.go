package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	})

	log.Println("WS server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
