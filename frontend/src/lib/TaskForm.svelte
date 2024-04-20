<script lang="ts">
	import { apiRequest } from "./api";
	import { NewTask, type Task } from "$lib/types";

    export let handleCancel: () => void = () => console.log('reset new task form')
    export let handleSave: (task: Task, i?: number) => void = () => console.log('task form saved')
    export let handleDeleteCallback: (task: Task, i?: number) => void = () => console.log('delete task called but no callback')

    export let task: Task | undefined = undefined
    export let i: number | undefined = undefined
    export let jobId: number
    let taskForm: NewTask
    let edit: boolean = false

    const handleSubmit = async() => {
        if (taskForm.name === null || taskForm.name.length === 0) {
            alert(`Task name cannot be empty`)
            return
        }
        if (task !== undefined) {
            if (jobId !== task.job) {
                alert(`Job Id provided (${jobId}) does not match job id on task (${task.job})`)
                return
            }
            task = {
                ...task,
                ...taskForm
            }
            const res = await apiRequest(`/jobs/${jobId}/tasks/edit`, task, 'POST', true)
            if (res.ok) {
                const json = await res.json() 
                handleSave(json, i)
                return
            }
            alert(`Unable to edit task: HTTP ${res.status}`)
            return
        }
        const res = await apiRequest(`/jobs/${jobId}/tasks/create`, taskForm, 'POST', true)
        if (res.ok) {
            const json = await res.json() 
            handleSave(json)
            taskForm = new NewTask()
            return
        }
        alert(`Unable to create task: HTTP ${res.status}`)
        return
    }

    const handleDelete = async() => {
        if (task === undefined) {
            alert("Cannot delete undefined task")
            return 
        }
        const res = await apiRequest(`/jobs/${jobId}/tasks/${task.id}`, undefined, 'DELETE', true)
        if (res.ok) {
            handleDeleteCallback(task, i)
            return
        }
        alert(`Unable to delete task: HTTP ${res.status}`)
        return
    }

    const init = async() => {
        if (task !== undefined) {
            if (i === undefined) {
                alert("Task index must be provided if editing an existing task, please refresh and try again")
                return
            }
            edit = true
            taskForm = new NewTask(
                task.name,
                task.description,
                task.partName,
                task.partLink,
                task.dueDate,
            )
            return 
        }
        edit = false 
        taskForm = new NewTask()
    }
    init()
</script>
<form class="space-y-2" on:submit|preventDefault={handleSubmit} on:reset={handleCancel}>
    <input name="task-form-name" id="task-form-name" bind:value={taskForm.name} placeholder="Task name..." class="border border-black" />
    <textarea name="task-form-description" id="task-form-description" bind:value={taskForm.description} placeholder="Task description..." class="border border-black"></textarea>
    <input name="task-form-part-name" id="task-form-part-name" bind:value={taskForm.partName} placeholder="Part name..." class="border border-black" />
    <input name="task-form-part-link" id="task-form-part-link" bind:value={taskForm.partLink} placeholder="Part link..." class="border border-black" />
    <input name="task-form-due-date" id="task-form-due-date" bind:value={taskForm.dueDate} placeholder="Due date..." class="border border-black" />
    <button type="submit">Save</button>
    <button type="reset">Cancel</button>
    {#if task !== undefined}
        <button on:click|preventDefault={handleDelete}>Delete</button>
    {/if}
</form>