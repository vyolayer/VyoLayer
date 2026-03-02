import type { ReactNode } from 'react'

/**
 * Shared ambient-glow background used by both Login and Register pages.
 * Wrap page content with this to get consistent auth page aesthetics.
 */
export function AuthPageLayout({ children }: { children: ReactNode }) {
    return (
        <div className="relative min-h-[calc(100vh-4rem)] flex items-center justify-center bg-background px-4">
            {/* Ambient glow orbs */}
            <div className="absolute inset-0 overflow-hidden pointer-events-none" aria-hidden>
                <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/5 rounded-full blur-3xl" />
                <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-primary/5 rounded-full blur-3xl" />
            </div>

            <div className="relative w-full max-w-md">{children}</div>
        </div>
    )
}
