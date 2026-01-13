package main

import (
	"log"
	"net/http"

	db "multiplayerGame/database"
)

type LocalUserData struct {
	localUsername    string
	currentSessionID string
}

var LocalData LocalUserData = LocalUserData{localUsername: "", currentSessionID: ""}

func main() {
	err := db.RunFirstTimeShemas(db.GetDB())

	if err != nil {
		log.Println(err)
		panic(err)
	}

	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/session-resume", IsSessionResume)
	http.HandleFunc("/session", InitSession)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(hub, w, r)
	})

	log.Println("WS server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
