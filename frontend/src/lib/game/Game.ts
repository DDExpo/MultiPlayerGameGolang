import { Application, Assets, Container, Sprite } from "pixi.js"
import { Controller } from "./Controllers"
import phase1 from "$lib/assets/players/phase1.png"
import { uiHasFocus } from "$lib/stores/ui.svelte"

export class Game {
  app:   Application
  world: Container

  constructor() {
    this.app = new Application()
    this.world = new Container()
  }

  async init() {
    await this.app.init({
      background: "#000000",
      resizeTo: window,
    })

    const texture          = await Assets.load(phase1)
    const player           = new Sprite(texture)
    const ccc              = new Sprite(texture)
    const playerWorldCords = { worldX: 0, worldY: 0 }
    const screenCenterX    = this.app.renderer.width  / 2
    const screenCenterY    = this.app.renderer.height / 2
  
    ccc.anchor.set(0.3)
    player.anchor.set(0.5)
    player.scale.set(0.3)
    player.x = screenCenterX
    player.y = screenCenterY

    this.app.stage.addChild(this.world)
    this.app.stage.addChild(player)
    this.world.addChild(ccc)

    const controller = new Controller()
    
    this.app.ticker.add((t) => {
      let dx = 0, dy = 0
      
      if (controller.keys.space.pressed) {

      }
      if (!uiHasFocus.isFocused){
         if (controller.keys.left.pressed) dx -= 1
         if (controller.keys.right.pressed) dx += 1
         if (controller.keys.up.pressed) dy -= 1
         if (controller.keys.down.pressed) dy += 1
      }

      if (dx !== 0 || dy !== 0) {
        player.angle = Math.atan2(dy, dx) * (180 / Math.PI) + 90
        
        const speed = controller.checkDoubleTap.pressed ? 2 : 1

        playerWorldCords.worldX += dx * speed;
        playerWorldCords.worldY += dy * speed;

        this.world.x = -playerWorldCords.worldX + screenCenterX;
        this.world.y = -playerWorldCords.worldY + screenCenterY;
      }
    })
  }

  mount(el: HTMLElement) {
    el.appendChild(this.app.canvas)
    this.world.position.set(
      this.app.screen.width / 2,
      this.app.screen.height / 2
    )
  }

  destroy() {
    this.app.destroy(true)
  }
}