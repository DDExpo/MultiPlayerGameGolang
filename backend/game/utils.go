package game

import (
	"math"
	"time"
)

func NewPlayer(username, sessionID string) *Player {
	return &Player{
		Movements: &PLayerMovements{
			X:         4000,
			Y:         4000,
			Angle:     0,
			LastX:     4000,
			LastY:     4000,
			LastAngle: 0,
			Speed:     100,
			LastSpeed: 100,
		},

		Combat: &PlayerCombat{
			HP:              1,
			Damage:          1,
			Kills:           0,
			WeaponType:      0,
			WeaponSpeed:     1,
			WeaponWidth:     1,
			WeaponRange:     1,
			LastHP:          1,
			LastDamage:      1,
			LastKills:       0,
			LastWeaponType:  0,
			LastWeaponSpeed: 1,
			LastWeaponWidth: 1,
			LastWeaponRange: 1,
		},

		Meta: &PlayerMetadata{
			Username:  username,
			SessionID: sessionID,
		},

		Input: PlayerInput{
			MoveX: 0,
			MoveY: 0,
			Angle: 0,
			Dash:  false,
		},
	}
}

func ResetStats(p *Player) {
	p.Movements.X = 4000
	p.Movements.Y = 4000
	p.Movements.LastX = 4000
	p.Movements.LastY = 4000
	p.Movements.Angle = 0
	p.Movements.LastAngle = 0
	p.Movements.Speed = 1
	p.Movements.LastSpeed = 1

	p.Combat.HP = 1
	p.Combat.Kills = 0
	p.Combat.Damage = 1
	p.Combat.WeaponType = 1
	p.Combat.WeaponSpeed = 1
	p.Combat.WeaponWidth = 1
	p.Combat.WeaponRange = 1

	p.Input.MoveX = 0
	p.Input.MoveY = 0
	p.Input.Angle = 0
	p.Input.Dash = false
}

func CreateProjectile(p *Player) *Projectile {
	radians := (float64(p.Movements.Angle) - 90.0) * (math.Pi / 180.0)
	speed := ProjectileSpeed * float64(p.Combat.WeaponSpeed)
	vx := float32(math.Cos(radians) * speed)
	vy := float32(math.Sin(radians) * speed)

	prj := projectilePool.Get().(*Projectile)
	*prj = Projectile{
		X:              p.Movements.X,
		Y:              p.Movements.Y,
		VX:             vx,
		VY:             vy,
		OwnerId:        p.Meta.SessionID[:16],
		Damage:         p.Combat.Damage,
		Radius:         projectileRadius * float32(p.Combat.WeaponWidth),
		ProjectileType: p.Combat.WeaponType,
		ProjectileId:   projectileID.Add(1),
		LifeTime:       ProjectileLifetime * float32(p.Combat.WeaponRange),
		SpawnTime:      time.Now(),
		LastUpdateTime: time.Now(),
	}
	return prj
}
