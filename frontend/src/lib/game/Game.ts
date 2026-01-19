import { Application, Assets, Text, Container, BlurFilter, Graphics, Sprite, TextStyle } from "pixi.js"
import { Controller } from "./Controllers"
import phase1 from "$lib/assets/players/phase1.png"
import { uiHasFocus, userRegistered } from "$lib/stores/ui.svelte"
import { ClientData, otherProjectiles, playersState } from "$lib/stores/game.svelte"
import { getSocket } from "$lib/net/socket"
import { MsgType, PADDING_USERNAME, WORLD_HEIGHT, WORLD_WIDTH } from "$lib/Consts"
import { randomBrightColor } from "$lib/utils"
import { ProjectilePool } from "./Projectiles"
import { serializeInputMsg } from "$lib/net/serialize"


export class Game {
  app:                  Application
  world:                Container
  playersContainer:     Container
  starfield:            Container
  projectilesContainer: Container
  projectilePool:       ProjectilePool

  lastShotTime = 0;
  shootCooldown = 600;

  constructor() {
    this.app                  = new Application()
    this.world                = new Container()
    this.playersContainer     = new Container()
    this.projectilesContainer = new Container()
    this.starfield            = new Container()
    this.projectilePool       = new ProjectilePool(this.projectilesContainer)
  }
  
  async init() {
    await this.app.init({ 
      background: "#000000", 
      resizeTo: window,
      antialias: false,
      resolution: window.devicePixelRatio || 1,
    })
    const bounds = new Graphics()
    bounds.rect(0, 0, WORLD_WIDTH, WORLD_HEIGHT).stroke({ width: 12, color: 0xff0000 })
    this.world.addChild(bounds)
  
    const safeZone = new Graphics()
    safeZone.rect(3600, 3600, 800, 800).stroke({ width: 6, color: 0x1fff01 })
    this.world.addChild(safeZone)
    
    const safeZoneText = new Text({
      text: "SAFE ZONE",
      style: new TextStyle({ fontSize: 50, fontFamily: "PixelUI", fill: 0x1fff01 }),
      resolution: 2,
    })
    safeZoneText.position.set(3780, 4420)
    this.world.addChild(safeZoneText)

    this.createStarfield()
    this.world.addChild(this.starfield)
    this.world.addChild(this.projectilesContainer)
    this.world.addChild(this.playersContainer)

    const controller = new Controller()
    const otherPlayerSprites = new Map<string, [Sprite, Text]>()
    
    const texture = await Assets.load(phase1)
    const player  = new Sprite(texture)
    const playerWorldCords = { X: 0, Y: 0 }
    player.anchor.set(0.5)
    player.scale.set(0.3)
    let cx = this.app.screen.width  / 2
    let cy = this.app.screen.height / 2 
    player.x = cx
    player.y = cy

    this.app.stage.addChild(player)
    
    this.app.stage.addChild(this.world)
    
    const userName = new Text({
      text: ClientData.Username,
      style: new TextStyle({ fontSize: 10, fontFamily: "PixelUI", fill: ClientData.Color}),
      resolution: 2,
    })
    userName.anchor.set(0.5, 1)
    let userNameCreated = false
    
    this.app.ticker.add((t) => {
      const socket = getSocket()
      
      if (userRegistered.isRegistered && !uiHasFocus.isFocused) {
        if (!userNameCreated) {
          this.world.addChild(userName)
          userNameCreated = true
        }
        let dx = 0, dy = 0
        
        if (controller.keys.left.pressed)  dx -= 1
        if (controller.keys.right.pressed) dx += 1
        if (controller.keys.up.pressed)    dy -= 1
        if (controller.keys.down.pressed)  dy += 1
        
        if (dx != 0 || dy != 0) {
          player.angle = Math.atan2(dy, dx) * (180 / Math.PI) + 90
          playerWorldCords.X = dx
          playerWorldCords.Y = dy
        }
        
        if (socket && socket.readyState === WebSocket.OPEN) { 
          const movements = serializeInputMsg(dx, dy, controller.checkDoubleTap.pressed, player.angle)
          socket.send(movements)
        }
        cx = this.app.screen.width  / 2
        cy = this.app.screen.height / 2
        player.position.set(cx, cy)
     
        const [sx, sy, ang] = playersState[ClientData.Username]
        
        playerWorldCords.X = sx
        playerWorldCords.Y = sy
        userName.position.set(sx, sy + PADDING_USERNAME)
        player.angle = ang        
        this.world.position.set(-playerWorldCords.X + cx, -playerWorldCords.Y + cy)
        
        if (controller.keys.space.pressed && !this.isSafeZone(playerWorldCords.X, playerWorldCords.Y)) {
          const now = performance.now()
          if (now - this.lastShotTime >= this.shootCooldown) {
            if (socket && socket.readyState === WebSocket.OPEN) { 
              const buffer = new ArrayBuffer(1)
              const view = new DataView(buffer)
              view.setUint8(0, MsgType.USER_SHOOT)
              socket.send(buffer)
            }
            this.projectilePool.spawn(
              playerWorldCords.X, playerWorldCords.Y, player.angle,
              ClientData.Username, ClientData.Damage, ClientData.WeaponWidth, ClientData.WeaponRange,
            )
            this.lastShotTime = now
          }}
      }
      
      for (const [usr, [x, y, ang]] of Object.entries(playersState)) {
        if (usr !== ClientData.Username) {
          
          const entry = otherPlayerSprites.get(usr)
          if (!entry) {
            const sprite = new Sprite(texture)
            sprite.anchor.set(0.5)
            sprite.scale.set(0.3)
            
            const textUser = new Text({
              text: usr,
              style: {fontSize: 10, fontFamily: "PixelUI", fill: randomBrightColor()},
              resolution: 2,
            })
            textUser.anchor.set(0.5, 1)
            this.playersContainer.addChild(sprite)
            this.playersContainer.addChild(textUser)
            otherPlayerSprites.set(usr, [sprite, textUser])
          }
          
          const [sprite, textUser] = otherPlayerSprites.get(usr)!
          sprite.position.set(x, y)
          sprite.angle = ang
          textUser.position.set(x, y + PADDING_USERNAME)
        }}
        for (const [usr, [x, y, a, d, ww, wr]] of Object.entries(otherProjectiles)) {
          this.projectilePool.spawn(x, y, a, usr, d, ww, wr)
          delete otherProjectiles[usr]
        }

    this.projectilePool.update(t.deltaTime)
    })
  }
  
  createStarfield() {
    if (this.starfield.children.length > 0) return
    
    const starCount = 2000
    const stars     = new Graphics()
    
    for (let i = 0; i < starCount; i++) {
      const size  = Math.random() * 2 + 1
      const x     = Math.random() * WORLD_WIDTH
      const y     = Math.random() * WORLD_HEIGHT
      const color = randomBrightColor()
      const alpha = Math.random() * 0.6 + 0.3
      
      stars.circle(x, y, size).fill({ color, alpha })
    }
  
    stars.filters = [new BlurFilter({ strength: 0.5, quality: 2 })]
    
    this.starfield.addChild(stars)
    this.starfield.enableRenderGroup()
  }
  
  mount(el: HTMLElement) {
    el.appendChild(this.app.canvas)
    this.world.position.set(
      this.app.screen.width  / 2,
      this.app.screen.height / 2
    )
  }

  isSafeZone(x: number, y: number): boolean {
    return x >= 3590 && x <= 4410 && y >= 3590 && y <= 4410
  } 

  destroy() {
    this.app.destroy(true)
  }
}