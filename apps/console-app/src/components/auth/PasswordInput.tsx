import { useState } from 'react'
import { Eye, EyeOff } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import type { ComponentProps } from 'react'

type PasswordInputProps = Omit<ComponentProps<'input'>, 'type'>

/**
 * An `<input type="password">` with a visibility toggle button built in.
 * Uses shadcn Input + Button so it automatically picks up design-token styles.
 */
export function PasswordInput({ className, ...props }: PasswordInputProps) {
    const [visible, setVisible] = useState(false)

    return (
        <div className="relative">
            <Input type={visible ? 'text' : 'password'} className={cn('pr-10', className)} {...props} />
            <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                onClick={() => setVisible((v) => !v)}
                className="absolute right-1.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                aria-label={visible ? 'Hide password' : 'Show password'}
            >
                {visible ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
            </Button>
        </div>
    )
}
