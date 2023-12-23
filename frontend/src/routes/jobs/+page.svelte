<script lang="ts">
	import { getJobs } from "$lib/api";
	import type { Job } from "$lib/types";

    let jobs: Array<Job> = []
    let err: boolean = false
    let sortStr: 'last_updated' | 'newest' | 'oldest' | 'completed' | 'za' | 'az' | 'select-one' = 'select-one'
    let searchStr = ''
    let searchUserId = ''
    let searchVehicleId = ''
    let isTemplateCheck = false

    const init = async() => {
        const res = await getJobs()
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        jobs = await res.json()
        if (jobs.length === 0) {
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
        if (searchVehicleId.length > 0) {
            params["vehicle"] = searchStr
        }
        if (isTemplateCheck) {
            params["template"] = "true"
        }
        const res = await getJobs(params)
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        const json = await  res.json()
        jobs = json
        return
    }
</script>
<div>
    <div class="flex justify-between p-2">
        <input bind:value={searchStr} placeholder="Search..." />
        <input bind:value={searchUserId} placeholder="User ID..." />
        <input bind:value={searchVehicleId} placeholder="Vehicle ID..." />
        <div>
            <input type="checkbox" bind:value={isTemplateCheck} id="is-template-check" name="is-template-check" />
            <label for="is-template-check">Is template?</label>
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
        <button on:click={runSearch}>Search</button>
    </div>
    {#if err === true}
        <p>Looks a little empty, try adding a job <a class="underline text-link" href="/jobs/create"><b>here</b></a></p>
    {:else}
        <ul>
            {#each jobs as job}
                <li>{job.name}</li>
            {/each}
        </ul>
    {/if}
</div>