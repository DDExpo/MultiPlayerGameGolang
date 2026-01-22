package main

const MinUsernameLength = 1
const MaxUsernameLength = 128
const MaxUserMsgLength = 512
const MaxAgeSession = 60 * 60 * 12

const (
	UserStateDeltaPOS    byte = 1 << 0 // X, Y, Angle
	UserStateDeltaSTATS  byte = 1 << 1 // HP, Kills, Damage
	UserStateDeltaWEAPON byte = 1 << 2 // WeaponType, WeaponSpeed, WeaponWidth, WeaponRange
)

const (
	_ = iota
	MsgTypeUserChat
	MsgTypeUserState
	MsgTypeUserShootStatus
	MsgTypeUserResumedDeath
	StateTypeUserReg
	StateTypeUserDead
	StateTypeUserInput
	StateTypeUserResume
	StateTypeUserCurrentState
	StateTypeUserPressedShoot
)
