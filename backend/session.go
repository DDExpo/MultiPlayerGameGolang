package main

type UserState struct {
	Username         string
	CurrentSessionID string
	X, Y             float32
	Speed            int8
	Angle            float32
	HP               uint8
	Damage           uint8
	Kills            uint8
	WeaponType       uint8
	WeaponWidth      uint8
	WeaponRange      uint8
	Input            PlayerInput
	LastX, LastY     float32
	LastSpeed        int8
	LastAngle        float32
	LastHP           uint8
	LastKills        uint8
	LastDamage       uint8
	LastWeaponType   uint8
	LastWeaponWidth  uint8
	LastWeaponRange  uint8
}

type PlayerInput struct {
	Seq     uint16
	MoveX   int8
	MoveY   int8
	Dash    bool
	Angle   float32
	Actions uint8
}

var LocalUserData UserState = UserState{
	Username: "", CurrentSessionID: "", X: 4000.0, Y: 4000.0, HP: 1, Speed: 1, Angle: 90,
	Damage: 1, Kills: 0, WeaponType: 1, WeaponWidth: 1, WeaponRange: 1,
}
