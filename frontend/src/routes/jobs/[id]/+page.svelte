<script lang="ts">
    import { page } from "$app/stores";
	import TaskForm from "$lib/TaskForm.svelte";
	import { apiRequest, getTasks, getLabels, getVehicles, getToken, getJWTData } from "$lib/api";
	import type { Job, Label, Task, Vehicle } from "$lib/types";

    let job: Job | null
    let jobForm: Job
    let tasks: Array<Task> = []
    let vehicles: Array<Vehicle> = []
    let edit = false
    let editTaskIdx: number | null = null 
    let labels: Array<Label> = []
    let addLabelValue: string = 'select-one' // will be select-one or index of label to add in labels

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

    const handleToggleTask = async(task: Task, i: number) => {
        const res = await apiRequest(`/jobs/${task.job}/tasks/${task.id}/complete?incomplete=${task.isComplete === 1}`, undefined, "PATCH", true)
        if (!res.ok) {
            if (res.status === 404) {
                alert("Task not found")
                return 
            }
            const msg = await res.text() 
            alert(`Unable to get toggle task, please try again: \r\n${msg}`)
            return
        }
        let newTasks = tasks 
        newTasks[i] = {
            ...task,
            isComplete: task.isComplete === 1 ? 0 : 1
        }
        tasks = newTasks
    }

    const handleAddLabel = async() => {
        if (job === null) {
            alert("Job not found, please refresh and try again")
            return 
        }
        if (addLabelValue === 'select-one') {
            alert("Select a label before saving")
            return 
        }
        const addedLabel = labels[Number(addLabelValue)]
        const res = await apiRequest(`/jobs/${job.id}/assignLabel/${addedLabel.id}`, undefined, 'POST', true)
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Unable to get add label, please try again: \r\n${msg}`)
            return
        }
        const newJob = job 
        if (newJob.labels === null) {
            newJob.labels = []
        }
        newJob.labels = [addedLabel, ...newJob.labels]
        job = newJob
        addLabelValue = 'select-one'
    }

    const handleDeleteLabel = async(i: number, id: number) => {
        if (job === null) {
            alert("Job not found, please refresh and try again")
            return 
        }
        const res = await apiRequest(`/jobs/${job.id}/assignLabel/${id}?unassign=true`, undefined, 'POST', true)
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Unable to get remove label, please try again: \r\n${msg}`)
            return
        }
        const newJob = job 
        if (newJob.labels === null) {
            alert("Something went wrong, please try again")
            return
        }
        newJob.labels.splice(i, 1)
        job = newJob
    }

    const init = async() => {
        // get jwt data
        const jwt = await getToken()
        if (jwt === null) {
            window.location.href = "/login"
            return
        }
        const jwtData = await getJWTData(jwt)
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
        // get vehicles
        res = await getVehicles({
            'user': jwtData.id
        })
        if (!res.ok) {
            alert(`Could not get your vehicles, please refresh and try again: ${await res.text()}`)
            return 
        }
        vehicles = await res.json()
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
        // get labels for job
        res = await getLabels()
        // error if non200 response
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Unable to get labels, please try again: \r\n${msg}`)
            return
        }
        labels = await res.json()
    }
    init()
</script>
<div>
    {#if edit}
        <button on:click={handleDelete}>Delete</button>
        <br />
        <input name="edit-job-name" id="edit-job-name" bind:value={jobForm.name} />
        <div>
            <label for="edit-job-vehicle">Vehicle</label>
            <select id="edit-job-vehicle" name="edit-job-vehicle" bind:value={jobForm.vehicle}>
                <option value={null}>None</option>
                {#each vehicles as vehicle}
                    <option value={vehicle.id}>{vehicle.name}</option>
                {/each}
            </select>
        </div>
        <textbox contenteditable name="edit-job-description" id="edit-job-description" bind:textContent={jobForm.description} />
        <textbox contenteditable name="edit-job-instructions" id="edit-job-instructions" bind:textContent={jobForm.instructions} />
        <div>
            <label for="edit-job-status">Status</label>
            <select id="edit-job-status" name="edit-job-status" bind:value={jobForm.isComplete}>
                <option value={0}>Incomplete</option>
                <option value={1}>Complete</option>
            </select>
        </div>    {:else if job}
        <h1 class="inline-block">{job.name ?? ""}</h1>
        <p><b>Vehicle: </b><a href="/vehicles/{job.vehicle}">{job.vehicle}</a></p>
        <p><b>Description: </b>{job.description ?? ""}</p>
        <p><b>Instructions: </b>{job.instructions ?? ""}</p>
        <p><b>Status: </b>{job.isComplete === 1 ? "Completed" : "Incomplete"}</p>
    {/if}
    {#if job && job.labels !== null}
        <div class="inline-block">
            {#if edit}
                <div class="inline-block">
                    <select bind:value={addLabelValue}>
                        <option value="select-one" selected disabled>Select label</option>
                        {#each labels as label, i}
                            <option value={i}>{label.name}</option>
                        {/each}
                    </select>
                    <button on:click={handleAddLabel}>Save</button>
                </div>
            {/if}
            {#each job.labels as label, i}
                <div class="inline-block p-1 mr-1" style="background-color: {label.color}">
                    <span>{label.name}</span>
                    {#if edit}
                        <button on:click={() => handleDeleteLabel(i, label.id)}><i class="fa-solid fa-x border border-black bg-slate-300 p-1"></i></button>
                    {/if}
                </div>
            {/each}
        </div>
    {/if}
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
                        <input type="checkbox" on:change={() => handleToggleTask(task, i)} checked={task.isComplete === 1} id="task-{task.id}-checkbox" name="task-{task.id}-checkbox" />
                        <label for="task-{task.id}-checkbox">{task.name}</label>
                        <button on:click={() => editTaskIdx = i}>Edit</button>
                    </li>
                {/if}
            {/each}
        </ul>
    </div>
</div>