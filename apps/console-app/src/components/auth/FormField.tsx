import type { ReactNode } from 'react'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

interface FormFieldProps {
    /** The `id` linked to the input, also becomes the label's `htmlFor`. */
    id: string
    label: string
    error?: string
    hint?: string
    className?: string
    children: ReactNode
}

/**
 * Consistent form field wrapper used across all auth forms.
 * Renders: label → input slot → error message → hint text
 */
export function FormField({ id, label, error, hint, className, children }: FormFieldProps) {
    return (
        <div className={cn('flex flex-col gap-1.5', className)}>
            <Label htmlFor={id}>{label}</Label>
            {children}
            {error && (
                <p className="text-sm text-destructive" role="alert">
                    {error}
                </p>
            )}
            {hint && !error && <p className="text-xs text-muted-foreground">{hint}</p>}
        </div>
    )
}
