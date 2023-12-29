<script lang="ts">
	import { getUsers } from "$lib/api";
	import type { User } from "$lib/types";

    let users: Array<User> = []
    let filtered: Array<User> = []
    let searchStr: string = ''
    let err: boolean = false

    // runs at render
    const init = async() => {
        // get all users
        const res = await getUsers()
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
        }
        users = await res.json()
        if (users.length === 0) {
            err = true
        }
        // update filtered to match returned users
        filtered = users
    }

    // called by search btn
    const updateSearch = async() => {
        // if search string isnt empty
        if (searchStr.length > 0) {
            const res = await getUsers({"q":searchStr})
            if (!res.ok) {
                const msg = await res.text()
                alert(`Something went wrong, please try again:\r\n${msg}`)
            }
            const json = await res.json()
            if (json.length === 0) {
                err = true
            }
            // update filtered to match returned users
            filtered = json            
        } else {
            filtered = users
        }
        return
    }

    init()
</script>
<div class="p-2">
    <div  class="p-2">
        <input bind:value={searchStr} class="border border-black" placeholder="Search..." />
        <button on:click={updateSearch}><i class="fa-solid fa-search border border-black bg-slate-300 p-1"></i></button>
    </div>
    {#if err === true}
        <p>Looks a little empty, you should <a class="underline text-link" href="join"><b>join here</b></a></p>
    {:else}
        <ul class="p-2">
            {#each filtered as user}
                <a href="/users/{user.username}">
                    <li class="p-2 bg-slate-200 mb-1">{user.username}</li>
                </a>
            {/each}
        </ul>
    {/if}
</div>