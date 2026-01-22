import { Application, Assets, Text, Container, BlurFilter, Graphics, Sprite, TextStyle } from "pixi.js"
import { Controller } from "./Controllers"
import phase1 from "$lib/assets/players/phase1.png"
import { ClientData, playersState, projectileQueue } from "$lib/stores/game.svelte"
import { getSocket } from "$lib/net/socket"
import { PADDING_USERNAME, StateType, WORLD_HEIGHT, WORLD_WIDTH } from "$lib/Consts"
import { randomBrightColor } from "$lib/utils"
import { initProjectilePool, projectilePool } from "./Projectiles"
import { serializeInputMsg } from "$lib/net/serialize"
import { userUiState } from "$lib/stores/ui.svelte"


export class Game {
  app:                  Application
  world:                Container
  playersContainer:     Container
  starfield:            Container
  projectilesContainer: Container
  userName:             Text

  screenCenterX:        number
  screenCenterY:        number

  lastShotTime  = 0
  shootCooldown = 600

  constructor() {
    this.app                  = new Application()
    this.world                = new Container()
    this.playersContainer     = new Container()
    this.projectilesContainer = new Container()
    this.starfield            = new Container()
    this.userName             = new Text
    
    this.screenCenterX        = 0
    this.screenCenterY        = 0

    initProjectilePool(this.projectilesContainer)
  }
  
  async init() {
    await this.app.init({ 
      background: "#000000", 
      resizeTo: window,
      antialias: false,
      resolution: window.devicePixelRatio || 1,
    })
    this.screenCenterX = this.app.screen.width  / 2
    this.screenCenterY = this.app.screen.height / 2

    this.playersContainer.enableRenderGroup()
    this.projectilesContainer.enableRenderGroup()
    
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
    player.x = this.screenCenterX
    player.y = this.screenCenterY
    
    this.app.stage.addChild(player)
    this.app.stage.addChild(this.world)
    
    this.userName.anchor.set(0.5, 1)
    this.world.addChild(this.userName)
    
    this.app.ticker.add((t) => {
      const socket = getSocket()
      const socketReady = socket?.readyState === WebSocket.OPEN
      
      for (const [usr, state] of Object.entries(playersState)) {
        const { movement, combat } = state
        const [x, y, ang]          = movement
        const visible              = !combat.dead
        
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
          
          if (sprite.visible !== visible) {
            sprite.visible = visible
            textUser.visible = visible
          }

          if (!visible) continue
          
          sprite.position.set(x, y)
          sprite.angle = ang
          textUser.position.set(x, y + PADDING_USERNAME)
        }
      }
      for (const [usr, x, y, angle, ww, wr] of projectileQueue) {
        projectilePool!.spawn(x, y, angle, usr, ww, wr)
      }
      projectileQueue.length = 0
      projectilePool!.update()
      
      
      if (!userUiState.registered || userUiState.focused) return
      
      if (!this.userName.text) { this.setUsernameTextStyle() }


      let lastDx = 0, lastDy = 0
      
      const dx = (controller.keys.right.pressed ? 1 : 0) - (controller.keys.left.pressed ? 1 : 0)
      const dy = (controller.keys.down.pressed  ? 1 : 0) - (controller.keys.up.pressed   ? 1 : 0)
      
      if (dx !== lastDx || dy !== lastDy) {
        player.angle = Math.atan2(dy, dx) * (180 / Math.PI) + 90
        playerWorldCords.X = dx
        playerWorldCords.Y = dy
        lastDx = dx
        lastDy = dy
      }
      if (socketReady) { 
        const movements = serializeInputMsg(dx, dy, controller.checkDoubleTap.pressed, player.angle)
        socket.send(movements)
      }
      player.position.set(this.screenCenterX, this.screenCenterY)
      
      const { movement, combat } = playersState[ClientData.Username]
      const [sx, sy, ang] = movement
      const visible = !combat.dead
      
      if (player.visible !== visible) {
        player.visible = visible
        this.userName.visible = visible
      }

      if (!visible) return

      playerWorldCords.X = sx
      playerWorldCords.Y = sy
      player.angle = ang
      this.userName.position.set(sx, sy + PADDING_USERNAME)
        
      const newWorldX = -playerWorldCords.X + this.screenCenterX
      const newWorldY = -playerWorldCords.Y + this.screenCenterY

      if (this.world.x !== newWorldX || this.world.y !== newWorldY) {
        this.world.position.set(newWorldX, newWorldY)
      }
      
      if (controller.keys.space.pressed && !this.isSafeZone(playerWorldCords.X, playerWorldCords.Y)) {
        const now = performance.now()
        if (now - this.lastShotTime >= this.shootCooldown) {
          const projectileId = projectilePool!.spawn(
            playerWorldCords.X, playerWorldCords.Y, player.angle,
            ClientData.Username, ClientData.WeaponWidth, ClientData.WeaponRange,
          )
          if (socketReady) { 
            const buffer = new ArrayBuffer(3)
            const view = new DataView(buffer)
            view.setUint8(0, StateType.USER_PRESSED_SHOOT)
            view.setUint16(1, projectileId, true)
            socket.send(buffer)
          }
          this.lastShotTime = now
      }}
    })

    this.app.renderer.on('resize', () => {
      this.screenCenterX = this.app.screen.width  / 2
      this.screenCenterY = this.app.screen.height / 2
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
  
  public mount(el: HTMLElement) {
    el.appendChild(this.app.canvas)
    this.world.position.set(this.screenCenterX, this.screenCenterY)
  }
  
  public setUsernameTextStyle() {
    this.userName.text = ClientData.Username
    this.userName.style = new TextStyle({ fontSize: 10, fontFamily: "PixelUI", fill: ClientData.Color})
    this.userName.resolution = 2
  }
  
  public destroy() {
    this.app.destroy(true)
  }

  private isSafeZone(x: number, y: number): boolean {
    return x >= 3590 && x <= 4410 && y >= 3590 && y <= 4410
  } 
}