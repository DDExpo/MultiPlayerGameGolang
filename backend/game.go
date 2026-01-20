package main

import (
	"sync"
	"time"
)

var mu sync.Mutex
var Projectiles []*Projectile = []*Projectile{}
var FastProjectileCheck = NewSpatialHash(400)

func applyInput(p *Player) {

	if p.Input.Dash {
		p.Movements.Speed = 2
	} else {
		p.Movements.Speed = 1
	}

	moveSpeed := float32(p.Movements.Speed) * 100.0

	p.Movements.X += float32(p.Input.MoveX) * moveSpeed * DT
	p.Movements.Y += float32(p.Input.MoveY) * moveSpeed * DT

	if p.Movements.X < 0 || p.Movements.Y < 0 || p.Movements.X > WorldWidth || p.Movements.Y > WorldHeight {
		p.Movements.X = 4000
		p.Movements.Y = 4000
	}

	FastProjectileCheck.Update(p.Meta.SessionID[:16], p.Movements.X, p.Movements.Y)
	p.Movements.Angle = p.Input.Angle
}
func simulateProjectile(p *Projectile, now time.Time) (bool, string) {
	elapsed := float32(now.Sub(p.SpawnTime).Seconds())

	px := p.X + p.VX*elapsed
	py := p.Y + p.VY*elapsed

	if hit, ok := CheckCollision(px, py, p.Radius, p.OwnerId); ok {
		return false, hit
	}

	if px >= 3590 && px <= 4410 && py >= 3590 && py <= 4410 {
		return false, ""
	}

	if px < 0 || py < 0 || px > WorldWidth || py > WorldHeight {
		return false, ""
	}

	if elapsed >= p.LifeTime {
		return false, ""
	}

	return true, ""
}

func CheckCollision(px, py float32, ProjectileRadius float32, owner string) (string, bool) {

	cell := FastProjectileCheck.GetCell(px, py)
	if cell == nil {
		return "", false
	}

	for id, target := range cell {
		if id == owner {
			continue
		}

		if px >= target.X-ProjectileRadius &&
			px <= target.X+PlayerWidth+ProjectileRadius &&
			py >= target.Y-ProjectileRadius &&
			py <= target.Y+PlayerHeight+ProjectileRadius {
			return id, true
		}
	}

	return "", false
}

func ApplyDamage(target *Player, attacker *Player) (string, bool) {
	target.Combat.HP -= attacker.Combat.Damage
	if target.Combat.HP <= 0 {
		attacker.Combat.Kills += 1
		attacker.Combat.HP += 1
	}
	return "", false
}
