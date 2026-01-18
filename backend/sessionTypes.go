package main

type PLayerMovements struct {
	X         float32
	Y         float32
	Angle     float32
	LastX     float32
	LastY     float32
	LastAngle float32
	Speed     int8
	LastSpeed int8
	_         [2]byte
}

type PlayerCombat struct {
	HP              uint8
	Damage          uint8
	Kills           uint8
	WeaponType      uint8
	WeaponWidth     uint8
	WeaponRange     uint8
	LastHP          uint8
	LastDamage      uint8
	LastKills       uint8
	LastWeaponType  uint8
	LastWeaponWidth uint8
	LastWeaponRange uint8
}

type PlayerMetadata struct {
	Username  string
	SessionID string
}

type PlayerInput struct {
	Seq   uint16
	MoveX int8
	MoveY int8
	Angle float32
	Dash  bool
}

type Player struct {
	Movements *PLayerMovements
	Combat    *PlayerCombat
	Meta      *PlayerMetadata
	Input     PlayerInput

	IsConnected bool
}
