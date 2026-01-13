package main

type localUserData struct {
	localUsername    string
	currentSessionID string
}

var localData LocalUserData = LocalUserData{localUsername: "", currentSessionID: ""}
