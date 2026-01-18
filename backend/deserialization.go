package main

import (
	"fmt"
	"regexp"
	"strings"
)

func DeserializeUserMsg(msg []byte) (text string, color string, err error) {
	if len(msg) < 2 {
		return "", "", fmt.Errorf("packet too small")
	}

	offset := 0

	readField := func() (string, error) {
		textLen := int(msg[offset])
		offset++
		if offset+textLen > len(msg) {
			return "", fmt.Errorf("invalid length")
		}
		b := msg[offset : offset+textLen]
		offset += textLen
		return string(b), nil
	}

	text, err = readField()
	if err != nil {
		return
	}

	if len(text) > MaxUserMsgLength {
		return "", "", fmt.Errorf("text too long")
	}

	text = strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, text)

	color, err = readField()
	if err != nil {
		return
	}

	if !regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`).MatchString(color) {
		color = "#ffffff"
	}
	return
}
