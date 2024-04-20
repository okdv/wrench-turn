<script lang="ts">
	import { getAlerts, getJWTData, getJobs, getToken, getVehicles } from "$lib/api";
	import type { Alert, Job, Vehicle } from "$lib/types";

    let jobs: Array<Job> = [] 
    let vehicles: Array<Vehicle> = [] 
    let alerts: Array<Alert> = []

    const init = async() => {
        const jwt = await getToken()
        if (jwt === null) {
            alert(`Something went wrong, please refresh and try again`)
            return  
        }
        const jwtData = await getJWTData(jwt)
        // get jobs
        let res = await getJobs({
            "user": jwtData.id,
            "complete": "0",
            "sort": "due_date"
        })
        // error if non200 response
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Unable to get Jobs, please try again: \r\n${msg}`)
            return
        }
        jobs = await res.json()
        // get vehicles
        res = await getVehicles({
            "user": jwtData.id
        })
        // error if non200 response
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Unable to get Vehicles, please try again: \r\n${msg}`)
            return
        }
        vehicles = await res.json()
        // get alerts
        res = await getAlerts({
            "user": jwtData.id,
            "isAlerted": "1"
        })
        // error if non200 response
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Unable to get Alerts, please try again: \r\n${msg}`)
            return
        }
        alerts = await res.json()
    }
    init()
</script>
<div class="w-full flex justify-between">
    <div>
        <h1>Dash</h1>
        <div>
            <h2>Jobs</h2>
            <ul>
                {#each jobs as job}
                    <a href="/jobs/{job.id}">
                        <li>{job.name}</li>
                    </a>
                {/each}
            </ul>
            <a href="/jobs">See more</a>
        </div>
        <div>
            <h2>Vehicles</h2>
            <ul>
                {#each vehicles as vehicle}
                    <a href="/vehicles/{vehicle.id}">
                        <li>{vehicle.name}</li>
                    </a>
                {/each}
            </ul>
            <a href="/vehicles">See more</a>
        </div>
    </div>
    <div>
        <h2>Alerts</h2>
        <ul>
            {#each alerts.slice(0,5) as alert}
                <li>{alert.name}</li>
            {/each}
        </ul>
        <a href="/alerts">See more</a>
    </div>
</div>