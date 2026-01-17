import { WORLD_HEIGHT, WORLD_WIDTH } from "$lib/Consts"
import { Container, Graphics } from "pixi.js"

export class Projectile {
  sprite:   Graphics
  vx:       number
  vy:       number
  lifetime: number
  damage:   number
  ownerId:  string

  constructor(x: number, y: number, angle: number, speed: number, ownerId: string, damage: number) {
    this.sprite = new Graphics()
    this.sprite.circle(0, 0, 3).fill({ color: 0xffffff })
    this.sprite.position.set(x, y)
    this.sprite.angle = angle
  
    const radians = (angle - 90) * (Math.PI / 180)
    this.vx = Math.cos(radians) * speed
    this.vy = Math.sin(radians) * speed
    
    this.lifetime = 0
    this.damage   = damage
    this.ownerId  = ownerId
  }

  update(deltaTime: number): boolean {
    this.sprite.x += this.vx * deltaTime
    this.sprite.y += this.vy * deltaTime
    this.lifetime += deltaTime

    return this.lifetime > 1000 || this.isOutOfBounds()
  }

  isOutOfBounds(): boolean {
    return this.sprite.x < 0 || this.sprite.x > WORLD_WIDTH ||
           this.sprite.y < 0 || this.sprite.y > WORLD_HEIGHT
  }
}

export class ProjectilePool {
  active:    Projectile[] = []
  inactive:  Projectile[] = []
  container: Container

  constructor(container: Container) {
    this.container = container
  }

  spawn(x: number, y: number, angle: number, speed: number, ownerId: string, damage: number, width: number, range: number) {
    let projectile: Projectile
    
    if (this.inactive.length > 0) {
      projectile = this.inactive.pop()!
      projectile.sprite.position.set(x, y)
      projectile.sprite.angle = angle
      const radians = (angle - 90) * (Math.PI / 180)
      projectile.vx = Math.cos(radians) * speed
      projectile.vy = Math.sin(radians) * speed
      projectile.lifetime = 0
      projectile.ownerId = ownerId
      projectile.sprite.visible = true
    } else {
      projectile = new Projectile(x, y, angle, speed, ownerId, damage)
      this.container.addChild(projectile.sprite)
    }
    this.active.push(projectile)
  }

  update(deltaTime: number) {
    for (let i = this.active.length - 1; i >= 0; i--) {
      const projectile = this.active[i]
      const shouldRemove = projectile.update(deltaTime)
      
      if (shouldRemove) {
        this.active.splice(i, 1)
        projectile.sprite.visible = false
        this.inactive.push(projectile)
      }
    }
  }

  checkCollisions(targets: Container[]): { projectile: Projectile, target: Container }[] {
    const hits: { projectile: Projectile, target: Container }[] = []
    
    for (const projectile of this.active) {
      for (const target of targets) {
        if (this.circleCollision(projectile.sprite, target)) {
          hits.push({ projectile, target })
        }
      }
    }
    
    return hits
  }

  circleCollision(a: Container, b: Container): boolean {
    const dx = a.x - b.x
    const dy = a.y - b.y
    const distance = Math.sqrt(dx * dx + dy * dy)

    const aRadius = Math.max(a.width, a.height) / 2
    const bRadius = Math.max(b.width, b.height) / 2
    const radiusSum = aRadius + bRadius
    
    return distance < radiusSum
  }

  destroy(projectile: Projectile) {
    const index = this.active.indexOf(projectile)
    if (index > -1) {
      this.active.splice(index, 1)
      projectile.sprite.visible = false
      this.inactive.push(projectile)
    }
  }
}