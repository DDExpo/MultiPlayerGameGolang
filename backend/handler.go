package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func IsSessionResume(h *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("No session cookie: %v", err)
		http.Error(w, "no session found", http.StatusUnauthorized)
		return
	}

	username, timeRegistered, err := db.DBGetUser(db.GetDB(), cookie.Value)
	if err != nil {
		log.Printf("Something went wrong: %s.  err: %v", cookie.Value, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	age := time.Now().Unix() - timeRegistered
	if age > MaxAgeSession {
		err = db.DbDelete(db.GetDB(), []string{cookie.Value})
		if err != nil {
			log.Printf("Something went wrong: %s, err: %v", cookie.Value, err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		log.Printf("Session expired: %s (age: %d)", cookie.Value, age)
		http.Error(w, "session expired", http.StatusUnauthorized)
		return
	}

	key := cookie.Value[:16]
	h.mu.RLock()
	_, exists := h.players[key]
	h.mu.RUnlock()

	if !exists {
		h.mu.Lock()
		h.players[key] = NewPlayer(username, cookie.Value)
		h.activeUsernames[username] = true
		h.mu.Unlock()

		log.Printf("Recreated player for session resume: %s", username)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    cookie.Value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   MaxAgeSession,
	})

	resp := map[string]string{"username": username}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}

	log.Printf("Session resumed: %s", username)
}

func InitSession(h *Hub, w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	trimmedUsername := strings.TrimSpace(req.Username)
	if err := validateUsername(trimmedUsername); err != nil {
		log.Printf("Username validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.mu.RLock()
	isTaken := h.activeUsernames[trimmedUsername]
	h.mu.RUnlock()
	if isTaken {
		http.Error(w, "username already taken", http.StatusBadRequest)
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		log.Printf("Failed to generate session ID: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user := db.UserData{
		ID:             sessionID,
		Username:       trimmedUsername,
		TimeRegistered: time.Now().Unix(),
	}
	if err := db.DBSaveUser(db.GetDB(), user); err != nil {
		log.Printf("Failed to save user to database: %v", err)
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   MaxAgeSession,
	})

	log.Printf("Created new session for user '%s': %s", trimmedUsername, sessionID)
	h.mu.Lock()
	h.players[sessionID[:16]] = NewPlayer(trimmedUsername, sessionID)
	h.activeUsernames[trimmedUsername] = true
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("No session cookie in WebSocket request: %v", err)
		http.Error(w, "unauthorized - no session", http.StatusUnauthorized)
		return
	}

	hub.mu.RLock()
	var player *Player = hub.players[cookie.Value[:16]]
	hub.mu.RUnlock()

	if player == nil {
		log.Printf("Player not found for session: %s", cookie.Value)
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := NewClient(hub, conn, player)

	hub.mu.Lock()
	player.IsConnected = true
	hub.mu.Unlock()

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
