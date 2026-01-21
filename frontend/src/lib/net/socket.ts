import { messages, userUiState } from "$lib/stores/ui.svelte"
import { ClientData, projectileQueue, playersState } from "$lib/stores/game.svelte"
import { MsgType, StateType } from "$lib/Consts"
import { UserStateDelta } from "$lib/types/enums"
import { projectilePool } from "$lib/game/Projectiles"
import { deserializeCombat, deserializePosition, deserializeWeapon,
         initReader, readString, readUint16, readUint8 } from "./deserialzie"

let socket: WebSocket | null = null

function createSocket(url: string) {
  if (socket && socket.readyState === WebSocket.OPEN) return socket
  socket = new WebSocket(url)

  socket.onopen = () => { console.log("Connected")}
  
 socket.onmessage = (e) => {
  if (!(e.data instanceof ArrayBuffer)) {
    console.error("Expected binary message")
    return
  }

  const view = new DataView(e.data)
  initReader(view, e.data)

  const msgType = readUint8()

  switch (msgType) {

    case MsgType.USER_CHAT: {
      const username  = readString()
      const text      = readString()
      const timestamp = readString()
      const color     = readString()

      messages.push([color, [`${username}: ${text}`, timestamp]])
      break
    }

    case MsgType.USER_SHOOT_STATUS: {
      const alive = readUint8()
      const id    = readUint16()
      if (!alive && projectilePool) projectilePool.destroy(id)
      break
    }

    case MsgType.USER_STATE: {
      const stateType = readUint8()
      const username  = readString()
      const player    = playersState[username] ??= { movement: [4000, 4000, 0], combat: { dead: false }}

      switch (stateType) {
        case StateType.USER_REG: {
          const [x, y, angle] = deserializePosition()

          player.movement = [x, y, angle]
          player.combat.dead = false

          if (username === ClientData.Username) {
            deserializeCombat()
            deserializeWeapon()
          }

          break
        }

        // 🔹 death event
        case StateType.USER_DEAD: {
          player.combat.dead = true

          if (username === ClientData.Username) {
            userUiState.alive   = false
            userUiState.focused = true
          }
          break
        }

        case StateType.USER_CURRENT_STATE: {
          const deltaMask = readUint8()

          if (deltaMask & UserStateDelta.POS) {
            const [x, y, angle] = deserializePosition()
            player.movement = [x, y, angle]
          }

          if (username === ClientData.Username) {
            if (deltaMask & UserStateDelta.STATS)  deserializeCombat()
            if (deltaMask & UserStateDelta.WEAPON) deserializeWeapon()
          }

          break
        }
        case StateType.USER_PRESSED_SHOOT: {
          if (username !== ClientData.Username) {
            const [x, y, angle] = deserializePosition()
            const weaponWidth  = readUint8()
            const weaponRange  = readUint16()
            projectileQueue.push([username, x, y, angle, weaponWidth, weaponRange])
          }
          break
        }
      }

      break
    }

    default:
      console.error("Unknown message type:", msgType)
  }
}
  
  socket.onclose = e => {
    console.log('WS close', {
      code: e.code,
      reason: e.reason,
      wasClean: e.wasClean
    })
    
    delete playersState[ClientData.Username]
  }

  socket.onerror = (e) => { console.error("Socket error", e)}

  return socket
}

export function initSocket(): WebSocket {
  socket = createSocket("ws://localhost:8000/ws")
  socket.binaryType = "arraybuffer"
  return socket
}

export function getSocket(): WebSocket | null {
  return socket
}

export function closeSocket() {
  socket?.close()
  socket = null
}

export function waitForOpen(socket: WebSocket): Promise<void> {
  if (socket.readyState === WebSocket.OPEN) return Promise.resolve()

  return new Promise((resolve, reject) => {
    const onOpen = () => {
      socket.removeEventListener('open', onOpen)
      resolve()
    }
    const onClose = () => {
      socket.removeEventListener('open', onOpen)
      reject(new Error('Socket closed before opening'))
    }

    socket.addEventListener('open',  onOpen)
    socket.addEventListener('close', onClose)
  })
}