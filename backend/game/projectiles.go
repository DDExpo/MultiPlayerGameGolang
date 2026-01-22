package game

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	projectiles   []*Projectile
	projectilesMu sync.RWMutex
	projectileID  atomic.Uint32
)

var projectilePool = sync.Pool{
	New: func() any {
		return &Projectile{}
	},
}

func AddProjectile(p *Projectile) {
	projectilesMu.Lock()
	defer projectilesMu.Unlock()

	projectiles = append(projectiles, p)
}

func TickProjectiles(now time.Time) (map[string][]Hit, []uint32) {
	projectilesMu.Lock()
	defer projectilesMu.Unlock()

	hits := make(map[string][]Hit)
	deadIDs := make([]uint32, 0)
	alive := projectiles[:0]

	for _, prj := range projectiles {
		dt := float32(now.Sub(prj.LastUpdateTime).Seconds())

		prj.X += prj.VX * dt
		prj.Y += prj.VY * dt
		dead := false

		if targetId, collided := CheckCollision(prj.X, prj.Y, prj.Radius, prj.OwnerId); collided {
			hits[targetId] = append(hits[targetId], Hit{
				ProjectileId: prj.ProjectileId,
				OwnerId:      prj.OwnerId,
				Damage:       prj.Damage,
			})
			dead = true

		} else if prj.X >= 3590 && prj.X <= 4410 && prj.Y >= 3590 && prj.Y <= 4410 {
			dead = true

		} else if prj.X < 0 || prj.Y < 0 || prj.X > WorldWidth || prj.Y > WorldHeight {
			dead = true

		} else if dt >= prj.LifeTime {
			dead = true
		}

		prj.LastUpdateTime = now
		if dead {
			deadIDs = append(deadIDs, prj.ProjectileId)
			projectilePool.Put(prj)
		} else {
			alive = append(alive, prj)
		}

	}

	projectiles = alive
	return hits, deadIDs
}
