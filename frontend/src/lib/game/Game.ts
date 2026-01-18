import { Application, Assets, Text, Container, BlurFilter, Graphics, Sprite, TextStyle } from "pixi.js"
import { Controller } from "./Controllers"
import phase1 from "$lib/assets/players/phase1.png"
import { uiHasFocus, userRegistered } from "$lib/stores/ui.svelte"
import { ClientData, playersState } from "$lib/stores/game.svelte"
import { getSocket } from "$lib/net/socket"
import { WORLD_HEIGHT, WORLD_WIDTH } from "$lib/Consts"
import { randomBrightColor } from "$lib/utils"
import { ProjectilePool } from "./Projectiles"
import { serializeInputMsg } from "$lib/net/serialize"


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
    const otherPlayerSprites = new Map<string, [Sprite, Text]>()
    
    this.createStarfield()
    this.world.addChild(this.starfield)
    
    const bounds = new Graphics()
    bounds.rect(0, 0, WORLD_WIDTH, WORLD_HEIGHT).stroke({ width: 12, color: 0xff0000 })
    this.world.addChild(bounds)
    
    const texture = await Assets.load(phase1)
    const player  = new Sprite(texture)
    const playerWorldCords = { X: 0, Y: 0 }
    player.anchor.set(0.5)
    player.scale.set(0.3)
    player.x = this.app.screen.width  / 2
    player.y = this.app.screen.height / 2
    this.app.stage.addChild(player)

    this.app.stage.addChild(this.world)
    
    const userName = new Text({
      text: ClientData.Username,
      style: new TextStyle({ fontSize: 10, fontFamily: "PixelUI", fill: ClientData.Color}),
      resolution: 2,
    })
    this.world.addChild(userName)
    
    this.app.ticker.add((t) => {
      if (userRegistered.isRegistered) {
        const socket = getSocket()
        let dx = 0, dy = 0
        
        if (controller.keys.space.pressed ) {
          const now = performance.now();
          if (now - this.lastShotTime >= this.shootCooldown) {
            this.projectilePool.spawn(
              playerWorldCords.X, 
              playerWorldCords.Y, 
              player.angle,
              10,
              ClientData.Username,
              ClientData.Damage,
              ClientData.WeaponRange,
              ClientData.WeaponWidth,
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

        player.angle = Math.atan2(dy, dx) * (180 / Math.PI) + 90

        if (socket && socket.readyState === WebSocket.OPEN) {
          const msg = serializeInputMsg(dx, dy, controller.checkDoubleTap.pressed, player.angle)
          socket.send(msg)
        }
        const cx = this.app.screen.width  / 2
        const cy = this.app.screen.height / 2
        player.position.set(cx, cy)

        for (const [usr, [x, y, a]] of Object.entries(playersState)) {
          if (usr !== ClientData.Username) {
            
            const entry = otherPlayerSprites.get(usr)
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
              otherPlayerSprites.set(usr, [sprite, textUser])
            }

            const [sprite, textUser] = otherPlayerSprites.get(usr)!
            sprite.position.set(x, y)
            sprite.angle = a
            textUser.position.set(x, y + 16)
        }}
        if (!playersState) return
        const [sx, sy, spd, ang] = playersState[ClientData.Username]

        playerWorldCords.X = sx
        playerWorldCords.Y = sy
        userName.x = sx - userName.width / 2
        userName.y = sy + 16
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