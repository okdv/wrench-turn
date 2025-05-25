<script lang="ts">
	import { getEnv, verifyToken } from "$lib/api";
    import "../app.css";
    import {page} from '$app/stores';

    let isLoggedIn = false
    let version = ''

    const init = async() => {
      isLoggedIn = await verifyToken()
      const res = await getEnv() 
      if (!res.ok) {
        alert("Unable to get env data from API")
        return 
      }
      const envJson = await res.json() 
      version = envJson["API_VERSION"]
    }
    init()
  </script>
  
  <svelte:head>
    <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200&icon_names=home">
  </svelte:head>


  <div class="flex justify-between p-2 m-2 border-2 border-blue-500 rounded-md">
    {#if $page.url.pathname !== "/"}
    <a href = '/'>
      <span class="material-symbols-outlined">home</span>
    </a>
    {/if}

    <a href="{isLoggedIn ? "/dash" : "/"}">
      <h1 class="text-xl" class:hide = {$page.url.pathname !== '/'}>WrenchTurn</h1>
    </a>
    <div class="flex justify-around bg-white space-x-1">
      <a href="/users" class="p-2 border-2 border-blue-500 rounded-md" >Users</a>
      <a href="/jobs" class="p-2 border-2 border-blue-500 rounded-md">Jobs</a>
      <a href="/vehicles" class="p-2 border-2 border-blue-500 rounded-md">Vehicles</a>
      {#if isLoggedIn}
        <a href="/settings" class="p-2 border-2 border-blue-500 rounded-md">Settings</a>
      {:else}
        <a href="/login" class="p-2 border-2 border-blue-500 rounded-md">Login</a>
        <a href="/join" class="p-2 border-2 border-blue-500 rounded-md">Join</a>
      {/if}
    </div>
  </div>
  <slot />
  <div class="text-center fixed bottom-0 right-0 left-0">
    <p class="text-center inline-block p-2 mr-2">Powered by WrenchTurn {version}</p>
    &bull;
    <a href="https://github.com/okdv/wrench-turn" class="inline-block p-2 ml-2 text-link">
      <i class="fa-solid fa-star"></i>
      &nbsp;or&nbsp;
      <i class="fa-solid fa-code-fork"></i>
      &nbsp;on&nbsp;
      <i class="fa-brands fa-github"></i>
    </a>
  </div>
  