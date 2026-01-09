<script lang="ts">
	import favicon from '$lib/assets/favicon.svg'
	import { onMount, setContext } from 'svelte'
	import { messages } from '$lib/stores/messages.svelte';
	import '../app.css'

	let { children } = $props()
  const socket: WebSocket = new WebSocket("ws://localhost:8000/ws")
  setContext('socket', socket)
  
  onMount(async() => {
    socket.onopen = () => {
      console.log("Connected to server")
    }
    socket.onmessage = (event) => {
      const msg = JSON.parse(event.data) as WSMessage
      switch (msg.type) {
        case "chat":
          messages.push(msg.data)
          break;
        default:
          break;
      }
    }
    socket.onclose = () => {
      console.log("Disconnected")
    }
    socket.onerror = (err) => {
      console.error("Socket error:", err)
    }
  })
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
