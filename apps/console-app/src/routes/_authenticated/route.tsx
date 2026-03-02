import { createFileRoute, Outlet, Navigate } from '@tanstack/react-router'
import { useAuth } from '@/contexts/AuthContext'
import { Loader2 } from 'lucide-react'

export const Route = createFileRoute('/_authenticated')({
    component: AuthenticatedLayout,
})

function AuthenticatedLayout() {
    const { isAuthenticated, isLoading } = useAuth()

    if (isLoading) {
        return (
            <div className="flex min-h-[calc(100vh-3.5rem)] items-center justify-center">
                <Loader2 className="w-8 h-8 text-primary animate-spin" />
            </div>
        )
    }

    if (!isAuthenticated) {
        return <Navigate to="/auth/login" />
    }

    return <Outlet />
}
