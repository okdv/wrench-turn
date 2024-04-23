import { PUBLIC_API_URL } from '$env/static/public'
import type { JWTPayload } from '$lib/types'

// getJWTData
// return json object of JWT payload
export const getJWTData = async(jwt: string): Promise<JWTPayload> => JSON.parse(atob(jwt.split('.')[1]))
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
        localStorage.removeItem('wrenchturn-jwt-expiration')
    } else {
        localStorage.setItem('wrenchturn-jwt', jwt)
        // extract payload portion of jwt, decode base64, parse into json, save expiration to localstorage
        const jwtPayload = await getJWTData(jwt)
        localStorage.setItem('wrenchturn-jwt-expiration', jwtPayload.exp.toString())
    }
    const token = await getToken() 
    if (token === jwt) {
        return true
    }
    return false
}

// apiRequest
// fetch proxy purpose built for api requests
export const apiRequest = async(endpoint: string, body?: unknown, method?: string, redirectOnFail?: boolean, preventRefresh?: boolean): Promise<Response> => {
    // initialize mark for refresh, default to false - dont this way so refresh can take place after api call
    let refresh = false
    // get jwt from localstorage
    const jwt = await getToken()
    // initialize headers for fetch
    const headers: {[key:string]: string} = {}
    // if req body, add json content type
    if (body) {
        headers['Content-Type'] = 'application/json'
    }
    // if jwt isnt null, add auth header with it as bearer, check if it expires within 24 hours, mark for refresh
    if (jwt) {
        headers['Authorization'] = `Bearer ${jwt}`
        // if refreshing isnt prevented
        if (preventRefresh !== true) {
            // get expire time from localstorage and current unix timestamp
            const expireUnixTimestamp = Number(localStorage.getItem('wrenchturn-jwt-expiration'))
            const currentUnixTimestamp = Date.now() 
            const timeLeft = expireUnixTimestamp - currentUnixTimestamp
            // mark for refresh if less than 24 hours from each other
            if (timeLeft < 86400 && timeLeft > 0) {            
                refresh = true
            }
        }
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
    // if marked for refresh 
    if (refresh === true) {
        // call refresh api endpoint
        // preventrefresh must be true to prevent endless refresh loop
        const res = await apiRequest("/refresh", null, "GET", false, true)
        // if 200 response, save new jwt via setToken
        if (res.status === 200) {
            const json = await res.json()
            await setToken(json.Value)
        }
    }
    return res
}
// verifyToken 
// return boolean that checks if JWT is valid via api
export const verifyToken =async (): Promise<boolean> => {
    const res = await apiRequest("/verify", null, "GET", false, false)
    if (res.status === 200) {
        return true
    }
    return false 
}
// paramStrConstruct
// generate param string from array, e.g. ["foo":"bar","bar":"foo"] => ?foo=bar&bar=foo
export const paramStrConstruct = async(params?: {[key:string]:string}): Promise<string> => {
    let paramStr: string = ''
    if (params) {
        Object.keys(params).forEach((key, i) => {
            if (i === 0) {
                paramStr += `?`
            } else {
                paramStr += '&'
            }
            paramStr += `${key}=${params[key]}`
        })
    }
    return paramStr
}
// getEnv
// apiRequest proxy purpose built for get env endpoint
export const getEnv = async(): Promise<Response> => {
    return apiRequest(`/env`, undefined, 'GET')
}
// getVehicles
// apiRequest proxy purpose built for get vehicles endpoint
export const getVehicles = async(params?: {[key:string]:string}): Promise<Response> => {
    const paramStr = await paramStrConstruct(params)
    return apiRequest(`/vehicles${paramStr}`, undefined, 'GET')
}
// getJobs
// apiRequest proxy purpose built for get jobs endpoint
export const getJobs = async(params?: {[key:string]:string}): Promise<Response> => {
    const paramStr = await paramStrConstruct(params)
    return apiRequest(`/jobs${paramStr}`, undefined, 'GET')
}
// getTasks
// apiRequest proxy purpose built for get tasks by job endpoint
export const getTasks = async(jobId: number, params?: {[key:string]:string}): Promise<Response> => {
    const paramStr = await paramStrConstruct(params)
    return apiRequest(`/jobs/${jobId}/tasks${paramStr}`, undefined, 'GET')
}
// getUsers
// apiRequest proxy purpose built for get users endpoint
export const getUsers = async(params?: {[key:string]:string}): Promise<Response> => {
    const paramStr = await paramStrConstruct(params)
    return apiRequest(`/users${paramStr}`, undefined, 'GET')
}
// getUser 
// apiRequest proxy purpose built for get user endpoint
export const getUser = async(username?: string): Promise<Response> => {
    if (!username) {
        const jwt = await getToken()
        if (jwt) {
            const jwtData = await getJWTData(jwt)
            username = jwtData.username
        }
    }
    return apiRequest(`/users/${username}`, undefined, 'GET', true, false)
}
// getAlerts
// apiRequest proxy purpose built for get alerts endpoint
export const getAlerts = async(params?: {[key:string]:string}): Promise<Response> => {
    const paramStr = await paramStrConstruct(params)
    return apiRequest(`/alerts${paramStr}`, undefined, 'GET')
}
// updateAlertReadStatus
// apiRequest proxy purpose build for updated read status of an alert
export const updateAlertReadStatus = async(alertId: number, unread?: boolean): Promise<Response> => {
    if (unread === true) {
        return apiRequest(`/alerts/${alertId}/read?unread=true`, undefined, 'PATCH')
    }
    return apiRequest(`/alerts/${alertId}/read`, undefined, 'PATCH')
}
// getLabels
// apiRequest proxy purpose built for list labels endpoint
export const getLabels = async(params?: {[key:string]:string}): Promise<Response> => {
    const paramStr = await paramStrConstruct(params)
    return apiRequest(`/labels${paramStr}`, undefined, 'GET')
}