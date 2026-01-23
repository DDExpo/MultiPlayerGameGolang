package main

import "multiplayerGame/game"

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

type GameCommand interface {
	isGameCommand()
}

type SpawnProjectileCmd struct {
	Player *game.Player
}
type UserRegistrationCmd struct {
	Player *game.Player
}
type UserResumeSession struct {
	Player *game.Player
}
type UserResumedDeathCmd struct {
	Player *game.Player
}
type UserInputCmd struct {
	Player *game.Player
	Input  game.PlayerInput
}

func (UserRegistrationCmd) isGameCommand() {}
func (UserResumeSession) isGameCommand()   {}
func (UserResumedDeathCmd) isGameCommand() {}
func (SpawnProjectileCmd) isGameCommand()  {}
func (UserInputCmd) isGameCommand()        {}
