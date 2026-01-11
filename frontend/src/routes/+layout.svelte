<script lang="ts">
  import favicon from '$lib/assets/favicon.svg'
  import { onMount, onDestroy, setContext } from 'svelte'
  import '../app.css'

  import { Game } from '$lib/game/Game'
  import { socket } from '$lib/net/socket'

  let game: Game
  let container: HTMLDivElement

  let { children } = $props()

  setContext("socket", socket)

  onMount(async () => {
    game = new Game()
    await game.init()
    game.mount(container)
  })

  onDestroy(() => {
    game?.destroy()
    socket?.close()
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