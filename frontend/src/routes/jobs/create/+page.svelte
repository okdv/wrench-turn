<script lang="ts">
	import { apiRequest, getToken } from "$lib/api";
	import { NewJob } from "$lib/types";

    let job = new NewJob()

    const handleSubmit = async() => {
        const res = await apiRequest("/jobs/create", job, 'POST', true)
        if (res.ok) {
            const json = await res.json() 
            window.location.href = `/jobs/${json.id}`
            return
        }
        const msg = await res.text()
        alert(`Something went wrong, please try again:\r\n${msg}`)
        return
    }

    const init = async() => {
        const jwt = await getToken()
        if (jwt === null) {
            window.location.href = "/login"
        }
    }
    init()
</script>
<div>
    <form name="create-job" id="create-job" on:submit|preventDefault={handleSubmit}>
        <div>
            <label for="new-job-name">Name</label>
            <input name="new-job-name" id="new-job-name" placeholder="Oil change" bind:value={job.name} />
        </div>
        <div>
            <button type="reset">Reset</button>
            <button type="submit">Submit</button>
        </div>
    </form>
</div>