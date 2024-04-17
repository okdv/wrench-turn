<script lang="ts">
	import { getLabels, apiRequest, getJWTData, getToken } from "$lib/api";
	import { NewLabel, type Label } from "$lib/types";

    let labels: Array<Label> = []
    let err: boolean = false
    let sortStr: 'last_updated' | 'newest' | 'oldest' | 'completed' | 'za' | 'az' | 'select-one' = 'select-one'
    let searchStr = ''
    let searchUserId = ''
    let searchJobId = ''
    let newLabel = new NewLabel()

    // runs on render
    const init = async() => {
        // get all labels
        const res = await getLabels()
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        // assign labels to state
        labels = await res.json()
        if (labels.length === 0) {
            err = true
        }
        return
    }

    // runs on search btn
    const updateSearch = async() => {
        // create empty string obj, add items to obj for each filter 
        const params: {[key:string]: string}= {}
        if (searchStr.length > 0) {
            params["q"] = encodeURIComponent(searchStr)
        }
        if (sortStr !== 'select-one') {
            params["sort"] = sortStr
        }
        if (searchUserId.length > 0) {
            params["user"] = searchStr
        }
        if (searchJobId.length > 0) {
            params["job"] = searchJobId
        }        
        const res = await getLabels(params)
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return
        }
        // update labels
        const json = await  res.json()
        labels = json
        return
    }

    // runs on save new label
    const handleAddLabel = async() => {
        // check if name is blank
        if (newLabel.name.length === 0) {
            alert("Name cannot be blank!")
            return
        }
        // add user id to label body 
        const jwt = await getToken()
        if (jwt === null) {
            alert("Something went wrong, please refresh and try again")
            return
        }
        const jwtData = await getJWTData(jwt)
        newLabel.user = Number(jwtData.id)
        // create via api
        let res = await apiRequest(`/labels/create`, newLabel, 'POST', true, false)
        // throw error if non 2xx
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again:\r\n${msg}`)
            return        
        }
        // add label to list
        const json = await res.json()
        const newLabels = labels
        newLabels.unshift(json)
        labels = newLabels    
        // reset form
        newLabel = new NewLabel()
        return 
    }

    init()
</script>
<div class="p-2">
    <div class="flex justify-between p-2">
        <input bind:value={searchStr} placeholder="Search..." class="border border-black" />
        <input bind:value={searchJobId} placeholder="Job ID..." class="border border-black" />
        <input bind:value={searchUserId} placeholder="User ID..." class="border border-black" />
        <select bind:value={sortStr}>
            <option value="select-one">Sort by</option>
            <option value="az">A -&gt; Z</option>
            <option value="za">A &lt;- Z</option>
            <option value="newest">Newest</option>
            <option value="oldest">Oldest</option>
            <option value="last_updated">Last updated</option>
        </select>
        <button on:click={updateSearch}><i class="fa-solid fa-search border border-black bg-slate-300 p-1"></i></button>
    </div>
    <div class="flex justify-between">
        <h3><b>Add label</b></h3>
        <div>
            <label for="new-label-form-name">Name</label>
            <input placeholder="Routine" name="new-label-form-name" id="new-label-form-name" bind:value={newLabel.name} />
        </div>
        <div>
            <label for="new-label-form-color">Color</label>
            <input type="color" name="new-label-form-color" id="new-label-form-color" bind:value={newLabel.color} />
        </div>
        <button on:click={handleAddLabel}>Save</button>
    </div>
    {#if err === true}
        <p>Looks a little empty, try adding a label <a class="underline text-link" href="/labels/create"><b>here</b></a></p>
    {:else}
        <ul class="p-2">
            {#each labels as label}
                <li class="p-2 bg-slate-200 mb-1">{label.name}</li>
            {/each}
        </ul>
    {/if}
</div>