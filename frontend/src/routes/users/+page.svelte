<script lang="ts">
	import { getUsers } from "$lib/api";
	import type { User } from "$lib/types";
    import List from "$lib/List.svelte"

    let users: Array<User> = []
    let filtered: Array<User> = []
    let searchStr: string = ''
    let err: boolean = false

    const init = async() => {
        const res = await getUsers()
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
        }
        users = await res.json()
        if (users.length === 0) {
            err = true
        }
        filtered = users
    }

    const updateSearch = async() => {
        console.log(searchStr)
        if (searchStr.length > 0) {
            const matched = users.filter(user => {
                if (user.username.toLowerCase().includes(searchStr.toLowerCase())) {
                    return true 
                } else if (user.description && user.description.toLowerCase().includes(searchStr.toLowerCase())) {
                    return true 
                }
                return false 
            })
            console.log(matched)
            filtered = matched
        } else {
            filtered = users
        }
    }

    init()
</script>
<div>
    <div>
        <input bind:value={searchStr} />
        <button on:click={updateSearch}>Search</button>
    </div>
    {#if err === true}
        <p>Looks a little empty, you should <a class="underline text-link" href="join"><b>join here</b></a></p>
    {:else}
        <List users={filtered} />
    {/if}
</div>