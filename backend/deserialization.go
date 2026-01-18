package main

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

func DeserializeUserMsg(msg []byte) (text string, color string, err error) {
	if len(msg) < 3 {
		return "", "", fmt.Errorf("packet too small")
	}

	offset := 0

	readText := func() (string, error) {
		if offset+2 > len(msg) {
			return "", fmt.Errorf("invalid length prefix")
		}
		textLen := int(binary.LittleEndian.Uint16(msg[offset:]))
		offset += 2

		if offset+textLen > len(msg) {
			return "", fmt.Errorf("invalid length")
		}

		b := msg[offset : offset+textLen]
		offset += textLen
		return string(b), nil
	}

	readColor := func() (string, error) {
		if offset+1 > len(msg) {
			return "", fmt.Errorf("invalid length prefix")
		}
		colorLen := int(msg[offset])
		offset++

		if offset+colorLen > len(msg) {
			return "", fmt.Errorf("invalid length")
		}

		b := msg[offset : offset+colorLen]
		offset += colorLen
		return string(b), nil
	}

	text, err = readText()
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

	color, err = readColor()
	if err != nil {
		return
	}

	if !regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`).MatchString(color) {
		color = "#ffffff"
	}

	return
}
