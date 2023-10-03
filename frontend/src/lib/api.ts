import { PUBLIC_API_URL } from '$env/static/public'

type ApiResponse = {
    status: number,
    body: {[key:string]: unknown},
    message?: string
}

export const apiRequest = async(endpoint: string, body?: unknown, method?: string): Promise<ApiResponse> => {
    let status = 200
    try {
        const res = await fetch(`${PUBLIC_API_URL ?? 'http://localhost:8080'}${endpoint}`, {
            method: method ?? (body ? 'POST' : 'GET'),
            body: body ? JSON.stringify(body) : undefined,
            headers: {
                'Content-Type': 'application/json'
            }
        })
        status = res.status
        if (res.ok) {
            const json = await res.json() 
            return {
                status,
                body: json
            } 
        } 
        const errText = await res.text()
        throw `API Call Error: ${endpoint} responded with HTTP ${res.status}: ${errText}`
    } catch (error) {
        console.error(error)
        return {
            status,
            body: {},
            message: error as string
        }
    }
}