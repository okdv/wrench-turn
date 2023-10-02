<script lang="ts">
	import { apiRequest } from "$lib/api";

    let newUserForm = {
        username: null,
        password: null,
        confirmPassword: null,
        email: null,
    }

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
        console.log(res)
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