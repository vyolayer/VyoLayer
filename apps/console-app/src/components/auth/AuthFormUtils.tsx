import type { ReactNode } from 'react'

interface ServerErrorProps {
    message: string | null
}

/**
 * Displays an API / server-level error banner inside auth forms.
 * Returns null when there is no message so it can be safely rendered unconditionally.
 */
export function ServerError({ message }: ServerErrorProps) {
    if (!message) return null

    return (
        <div
            className="mb-5 rounded-lg border border-destructive/20 bg-destructive/10 px-4 py-3 text-sm text-destructive"
            role="alert"
        >
            {message}
        </div>
    )
}

interface LoadingButtonChildrenProps {
    isLoading: boolean
    loadingText: string
    children: ReactNode
}

/**
 * Renders either a spinner + loading text or the normal button children.
 * Use inside a `<Button type="submit" disabled={isLoading}>`.
 */
export function LoadingButtonChildren({ isLoading, loadingText, children }: LoadingButtonChildrenProps) {
    if (!isLoading) return <>{children}</>

    return (
        <>
            <svg className="w-4 h-4 animate-spin" viewBox="0 0 24 24" fill="none" aria-hidden>
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
            </svg>
            {loadingText}
        </>
    )
}
