import { cn } from '@/lib/utils'
import { Link } from '@tanstack/react-router'
import * as React from 'react'

export const Logo: React.FC<{ className?: string }> = ({ className }) => {
    return (
        <Link to="/">
            <img src="/wl_logo.png" alt="Logo" className={cn('h-12 w-auto', className)} />
        </Link>
    )
}
