<script lang="ts">
  import '../app.css'
  import favicon from '$lib/assets/favicon.svg'
  import { onMount, onDestroy } from 'svelte'
  import { Game } from '$lib/game/Game'
	import { closeSocket, initSocket, waitForOpen } from '$lib/net/socket';
	import { userRegistered } from '$lib/stores/ui.svelte';
	import { ClientData, playersState } from '$lib/stores/game.svelte';
	import { MsgType } from '$lib/Consts';
	import { randomBrightColor } from '$lib/utils';

  let game: Game
  let container: HTMLDivElement

  let { children } = $props()
  
  onMount(async () => {
    try {
      const res = await fetch("http://localhost:8000/session-resume", { credentials: "include" })
      if (res.ok) {
        
        const socket = initSocket()
        await waitForOpen(socket)
        
        const buffer = new ArrayBuffer(1)
        const view = new DataView(buffer)
        view.setUint8(0, MsgType.USER_RESUME)
        socket.send(buffer)

        const data = await res.json()
        ClientData.Username = data.username
        ClientData.Color = randomBrightColor()
        playersState[ClientData.Username] = [4000, 4000, 0, 0]
        console.log("Session resumed")
        userRegistered.isRegistered = true
      } else {
        console.log("No valid session - user needs to login")
      }
    } catch (err) {
      console.error("Session resume failed:", err)
    }
    
    game = new Game()
    await game.init()
    game.mount(container)
  })

  onDestroy(() => {
    game?.destroy()
    closeSocket()
  })


</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<div bind:this={container} class="game-root"></div>
<div class="ui-layer"> {@render children()} </div>
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