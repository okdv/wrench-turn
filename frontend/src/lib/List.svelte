<script lang="ts">
	import type { User } from "$lib/types";

    export let users: Array<User> = []
</script>
<div>
    <ul>
        {#each users as user}
            <li>{user.username}</li>
        {/each}
    </ul>
</div>