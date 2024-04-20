<script lang="ts">
    import { apiRequest, verifyToken, setToken } from '$lib/api'

    const checkLogin = async() => {
        const isLoggedIn = await verifyToken() 
        if (isLoggedIn === true) {
            console.log("already logged in, redirecting to homepage")
            window.location.href = "/dash"
        }
    }

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
        if (!res.ok) {
            credentials = new Credentials()
            if (res.status === 401) {
                alert("Username or password incorrect, please try again")
                return
            }
            const msg = await res.text()
            alert(`Login error, please try again: \r\n${msg}`)  
            return 
        }
        const json = await res.json() 
        const tokenAdded = await setToken(json["Value"])
        if (!tokenAdded) {
            credentials = new Credentials()
            alert("Had trouble saving your login token, please try again")
            return 
        }
        window.location.href = "/dash"
        return

    }

    checkLogin()
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