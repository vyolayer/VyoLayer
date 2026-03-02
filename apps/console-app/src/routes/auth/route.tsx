import { createFileRoute, Outlet, Navigate } from '@tanstack/react-router'
import { useAuth } from '@/contexts/AuthContext'
import { Loader2 } from 'lucide-react'

export const Route = createFileRoute('/auth')({
    component: AuthLayout,
})

function AuthLayout() {
    const { isAuthenticated, isLoading } = useAuth()

    // Show loading while checking auth state
    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-slate-950 via-slate-900 to-slate-950">
                <Loader2 className="w-8 h-8 text-cyan-400 animate-spin" />
            </div>
        )
    }

    // Redirect to home if already authenticated
    if (isAuthenticated) {
        return <Navigate to="/" />
    }

    return <Outlet />
}
