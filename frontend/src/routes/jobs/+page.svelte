<script lang="ts">
	import { getJobs } from "$lib/api";
	import type { Job } from "$lib/types";

    let jobs: Array<Job> = []
    let err: boolean = false
    let sortStr: 'last_updated' | 'newest' | 'oldest' | 'completed' | 'za' | 'az' | 'select-one' = 'select-one'
    let searchStr = ''
    let searchUserId = ''
    let searchVehicleId = ''
    let searchLabelId = ''
    let isTemplateCheck = false

    // runs on render
    const init = async() => {
        // get all jobs
        const res = await getJobs()
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        // assign jobs to state
        jobs = await res.json()
        if (jobs.length === 0) {
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
        if (searchUserId.length > 0) {
            params["user"] = searchUserId
        }
        if (searchLabelId.length > 0) {
            params["label"] = searchLabelId
        }
        if (searchVehicleId.length > 0) {
            params["vehicle"] = searchVehicleId
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
        // update jobs
        const json = await  res.json()
        jobs = json
        return
    }

    init()
</script>
<div class="p-2">
    <div class="flex justify-between p-2">
        <input bind:value={searchStr} placeholder="Search..." class="border border-black" />
        <input bind:value={searchUserId} placeholder="User ID..." class="border border-black" />
        <input bind:value={searchLabelId} placeholder="Label ID..." class="border border-black" />
        <input bind:value={searchVehicleId} placeholder="Vehicle ID..." class="border border-black" />
        <div>
            <input type="checkbox" bind:value={isTemplateCheck} id="is-template-check" name="is-template-check" />
            <label for="is-template-check" class="select-none">Is template?</label>
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
        <p>Looks a little empty, try adding a job <a class="underline text-link" href="/jobs/create"><b>here</b></a></p>
    {:else}
        <ul class="p-2">
            {#each jobs as job}
                <a href="/jobs/{job.id}">
                    <li class="p-2 bg-slate-200 mb-1">{job.name}</li>
                </a>
            {/each}
        </ul>
    {/if}
</div>