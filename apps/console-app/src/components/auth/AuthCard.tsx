import type { ReactNode } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

interface AuthCardProps {
    title: string
    subtitle: string
    children: ReactNode
    footer?: ReactNode
}

/**
 * Card shell shared between Login and Register pages.
 * Provides consistent branding header + form slot + footer link.
 */
export function AuthCard({ title, subtitle, children, footer }: AuthCardProps) {
    return (
        <Card className="shadow-2xl">
            <CardHeader className="text-center">
                <CardTitle className="text-3xl font-bold tracking-tight">{title}</CardTitle>
                <CardDescription className="text-muted-foreground mt-2">{subtitle}</CardDescription>
            </CardHeader>
            <CardContent className="pt-6">{children}</CardContent>

            {footer && (
                <div className="px-6 pb-6 pt-4 border-t border-border text-center text-sm text-muted-foreground">
                    {footer}
                </div>
            )}
        </Card>
    )
}
