package main

import "time"

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
	HP              int16
	Damage          int16
	Kills           uint16
	WeaponType      uint8
	WeaponWidth     uint8
	WeaponRange     uint16
	LastHP          int16
	LastDamage      int16
	LastKills       uint16
	LastWeaponType  uint8
	LastWeaponWidth uint8
	LastWeaponRange uint16
}

type PlayerMetadata struct {
	Username  string
	SessionID string
}

type PlayerInput struct {
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

type Projectile struct {
	X, Y           float32
	VX, VY         float32
	OwnerId        string
	Damage         int16
	Radius         float32
	Width          uint8
	Range          uint16
	ProjectileType uint8
	ProjectileId   uint16
	LifeTime       float32
	SpawnTime      time.Time
}
