import { Application, Assets, Container, Sprite, Texture } from "pixi.js"
import { Controller } from "./Controllers"
import phase1 from "$lib/assets/players/phase1.png"
import { uiHasFocus, userRegistered } from "$lib/stores/ui.svelte"
import { localUser, playersState } from "$lib/stores/game.svelte"
import { socket } from '$lib/net/socket'


export class Game {
  app:   Application
  world: Container
  playersContainer: Container

  constructor() {
    this.app              = new Application()
    this.world            = new Container()
    this.playersContainer = new Container()
  }

async init() {
  await this.app.init({
    background: "#000000",
    resizeTo: window,
  })

  const texture  = await Assets.load(phase1)
  const ccc      = new Sprite(texture)
  ccc.anchor.set(0.3)

  const player           = new Sprite(texture)
  player.anchor.set(0.5)
  player.scale.set(0.3)
  player.x = this.app.screen.width  / 2
  player.y = this.app.screen.height / 2
  const playerWorldCords = { worldX: 0, worldY: 0 }

  this.app.stage.addChild(this.world)
  this.app.stage.addChild(player)

  this.world.addChild(this.playersContainer)
  this.world.addChild(ccc)

  const controller = new Controller()
  const otherPlayerSprites = new Map<string, Sprite>()

  this.app.ticker.add((t) => {
    if (userRegistered.isRegistered) {
      let dx = 0, dy = 0
      
      if (controller.keys.space.pressed) {}
      
      if (!uiHasFocus.isFocused){
        if (controller.keys.left.pressed)  dx -= 1
        if (controller.keys.right.pressed) dx += 1
        if (controller.keys.up.pressed)    dy -= 1
        if (controller.keys.down.pressed)  dy += 1
      }

      if (dx !== 0 || dy !== 0) {
        player.angle = Math.atan2(dy, dx) * (180 / Math.PI) + 90
        const speed = controller.checkDoubleTap.pressed ? 2 : 1
        playerWorldCords.worldX += dx * speed
        playerWorldCords.worldY += dy * speed
      }

      const cx = this.app.screen.width  / 2
      const cy = this.app.screen.height / 2
      player.position.set(cx, cy)
      
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'userCords', data: {
          username: localUser.Username,
          x: playerWorldCords.worldX,
          y: playerWorldCords.worldY,
          angle: player.angle
        }}))
      }

      const currentPlayers = new Set<string>()
      
      for (const [username, [x, y, a]] of Object.entries(playersState)) {
        if (username !== localUser.Username) {
          currentPlayers.add(username)
          
          let sprite = otherPlayerSprites.get(username)
          if (!sprite) {
            sprite = createPlayer(texture)
            otherPlayerSprites.set(username, sprite)
            this.playersContainer.addChild(sprite)
          }
          sprite.position.set(x, y)
          sprite.angle = a
        }
      }

      for (const [username, sprite] of otherPlayerSprites.entries()) {
        if (!currentPlayers.has(username)) {
          this.playersContainer.removeChild(sprite)
          sprite.destroy()
          otherPlayerSprites.delete(username)
        }
      }
      this.world.position.set(-playerWorldCords.worldX + cx, -playerWorldCords.worldY + cy)
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

function createPlayer(texture: Texture): Sprite {
  const sprite = new Sprite(texture)
  sprite.anchor.set(0.5)
  sprite.scale.set(0.3)
  return sprite
}

