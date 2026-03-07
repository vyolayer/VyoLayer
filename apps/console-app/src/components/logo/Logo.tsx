import { useAuth } from '@/contexts/AuthContext'
import { cn } from '@/lib/utils'
import { Link } from '@tanstack/react-router'
import * as React from 'react'

export const Logo: React.FC<{ className?: string }> = ({ className }) => {
    // is authenticated
    const { isAuthenticated } = useAuth()

    const to = isAuthenticated ? '/dashboard' : '/'

    return (
        <Link to={to} className="inline-flex items-center select-none">
            <img src="/wl_logo.png" alt="Logo" className={cn('h-12 w-auto', className)} />
            <span className="text-xl font-bold tracking-tighter text-emerald-500">WorkLayer</span>
        </Link>
    )
}
