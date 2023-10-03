<script lang="ts">
    import { apiRequest } from '$lib/api'

    class Credentials {
        username: string | null
        password: string | null
        constructor(username?: string | null, password?: string | null) {
            this.username = username ?? null 
            this.password = password ?? null 
        }
    }

    let credentials = new Credentials()

    const handleSubmit = async() => {
        if (credentials.username === null || credentials.username === "") {
            alert("Username required")
            return
        }
        const res = await apiRequest('/auth', credentials)
        if (res.status < 299) {
            window.location.href = "/dash"
            return
        }
        credentials = new Credentials()
        alert(`Login error, please try again: \r\n${res.message}`)
        return
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