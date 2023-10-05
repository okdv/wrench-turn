<script lang="ts">
	import { getJobs } from "$lib/api";
	import type { Job } from "$lib/types";

    let jobs: Array<Job> = []
    let err: boolean = false

    const init = async() => {
        const res = await getJobs()
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
        }
        jobs = await res.json()
        if (jobs.length === 0) {
            err = true
        }
    }
    init()
</script>
<div>
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