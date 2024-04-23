<script lang="ts">
	import { apiRequest, getJWTData, getToken, getVehicles } from "$lib/api";
	import { NewJob, type Vehicle } from "$lib/types";

    let job = new NewJob()
    let vehicles: Array<Vehicle> = []
    let isTemplate = false 
    let selectedVehicleId: string = "select-one"
    let repeats = false 
    let timeIntervalUnit: "select-one" | "month" | "day" | "year" | "week" = "select-one"
    
    // runs on render, redirect if not logged in
    const init = async() => {
        const jwt = await getToken()
        if (jwt === null) {
            window.location.href = "/login"
            return
        }
        const jwtData = await getJWTData(jwt)
        const vehiclesRes = await getVehicles({
            'user': jwtData.id
        })
        if (!vehiclesRes.ok) {
            alert(`Could not get your vehicles, please refresh and try again: ${await vehiclesRes.text()}`)
            return 
        }
        vehicles = await vehiclesRes.json()
    }
    // runs on submit btn, create via api
    const handleSubmit = async() => {
        if (selectedVehicleId !== "select-one") {
            job.vehicle = Number(selectedVehicleId)
        }
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
    // runs on change of template input, converts to int
    const updateIsTemplate = async() => job.isTemplate = isTemplate ? 1 : 0
    // runs on change of repeats input, converts to int
    const updateRepeats = async() => job.repeats = repeats ? 1 : 0
    // runs on change of time unit select, replaces "select-one" with null
    const updateTimeIntervalUnit = async() => {
        if (timeIntervalUnit === "select-one") {
            job.timeIntervalUnit = null
            return
        }
        job.timeIntervalUnit = timeIntervalUnit
    }
    init()
</script>
<div>
    <form name="create-job" id="create-job" on:submit|preventDefault={handleSubmit} class="space-y-1">
        <div>
            <label for="new-job-name">Name</label>
            <input name="new-job-name" id="new-job-name" placeholder="Oil change" bind:value={job.name} class="border border-black" />
        </div>
        <div>
            <label for="new-job-description">Description</label>
            <textarea name="new-job-description" id="new-job-description" placeholder="Replace oil and oil filter" bind:value={job.description} class="border border-black" />
        </div>
        <div>
            <label for="new-job-instructions">Instructions</label>
            <textarea name="new-job-instructions" id="new-job-instructions" placeholder="Undo 12mm oil drain plug, drain oil, replace plug gasket and reinstall, remove and replace oil filter, refill oil" bind:value={job.instructions} class="border border-black" />
        </div>
        <div>
            <label for="new-job-vehicle">Attach vehicle</label>
            <select id="new-job-vehicle" name="new-job-vehicle" bind:value={selectedVehicleId}>
                <option value="select-one" disabled selected>Select one</option>
                {#each vehicles as vehicle}
                    <option value={vehicle.id}>{vehicle.name}</option>
                {/each}
            </select>
        </div>
        <div>
            <input name="new-job-is-template" id="new-job-is-template" type="checkbox" bind:value={isTemplate} class="border border-black" on:change={updateIsTemplate} />
            <label for="new-job-is-template" class="select-none">Is template?</label>
        </div>
        <div>
            <input name="new-job-repeats" id="new-job-repeats" type="checkbox" bind:value={repeats} class="border border-black" on:change={updateRepeats} />
            <label for="new-job-repeats" class="select-none">Repeats?</label>
        </div>
        <div>
            <label for="new-job-odo-interval">Odometer interval</label>
            <input name="new-job-odo-interval" id="new-job-odo-interval" type="number" placeholder="4000" bind:value={job.odoInterval} class="border border-black" />
        </div>
        <div>
            <label for="new-job-time-interval">Time interval</label>
            <input name="new-job-time-interval" id="new-job-time-interval" type="number" placeholder="6" bind:value={job.timeInterval} class="border border-black" />
        </div>
        <div>
            <label for="new-job-time-interval-unit">Time interval</label>
            <select name="new-job-time-interval-unit" id="new-job-time-interval-unit" bind:value={timeIntervalUnit} on:change={updateTimeIntervalUnit}>
                <option value="select-one">Select one</option>
                <option value="day">Days</option>
                <option value="week">Weeks</option>
                <option value="month">Months</option>
                <option value="year">Years</option>
            </select>
        </div>
        <div>
            <label for="new-job-due-date">Due date (format: 2023-10-05T19:32:20Z)</label>
            <input name="new-job-due-date" id="new-job-due-date" placeholder="2023-10-05T19:32:20Z" bind:value={job.dueDate} class="border border-black" />
        </div>
        <div>
            <button type="reset">Reset</button>
            <button type="submit">Submit</button>
        </div>
    </form>
</div>