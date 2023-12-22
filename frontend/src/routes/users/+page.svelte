<script lang="ts">
	import { getUsers } from "$lib/api";
	import type { User } from "$lib/types";

    let users: Array<User> = []
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
    }
    init()
</script>
<div>
    {#if err === true}
        <p>Looks a little empty, you should <a class="underline text-link" href="join"><b>join here</b></a></p>
    {:else}
        <ul>
            {#each users as user}
                <li>{user.username}</li>
            {/each}
        </ul>
    {/if}
</div>