package main

import (
	"bytes"
	"encoding/binary"
)

func SerializeUserDead(username string) []byte {
	nameBytes := []byte(username)
	b := make([]byte, 0, 1+1+len(nameBytes))
	b = append(b, MsgTypeUserDead)
	b = append(b, byte(len(nameBytes)))
	b = append(b, nameBytes...)
	return b
}

func SerializeUserShootStatus(alive bool, id uint32) []byte {

	b := make([]byte, 0, 1+1+4)
	b = append(b, MsgTypeShootStatus)
	if alive {
		b = append(b, 1)
	} else {
		b = append(b, 0)
	}
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, id)
	b = append(b, idBytes...)
	return b
}

func SerializeUserStateDelta(msgType uint8, p *Player, deltaMask byte) []byte {
	buf := &bytes.Buffer{}
	buf.WriteByte(msgType)
	buf.WriteByte(deltaMask)

	usernameBytes := []byte(p.Meta.Username)
	usernameBytes = bytes.TrimRight(usernameBytes, "\x00")
	buf.WriteByte(byte(len(usernameBytes)))
	buf.Write(usernameBytes)

	if deltaMask&UserStateDeltaPOS != 0 {
		binary.Write(buf, binary.LittleEndian, p.Movements.X)
		binary.Write(buf, binary.LittleEndian, p.Movements.Y)
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
