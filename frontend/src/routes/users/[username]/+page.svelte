<script lang="ts">
    import { page } from "$app/stores";
	import { apiRequest, getJWTData, getToken } from "$lib/api";
	import type { User } from "$lib/types";

    let user: User | null
    let userForm: User
    let edit = false
    let editEligible = false

    const toggleEdit = async() => {
        if (!user) {
            alert("Something went wrong, please try again")
            return    
        }
        userForm = user
        edit = !edit
        return
    }

    const handleEdit = async() => {
        const res = await apiRequest("/users/edit", userForm, 'POST', true)
        if (!res.ok) {
            const msg = await res.text() 
            alert(`Login error, please try again: \r\n${msg}`)
            return
        }
        user = await res.json()
        if (!user) {
            alert("Something went wrong, please try again")
            return
        }
        await toggleEdit()
        userForm = user
    }

    const init = async() => {
        const jwt = await getToken()
        if (jwt !== null) {
            const json = await getJWTData(jwt)
            if (json.username === $page.params.username) {
                editEligible = true
            }
        }
        const res = await apiRequest(`/users/${$page.params.username}`)
        if (!res.ok) {
            if (res.status === 404) {
                alert("User not found")
                return 
            }
            const msg = await res.text() 
            alert(`Login error, please try again: \r\n${msg}`)
            return
        }
        user = await res.json()
    }
    init()
</script>
<div>
    <p>{user?.username}</p>
    <p>{!user ? "..." : user.email}</p>
    <p>{!user ? "..." : (user.description ?? "")}</p>
</div>