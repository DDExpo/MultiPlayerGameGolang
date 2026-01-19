import { PROJECTILE_LIFE_FRAMES, PROJECTILE_SPEED, WORLD_HEIGHT, WORLD_WIDTH } from "$lib/Consts"
import { Container, Graphics } from "pixi.js"
import { FastProjectileCheck } from "./optimizations"
import { getSocket } from "$lib/net/socket"

export class Projectile {
  sprite:   Graphics
  vx:       number
  vy:       number
  lifetime: number
  ownerId:  string
  damage:   number
  width:  number
  range:  number


  constructor(x: number, y: number, angle: number, speed: number, ownerId: string, damage: number, width: number, range: number) {
    this.sprite = new Graphics()
    this.sprite.circle(0, 0, 3 * width).fill({ color: 0xffffff })
    this.sprite.position.set(x, y)
    this.sprite.angle = angle
  
    const radians = (angle - 90) * (Math.PI / 180)
    this.vx = Math.cos(radians) * speed
    this.vy = Math.sin(radians) * speed
    
    this.lifetime = 0
    this.damage   = damage
    this.ownerId  = ownerId
    this.width    = width
    this.range    = range
  }

  update(deltaTime: number): boolean {
    this.sprite.x += this.vx * deltaTime
    this.sprite.y += this.vy * deltaTime
    this.lifetime += deltaTime
    return this.isCollided() || this.lifetime > (PROJECTILE_LIFE_FRAMES * this.range) || this.isOutOfBounds() || this.isSafeZone()
  }

  isOutOfBounds(): boolean {
    return this.sprite.x < 0 || this.sprite.x > WORLD_WIDTH ||
           this.sprite.y < 0 || this.sprite.y > WORLD_HEIGHT
  }

  isSafeZone(): boolean {
    return this.sprite.x >= 3590 && this.sprite.x <= 4410 && this.sprite.y >= 3590 && this.sprite.y <= 4410
  } 

  isCollided(): boolean {
    for (const [usr , [userX, userY]] of FastProjectileCheck.getCell(this.sprite.x, this.sprite.y)) {
     if (usr != this.ownerId) {
        const dx = userX - this.sprite.x
        const dy = userY - this.sprite.y
        if (dx * dx + dy * dy < 50)  {
          return true
        }
     } 
    }
    return false 
  }

}

export class ProjectilePool {
  active:    Projectile[] = []
  inactive:  Projectile[] = []
  container: Container

  constructor(container: Container) {
    this.container = container
  }

  spawn(x: number, y: number, angle: number, ownerId: string, damage: number, width: number, range: number) {
    let projectile: Projectile
    
    if (this.inactive.length > 0) {
      projectile = this.inactive.pop()!
      projectile.sprite.position.set(x, y)
      projectile.sprite.angle = angle
      const radians = (angle - 90) * (Math.PI / 180)
      projectile.vx = Math.cos(radians) * PROJECTILE_SPEED
      projectile.vy = Math.sin(radians) * PROJECTILE_SPEED
      projectile.lifetime = 0
      projectile.ownerId = ownerId
      projectile.sprite.visible = true
    } else {
      projectile = new Projectile(x, y, angle, PROJECTILE_SPEED, ownerId, damage, width, range)
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
}