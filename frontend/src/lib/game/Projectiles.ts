import { PROJECTILE_LIFE, PROJECTILE_SPEED} from "$lib/Consts"
import { Container, Graphics } from "pixi.js"

export class Projectile {
  sprite:    Graphics
  x:         number
  y:         number
  vx:        number
  vy:        number
  lifetime:  number
  ownerId:   string
  width:     number
  range:     number
  id:        number
  spawnTime: number


  constructor(id: number, x: number, y: number, angle: number, speed: number, ownerId: string, width: number, range: number) {
    this.id = id
    this.x = x
    this.y = y
    this.spawnTime = performance.now()
    this.sprite = new Graphics()
    this.sprite.circle(0, 0, 3 * width).fill({ color: 0xffffff })
    this.sprite.position.set(x, y)
    this.sprite.angle = angle
    this.ownerId = ownerId

    const radians = (angle - 90) * (Math.PI / 180)
    this.vx = Math.cos(radians) * speed
    this.vy = Math.sin(radians) * speed

    this.range = range
    this.lifetime = PROJECTILE_LIFE * this.range
    this.width = width
  }

  update(currentTime: number) {
    const t = (currentTime - this.spawnTime) / 1000
    this.sprite.x = this.x + this.vx * t
    this.sprite.y = this.y + this.vy * t
  }


  hasExpired(currentTime: number): boolean {
    const t = (currentTime - this.spawnTime) / 1000
    return t >= this.lifetime
  }
}

export class ProjectilePool {
  active: Map<number, Projectile> = new Map()
  inactive: Projectile[] = []
  container: Container
  private nextId = 0

  constructor(container: Container) {
    this.container = container
  }

  spawn(x: number, y: number, angle: number, ownerId: string, width: number, range: number): number {
    const id = this.nextId++
    let projectile: Projectile

    if (this.inactive.length > 0) {
      projectile = this.inactive.pop()!
      projectile.x = x
      projectile.y = y
      projectile.sprite.position.set(x, y)
      projectile.sprite.angle = angle
      const radians = (angle - 90) * (Math.PI / 180)
      projectile.vx = Math.cos(radians) * PROJECTILE_SPEED
      projectile.vy = Math.sin(radians) * PROJECTILE_SPEED
      projectile.spawnTime = performance.now()
      projectile.lifetime = PROJECTILE_LIFE * range
      projectile.ownerId = ownerId
      projectile.width = width
      projectile.range = range
      projectile.sprite.visible = true
      projectile.id = id
    }  else {
      projectile = new Projectile(id, x, y, angle, PROJECTILE_SPEED, ownerId, width, range)
      this.container.addChild(projectile.sprite)
    }

    this.active.set(id, projectile)
    return id
  }

  update() {
    const now = performance.now()
    for (const [id, projectile] of this.active.entries()) {
        projectile.update(now)
        if (projectile.hasExpired(now)) {
            this.destroyProjectile(id)
        }
    }
  }

  destroyProjectile(id: number) {
    const projectile = this.active.get(id)
    if (!projectile) return
    projectile.sprite.visible = false
    this.inactive.push(projectile)
    this.active.delete(id)
  }
}

export let projectilePool: ProjectilePool | null = null

export function initProjectilePool(container: Container) {
  if (!projectilePool) {
    projectilePool = new ProjectilePool(container)
  }
}