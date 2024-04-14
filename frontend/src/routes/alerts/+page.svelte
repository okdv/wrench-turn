<script lang="ts">
    import { getAlerts, updateAlertReadStatus } from "$lib/api";
	import type { Alert } from "$lib/types";

    let alerts: Array<Alert> = [] 
    let err = false 
    let sortStr: 'last_updated' | 'newest' | 'oldest' | 'completed' | 'za' | 'az' | 'select-one' = 'select-one'
    let searchStr = ''
    let searchVehicleId = ''
    let searchJobId = ''
    let searchTaskId = ''
    let isReadCheck = false

    // runs on render
    const init = async() => {
        // get all alerts
        const res = await getAlerts() 
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        // assign alerts to state
        alerts = await res.json()
        if (alerts.length === 0) {
            err = true
        }
        return
    }

    // runs on search btn
    const updateSearch = async() => {
        // create empty string obj, add items to obj for each filter 
        const params: {[key:string]: string}= {}
        if (searchStr.length > 0) {
            params["q"] = encodeURIComponent(searchStr)
        }
        if (sortStr !== 'select-one') {
            params["sort"] = sortStr
        }
        if (searchVehicleId.length > 0) {
            params["vehicle"] = searchStr
        }
        if (searchJobId.length > 0) {
            params["job"] = searchStr
        }
        if (searchTaskId.length > 0) {
            params["task"] = searchStr
        }
        if (isReadCheck) {
            params["read"] = "true"
        }
        const res = await getAlerts(params)
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        // update alerts
        const json = await  res.json()
        alerts = json
        return
    }

    // runs on update read status
    const updateReadStatus = async(i: number, alertId: number, unread?: boolean) => {
        let alertItem = alerts[i]
        unread = unread ?? false 
        const res = await updateAlertReadStatus(alertId, unread)
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return        
        }
        const newAlerts = alerts 
        newAlerts[i].isRead = Number(!unread)
        alerts = newAlerts
    }

    init()
</script>
<div class="p-2">
    <div class="flex justify-between p-2">
        <input bind:value={searchStr} placeholder="Search..." class="border border-black" />
        <input bind:value={searchVehicleId} placeholder="Vehicle ID..." class="border border-black" />
        <input bind:value={searchJobId} placeholder="Job ID..." class="border border-black" />
        <input bind:value={searchTaskId} placeholder="Task ID..." class="border border-black" />

        <div>
            <input type="checkbox" bind:value={isReadCheck} id="is-read-check" name="is-read-check" />
            <label for="is-read-check" class="select-none">Is read?</label>
        </div>
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
        <p>Looks a little empty, alerts are generated automatically, so try using the service to see some alerts here!</p>
    {:else}
        <ul class="p-2">
            {#each alerts as alert, i}
                    <li class="p-2 bg-slate-200 mb-1">
                        <div class="flex justify-between pb-1">
                            <h4>{alert.name}</h4>
                            {#if alert.vehicle !== null}
                                <a href="/vehicle/{alert.vehicle}">Vehicle <i class="fa-solid fa-square-arrow-up-right"></i></a>
                            {/if}
                            {#if alert.job !== null}
                                <a href="/job/{alert.job}">Job <i class="fa-solid fa-square-arrow-up-right"></i></a>
                            {/if}
                            {#if alert.task !== null}
                                <a href="/task/{alert.task}">Task <i class="fa-solid fa-square-arrow-up-right"></i></a>
                            {/if}
                            {#if alert.isRead === 1}
                                <button class="p-1 border border-black bg-slate-100" on:click={() => updateReadStatus(i, alert.id, true)}>Mark unread</button>
                            {:else}
                                <button class="p-1 border border-black bg-slate-100" on:click={() => updateReadStatus(i, alert.id, false)}>Mark read</button>
                            {/if}
                        </div>
                        <hr class="h-1 w-full bg-slate-100" />
                        <p>{alert.description}</p>
                    </li>
            {/each}
        </ul>
    {/if}
</div>