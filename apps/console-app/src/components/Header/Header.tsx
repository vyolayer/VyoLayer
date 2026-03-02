import { useRouter } from '@tanstack/react-router'
import { useAuth } from '@/contexts/AuthContext'
import { HeaderAuth } from './HeaderAuth'
import { Logo } from '../logo/Logo'
import { useCallback } from 'react'

/**
 * App-wide header: brand logo + auth controls + slide-in sidebar.
 * State management (sidebar open/closed, logout handler) lives here
 * and is passed down as props to the sub-components.
 */
export default function Header() {
    const { isAuthenticated, user, logout, isLoading } = useAuth()
    const router = useRouter()

    const handleLogout = useCallback(async () => {
        await logout()
        router.navigate({ to: '/auth/login', replace: true })
    }, [logout, router])

    return (
        <header className="sticky top-0 z-30 w-full border-b bg-card backdrop-blur supports-[backdrop-filter]:bg-card/80">
            <div className="mx-auto flex h-14 w-full max-w-7xl items-center justify-between px-6">
                {/* Left: Brand */}
                <Logo />

                {/* Right: Auth Controls */}
                <div className="flex items-center gap-3">
                    {!isLoading ? (
                        <HeaderAuth isAuthenticated={isAuthenticated} user={user} onLogout={handleLogout} />
                    ) : (
                        <div className="h-9 w-9 animate-pulse rounded-full bg-muted" />
                    )}
                </div>
            </div>
        </header>
    )
}
