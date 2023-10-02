<script lang="ts">
    import { apiRequest } from '$lib/api'

    const credentials = {
        username: null, 
        password: null
    }

    const handleSubmit = async() => {
        if (credentials.username === null || credentials.username === "") {
            alert("Username required")
            return
        }
        const res = await apiRequest('/auth', credentials)
        console.log(res)
    }
</script>
<form name="login" id="login" on:submit|preventDefault={handleSubmit}>
    <h1 class="text-xl">Login</h1>
    <div>
        <label for="username">Username</label>
        <input name="username" id="username" placeholder="TheStig420" bind:value={credentials.username} />
    </div>
    <div>
        <label for="password">Password</label>
        <input name="password" id="password" type="password" placeholder="••••••••••" bind:value={credentials.password} />
    </div>
    <div>
        <button type="submit">Submit</button>
        <button type="reset">Reset</button>
    </div>
</form>