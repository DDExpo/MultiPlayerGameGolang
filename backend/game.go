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

	FastProjectileCheck.Update(p.Meta.Username, p.Movements.X, p.Movements.Y)
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

func CheckCollision(px, py float32, radius float32, owner string) (string, bool) {

	cell := FastProjectileCheck.GetCell(px, py)
	if cell == nil {
		return "", false
	}

	for id, target := range cell {
		if id == owner {
			continue
		}

		closestX := px
		if px < target.X {
			closestX = target.X
		} else if px > target.X+PlayerWidth {
			closestX = target.X + PlayerWidth
		}

		closestY := py
		if py < target.Y {
			closestY = target.Y
		} else if py > target.Y+PlayerHeight {
			closestY = target.Y + PlayerHeight
		}

		dx := px - closestX
		dy := py - closestY

		if dx*dx+dy*dy < radius*radius {
			return id, true
		}
	}

	return "", false
}

func ApplyDamage(target *Player, attacker *Player) (string, bool) {
	target.Combat.HP -= attacker.Combat.Damage
	if target.Combat.HP <= 0 {
		ResetStats(target)
	}
	attacker.Combat.Kills += 1
	return "", false
}
