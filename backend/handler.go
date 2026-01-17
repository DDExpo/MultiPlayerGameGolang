package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "multiplayerGame/database"

	"github.com/gorilla/websocket"
)

type SessionRequestUsername struct {
	Username string `json:"username"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func IsSessionResume(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("No session cookie found: %v", err)
		http.Error(w, "internal error", 400)
		return
	}

	log.Printf("Found session cookie: %s", cookie.Value)

	username, timeRegistered, dbErr := db.DBGetUser(db.GetDB(), cookie.Value)
	if dbErr != nil {
		log.Printf("DB lookup failed for session %s: %v", cookie.Value, dbErr)
		LocalUserData.CurrentSessionID = ""
		http.Error(w, "internal error", 400)
		return
	}

	age := time.Now().Unix() - timeRegistered
	log.Printf("Session age: %d seconds, max age: %d", age, MaxAgeSession)

	if age > MaxAgeSession {
		log.Printf("Session expired")
		LocalUserData.CurrentSessionID = ""
		http.Error(w, "internal error", 400)
		return
	}

	LocalUserData.CurrentSessionID = cookie.Value
	LocalUserData.Username = username

	resp := SessionRequestUsername{Username: username}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}

	log.Printf("Session valid, reusing")
}

func InitSession(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SessionRequestUsername
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if LocalUserData.CurrentSessionID == "" {
		sessionID, err := generateSessionID()
		if err != nil {
			http.Error(w, "internal error", 500)
			return
		}

		LocalUserData.CurrentSessionID = sessionID
		log.Printf("Creating new session (HTTP): %s", sessionID)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   MaxAgeSession,
		})
	}

	user := db.UserData{
		ID:             LocalUserData.CurrentSessionID,
		Username:       req.Username,
		TimeRegistered: time.Now().Unix(),
	}
	err := db.DBSaveUser(db.GetDB(), user)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade failed:", err)
		return
	}

	client := NewClient(hub, conn, LocalUserData.CurrentSessionID)
	hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
