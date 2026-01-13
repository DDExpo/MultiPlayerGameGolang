import { messages } from "$lib/stores/ui.svelte"
import { localUser, playersState } from "$lib/stores/game.svelte"
import { MsgType } from "$lib/Consts"
import { readFloat32, readString } from "./binaryEncodingDecoding"

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
        const username = readString(view, buffer, off)
        const x        = readFloat32(view, off)
        const y        = readFloat32(view, off)
        const angle    = readFloat32(view, off)

        playersState[username] = [x, y, angle]
        break
      }
      case MsgType.CHAT: {
        const username  = readString(view, buffer, off)
        const text      = readString(view, buffer, off)
        const timestamp = readString(view, buffer, off)
        messages.push(`${username}: ${text} ${timestamp}`)
        break
      }

      case MsgType.USER: {
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
  const socket = createSocket("ws://localhost:8000/ws")
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