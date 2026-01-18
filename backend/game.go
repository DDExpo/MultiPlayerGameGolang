package main

func applyInput(p *Player, dt float32) {

	if p.Input.Dash {
		p.Movements.Speed = 2
	} else {
		p.Movements.Speed = 1
	}

	moveSpeed := float32(p.Movements.Speed) * 60.0

	p.Movements.X += float32(int8(p.Input.MoveX)) * moveSpeed * dt
	p.Movements.Y += float32(int8(p.Input.MoveY)) * moveSpeed * dt

	p.Movements.Angle = p.Input.Angle
}

func clampToWorld(p *Player) {

	if p.Movements.X < 0 || p.Movements.Y < 0 || p.Movements.X > WorldWidth || p.Movements.Y > WorldHeight {
		p.Movements.X = 4000
		p.Movements.Y = 4000
	}
}
