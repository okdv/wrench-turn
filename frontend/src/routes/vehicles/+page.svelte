<script lang="ts">
	import { getVehicles } from "$lib/api";
	import type { Vehicle } from "$lib/types";

    let vehicles: Array<Vehicle> = []
    let err: boolean = false
    let sortStr: 'last_updated' | 'newest' | 'oldest' | 'completed' | 'za' | 'az' | 'select-one' = 'select-one'
    let searchStr = ''
    let searchUserId = ''
    let searchJobId = ''

    // runs on render
    const init = async() => {
        // get all vehicles
        const res = await getVehicles()
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        // assign vehicles
        vehicles = await res.json()
        if (vehicles.length === 0) {
            err = true
        }
        return
    }
    
    // runs on search btn
    const updateSearch = async() => {
        // create empty param obj, add items to obj for each provided filter
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

    init()
</script>
<div class="p-2">
    <div class="flex justify-between p-2">
        <input bind:value={searchStr} placeholder="Search..." class="border border-black" />
        <input bind:value={searchUserId} placeholder="User ID..." class="border border-black" />
        <input bind:value={searchJobId} placeholder="Job ID..." class="border border-black" />
        <select bind:value={sortStr}>
            <option value="select-one">Sort by</option>
            <option value="az">A -&gt; Z</option>
            <option value="za">A &lt;- Z</option>
            <option value="completed">Completed</option>
            <option value="newest">Newest</option>
            <option value="oldest">Oldest</option>
            <option value="last_updated">Last updated</option>
        </select>
        <button on:click={updateSearch}><i class="fa-solid fa-search border border-black bg-slate-300 p-1"></i></button>
    </div>
    {#if err === true}
        <p>Looks a little empty, try adding a vehicle <a class="underline text-link" href="/vehicles/create"><b>here</b></a></p>
    {:else}
        <ul class="p-2">
            {#each vehicles as vehicle}
                <a href="/vehicles/{vehicle.id}">
                    <li class="p-2 bg-slate-200 mb-1">{vehicle.name}</li>
                </a>   
            {/each}
        </ul>
    {/if}
</div>