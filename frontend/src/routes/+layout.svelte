<script lang="ts">
  import '../app.css'
  import favicon from '$lib/assets/favicon.svg'
  import { onMount, onDestroy } from 'svelte'
  import { Game } from '$lib/game/Game'
	import { closeSocket, initSocket } from '$lib/net/socket';
	import { userRegistered } from '$lib/stores/ui.svelte';
	import { localUser } from '$lib/stores/game.svelte';

  let game: Game
  let container: HTMLDivElement

  let { children } = $props()
  
  onMount(async () => {
    
    fetch("http://localhost:8000/session-resume", { credentials: "include" })
    .then(res => {
      if (!res.ok) console.log("Session resume failed")
      return res.json()
    })
    .then(data => {
      localUser.Username = data.username
      initSocket()
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

<style>
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