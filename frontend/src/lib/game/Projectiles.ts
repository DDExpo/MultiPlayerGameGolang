import { PROJECTILE_LIFE, PROJECTILE_SPEED } from "$lib/Consts"
import { Container, Graphics } from "pixi.js"

export class Projectile {
  sprite: Graphics
  x: number
  y: number
  vx: number
  vy: number
  lifetime: number
  ownerId: string
  width: number
  range: number
  id: number
  spawnTime: number
  lastUpdateTime: number

  constructor(id: number, x: number, y: number, angle: number, speed: number,
              ownerId: string, width: number,  range: number )
  {
    this.id = id
    this.x = x
    this.y = y
    this.spawnTime = performance.now()
    this.lastUpdateTime = this.spawnTime
    this.ownerId = ownerId
    this.width = width
    this.range = range
    this.lifetime = PROJECTILE_LIFE * range

    this.sprite = new Graphics()
    this.sprite.circle(0, 0, 3 * width).fill({ color: 0xffffff })
    this.sprite.position.set(x, y)
    this.sprite.angle = angle

    const radians = (angle - 90) * (Math.PI / 180)
    this.vx = Math.cos(radians) * speed
    this.vy = Math.sin(radians) * speed
  }

  reset(id: number, x: number, y: number, angle: number, speed: number,
        ownerId: string, width: number, range: number)
  {
    this.id = id
    this.x = x
    this.y = y
    this.spawnTime = performance.now()
    this.lastUpdateTime = this.spawnTime
    this.ownerId = ownerId
    this.width = width
    this.range = range
    this.lifetime = PROJECTILE_LIFE * range
    this.sprite.visible = true
    this.sprite.position.set(x, y)
    this.sprite.angle = angle

    const radians = (angle - 90) * (Math.PI / 180)
    this.vx = Math.cos(radians) * speed
    this.vy = Math.sin(radians) * speed
  }

  update(currentTime: number) {
    const dt = (currentTime - this.lastUpdateTime) / 1000
    this.x += this.vx * dt
    this.y += this.vy * dt
    this.sprite.position.set(this.x, this.y)
    this.lastUpdateTime = currentTime
  }

  hasExpired(currentTime: number): boolean {
    return (currentTime - this.spawnTime) / 1000 >= this.lifetime
  }

}

export class ProjectilePool {
  private active = new Map<number, Projectile>()
  private inactive: Projectile[] = []
  private container: Container
  private availableIds: number[] = []

  constructor(container: Container) { this.container = container }

  spawn(
    id: number, x: number, y: number, angle: number, ownerId: string,
    speed: number, width: number, range: number): number
  {
    let projectile: Projectile

    if (this.inactive.length > 0) {
      projectile = this.inactive.pop()!
      projectile.reset(id, x, y, angle, PROJECTILE_SPEED * speed, ownerId, width, range)
    } else {
      projectile = new Projectile(id, x, y, angle, PROJECTILE_SPEED * speed, ownerId, width, range)
      this.container.addChild(projectile.sprite)
    }

    this.active.set(id, projectile)
    return id
  }

  update() {
    const now = performance.now()
    const expired: number[] = []

    for (const [id, projectile] of this.active.entries()) {
      projectile.update(now)
      if (projectile.hasExpired(now)) expired.push(id)
    }

    for (const id of expired) this.destroy(id)
  }

  destroy(id: number) {
    const projectile = this.active.get(id)
    if (!projectile) return
    projectile.sprite.visible = false
    this.inactive.push(projectile)
    this.active.delete(id)
    this.availableIds.push(id)
  }

  clear() {
    for (const [id] of this.active.entries()) {
      this.destroy(id)
    }
  }

  getActive(): Map<number, Projectile> {
    return this.active
  }
}

export let projectilePool: ProjectilePool | null = null

export function initProjectilePool(container: Container) {
  if (!projectilePool) {
    projectilePool = new ProjectilePool(container)
  }
}