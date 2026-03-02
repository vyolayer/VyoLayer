// import { useState, type ReactNode } from 'react'
// import {
//     AlertDialog,
//     AlertDialogContent,
//     AlertDialogHeader,
//     AlertDialogTitle,
//     AlertDialogDescription,
//     AlertDialogFooter,
//     AlertDialogCancel,
//     AlertDialogAction,
// } from '@/components/ui/alert-dialog'
// import { Input } from '@/components/ui/input'
// import { Loader2 } from 'lucide-react'
// import { Button } from '../ui/button'

// interface ConfirmDialogProps {
//     open: boolean
//     onOpenChange: (open: boolean) => void
//     title: string
//     description: string
//     confirmLabel?: string
//     confirmVariant?: 'default' | 'destructive'
//     /** If provided, user must type this text to confirm */
//     confirmText?: string
//     confirmPlaceholder?: string
//     isLoading?: boolean
//     onConfirm: () => void
//     children?: ReactNode
// }

// export function ConfirmDialog({
//     open,
//     onOpenChange,
//     title,
//     description,
//     confirmLabel = 'Confirm',
//     confirmVariant = 'destructive',
//     confirmText,
//     confirmPlaceholder,
//     isLoading,
//     onConfirm,
//     children,
// }: ConfirmDialogProps) {
//     const [inputValue, setInputValue] = useState('')

//     const handleOpenChange = (v: boolean) => {
//         if (!v) setInputValue('')
//         onOpenChange(v)
//     }
//     const canConfirm = confirmText ? inputValue.trim() === confirmText : true

//     return (
//         <AlertDialog open={open} onOpenChange={handleOpenChange}>
//             <AlertDialogContent>
//                 <AlertDialogHeader>
//                     <AlertDialogTitle>{title}</AlertDialogTitle>
//                     <AlertDialogDescription>{description}</AlertDialogDescription>
//                 </AlertDialogHeader>

//                 {children}

//                 {confirmText && (
//                     <div className="space-y-2 py-2">
//                         <p className="text-sm text-muted-foreground">
//                             Type <span className="font-mono font-semibold text-foreground">{confirmText}</span> to
//                             confirm.
//                         </p>
//                         <Input
//                             value={inputValue}
//                             onChange={(e) => setInputValue(e.target.value)}
//                             placeholder={confirmPlaceholder ?? confirmText}
//                             autoFocus
//                         />
//                     </div>
//                 )}

//                 <AlertDialogFooter>
//                     <AlertDialogCancel disabled={isLoading}>Cancel</AlertDialogCancel>
//                     <AlertDialogAction asChild>
//                         <Button
//                             type="button"
//                             variant="destructive"
//                             onClick={onConfirm}
//                             disabled={!canConfirm || isLoading}
//                         >
//                             {isLoading && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
//                             {confirmLabel}
//                         </Button>
//                     </AlertDialogAction>
//                 </AlertDialogFooter>
//             </AlertDialogContent>
//         </AlertDialog>
//     )
// }
import { useState, type ReactNode } from 'react'
import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogCancel,
    AlertDialogAction,
} from '@/components/ui/alert-dialog'
import { Input } from '@/components/ui/input'
import { Loader2 } from 'lucide-react'

interface ConfirmDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void

    title: string
    description?: string

    confirmLabel?: string
    confirmVariant?: 'default' | 'destructive'

    /** If provided, user must type this exact text */
    confirmText?: string
    confirmPlaceholder?: string

    /** Async-safe confirm handler */
    isLoading?: boolean
    onConfirm: () => Promise<void> | void

    children?: ReactNode
}

export function ConfirmDialog({
    open,
    onOpenChange,
    title,
    description,
    confirmLabel = 'Confirm',
    confirmVariant = 'destructive',
    confirmText,
    confirmPlaceholder,
    isLoading,
    onConfirm,
    children,
}: ConfirmDialogProps) {
    const [inputValue, setInputValue] = useState('')
    const [error, setError] = useState<string | null>(null)

    const handleOpenChange = (value: boolean) => {
        if (!value) {
            setInputValue('')
            setError(null)
        }
        onOpenChange(value)
    }
    const canConfirm = confirmText ? inputValue.trim() === confirmText.trim() : true

    const handleConfirm = async () => {
        try {
            setError(null)

            await onConfirm()

            handleOpenChange(false) // auto-close on success
        } catch (err: any) {
            setError(err?.message ?? 'Something went wrong. Please try again.')
        } finally {
        }
    }

    return (
        <AlertDialog open={open} onOpenChange={handleOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>{title}</AlertDialogTitle>
                    {description && <AlertDialogDescription>{description}</AlertDialogDescription>}
                </AlertDialogHeader>

                {children}

                {confirmText && (
                    <div className="space-y-2 py-2">
                        <p className="text-sm text-muted-foreground">
                            Type <span className="font-mono font-semibold text-foreground">{confirmText}</span> to
                            confirm.
                        </p>

                        <Input
                            value={inputValue}
                            onChange={(e) => setInputValue(e.target.value)}
                            placeholder={confirmPlaceholder ?? confirmText}
                            autoFocus
                        />
                    </div>
                )}

                {error && (
                    <div className="rounded-md border border-destructive/40 bg-destructive/10 p-3 text-sm text-destructive">
                        {error}
                    </div>
                )}

                <AlertDialogFooter>
                    <AlertDialogCancel disabled={isLoading}>Cancel</AlertDialogCancel>

                    <AlertDialogAction asChild>
                        <button
                            type="button"
                            disabled={!canConfirm || isLoading}
                            onClick={handleConfirm}
                            className={
                                confirmVariant === 'destructive'
                                    ? 'inline-flex items-center justify-center rounded-md bg-destructive px-4 py-2 text-sm font-medium text-destructive-foreground hover:bg-destructive/90 disabled:opacity-50'
                                    : 'inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-50'
                            }
                        >
                            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {confirmLabel}
                        </button>
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    )
}
