import type { ReactNode } from 'react'
import type { User, LoginPayload, RegisterPayload } from '@/lib/services/auth.service'
import { createContext, useContext, useState, useCallback, useEffect } from 'react'
import {
    loginUser as apiLogin,
    registerUser as apiRegister,
    logoutUser as apiLogout,
    validateSession as apiValidate,
} from '@/lib/services/auth.service'

// ─── Context Types ─────────────────────────────────────────────────

interface AuthState {
    user: User | null
    isAuthenticated: boolean
    isLoading: boolean
}

interface AuthContextType extends AuthState {
    login: (payload: LoginPayload) => Promise<void>
    register: (payload: RegisterPayload) => Promise<User>
    logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

// ─── Provider ──────────────────────────────────────────────────────

export function AuthProvider({ children }: { children: ReactNode }) {
    const [state, setState] = useState<AuthState>({
        user: null,
        isAuthenticated: false,
        isLoading: true,
    })

    // Validate session on mount (cookies are sent automatically)
    useEffect(() => {
        const checkSession = async () => {
            try {
                await apiValidate()
                // If validation succeeds the user has a valid cookie session.
                // /auth/validate only returns a success message, not the user object.
                // We cache user data in localStorage from login response.
                const cachedUser = localStorage.getItem('wl_user')
                if (cachedUser) {
                    setState({
                        user: JSON.parse(cachedUser) as User,
                        isAuthenticated: true,
                        isLoading: false,
                    })
                } else {
                    // Session is valid but no cached user — still mark as authenticated
                    setState({
                        user: null,
                        isAuthenticated: true,
                        isLoading: false,
                    })
                }
            } catch {
                // Session invalid or expired
                localStorage.removeItem('wl_user')
                setState({
                    user: null,
                    isAuthenticated: false,
                    isLoading: false,
                })
            }
        }

        checkSession()
    }, [])

    const login = useCallback(async (payload: LoginPayload) => {
        const response = await apiLogin(payload)
        localStorage.setItem('wl_user', JSON.stringify(response.user))
        setState({
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
        })
    }, [])

    const register = useCallback(async (payload: RegisterPayload) => {
        const user = await apiRegister(payload)
        return user
    }, [])

    const logout = useCallback(async () => {
        try {
            await apiLogout()
        } catch {
            // Even if logout API fails, clear local state
        }
        localStorage.removeItem('wl_user')
        localStorage.removeItem('wl_avatar_bg_color')
        setState({
            user: null,
            isAuthenticated: false,
            isLoading: false,
        })
    }, [])

    return <AuthContext.Provider value={{ ...state, login, register, logout }}>{children}</AuthContext.Provider>
}

// ─── Hook ──────────────────────────────────────────────────────────

export function useAuth() {
    const context = useContext(AuthContext)
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}
