<script lang="ts">
	import { getContext } from "svelte"
  import { messages, uiHasFocus, userRegistered } from "$lib/stores/ui.svelte"
	import { invalid } from "@sveltejs/kit"
	import { MAX_MESSAGE_LENGTH, MAX_USERNAME_LENGTH, MESSAGE_COOLDOWN_MS } from "$lib/Consts";
	import { localUser, playersState } from "$lib/stores/game.svelte";

  let hide: boolean    = $state(false)
  let message: string  = $state("")
  let lastTimeMessageSent: Date | null = null
  const socket: WebSocket = getContext('socket')
  
  function sendMessage() {
    if (!socket || socket.readyState !== WebSocket.OPEN) return invalid("Connection not ready")
    if (!userRegistered.isRegistered) return invalid("User is not authorized")
  
    const now = new Date()
    let trimmed = message.trim()
    if (!trimmed) return invalid("Message is empty")
    if (trimmed.length > MAX_MESSAGE_LENGTH) return invalid("Message is too long")
    if (lastTimeMessageSent && (now.getTime() - lastTimeMessageSent.getTime()) >= MESSAGE_COOLDOWN_MS) return invalid("Message can be sent after 15 seconds")

    try {
      const iso = now.toISOString()
      const formatted = iso.replace(/-/g, ':').replace('T', ' ').slice(0, 16)
      socket.send(JSON.stringify({ type: 'chat', data: `${localUser.Username}: ${trimmed} ${formatted}`}))
      lastTimeMessageSent = now
      message = ''
    } catch {
      invalid("Failed to send message")
    }
  }

  function createUser() {
    if (!socket || socket.readyState !== WebSocket.OPEN) return invalid("Connection not ready")

    let trimmed = localUser.Username.trim()
    if (trimmed in playersState) return invalid("Username already created")
    if (!trimmed) return invalid("Username is empty")
    if (localUser.Username.length > MAX_USERNAME_LENGTH) return invalid("Username is too big")
    
    try {
      socket.send(JSON.stringify({ type: 'user', data: localUser.Username}))
      userRegistered.isRegistered = true
    } catch {
      invalid("Failed to register user")
    }
  }

</script>

<div class="main">
  {#if !userRegistered.isRegistered}
    <div class="game-enter-screen">
      <div style:color="White">Enter username</div>
      <input type="text"
             onfocus={() => uiHasFocus.isFocused = true} 
             onblur={() => uiHasFocus.isFocused = false} 
             bind:value={localUser.Username}>
      <button style:font-size="1 rem" onclick={ createUser }>Enter</button>
    </div>
  {/if}
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
    font-size: 0.7em;
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