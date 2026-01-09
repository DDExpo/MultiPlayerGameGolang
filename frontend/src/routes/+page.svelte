<script lang="ts">
	import { getContext } from "svelte"
  import { messages } from "$lib/stores/messages.svelte";

  let hide: boolean = $state(false)
  let message: string = $state("")
  const socket: WebSocket = getContext('socket')
  
  function sendMessage() {
    if (socket && socket.readyState === WebSocket.OPEN) {
      if (message) {
        const now = new Date();
        const iso = now.toISOString();
        const formatted = iso.replace(/-/g, ':').replace('T', ' ').slice(0, 16);

        socket.send(JSON.stringify({ type: "chat", data: "user:" + " " + message + " " + formatted }))
        message = ""
      }
    }
  }

</script>

<div class="main">
  <div class="game-enter-screen">
    <button>Enter</button>
  </div>
  <div class="chat-block">
    {#if !hide}
      <div class="messages-screen">
        {#each messages as msg} <p>{ msg }</p> {/each}
      </div>
      <div class="input-row">
        <input type="text" bind:value={message}>
        <button class="btn-send-message" onclick={ sendMessage }>Enter</button>
      </div>
      <button class="btn-hide-chat" onclick={() => hide = !hide}>-</button>
    {:else}
      <button onclick={() => hide = !hide} style:width="45px">Enter</button>
    {/if}
  </div>

</div>

<style>
  .main {
    height: 100%;
    position: relative;
    justify-items: center;
    align-content: center;
  }
  .chat-block {
    max-width: 240px;
    min-width: 240px;
    left: 0;
    bottom: 0;
    margin: 8px;
    position: absolute;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .messages-screen {
    max-width: 232px;
    min-height: 100px;
    padding: 4px;
  
    text-overflow: clip;
    white-space: normal;
    overflow: auto;
    background-color: darkgray;
  }
  .input-row {
    display: flex;
    width: 100%;
    gap: 4px;
  }
  .input-row input {
    flex: 1;
  }

  .btn-hide-chat {
    position: absolute;
    top: -20px;
    right: -20px;
  }

</style>