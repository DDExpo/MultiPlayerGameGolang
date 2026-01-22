package game

import "time"

const TickRate = 60
const TickDuration = time.Second / TickRate

var DT = float32(TickDuration.Seconds())

const PlayerWidth = 21
const PlayerHeight = 24

const ProjectileSpeed = 358
const projectileRadius = 3
const ProjectileLifetime = 5.0 // Seconds

const WorldWidth = 8000
const WorldHeight = 8000
