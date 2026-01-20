<script lang="ts">
  import { messages, userUiState } from "$lib/stores/ui.svelte"
  import { MAX_MESSAGE_LENGTH, MAX_USERNAME_LENGTH, MESSAGE_COOLDOWN_MS, MsgType } from "$lib/Consts"
  import { ClientData, playersState } from "$lib/stores/game.svelte"
  import { getSocket, initSocket, waitForOpen } from "$lib/net/socket"
	import { randomBrightColor } from "$lib/utils"
	import { HttpStatus } from "$lib/types/enums"
	import { serializeChatMsg } from "$lib/net/serialize";


  let socket: WebSocket | null
  let hide: boolean                    = $state(true)
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
    if (!userUiState.registered)                         { showError('User is not authorized', 'chat'); return }

    const now = new Date()
    const trimmed = message.trim()
    if (!trimmed)                            { showError('Message is empty', 'chat'); return }
    if (trimmed.length > MAX_MESSAGE_LENGTH) { showError(`Message is too long (max ${MAX_MESSAGE_LENGTH} characters)`, 'chat'); return}
    if (lastTimeMessageSent && (now.getTime() - lastTimeMessageSent.getTime()) < MESSAGE_COOLDOWN_MS) {showError('Please wait 15 seconds between messages', 'chat'); return }
    
    try {
      const msg = serializeChatMsg(trimmed, ClientData.Color)
      socket.send(msg)
      lastTimeMessageSent = now
      message = ""
    } catch (err) { console.error(`Failed: ${err}`)}
  }

  async function createUser() {

    ClientData.Username = ClientData.Username.trim()

    if (!ClientData.Username)                             { showError('Username is empty', 'login'); return }
    if (ClientData.Username.length > MAX_USERNAME_LENGTH) { showError(`Username is too long (max ${MAX_USERNAME_LENGTH} characters)`, 'login'); return }

    try {
      const res = await fetch("http://localhost:8000/initialize-session", {
        method:      "POST",
        credentials: "include",
        headers:   { "Content-Type": "application/json" },
        body:        JSON.stringify({ username: ClientData.Username }),
      })
      
      if (res.status === HttpStatus.BAD_REQUEST) { showError(await res.text(), 'login'); return }
      if (res.status !== HttpStatus.CREATED)     { showError('Something went wrong', 'login'); console.error(await res.text()); return } 
      
      socket = initSocket()
      const buf = new ArrayBuffer(1);
      new DataView(buf).setUint8(0, MsgType.USER_REG);
      
      await waitForOpen(socket)
      socket.send(buf)

      userUiState.registered = true
      userUiState.alive      = true
      playersState[ClientData.Username] = [4000, 4000, 0, false]
      ClientData.Color = randomBrightColor()
    
    } catch (err) { console.error(`Failed: ${err}`)}
  }

</script>


<div class="main">

  {#if !userUiState.registered}
    <div class="game-enter-screen">
      <div style:color="White">Enter username</div>
      <input type="text"
             onfocus={() => userUiState.focused = true} 
             onblur={()  => userUiState.focused = false} 
             bind:value={ClientData.Username}>
      <button style:font-size="1rem" onclick={createUser}>Enter</button>
      {#if place === "login"}
        <div class="error-message">{error}</div>
      {/if}
    </div>
  {/if}
  
  <div class="chat-block">
    {#if !hide}
      <div class="messages-screen">
        {#each messages as [color, [header, msgTime]]}
          <p>
            <span style:color={color}> {header} </span> {msgTime}
          </p>
        {/each}
      </div>

      <div class="input-row">
        <input type="text"
               onfocus={() => userUiState.focused = true} 
               onblur={()  => userUiState.focused = false} 
               bind:value={message}
               onkeydown={(e) => e.key === 'Enter' && sendMessage()}>
        <button class="btn-send-message" onclick={sendMessage}>send</button>
      </div>

      {#if place === "chat"}
        <div class="error-message">{error}</div>
      {/if}
      <button class="btn-hide-chat" onclick={() => hide = !hide} style:width="16px" style:height="16px" aria-label="open"></button>

    {:else}
      <button onclick={() => hide = !hide} style:width="16px" style:height="16px" aria-label="open"></button>
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
    max-width: 350px;
    min-width: 350px;
    left: 0;
    bottom: 0;
    margin: 8px;
    position: absolute;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .messages-screen {
    max-width: 350px;
    min-height: 150px;
    padding: 4px;
    color: white;
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