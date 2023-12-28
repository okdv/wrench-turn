<script lang="ts">
	import { apiRequest, getToken } from "$lib/api";
	import { NewVehicle } from "$lib/types";

    let vehicle = new NewVehicle()

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

    const init = async() => {
        const jwt = await getToken()
        if (jwt === null) {
            window.location.href = "/login"
        }
    }
    init()
</script>
<div>
    <form name="create-vehicle" id="create-vehicle" on:submit|preventDefault={handleSubmit}>
        <div>
            <label for="new-vehicle-name">Name</label>
            <input name="new-vehicle-name" id="new-vehicle-name" placeholder="Ole Bessy" bind:value={vehicle.name} />
        </div>
        <div>
            <button type="reset">Reset</button>
            <button type="submit">Submit</button>
        </div>
    </form>
</div>