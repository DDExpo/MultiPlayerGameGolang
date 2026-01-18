package main

import (
	"bytes"
	"encoding/binary"
)

func SerializeUserStateDelta(p *Player, deltaMask byte) []byte {
	buf := &bytes.Buffer{}
	buf.WriteByte(MsgTypeUserState)
	buf.WriteByte(deltaMask)

	usernameBytes := []byte(p.Meta.Username)
	usernameBytes = bytes.TrimRight(usernameBytes, "\x00")
	buf.WriteByte(byte(len(usernameBytes)))
	buf.Write(usernameBytes)

	if deltaMask&UserStateDeltaPOS != 0 {
		binary.Write(buf, binary.LittleEndian, p.Movements.X)
		binary.Write(buf, binary.LittleEndian, p.Movements.Y)
		binary.Write(buf, binary.LittleEndian, p.Movements.Speed)
		binary.Write(buf, binary.LittleEndian, p.Movements.Angle)
	}
	if deltaMask&UserStateDeltaSTATS != 0 {
		binary.Write(buf, binary.LittleEndian, p.Combat.HP)
		binary.Write(buf, binary.LittleEndian, p.Combat.Kills)
		binary.Write(buf, binary.LittleEndian, p.Combat.Damage)
	}
	if deltaMask&UserStateDeltaWEAPON != 0 {
		binary.Write(buf, binary.LittleEndian, p.Combat.WeaponType)
		binary.Write(buf, binary.LittleEndian, p.Combat.WeaponWidth)
		binary.Write(buf, binary.LittleEndian, p.Combat.WeaponRange)
	}

	return buf.Bytes()
}

func SerializeUserMsg(username string, text string, timestamp string, color string) []byte {
	userBytes := []byte(username)
	textBytes := []byte(text)
	tsBytes := []byte(timestamp)
	colorBytes := []byte(color)

	size := 1 +
		1 + len(userBytes) +
		1 + len(textBytes) +
		1 + len(tsBytes) +
		1 + len(colorBytes)

	buf := make([]byte, size)
	offset := 0

	buf[offset] = byte(MsgTypeChat)
	offset++

	buf[offset] = byte(len(userBytes))
	offset++
	copy(buf[offset:], userBytes)
	offset += len(userBytes)

	buf[offset] = byte(len(textBytes))
	offset++
	copy(buf[offset:], textBytes)
	offset += len(textBytes)

	buf[offset] = byte(len(tsBytes))
	offset++
	copy(buf[offset:], tsBytes)
	offset += len(tsBytes)

	buf[offset] = byte(len(colorBytes))
	offset++
	copy(buf[offset:], colorBytes)
	offset += len(colorBytes)

	return buf
}
