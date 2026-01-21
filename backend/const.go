package main

import "time"

const MinUsernameLength = 1
const MaxUsernameLength = 128
const MaxUserMsgLength = 512
const MaxAgeSession = 60 * 60 * 12

const TickRate = 60
const TickDuration = time.Second / TickRate

var DT = float32(TickDuration.Seconds())

const PlayerWidth = 21
const PlayerHeight = 24

const ProjectileSpeed = 358
const projectileRadius = 3
const ProjectileLifetime = 5.0 // Seconds

const WorldWidth = 8000
const WorldHeight = 8000

const (
	UserStateDeltaPOS    byte = 1 << 0 // X, Y, Angle
	UserStateDeltaSTATS  byte = 1 << 1 // HP, Kills, Damage
	UserStateDeltaWEAPON byte = 1 << 2 // WeaponType, WeaponWidth, WeaponRange
)

const (
	MsgTypeUserChat           = 1
	MsgTypeUserState          = 2
	MsgTypeUserShootStatus    = 3
	MsgTypeUserResumedDeath   = 4
	StateTypeUserReg          = 5
	StateTypeUserDead         = 6
	StateTypeUserInput        = 7
	StateTypeUserResume       = 8
	StateTypeUserCurrentState = 9
	StateTypeUserPressedShoot = 10
)
