import { useAuth } from '@/contexts/AuthContext'
import { createFileRoute, Navigate, Outlet } from '@tanstack/react-router'
import { Loader2Icon } from 'lucide-react'

export const Route = createFileRoute('/_landing')({
    component: LandingPageLayout,
})

function LandingPageLayout() {
    const { isAuthenticated, isLoading } = useAuth()

    if (isLoading) {
        return (
            <div className="flex min-h-[calc(100vh-3.5rem)] items-center justify-center">
                <Loader2Icon className="w-8 h-8 text-primary animate-spin" />
            </div>
        )
    }

    if (isAuthenticated) {
        return <Navigate to="/dashboard" />
    }

    return (
        <div className="landing-page relative min-h-screen overflow-hidden">
            {/* ── Ambient glow blobs ── */}
            <div className="pointer-events-none absolute inset-0 -z-10">
                <div className="landing-blob-1 absolute -top-40 left-1/4 h-[600px] w-[600px] rounded-full opacity-20 blur-[120px]" />
                <div className="landing-blob-2 absolute top-1/3 -right-20 h-[500px] w-[500px] rounded-full opacity-15 blur-[100px]" />
                <div className="landing-blob-3 absolute bottom-1/4 -left-20 h-[400px] w-[400px] rounded-full opacity-10 blur-[80px]" />
            </div>

            <Outlet />
        </div>
    )
}
