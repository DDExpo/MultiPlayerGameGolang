import { messages } from "$lib/stores/ui.svelte"

export function createSocket(url: string) {
  const socket = new WebSocket(url)

  socket.onopen = () => {
    console.log("Connected")
  }

  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data)

    if (msg.type === "chat") {
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