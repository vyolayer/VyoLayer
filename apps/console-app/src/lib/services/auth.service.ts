import { api } from '../api'

export interface User {
    id: string
    fullName: string
    email: string
    status: 'active' | 'inactive'
    isEmailVerified: boolean
    joinedAt: string
}

export interface LoginResponse {
    accessToken: string
    refreshToken: string
    user: User
}

export interface RegisterPayload {
    email: string
    password: string
    fullName: string
}

export interface LoginPayload {
    email: string
    password: string
}

// ─── Auth API Functions ────────────────────────────────────────────

export async function registerUser(payload: RegisterPayload): Promise<User> {
    return api.post<User>('/auth/register', payload)
}

export async function loginUser(payload: LoginPayload): Promise<LoginResponse> {
    return api.post<LoginResponse>('/auth/login', payload)
}

export async function logoutUser(): Promise<void> {
    return api.post<void>('/auth/logout')
}

export async function validateSession(): Promise<void> {
    return api.post<void>('/auth/validate')
}

export async function refreshSession(): Promise<void> {
    return api.post<void>('/auth/refresh')
}

export async function changeUserPassword(payload: { currentPassword: string; newPassword: string }): Promise<void> {
    return api.post<void>('/auth/change-password', payload)
}
