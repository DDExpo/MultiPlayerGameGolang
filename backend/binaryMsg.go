package main

import (
	"bytes"
	"encoding/binary"
)

func SerializeUserStateDelta(p *UserState, deltaMask byte) []byte {
	buf := &bytes.Buffer{}
	buf.WriteByte(MsgTypeUserState)
	buf.WriteByte(deltaMask)

	usernameBytes := []byte(p.Username)
	buf.WriteByte(byte(len(usernameBytes)))
	buf.Write(usernameBytes)

	var out [16]byte
	copy(out[:], []byte(p.CurrentSessionID))
	buf.WriteByte(16)
	buf.Write(out[:])

	if deltaMask&UserStateDeltaPOS != 0 {
		binary.Write(buf, binary.LittleEndian, p.X)
		binary.Write(buf, binary.LittleEndian, p.Y)
		binary.Write(buf, binary.LittleEndian, p.Speed)
		binary.Write(buf, binary.LittleEndian, p.Angle)
	}
	if deltaMask&UserStateDeltaSTATS != 0 {
		binary.Write(buf, binary.LittleEndian, p.HP)
		binary.Write(buf, binary.LittleEndian, p.Kills)
		binary.Write(buf, binary.LittleEndian, p.Damage)
	}
	if deltaMask&UserStateDeltaWEAPON != 0 {
		binary.Write(buf, binary.LittleEndian, p.WeaponType)
		binary.Write(buf, binary.LittleEndian, p.WeaponWidth)
		binary.Write(buf, binary.LittleEndian, p.WeaponRange)
	}

	return buf.Bytes()
}

func computeDeltaMask(p *UserState) byte {
	var mask byte = 0

	// Position / movement
	if p.X != p.LastX || p.Y != p.LastY || p.Speed != p.LastSpeed || p.Angle != p.LastAngle {
		mask |= UserStateDeltaPOS
	}

	// Stats
	if p.HP != p.LastHP || p.Kills != p.LastKills || p.Damage != p.LastDamage {
		mask |= UserStateDeltaSTATS
	}

	// Weapon
	if p.WeaponType != p.LastWeaponType || p.WeaponWidth != p.LastWeaponWidth || p.WeaponRange != p.LastWeaponRange {
		mask |= UserStateDeltaWEAPON
	}

	return mask
}

func updateLastSent(p *UserState, mask byte) {
	if mask&UserStateDeltaPOS != 0 {
		p.LastX = p.X
		p.LastY = p.Y
		p.LastSpeed = p.Speed
		p.LastAngle = p.Angle
	}
	if mask&UserStateDeltaSTATS != 0 {
		p.LastHP = p.HP
		p.LastKills = p.Kills
		p.LastDamage = p.Damage
	}
	if mask&UserStateDeltaWEAPON != 0 {
		p.LastWeaponType = p.WeaponType
		p.LastWeaponWidth = p.WeaponWidth
		p.LastWeaponRange = p.WeaponRange
	}
}
