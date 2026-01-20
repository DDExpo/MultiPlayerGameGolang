
const keyMap = {
  KeyW:       "up",
  ArrowUp:    "up",
  KeyA:       "left",
  ArrowLeft:  "left",
  KeyS:       "down",
  ArrowDown:  "down",
  KeyD:       "right",
  ArrowRight: "right",
  Space:      "space",
  ShiftLeft:  "shift",
  ShiftRight: "shift"
}

export class Controller {
  keys: Record<string, { pressed: boolean, timestamp?: number }> = {
    up:        { pressed: false, timestamp: 0 },
    down:      { pressed: false, timestamp: 0 },
    left:      { pressed: false, timestamp: 0 },
    right:     { pressed: false, timestamp: 0 },
    shift:     { pressed: false, timestamp: 0 },
    space:     { pressed: false, timestamp: 0 },
  }
  checkDoubleTap = { pressed: false, lastPressedKey: ""}

  constructor() {
    window.addEventListener('keydown', (event) => this.keydownHandler(event))
    window.addEventListener('keyup',   (event) => this.keyupHandler(event))
  }

  keydownHandler(event: KeyboardEvent) {
    const key = keyMap[event.code as keyof typeof keyMap]
    if (!key) return
    const now = Date.now()
    if (now - this.keys[key].timestamp! < 200 || key == "shift") {
      this.checkDoubleTap.pressed = true
      this.checkDoubleTap.lastPressedKey = key
    }
    this.keys[key].pressed = true
  }

  keyupHandler(event: KeyboardEvent) {
    const key = keyMap[event.code as keyof typeof keyMap]
    if (!key) return
    if (key == this.checkDoubleTap.lastPressedKey) this.checkDoubleTap.pressed = false
 
    this.keys[key].pressed = false
    this.keys[key].timestamp = Date.now()
  }
}
