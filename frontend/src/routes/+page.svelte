<script lang="ts">
  import { messages, uiHasFocus, userRegistered } from "$lib/stores/ui.svelte"
  import { MAX_MESSAGE_LENGTH, MAX_USERNAME_LENGTH, MESSAGE_COOLDOWN_MS } from "$lib/Consts";
  import { localUser, playersState } from "$lib/stores/game.svelte";
  import { createBinaryChatMsg, createBinaryUserMsg } from "$lib/net/binaryEncodingDecoding";
  import { getSocket, initSocket } from "$lib/net/socket";


  let socket: WebSocket | null
  let hide: boolean                    = $state(false)
  let message: string                  = $state("")
  let error: string                    = $state("")
  let place: string                    = $state("")
  let errorTimeout: number | null      = null
  let lastTimeMessageSent: Date | null = null


  function showError(msg: string, pl: string) {
    error = msg
    place = pl
    if (errorTimeout) clearTimeout(errorTimeout)
    errorTimeout = window.setTimeout(() => {
      error = ""
      place = ""
      errorTimeout = null
    }, 3000)
  }

  function sendMessage() {
    socket = getSocket()
    if (!socket || socket.readyState !== WebSocket.OPEN) { showError('Connection not ready',   'chat'); return }
    if (!userRegistered.isRegistered)                    { showError('User is not authorized', 'chat'); return }

    const now = new Date()
    const trimmed = message.trim()
    if (!trimmed)                            { showError('Message is empty', 'chat'); return }
    if (trimmed.length > MAX_MESSAGE_LENGTH) { showError(`Message is too long (max ${MAX_MESSAGE_LENGTH} characters)`, 'chat'); return}
    if (lastTimeMessageSent && (now.getTime() - lastTimeMessageSent.getTime()) < MESSAGE_COOLDOWN_MS) {
      showError('Please wait 15 seconds between messages', 'chat')
      return
    }

    try {
      const iso = now.toISOString()
      const formatted = iso.replace(/-/g, ':').replace('T', ' ').slice(0, 16)
      const msg = createBinaryChatMsg(localUser.Username, trimmed, formatted)
      socket.send(msg)
      lastTimeMessageSent = now
      message = ""
    } catch (err) { console.error(`Failed: ${err}`)}
  }

  async function createUser() {
    const trimmed = localUser.Username.trim()
    
    if (!trimmed)                             { showError('Username is empty', 'login'); return }
    if (trimmed.length > MAX_USERNAME_LENGTH) { showError(`Username is too long (max ${MAX_USERNAME_LENGTH} characters)`, 'login'); return }
    if (trimmed in playersState)              { return }
    
    localUser.Username = trimmed

    try {
      const msg = createBinaryUserMsg(localUser.Username)
      const res = await fetch("http://localhost:8000/session", {
        method:      "POST",
        credentials: "include",
        headers:   { "Content-Type": "application/json" },
        body:        JSON.stringify({ username: trimmed }),
      })

      if (!res.ok) { console.log("Session creation failed"); return }
      
      socket = initSocket()
      await waitForOpen(socket)
      socket.send(msg)
      userRegistered.isRegistered = true
    
    } catch (err) { console.error(`Failed: ${err}`)}
  }

  function waitForOpen(socket: WebSocket): Promise<void> {
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
</script>


<div class="main">
  {#if !userRegistered.isRegistered}
    <div class="game-enter-screen">
      <div style:color="White">Enter username</div>
      <input type="text"
             onfocus={() => uiHasFocus.isFocused = true} 
             onblur={() => uiHasFocus.isFocused = false} 
             bind:value={localUser.Username}>
      <button style:font-size="1rem" onclick={createUser}>Enter</button>
      {#if place === "login"}
        <div class="error-message">{error}</div>
      {/if}
    </div>
  {/if}
  <div class="chat-block">
    {#if !hide}
      <div class="messages-screen">
        {#each messages as msg} <p>{msg}</p> {/each}
      </div>
      <div class="input-row">
        <input type="text"
               onfocus={() => uiHasFocus.isFocused = true}
               onblur={() => uiHasFocus.isFocused = false}
               bind:value={message}
               onkeydown={(e) => e.key === 'Enter' && sendMessage()}>
        <button class="btn-send-message" onclick={sendMessage}>send</button>
      </div>
      {#if place === "chat"}
        <div class="error-message">{error}</div>
      {/if}
      <button class="btn-hide-chat" onclick={() => hide = !hide}>-</button>
    {:else}
      <button onclick={() => hide = !hide} style:max-width="45px">Chat</button>
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

  .error-message {
    padding: 6px 8px;
    background-color: rgba(220, 53, 69, 0.9);
    color: white;
    border-radius: 4px;
    font-size: 0.75em;
    text-align: center;
    animation: slideIn 0.3s ease-out;
  }

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>