<script lang="ts">
	import { apiRequest } from "$lib/api";

    class NewUser {
        username: string | null
        password: string | null
        confirmPassword: string | null
        email: string | null

        constructor(username?: string | null, password?: string | null, confirmPassword?: string | null, email?: string | null) {
            this.username = username ?? null 
            this.password = password ?? null 
            this.confirmPassword = confirmPassword ?? null 
            this.email = email ?? null 
        }
    }

    let newUserForm = new NewUser()

    const handleSubmit = async() => {
        if (newUserForm.username == null || newUserForm.username == "") {
            alert("Username required")
            return
        }
        if (newUserForm.password !== newUserForm.confirmPassword) {
            alert("Passwords do not match")
            return 
        }
        const newUser = {
            username: newUserForm.username,
            password: newUserForm.password,
            email: newUserForm.email
        }

        const res = await apiRequest('/users/create', newUser)
        if (!res.ok) {
            const msg = await res.text()
            newUserForm = new NewUser()
            alert(`Login error, please try again: \r\n${msg}`)
            return
        }
        window.location.href = '/login'            
        return
    }
</script>
<form name="register" id="register" on:submit|preventDefault={handleSubmit}>
    <h1 class="text-xl text-red-500">Join</h1>
    <div>
        <label for="username"><span class="text-red">*&nbsp;</span>Username</label>
        <input name="username" id="username" placeholder="TheStig420" bind:value={newUserForm.username} />
    </div>
    <div>
        <label for="password">Password</label>
        <input name="password" id="password" type="password" placeholder="••••••••••" bind:value={newUserForm.password} />
        <br />
        <label for="confirm-password">Confirm Password</label>
        <input name="confirm-password" id="confirm-password" type="password" placeholder="••••••••••" bind:value={newUserForm.confirmPassword} />
    </div>
    <div>
        <label for="email">Email</label>
        <input name="email" id="email" placeholder="TheStig420@aol.com" bind:value={newUserForm.email} />
    </div>
    <div>
        <button type="submit">Submit</button>
        <button type="reset">Reset</button>
    </div>
</form>