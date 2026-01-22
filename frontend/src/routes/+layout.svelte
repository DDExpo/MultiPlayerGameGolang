<script lang="ts">
  import '../app.css'
  import favicon from '$lib/assets/favicon.svg'
  import { onMount, onDestroy } from 'svelte'
  import { Game } from '$lib/game/Game'
	import { closeSocket, getSocket, initSocket, waitForOpen } from '$lib/net/socket';
	import { ClientData, playersState } from '$lib/stores/game.svelte';
	import { MsgType, StateType } from '$lib/Consts';
	import { randomBrightColor } from '$lib/utils';
	import { userUiState } from '$lib/stores/ui.svelte';
  
  let container: HTMLDivElement
  let game: Game
  let { children } = $props()

  onMount(async () => {
    try {
      const res = await fetch("http://localhost:8000/session-resume", { credentials: "include" })
      if (res.ok) {
        
        const socket = initSocket()
        await waitForOpen(socket)
        
        const buffer = new ArrayBuffer(1)
        const view = new DataView(buffer)
        view.setUint8(0, StateType.USER_RESUME)
        socket.send(buffer)
        
        const data = await res.json()
        ClientData.Username = data.username
        ClientData.Color = randomBrightColor()

        playersState[ClientData.Username] = { movement: [4000, 4000, 0], combat: { dead: false } }
        userUiState.registered = true
        userUiState.alive      = true
        
        if (ClientData.Hp <= 0) {
          playersState[ClientData.Username].combat.dead = true
          userUiState.alive   = false
          userUiState.focused = true
        }
        
        console.log("Session resumed")
      } else {
        console.log("No valid session - user needs to login")
      }
    } catch (err) {
      console.error("Session resume failed:", err)
    }
    
    game = new Game()
    await game.init()
    game.mount(container)
    game.setUsernameTextStyle()
  })
  
  onDestroy(() => {
    game?.destroy()
    closeSocket()
  })

  async function resumeGame() {
    const socket = getSocket()
    const buffer = new ArrayBuffer(1)
    const view = new DataView(buffer)
    view.setUint8(0, MsgType.USER_RESUMED_DEATH)
  
    socket?.send(buffer)
    userUiState.alive   = true
    userUiState.focused = false
    playersState[ClientData.Username].movement    = [4000, 4000, 0]
    playersState[ClientData.Username].combat.dead = false
  }

</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<div bind:this={container} class="game-root"></div>

<div class="ui-layer"> {@render children()} </div>

{#if !userUiState.alive && userUiState.registered}
  <div class="resume"><button onclick={ resumeGame } > RESUME </button></div>
{/if}

<div class="player-stats">
  <div>
    HP { ClientData.Hp }
  </div>
  <div>
    <img src="/skull.webp" width="48" alt="skull-webp">
    { ClientData.Kills }
  </div>
</div>


<style>
  .player-stats {
    position: absolute;
    display: grid;
    bottom: 22%;
    margin-left: 8px;
    gap: 20px;
    z-index: 9;
    color: white;
    font-size: x-large;
  }

  .resume {
    height: 100%;
    position: relative;
    background: transparent;
    justify-self: center;
    align-content: center;
    inset: 0;
    z-index: 11;
  }

  .resume button {
    width: 136px;
    height: 48px;
    font-size: large;
  }

  .game-root {
    position: fixed;
    inset: 0;
    z-index: 0;
  }
  .ui-layer {

    position: fixed;
    inset: 0;
    z-index: 10;
  }
</style>