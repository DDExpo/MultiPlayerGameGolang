package main

import "time"

const MinUsernameLength = 1
const MaxUsernameLength = 128
const MaxUserMsgLength = 512
const MaxAgeSession = 60 * 60 * 12

const TickRate = 60
const TickDuration = time.Second / TickRate

const WorldWidth = 8000
const WorldHeight = 8000

const (
	UserStateDeltaPOS    byte = 1 << 0 // X, Y, Speed, Angle
	UserStateDeltaSTATS  byte = 1 << 1 // HP, Kills, Damage
	UserStateDeltaWEAPON byte = 1 << 2 // WeaponType, WeaponWidth, WeaponRange
)
