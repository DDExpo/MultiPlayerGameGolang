package game

import "time"

type Hit struct {
	ProjectileId uint32
	OwnerId      string
	Damage       int16
}

type Projectile struct {
	X, Y           float32
	VX, VY         float32
	OwnerId        string
	Damage         int16
	Radius         float32
	ProjectileType uint8
	ProjectileId   uint32
	LifeTime       float32
	SpawnTime      time.Time
	LastUpdateTime time.Time
}

type PLayerMovements struct {
	X         float32
	Y         float32
	Angle     float32
	LastX     float32
	LastY     float32
	LastAngle float32
	Speed     float32
	LastSpeed float32
}

type PlayerCombat struct {
	HP              int16
	Damage          int16
	Kills           uint16
	WeaponType      uint8
	WeaponSpeed     uint8
	WeaponWidth     uint8
	WeaponRange     uint16
	LastHP          int16
	LastDamage      int16
	LastKills       uint16
	LastWeaponType  uint8
	LastWeaponSpeed uint8
	LastWeaponWidth uint8
	LastWeaponRange uint16
}

type PlayerMetadata struct {
	Username  string
	SessionID string
}

type PlayerInput struct {
	MoveX float32
	MoveY float32
	Angle float32
	Dash  bool
}

type Player struct {
	Movements *PLayerMovements
	Combat    *PlayerCombat
	Meta      *PlayerMetadata
	Input     PlayerInput
}
