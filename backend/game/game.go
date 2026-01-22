package game

var FastProjectileCheck = NewSpatialHash(400)

func ApplyInput(p *Player) {

	p.Movements.Speed = 100.0
	if p.Input.Dash {
		p.Movements.Speed = 200.0
	}

	p.Movements.X += p.Input.MoveX * p.Movements.Speed * DT
	p.Movements.Y += p.Input.MoveY * p.Movements.Speed * DT

	if p.Movements.X < 0 || p.Movements.Y < 0 || p.Movements.X > WorldWidth || p.Movements.Y > WorldHeight {
		p.Movements.X = 4000
		p.Movements.Y = 4000
	}

	FastProjectileCheck.Update(p.Meta.SessionID[:16], p.Movements.X, p.Movements.Y)
	p.Movements.Angle = p.Input.Angle
}

func CheckCollision(px, py float32, projectileRadius float32, owner string) (string, bool) {

	cell := FastProjectileCheck.GetCell(px, py)
	if cell == nil {
		return "", false
	}

	for id, target := range cell {
		if id == owner {
			continue
		}

		hitboxLeft := target.X - PlayerWidth/2
		hitboxRight := target.X + PlayerWidth/2
		hitboxTop := target.Y - PlayerHeight/2
		hitboxBottom := target.Y + PlayerHeight/2

		closestX := clamp(px, hitboxLeft, hitboxRight)
		closestY := clamp(py, hitboxTop, hitboxBottom)

		dx := px - closestX
		dy := py - closestY
		distanceSquared := dx*dx + dy*dy
		radiusSquared := projectileRadius * projectileRadius

		if distanceSquared < radiusSquared {
			return id, true
		}
	}

	return "", false
}

func ApplyDamage(target *Player, attacker *Player) {
	target.Combat.HP -= attacker.Combat.Damage
	if target.Combat.HP <= 0 {
		attacker.Combat.Kills += 1
		attacker.Combat.HP += 1
		target.Movements.X = 4000
		target.Movements.Y = 4000
	}
}

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
