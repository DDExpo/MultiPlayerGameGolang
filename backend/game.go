package main

func applyInput(p *UserState, dt float32) {

	if p.Input.Dash {
		p.Speed = 2
	} else {
		p.Speed = 1
	}

	moveSpeed := float32(p.Speed) * 60.0

	p.X += float32(int8(p.Input.MoveX)) * moveSpeed * dt
	p.Y += float32(int8(p.Input.MoveY)) * moveSpeed * dt

	p.Angle = p.Input.Angle
}

func clampToWorld(p *UserState) {

	if p.X < 0 || p.Y < 0 || p.X > WorldWidth || p.Y > WorldHeight {
		p.X = 4000
		p.Y = 4000
	}
}
