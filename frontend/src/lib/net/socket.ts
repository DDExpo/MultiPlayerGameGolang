import { messages } from "$lib/stores/ui.svelte"
import { playersState } from "$lib/stores/game.svelte"

export function createSocket(url: string) {
  const socket = new WebSocket(url)

  socket.onopen = () => {
    console.log("Connected")
  }

  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data) as WSMessage
    if (msg.type === "userCords") {
      const { username, x, y, angle } = msg.data as PlayerMoveMessage
      playersState[username] = [x, y, angle]
    }
    else if (msg.type === "chat") {
      messages.push(msg.data)
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

export const socket = createSocket("ws://localhost:8000/ws")
