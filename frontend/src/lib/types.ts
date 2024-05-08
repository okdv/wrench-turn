import { apiRequest } from "$lib/api"
export class NewJob {
    name: string|null
    description: string|null
    instructions: string|null
    isTemplate: number
    vehicle: number|null
    originJob: number|null
    repeats: number
    odoInterval: number|null
    timeInterval: number|null
    timeIntervalUnit: 'month' | 'day' | 'week' | 'hour' | 'year' | null
    dueDate: string|null
    constructor() {
        this.name = null
        this.description = null 
        this.instructions = null 
        this.isTemplate = 0 
        this.vehicle = null 
        this.originJob = null 
        this.repeats = 0 
        this.odoInterval = null 
        this.timeInterval = null 
        this.timeIntervalUnit = null
        this.dueDate = null
    }
}

export type Job = 	{
    id: number,
    name: string,
    description: string|null,
    instructions: string|null,
    isTemplate: number,
    isComplete: number,
    vehicle: number|null,
    user: number,
    originJob: number|null,
    labels: Array<Label> | null
    repeats: number,
    odoInterval: number,
    timeInterval: number,
    timeIntervalUnit: 'month' | 'day' | 'week' | 'hour' | 'year',
    dueDate: number|null,
    completedAt: string|null,
    createdAt: string,
    updatedAt: string,
}

export class NewTask {
    name: string
    description: string|null
    partName: string|null
    partLink: string|null
    dueDate: string|null
    constructor(name?: string, description?: string | null, partName?: string | null, partLink?: string | null, dueDate?: string | null) {
        this.name = name ?? ""
        this.description = description ?? null 
        this.partName = partName ?? null 
        this.partLink = partLink ?? null
        this.dueDate = dueDate ?? null
    }
}

// Class for handling Tasks, new and existing
export class Task {
    id?: number
    name: string
    description: string|null
    isComplete: number
    job: number
    partName: string|null
    partLink: string|null
    dueDate: string|null
    readonly completedAt?: string|null
    readonly createdAt?: string
    readonly updatedAt?: string
    constructor(task: Task) {
        this.id = task.id
        this.job = task.job
        this.name = task.name ?? ""
        this.description = task.description ?? null 
        this.isComplete = task.isComplete ?? 0
        this.partName = task.partName ?? null 
        this.partLink = task.partLink ?? null
        this.dueDate = task.dueDate ?? null
        this.completedAt = task.completedAt
        this.createdAt = task.createdAt
        this.updatedAt = task.updatedAt
    }
    // method for marking task complete via api
    async markComplete(): Promise<void> {
        try {
            console.log(this)
            // throw error if it has no id
            if (!this.id) {
                throw new Error("This task does not have an ID and likely needs to be created before it can be updated")
            }
            // update via api
            const res = await apiRequest(`/jobs/${this.job}/tasks/${this.id}/complete?incomplete=${this.isComplete === 1}`, undefined, "PATCH", true)
            if (!res.ok) {
                throw new Error(`Unable to mark task complete: HTTP ${res.status} - ${res.text()}`)
            }
            // update locally instead of calling new data from api
            this.isComplete = 1
        } catch (error) {
            console.error("Error marking task complete: ", error)
            throw error
        }
    }
}

export class NewVehicle {
    name: string
    description: string|null
    type: string|null
    isMetric: number
    vin: string|null
    year: number|null
    make: string|null
    model: string|null
    trim: string|null
    odometer: number|null
    user: number|null
    constructor() {
        this.name = ""
        this.description = null 
        this.type = null 
        this.isMetric = 0 
        this.vin = null 
        this.year = null 
        this.make = null
        this.model = null 
        this.trim = null 
        this.odometer = null
        this.user = null
    }
}

export type Vehicle = 	{
    id: number,
    name: string,
    description: string|null,
    type: string|null,
    isMetric: number,
    vin: string|null,
    year: number|null,
    make: string|null,
    model: string|null,
    trim: string|null,
    odometer: number|null,
    user: number|null,
    createdAt: string,
    updatedAt: string,
}

type AlertType = 'notification' | 'reminder'
export class NewAlert {
    name: string
    description: string|null
    type: AlertType
    user: number
    vehicle: number|null
    job: number|null
    task: number|null
    constructor(name?: string, description?: string | null, type?: AlertType, user?: number, vehicle?: number | null, job?: number | null, task?: number | null) {
        this.name = name ?? ""
        this.description = description ?? null 
        this.type = type ?? 'notification' 
        this.user = user ?? 0
        this.vehicle = user ?? null
        this.job = user ?? null
        this.task = user ?? null
    }
}

export type Alert = 	{
    id: number,
    name: string,
    color: string|null,
    user: number,
    createdAt: string,
    updatedAt: string,
}

export class NewLabel {
    name: string;
    color: string|null;
    user: number|null;
    constructor(name?: string, color?: string | null, user?: number|null) {
        this.name = name ?? ""
        this.color = color ?? "#1d4ed8" 
        this.user = user ?? null
    }
}

export type Label = 	{
    id: number,
    name: string;
    color: string|null;
    user: number|null;
    createdAt: string,
    updatedAt: string,
}


export class User {
    id: number
    username: string
    email: string|null
    description: string|null
    hashedPw: string|null
    isAdmin: number
    createdAt: string
    updatedAt: string 
    constructor() {
        this.id = 0
        this.username = ''
        this.email = '' 
        this.description = '' 
        this.hashedPw = null 
        this.isAdmin = 0 
        this.createdAt = ''
        this.updatedAt = ''
    }
}

export type JWTPayload = {
    id: string,
    username: string,
    isAdmin: string,
    exp: number
}

export class UpdatePassword {
    password: string 
    newPassword: string
    confirmPassword: string 
    constructor() {
        this.password = ''
        this.newPassword = ''
        this.confirmPassword = ''
    }
}