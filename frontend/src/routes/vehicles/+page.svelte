<script lang="ts">
	import { getVehicles } from "$lib/api";
	import type { Vehicle } from "$lib/types";

    let vehicles: Array<Vehicle> = []
    let err: boolean = false
    let sortStr: 'last_updated' | 'newest' | 'oldest' | 'completed' | 'za' | 'az' | 'select-one' = 'select-one'
    let searchStr = ''
    let searchUserId = ''
    let searchJobId = ''

    const init = async() => {
        const res = await getVehicles()
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        vehicles = await res.json()
        if (vehicles.length === 0) {
            err = true
        }
        return
    }
    init()

    const runSearch = async() => {
        const params: {[key:string]: string}= {}
        if (searchStr.length > 0) {
            params["q"] = encodeURIComponent(searchStr)
        }
        if (sortStr !== 'select-one') {
            params["sort"] = sortStr
        }
        if (searchUserId.length > 0) {
            params["user"] = searchStr
        }
        if (searchJobId.length > 0) {
            params["job"] = searchStr
        }

        const res = await getVehicles(params)
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        const json = await  res.json()
        vehicles = json
        return
    }
</script>
<div>
    <div class="flex justify-between p-2">
        <input bind:value={searchStr} placeholder="Search..." />
        <input bind:value={searchUserId} placeholder="User ID..." />
        <input bind:value={searchJobId} placeholder="Job ID..." />
        <select bind:value={sortStr}>
            <option value="select-one">Sort by</option>
            <option value="az">A -&gt; Z</option>
            <option value="za">A &lt;- Z</option>
            <option value="completed">Completed</option>
            <option value="newest">Newest</option>
            <option value="oldest">Oldest</option>
            <option value="last_updated">Last updated</option>
        </select>
        <button on:click={runSearch}>Search</button>
    </div>
    {#if err === true}
        <p>Looks a little empty, try adding a vehicle <a class="underline text-link" href="/vehicles/create"><b>here</b></a></p>
    {:else}
        <ul>
            {#each vehicles as vehicle}
                <li>{vehicle.name}</li>
            {/each}
        </ul>
    {/if}
</div>