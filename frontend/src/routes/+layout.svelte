<script lang="ts">
  import '../app.css'
  import favicon from '$lib/assets/favicon.svg'
  import { onMount, onDestroy } from 'svelte'
  import { Game } from '$lib/game/Game'
	import { closeSocket, initSocket, getSocket, waitForOpen } from '$lib/net/socket';
	import { userRegistered } from '$lib/stores/ui.svelte';
	import { localUser } from '$lib/stores/game.svelte';
	import { MsgType } from '$lib/Consts';
	import { randomBrightColor } from '$lib/utils';

  let game: Game
  let container: HTMLDivElement

  let { children } = $props()
  
  onMount(async () => {
    
    fetch("http://localhost:8000/session-resume", { credentials: "include" })
    .then(res => {
      if (!res.ok) console.log("Session resume failed")
      return res.json()
    })
    .then(async data => {
      initSocket()
      const buffer = new ArrayBuffer(1)
      const view = new DataView(buffer)
      view.setUint8(0, MsgType.USER_RESUME)
      
      const socket = getSocket()
      await waitForOpen(socket!)
      socket!.send(view)
      localUser.Color = randomBrightColor()
      localUser.Username = data.username
      userRegistered.isRegistered = true
    })
    .catch(err => console.error(err))
  
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
    HP { localUser.Hp }
  </div>
  <div>
    <img src="/skull.webp" width="48" alt="skull-webp">
    { localUser.Kills }
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