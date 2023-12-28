<script lang="ts">
	import { verifyToken } from "$lib/api";
    import "../app.css";

    let isLoggedIn = false

    const init = async() => {
      isLoggedIn = await verifyToken()
    }
    init()
  </script>
  
  <div class="flex justify-between p-2">
    <a href="{isLoggedIn ? "/dash" : "/"}">
      <h1 class="text-xl">WrenchTurn</h1>
    </a>
    <div class="flex justify-around">
      <a href="/users" class="p-2">Users</a>
      <a href="/jobs" class="p-2">Jobs</a>
      <a href="/vehicles" class="p-2">Vehicles</a>
      {#if isLoggedIn}
        <a href="/settings" class="p-2">Settings</a>
      {:else}
        <a href="/login" class="p-2">Login</a>
        <a href="/join" class="p-2">Join</a>
      {/if}
    </div>
  </div>
  <slot />
  <div class="text-center">
    <p class="text-center inline-block p-2 mr-2">Powered by WrenchTurn v0.0.0-alpha</p>
    &bull;
    <a href="https://github.com/okdv/wrench-turn" class="inline-block p-2 ml-2 text-link">
      <i class="fa-solid fa-star"></i>
      &nbsp;or&nbsp;
      <i class="fa-solid fa-code-fork"></i>
      &nbsp;on&nbsp;
      <i class="fa-brands fa-github"></i>
    </a>
  </div>
  