<script lang="ts">
    import { page } from "$app/stores";
	import { apiRequest } from "$lib/api";
	import type { Vehicle } from "$lib/types";

    let vehicle: Vehicle | null
    let vehicleForm: Vehicle
    let edit = false

    const toggleEdit = async() => {
        if (!vehicle) {
            alert("Something went wrong, please try again")
            return    
        }
        vehicleForm = vehicle
        edit = !edit
        return
    }

    const handleEdit = async() => {
        const res = await apiRequest("/vehicles/edit", vehicleForm, 'POST', true)
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Login error, please try again: \r\n${msg}`)
            return
        }
        vehicle = await res.json()
        if (!vehicle) {
            alert("Something went wrong, please try again")
            return
        }
        await toggleEdit()
        vehicleForm = vehicle
    }

    const handleDelete = async() => {
        const res = await apiRequest(`/vehicles/${vehicle?.id}`, null, 'DELETE', true)
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Login error, please try again: \r\n${msg}`)
            return
        }
        alert("Vehicle successfully delete, taking you to the dashboard")
        window.location.href="/dash"
        return
    }

    const init = async() => {
        const res = await apiRequest(`/vehicles/${$page.params.id}`)
        if (!res.ok) {
            if (res.status === 404) {
                alert("Vehicle not found")
                return 
            }
            const msg = await res.text() 
            alert(`Login error, please try again: \r\n${msg}`)
            return
        }
        vehicle = await res.json()
    }
    init()
</script>
<div>
    {#if edit}
        <button on:click={handleDelete}>Delete</button>
        <br />
        <input name="edit-vehicle-name" id="edit-vehicle-name" bind:value={vehicleForm.name} />
    {:else}
        <h1>{!vehicle ? "..." : vehicle.name}</h1>
    {/if}
    {#if edit}
        <textbox contenteditable name="edit-vehicle-description" id="edit-vehicle-description" bind:textContent={vehicleForm.description} />
    {:else}
        <p>{!vehicle ? "..." : (vehicle.description ?? "")}</p>
    {/if}
    <p>{!vehicle ? "..." : (vehicle.description ?? "")}</p>
    <div>
        <button on:click|preventDefault={toggleEdit}>{edit ? "Cancel" : "Edit"}</button>
        {#if edit}
            <button type="submit" on:click|preventDefault={handleEdit}>Save</button>
        {/if}    
    </div>
</div>