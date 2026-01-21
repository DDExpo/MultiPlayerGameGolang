package main

import (
	"log"
	"net/http"

	db "multiplayerGame/database"
)

func main() {
	err := db.RunFirstTimeShemas(db.GetDB())

	if err != nil {
		log.Println(err)
		panic(err)
	}

	hub := NewHub()
	go hub.Run()
	go hub.RunGameLoop()

	http.HandleFunc("/session-resume", func(w http.ResponseWriter, r *http.Request) { IsSessionResume(hub, w, r) })
	http.HandleFunc("/initialize-session", func(w http.ResponseWriter, r *http.Request) { InitSession(hub, w, r) })
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { ServeWS(hub, w, r) })

	log.Println("WS server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
