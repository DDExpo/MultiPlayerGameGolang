import { Application, Assets, Text, Container, BlurFilter, Graphics, Sprite, TextStyle } from "pixi.js"
import { Controller } from "./Controllers"
import phase1 from "$lib/assets/players/phase1.png"
import { uiHasFocus, userRegistered } from "$lib/stores/ui.svelte"
import { localUser, playersState } from "$lib/stores/game.svelte"
import { createInputMsg } from "$lib/net/binaryEncodingDecoding"
import { getSocket } from "$lib/net/socket"
import { WORLD_HEIGHT, WORLD_WIDTH } from "$lib/Consts"
import { randomBrightColor } from "$lib/utils"
import { ProjectilePool } from "./Projectiles"


export class Game {
  app:                  Application
  world:                Container
  playersContainer:     Container
  starfield:            Container
  projectilePool:       ProjectilePool
  projectilesContainer: Container
  lastShotTime = 0;
  shootCooldown = 1000;

  constructor() {
    this.app                  = new Application()
    this.world                = new Container()
    this.playersContainer     = new Container()
    this.projectilesContainer = new Container()
    this.starfield            = new Container()
    this.projectilePool       = new ProjectilePool(this.projectilesContainer)
  }
  
  createStarfield() {
    if (this.starfield.children.length > 0) return
    
    const starCount = 2000
    const stars = new Graphics()
    
    for (let i = 0; i < starCount; i++) {
      const size = Math.random() * 2 + 1
      const x = Math.random() * WORLD_WIDTH
      const y = Math.random() * WORLD_HEIGHT
      const color = randomBrightColor()
      const alpha = Math.random() * 0.6 + 0.3
      
      stars.circle(x, y, size).fill({ color, alpha })
    }
  
    stars.filters = [new BlurFilter({ strength: 0.5, quality: 2 })]
    
    this.starfield.addChild(stars)
    this.starfield.enableRenderGroup()
  }

  async init() {
    await this.app.init({ 
      background: "#000000", 
      resizeTo: window,
      antialias: false,
      resolution: window.devicePixelRatio || 1,
    })
    
    this.world.addChild(this.projectilesContainer)
    this.world.addChild(this.playersContainer)
    const controller = new Controller()
    const socket = getSocket()
    
    this.createStarfield()
    this.world.addChild(this.starfield)
    
    const bounds = new Graphics()
    bounds.rect(0, 0, WORLD_WIDTH, WORLD_HEIGHT).stroke({ width: 12, color: 0xff0000 })
    this.world.addChild(bounds)
    
    const texture  = await Assets.load(phase1)
    const player = new Sprite(texture)
    const playerWorldCords = { X: 0, Y: 0 }
    player.anchor.set(0.5)
    player.scale.set(0.3)
    player.x = this.app.screen.width  / 2
    player.y = this.app.screen.height / 2
    this.app.stage.addChild(player)
    
    this.world.addChild(this.playersContainer)
    
    this.app.stage.addChild(this.world)
    const otherPlayerSprites = new Map<string, [Sprite, Text]>()

    const userName = new Text({
      text: localUser.Username,
      style: new TextStyle({ fontSize: 10, fontFamily: "PixelUI", fill: localUser.Color}),
      resolution: 2,
    })
    this.app.stage.addChild(userName)

    this.app.ticker.add((t) => {
      if (userRegistered.isRegistered) {
        let dx = 0, dy = 0
        
        if (controller.keys.space.pressed ) {
          const now = performance.now();
          if (now - this.lastShotTime >= this.shootCooldown) {
            this.projectilePool.spawn(
              playerWorldCords.X, 
              playerWorldCords.Y, 
              player.angle,
              10,
              localUser.ID,
              localUser.Damage,
              1,
              1,
            )
            this.lastShotTime = now
          }
        }
        this.projectilePool.update(t.deltaTime)
        
        if (!uiHasFocus.isFocused){
          if (controller.keys.left.pressed)  dx -= 1
          if (controller.keys.right.pressed) dx += 1
          if (controller.keys.up.pressed)    dy -= 1
          if (controller.keys.down.pressed)  dy += 1
        }

        const speed = controller.checkDoubleTap.pressed ? 2 : 1
        if (dx !== 0 || dy !== 0) {
          player.angle = Math.atan2(dy, dx) * (180 / Math.PI) + 90
          playerWorldCords.X += dx * speed
          playerWorldCords.Y += dy * speed
        }

        if (socket && socket.readyState === WebSocket.OPEN) {
          const msg = createInputMsg(dx, dy, controller.checkDoubleTap.pressed, player.angle)
          socket.send(msg)
        }
        const cx = this.app.screen.width  / 2
        const cy = this.app.screen.height / 2
        player.position.set(cx, cy)

        for (const [id, [usr, x, y, a]] of Object.entries(playersState)) {
          if (id !== localUser.ID) {
            
            const entry = otherPlayerSprites.get(id)
            if (!entry) {
              const sprite = new Sprite(texture)
              sprite.anchor.set(0.5)
              sprite.scale.set(0.3)
              
              const textUser = new Text({
                text: usr,
                style: {fontSize: 12, fontFamily: "PixelUI", fill: randomBrightColor()},
                resolution: 2,
              })
              textUser.anchor.set(0.5, 1)
              this.playersContainer.addChild(sprite)
              this.playersContainer.addChild(textUser)
              otherPlayerSprites.set(id, [sprite, textUser])
            }

            const [sprite, textUser] = otherPlayerSprites.get(id)!
            sprite.position.set(x, y)
            sprite.angle = a
            textUser.position.set(x, y + 16)
        }}
        
        const [usr, sx, sy, spd, ang] = playersState[localUser.ID]
        playerWorldCords.X = sx
        playerWorldCords.Y = sy
        userName.x = cx - userName.width / 2
        userName.y = cy + 16
        player.angle = ang
        
        this.world.position.set(-playerWorldCords.X + cx, -playerWorldCords.Y + cy)
      }
    })
  }

  mount(el: HTMLElement) {
    el.appendChild(this.app.canvas)
    this.world.position.set(
      this.app.screen.width  / 2,
      this.app.screen.height / 2
    )
  }

  destroy() {
    this.app.destroy(true)
  }
}