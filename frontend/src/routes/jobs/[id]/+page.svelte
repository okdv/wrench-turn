<script lang="ts">
    import { page } from "$app/stores";
	import TaskForm from "$lib/TaskForm.svelte";
	import { apiRequest, getTasks } from "$lib/api";
	import type { Job, Task } from "$lib/types";

    let job: Job | null
    let jobForm: Job
    let tasks: Array<Task> = []
    let edit = false
    let editTaskIdx: number | null = null 

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
            alert(`Unable to edit job, please try again: \r\n${msg}`)
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

    const handleNewTaskSave = async(task: Task) => {
        const newTasks = tasks
        newTasks.unshift(task)
        tasks = newTasks
    }

    const handleEditTaskSave = async(task: Task, i?: number) => {
        if (!i) {
            alert("Task was saved, but its index is unknown, please refresh")
            return 
        }
        const newTasks = tasks
        newTasks[i] = task
        tasks = newTasks
    }

    const init = async() => {
        // get job
        let res = await apiRequest(`/jobs/${$page.params.id}`)
        // error if non200 response
        if (!res.ok) {
            if (res.status === 404) {
                alert("Job not found")
                return 
            }
            const msg = await res.text() 
            alert(`Unable to get Jobs, please try again: \r\n${msg}`)
            return
        }
        // update global job
        job = await res.json()
        if (!job?.id) {
            alert("Job ID not available, please refresh and try again")
            return
        }
        // get tasks for job
        res = await getTasks(job.id)
        // error if non200 response
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Unable to get tasks, please try again: \r\n${msg}`)
            return
        }
        // update global tasks
        tasks = await res.json()
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
    <div>
        <h2>Tasks</h2>
        <ul>
            <li>
                <button>Add task</button>
            </li>
            <TaskForm jobId={Number($page.params.id)} handleSave={handleNewTaskSave} />
            {#each tasks as task, i}
                {#if editTaskIdx === i} 
                    <TaskForm jobId={task.job} task={task} i={i} handleSave={handleEditTaskSave} handleCancel={() => editTaskIdx = null} handleDeleteCallback={() => {
                        tasks.splice(i, 1)
                        tasks = tasks
                        editTaskIdx = null;
                        return;
                    }} />
                {:else}
                    <li>
                        <p>{task.name}</p>
                        <button on:click={() => editTaskIdx = i}>Edit</button>
                    </li>
                {/if}
            {/each}
        </ul>
    </div>
</div>