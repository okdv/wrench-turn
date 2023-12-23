<script lang="ts">
    import { page } from "$app/stores";
	import { apiRequest } from "$lib/api";
	import type { Job } from "$lib/types";

    let job: Job | null
    let jobForm: Job
    let edit = false

    const toggleEdit = async() => {
        if (!job) {
            alert("Something went wrong, please try again")
            return    
        }
        jobForm = job
        edit = !edit
        return
    }

    const handleEdit = async() => {
        const res = await apiRequest("/jobs/edit", jobForm, 'POST', true)
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Login error, please try again: \r\n${msg}`)
            return
        }
        job = await res.json()
        if (!job) {
            alert("Something went wrong, please try again")
            return
        }
        await toggleEdit()
        jobForm = job
    }

    const handleDelete = async() => {
        const res = await apiRequest(`/jobs/${job?.id}`, null, 'DELETE', true)
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Login error, please try again: \r\n${msg}`)
            return
        }
        alert("Job successfully delete, taking you to the dashboard")
        window.location.href="/dash"
        return
    }

    const init = async() => {
        const res = await apiRequest(`/jobs/${$page.params.id}`)
        if (!res.ok) {
            if (res.status === 404) {
                alert("Job not found")
                return 
            }
            const msg = await res.text() 
            alert(`Login error, please try again: \r\n${msg}`)
            return
        }
        job = await res.json()
    }
    init()
</script>
<div>
    {#if edit}
        <button on:click={handleDelete}>Delete</button>
        <br />
        <input name="edit-job-name" id="edit-job-name" bind:value={jobForm.name} />
    {:else}
        <h1>{!job ? "..." : job.name}</h1>
    {/if}
    {#if edit}
        <textbox contenteditable name="edit-job-description" id="edit-job-description" bind:textContent={jobForm.description} />
    {:else}
        <p>{!job ? "..." : (job.description ?? "")}</p>
    {/if}
    <p>{!job ? "..." : (job.description ?? "")}</p>
    <div>
        <h2>Instructions</h2>
        {#if edit}
            <textbox contenteditable name="edit-job-instructions" id="edit-job-instructions" bind:textContent={jobForm.instructions} />
        {:else}
            <p>{!job ? "..." : (job.instructions ?? "")}</p>
        {/if}
    </div>
    <div>
        <button on:click|preventDefault={toggleEdit}>{edit ? "Cancel" : "Edit"}</button>
        {#if edit}
            <button type="submit" on:click|preventDefault={handleEdit}>Save</button>
        {/if}    
    </div>
</div>