import { cn } from '@/lib/utils'
import * as React from 'react'

export const Container: React.FC<{ className?: string; children: React.ReactNode }> = ({ className, children }) => {
    return <div className={cn('mx-auto py-10 max-w-7xl px-4 space-y-8', className)}>{children}</div>
}
