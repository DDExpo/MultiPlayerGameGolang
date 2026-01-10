<script lang="ts">
	import { getContext } from "svelte"
  import { messages, uiHasFocus } from "$lib/stores/ui.svelte";

  let hide: boolean = $state(false)
  let message: string = $state("")
  let username: string = $state("user")
  const socket: WebSocket = getContext('socket')
  
  function sendMessage() {
    if (socket && socket.readyState === WebSocket.OPEN) {
      if (message) {
        const now = new Date();
        const iso = now.toISOString();
        const formatted = iso.replace(/-/g, ':').replace('T', ' ').slice(0, 16);

        socket.send(JSON.stringify({ type: 'chat', data: `${username}: ${message} ${formatted}`}))
        message = ''
      }
    }
  }

</script>

<div class="main">
  <div class="game-enter-screen">
    <div style:color="White">Enter username</div>
    <input type="text"
           onfocus={() => uiHasFocus.isFocused = true} 
           onblur={() => uiHasFocus.isFocused = false} 
           bind:value={username}>
    <button style:font-size="1 rem">Enter</button>
  </div>
  <div class="chat-block">
    {#if !hide}
      <div class="messages-screen">
        {#each messages as msg} <p>{ msg }</p> {/each}
      </div>
      <div class="input-row">
        <input type="text"
               onfocus={() => uiHasFocus.isFocused = true}
               onblur={() => uiHasFocus.isFocused = false}
               bind:value={message}>
        <button class="btn-send-message" onclick={ sendMessage }>send</button>
      </div>
      <button class="btn-hide-chat" onclick={() => hide = !hide}>-</button>
    {:else}
      <button onclick={() => hide = !hide} style:max-width="45px">Enter</button>
    {/if}
  </div>

</div>

<style>
  .main {
    height: 100%;
    position: relative;
    background: transparent;
    justify-items: center;
    align-content: center;
  }
  .game-enter-screen {
    display: grid;
    justify-items: center;
    gap: 4px;
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
    color: white;
    text-overflow: clip;
    white-space: normal;
    overflow: auto;
    background-color: rgba(169, 169, 169, 0.137);
  }
  .input-row {
    display: flex;
    width: 100%;
    gap: 4px;
  }
  .input-row input {
    flex: 1;
  }
  .btn-send-message {
    font-size: 0.5em;
  }
  .btn-hide-chat {
    position: absolute;
    top: -20px;
    right: -20px;
  }

</style>