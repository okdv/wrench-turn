<script lang="ts">
	import { apiRequest, getUser } from "$lib/api";
	import { UpdatePassword, User } from "$lib/types";

    let user: User | null
    let userForm: User = new User()
    let deleteInit = false
    let passwordForm = new UpdatePassword()

    // runs at render, sets user and userForm data
    const init = async() => {
        const res = await getUser() 
        user = await res.json()
        if (user !== null) {
            userForm = user
        }
    }
    // runs on click of submit btn
    const handleSubmit = async() => {
        const res = await apiRequest("/users/edit", userForm, "POST", true, false)
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again: \r\n${msg}`) 
            return
        }
        init()
    }
    // runs twice, on click of delete btn and confirm delete btn, first click sets deleteInit=true as a confirmation step
    const handleDelete = async() => {
        if (deleteInit === false) {
            deleteInit = true
            return
        }
        const res = await apiRequest(`/users/${user?.username}`, undefined, "DELETE", true, false)
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again: \r\n${msg}`) 
            return
        }
        alert("User deleted")
        window.location.href = "/logout"
    }
    // runs when update password form submit btn clicked
    const handleUpdatePassword = async() => {
        if (passwordForm.newPassword !== passwordForm.confirmPassword) {
            alert("Passwords do not match")
            return
        }
        const body = {
            username: user?.username,
            currentPassword: passwordForm.password,
            newPassword: passwordForm.confirmPassword
        }
        const res = await apiRequest("/users/updatePassword", body, 'POST', true, false)
        if (!res.ok) {
            const msg = await res.text()
            alert(`Something went wrong, please try again: \r\n${msg}`) 
            return
        }
        passwordForm = new UpdatePassword() 
        return 
    }
    init()
</script>
<div>
    <h1>Settings</h1>
    <form on:submit|preventDefault>
        <div>
            <label for="user-id">User ID</label>
            <input disabled type="number" name="user-id" id="user-id" value={user?.id} />
        </div>
        <div>
            <label for="username">Username</label>
            <input disabled name="username" id="username"  value={user?.username} />
        </div>
        <div>
            <label for="email">Email</label>
            <input type="email" name="email" id="email" bind:value={userForm.email} />
        </div>
        <div>
            <label for="description">Description</label>
            <br />
            <textarea name="description" id="description" bind:value={userForm.description}></textarea>
        </div>
        <div>
            <button on:click={init}>Reset</button>
            <button on:click={handleSubmit}>Save</button>
        </div>
        <button on:click={handleDelete}>{deleteInit ? "Click again to confirm" : "Delete"}</button>
        {#if deleteInit}
            <button on:click={() => deleteInit = false}>Cancel</button>
        {/if}
    </form>
    <hr />
    <h2>Update password</h2>
    <form on:submit|preventDefault>
        <div>
            <label for="password">Password</label>
            <input type="password" name="password" id="password" bind:value={passwordForm.password} />
        </div>
        <div>
            <label for="new-password">New password</label>
            <input type="password" name="new-password" id="new-password" bind:value={passwordForm.newPassword} />
        </div>
        <div>
            <label for="confirm-password">Confirm new password</label>
            <input type="password" name="confirm-password" id="confirm-password" bind:value={passwordForm.confirmPassword} />
        </div>
        <div>
            <button on:click={() => passwordForm = new UpdatePassword()}>Reset</button>
            <button on:click={handleUpdatePassword}>Update</button>
        </div>
    </form>
    <hr />
    <h3>Admin settings</h3>
</div>
