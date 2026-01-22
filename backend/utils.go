package main

import (
	"multiplayerGame/game"
)

func computeDeltaMask(p *game.Player) byte {
	var mask byte = 0

	if p.Movements.X != p.Movements.LastX || p.Movements.Y != p.Movements.LastY || p.Movements.Speed != p.Movements.LastSpeed || p.Movements.Angle != p.Movements.LastAngle {
		mask |= UserStateDeltaPOS
	}

	if p.Combat.HP != p.Combat.LastHP || p.Combat.Kills != p.Combat.LastKills || p.Combat.Damage != p.Combat.LastDamage {
		mask |= UserStateDeltaSTATS
	}

	if p.Combat.WeaponType != p.Combat.LastWeaponType || p.Combat.WeaponWidth != p.Combat.LastWeaponWidth || p.Combat.WeaponRange != p.Combat.LastWeaponRange {
		mask |= UserStateDeltaWEAPON
	}

	return mask
}

func updateLastSent(p *game.Player, mask byte) {
	if mask&UserStateDeltaPOS != 0 {
		p.Movements.LastX = p.Movements.X
		p.Movements.LastY = p.Movements.Y
		p.Movements.LastSpeed = p.Movements.Speed
		p.Movements.LastAngle = p.Movements.Angle
	}
	if mask&UserStateDeltaSTATS != 0 {
		p.Combat.LastHP = p.Combat.HP
		p.Combat.LastKills = p.Combat.Kills
		p.Combat.LastDamage = p.Combat.Damage
	}
	if mask&UserStateDeltaWEAPON != 0 {
		p.Combat.LastWeaponType = p.Combat.WeaponType
		p.Combat.LastWeaponWidth = p.Combat.WeaponWidth
		p.Combat.LastWeaponRange = p.Combat.WeaponRange
	}
}
