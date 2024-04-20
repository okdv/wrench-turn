<script lang="ts">
	import { apiRequest, getToken } from "$lib/api";
	import { NewVehicle } from "$lib/types";

    let vehicle = new NewVehicle()
    let isMetric = false
    
    // runs on render, redirect if not logged in
    const init = async() => {
        const jwt = await getToken()
        if (jwt === null) {
            window.location.href = "/login"
        }
    }
    // runs on submit btn, create via api
    const handleSubmit = async() => {
        const res = await apiRequest("/vehicles/create", vehicle, 'POST', true)
        if (res.ok) {
            const json = await res.json() 
            window.location.href = `/vehicles/${json.id}`
            return
        }
        const msg = await res.text()
        alert(`Something went wrong, please try again:\r\n${msg}`)
        return
    }
    // runs on change of metric input, converts it to int
    const updateIsMetric = async() => vehicle.isMetric = isMetric ? 1 : 0
    init()
</script>
<div>
    <form name="create-vehicle" id="create-vehicle" on:submit|preventDefault={handleSubmit} class="space-y-1">
        <div>
            <label for="new-vehicle-name">Name</label>
            <input name="new-vehicle-name" id="new-vehicle-name" placeholder="Ole Bessy" bind:value={vehicle.name} class="border border-black" />
        </div>
        <div>
            <label for="new-vehicle-description">Description</label>
            <textarea name="new-vehicle-description" id="new-vehicle-description" placeholder="..." bind:value={vehicle.description} class="border border-black" />
        </div>
        <div>
            <label for="new-vehicle-vin">Vin</label>
            <input name="new-vehicle-vin" id="new-vehicle-vin" placeholder="car" bind:value={vehicle.vin} class="border border-black" />
        </div>
        <div>
            <label for="new-vehicle-type">Type</label>
            <input name="new-vehicle-type" id="new-vehicle-type" placeholder="car" bind:value={vehicle.type} class="border border-black" />
        </div>
        <div>
            <label for="new-vehicle-year">Year</label>
            <input name="new-vehicle-year" id="new-vehicle-year" type="number" placeholder="1991" bind:value={vehicle.year} class="border border-black" />
        </div>
        <div>
            <label for="new-vehicle-make">Make</label>
            <input name="new-vehicle-make" id="new-vehicle-make" placeholder="Nissan" bind:value={vehicle.make} class="border border-black" />
        </div>
        <div>
            <label for="new-vehicle-model">Model</label>
            <input name="new-vehicle-model" id="new-vehicle-model" placeholder="Skyline" bind:value={vehicle.model} class="border border-black" />
        </div>
        <div>
            <label for="new-vehicle-trim">Trim</label>
            <input name="new-vehicle-trim" id="new-vehicle-trim" placeholder="GTR" bind:value={vehicle.trim} class="border border-black" />
        </div>
        <div>
            <input name="new-vehicle-is-metric" id="new-vehicle-is-metric" type="checkbox" bind:value={isMetric} class="border border-black" on:change={updateIsMetric} />
            <label for="new-vehicle-is-metric" class="select-none">Is metric?</label>
        </div>
        <div>
            <label for="new-vehicle-odo">Odometer</label>
            <input name="new-vehicle-odo" id="new-vehicle-odo" type="number" placeholder="25000" bind:value={vehicle.odometer} class="border border-black" />
        </div>
        <div>
            <button type="reset">Reset</button>
            <button type="submit">Submit</button>
        </div>
    </form>
</div>