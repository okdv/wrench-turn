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
    timeIntervalUnit: 'month' | 'day' | 'week' | 'hour' | null
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
    repeats: number,
    odoInterval: number,
    timeInterval: number,
    timeIntervalUnit: 'month' | 'day' | 'week' | 'hour',
    dueDate: number|null,
    completedAt: string|null,
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