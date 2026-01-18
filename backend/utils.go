package main

func NewPlayer(username, sessionID string) *Player {
	return &Player{
		Movements: &PLayerMovements{
			X:         4000,
			Y:         4000,
			Angle:     0,
			LastX:     4000,
			LastY:     4000,
			LastAngle: 0,
			Speed:     1,
			LastSpeed: 1,
		},

		Combat: &PlayerCombat{
			HP:              1,
			Damage:          1,
			Kills:           0,
			WeaponType:      0,
			WeaponWidth:     1,
			WeaponRange:     1,
			LastHP:          1,
			LastDamage:      1,
			LastKills:       0,
			LastWeaponType:  0,
			LastWeaponWidth: 1,
			LastWeaponRange: 1,
		},

		Meta: &PlayerMetadata{
			Username:  username,
			SessionID: sessionID,
		},

		Input: PlayerInput{
			Seq:   0,
			MoveX: 0,
			MoveY: 0,
			Angle: 0,
			Dash:  false,
		},

		IsConnected: true,
	}
}

func computeDeltaMask(p *Player) byte {
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

func updateLastSent(p *Player, mask byte) {
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
