package main

import (
	"encoding/binary"
	"math"
)

type BinaryWriter struct {
	buf []byte
}

func NewBinaryWriter() *BinaryWriter {
	return &BinaryWriter{buf: make([]byte, 0, 256)}
}

func (w *BinaryWriter) Bytes() []byte {
	return w.buf
}

func (w *BinaryWriter) WriteUint8(v uint8) {
	w.buf = append(w.buf, v)
}

func (w *BinaryWriter) WriteInt16(v int16) {
	w.buf = binary.LittleEndian.AppendUint16(w.buf, uint16(v))
}

func (w *BinaryWriter) WriteUint16(v uint16) {
	w.buf = binary.LittleEndian.AppendUint16(w.buf, v)
}

func (w *BinaryWriter) WriteUint32(v uint32) {
	w.buf = binary.LittleEndian.AppendUint32(w.buf, v)
}

func (w *BinaryWriter) WriteFloat32(v float32) {
	w.buf = binary.LittleEndian.AppendUint32(w.buf, math.Float32bits(v))
}

func (w *BinaryWriter) WriteString(s string) {
	w.WriteUint16(uint16(len(s)))
	w.buf = append(w.buf, []byte(s)...)
}

func SerializeUserReg(p *Player) []byte {
	w := NewBinaryWriter()
	w.WriteUint8(MsgTypeUserState)
	w.WriteUint8(StateTypeUserReg)
	w.WriteString(p.Meta.Username)
	w.WriteFloat32(p.Movements.X)
	w.WriteFloat32(p.Movements.Y)
	w.WriteFloat32(p.Movements.Angle)
	w.WriteInt16(p.Combat.HP)
	w.WriteUint16(p.Combat.Kills)
	w.WriteInt16(p.Combat.Damage)
	w.WriteUint8(p.Combat.WeaponType)
	w.WriteUint8(p.Combat.WeaponWidth)
	w.WriteUint16(p.Combat.WeaponRange)
	return w.Bytes()
}

func SerializeUserDead(username string) []byte {
	w := NewBinaryWriter()
	w.WriteUint8(MsgTypeUserState)
	w.WriteUint8(StateTypeUserDead)
	w.WriteString(username)
	return w.Bytes()
}

func SerializeUserCurrentState(deltaMask uint8, p *Player) []byte {
	w := NewBinaryWriter()
	w.WriteUint8(MsgTypeUserState)
	w.WriteUint8(StateTypeUserCurrentState)
	w.WriteString(p.Meta.Username)
	w.WriteUint8(deltaMask)

	if deltaMask&UserStateDeltaPOS != 0 {
		w.WriteFloat32(p.Movements.X)
		w.WriteFloat32(p.Movements.Y)
		w.WriteFloat32(p.Movements.Angle)
	}

	if deltaMask&UserStateDeltaSTATS != 0 {
		w.WriteInt16(p.Combat.HP)
		w.WriteUint16(p.Combat.Kills)
		w.WriteInt16(p.Combat.Damage)
	}

	if deltaMask&UserStateDeltaWEAPON != 0 {
		w.WriteUint8(p.Combat.WeaponType)
		w.WriteUint8(p.Combat.WeaponWidth)
		w.WriteUint16(p.Combat.WeaponRange)
	}

	return w.Bytes()
}

func SerializeUserPressedShoot(p *Player) []byte {
	w := NewBinaryWriter()
	w.WriteUint8(MsgTypeUserState)
	w.WriteUint8(StateTypeUserPressedShoot)
	w.WriteString(p.Meta.Username)
	w.WriteFloat32(p.Movements.X)
	w.WriteFloat32(p.Movements.Y)
	w.WriteFloat32(p.Movements.Angle)
	w.WriteUint8(p.Combat.WeaponWidth)
	w.WriteUint16(p.Combat.WeaponRange)
	return w.Bytes()
}

func SerializeUserShootStatus(alive bool, id uint16) []byte {
	w := NewBinaryWriter()
	w.WriteUint8(MsgTypeUserShootStatus)
	if alive {
		w.WriteUint8(1)
	} else {
		w.WriteUint8(0)
	}
	w.WriteUint16(id)
	return w.Bytes()
}

func SerializeUserChat(username, text, timestamp, color string) []byte {
	w := NewBinaryWriter()
	w.WriteUint8(MsgTypeUserChat)
	w.WriteString(username)
	w.WriteString(text)
	w.WriteString(timestamp)
	w.WriteString(color)
	return w.Bytes()
}
