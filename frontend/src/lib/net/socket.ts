import { messages } from "$lib/stores/ui.svelte"
import { ClientData, playersState } from "$lib/stores/game.svelte"
import { MsgType } from "$lib/Consts"
import { UserStateDelta } from "$lib/types/enums"
import { readFloat32, readString } from "./deserialzie"

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
      const buffer = e.data
      const view = new DataView(buffer)

      const off = { v: 0 }
      const msgType = view.getUint8(off.v++)

      switch (msgType) {
        case MsgType.USER_STATE: {
        const deltaMask = view.getUint8(off.v++)
        const username  = readString(view, buffer, off)
        let x, y, speed, angle

        if (deltaMask & UserStateDelta.POS) { 
          x     = readFloat32(view, off)
          y     = readFloat32(view, off)
          speed = view.getUint8(off.v++)
          angle = readFloat32(view, off)
          playersState[username] = [x, y, speed, angle]
          console.log(playersState)
        }

        if (username === ClientData.Username) {
          if (deltaMask & UserStateDelta.STATS) {
            ClientData.Hp     = view.getUint8(off.v++)
            ClientData.Kills  = view.getUint8(off.v++)
            ClientData.Damage = view.getUint8(off.v++)
          }

          if (deltaMask & UserStateDelta.WEAPON) {
            ClientData.WeaponType  = view.getUint8(off.v++)
            ClientData.WeaponWidth = view.getUint8(off.v++)
            ClientData.WeaponRange = view.getUint8(off.v++)
          }
        }
        break
      }
      case MsgType.USER_CHAT: {
        const username  = readString(view, buffer, off)
        const text      = readString(view, buffer, off)
        const timestamp = readString(view, buffer, off)
        const color     = readString(view, buffer, off)
        messages.push([color, [`${username}: ${text}`,`${timestamp}`]])
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
  }
  socket.onerror = (e) => {
    console.error("Socket error", e)
  }

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