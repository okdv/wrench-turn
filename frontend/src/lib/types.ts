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