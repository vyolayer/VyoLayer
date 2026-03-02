// API Client for the WorkLayer backend
// Uses cookies (HttpOnly) for authentication automatically

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:6999/api/v1'

export interface ApiSuccessResponse<T = unknown> {
    success: true
    statusCode: number
    message?: string
    data?: T
    meta?: {
        requestId?: string
        timestamp?: string
    }
}

export interface ApiErrorDetail {
    code: string
    message: string
    details?: unknown
    requestId?: string
    timestamp?: string
    metadata?: Record<string, unknown>
}

export interface ApiErrorResponse {
    success: false
    statusCode: number
    error: ApiErrorDetail
}

export type ApiResponse<T = unknown> = ApiSuccessResponse<T> | ApiErrorResponse

export class ApiError extends Error {
    public code: string
    public statusCode: number
    public details?: unknown

    constructor(errorResponse: ApiErrorResponse) {
        super(errorResponse.error.message)
        this.name = 'ApiError'
        this.code = errorResponse.error.code
        this.statusCode = errorResponse.statusCode
        this.details = errorResponse.error.details
    }
}

async function handleResponse<T>(response: Response): Promise<T> {
    // 204 No Content
    if (response.status === 204) {
        return undefined as T
    }

    const json = (await response.json()) as ApiResponse<T>

    if (!json.success) {
        throw new ApiError(json)
    }

    return json.data as T
}

async function refreshTokens(): Promise<boolean> {
    try {
        const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
            method: 'POST',
            credentials: 'include',
        })
        return response.ok
    } catch {
        return false
    }
}

async function fetchWithAuth<T>(url: string, options: RequestInit = {}): Promise<T> {
    const config: RequestInit = {
        ...options,
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
    }

    let response = await fetch(`${API_BASE_URL}${url}`, config)

    // If 401, try refreshing tokens and retry once
    if (response.status === 401) {
        const refreshed = await refreshTokens()
        if (refreshed) {
            response = await fetch(`${API_BASE_URL}${url}`, config)
        }
    }

    return handleResponse<T>(response)
}

export const api = {
    get: <T>(url: string) => fetchWithAuth<T>(url),

    post: <T>(url: string, body?: unknown) =>
        fetchWithAuth<T>(url, {
            method: 'POST',
            body: body ? JSON.stringify(body) : undefined,
        }),

    put: <T>(url: string, body?: unknown) =>
        fetchWithAuth<T>(url, {
            method: 'PUT',
            body: body ? JSON.stringify(body) : undefined,
        }),

    patch: <T>(url: string, body?: unknown) =>
        fetchWithAuth<T>(url, {
            method: 'PATCH',
            body: body ? JSON.stringify(body) : undefined,
        }),

    delete: <T>(url: string, body?: unknown) =>
        fetchWithAuth<T>(url, {
            method: 'DELETE',
            body: body ? JSON.stringify(body) : undefined,
        }),
}
