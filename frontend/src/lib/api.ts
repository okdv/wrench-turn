export const apiRequest = async(endpoint: string, body?: {[key:string]: unknown}, method?: string): Promise<{[key:string]:unknown}> => {
    try {
        const res = await fetch(`http://localhost:8080${endpoint}`, {
            method: method ?? (body ? 'POST' : 'GET'),
            body: body ? JSON.stringify(body) : undefined,
            headers: {
                'Content-Type': 'application/json'
            }
        })
        if (res.ok) {
            const json = await res.json() 
            return json 
        } 
        const errText = await res.text()
        throw `${res.status}, ${errText}`
    } catch (error) {
        console.error(`API Call Error: ${endpoint} responded with non-200 error: ${error}`)
        return {}
    }
}