import { PUBLIC_API_URL } from '$env/static/public'

// getToken
// get jwt from localstorage
export const getToken = async(): Promise<string | null> => {
    return localStorage.getItem('wrenchturn-jwt')
}
// setToken
// set jwt in localstorage (remove if jwt is null)
export const setToken = async(jwt?: string): Promise<boolean> => {
    if (!jwt) {
        localStorage.removeItem('wrenchturn-jwt')
    } else {
        localStorage.setItem('wrenchturn-jwt', jwt)
    }
    const token = await getToken() 
    if (token === jwt) {
        return true
    }
    return false
}
// apiRequest
// fetch proxy purpose built for api requests
export const apiRequest = async(endpoint: string, body?: unknown, method?: string, redirectOnFail?: boolean): Promise<Response> => {
    // get jwt from localstorage
    const jwt = await getToken()
    // initialize headers for fetch
    const headers: {[key:string]: string} = {}
    // if req body, add json content type
    if (body) {
        headers['Content-Type'] = 'application/json'
    }
    // if jwt isnt null, add auth header with it as bearer
    if (jwt) {
        headers['Authorization'] = `Bearer ${jwt}`
    }
    // create fetch
    const res = await fetch(`${PUBLIC_API_URL ?? 'http://localhost:8080'}${endpoint}`, {
        method: method ?? (body ? 'POST' : 'GET'),
        body: body ? JSON.stringify(body) : undefined,
        headers
    })
    if (res.status === 401 && redirectOnFail && redirectOnFail === true) {
        await setToken()
        window.location.href = '/login'
    }
    return res
}
// getJobs
// apiRequest proxy purpose built for get jobs endpoint
export const getJobs = async(params?: {[key:string]:string}): Promise<Response> => {
    let paramStr: string = ''
    if (params) {
        Object.keys(params).forEach((key, i) => {
            paramStr += `${key}=${params[key]}`
            if (i === 0) {
                paramStr += `?`
            } else {
                paramStr += '&'
            }
            paramStr += `${key}=${params[key]}`
        })
    }
    return apiRequest(`/jobs${paramStr}`, undefined, 'GET')
}