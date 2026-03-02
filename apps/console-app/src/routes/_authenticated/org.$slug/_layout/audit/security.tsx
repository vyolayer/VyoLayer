import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent } from '@/components/ui/card'
import { ShieldAlert } from 'lucide-react'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/audit/security')({
    component: AuditSecurityPage,
})

function AuditSecurityPage() {
    return (
        <Card>
            <CardContent className="py-16 text-center space-y-3">
                <ShieldAlert className="h-10 w-10 text-muted-foreground mx-auto" />
                <div>
                    <p className="text-sm font-medium">Security Events</p>
                    <p className="text-sm text-muted-foreground mt-1">
                        Security-related audit events will appear here.
                    </p>
                    <p className="text-xs text-muted-foreground mt-3">
                        Tracks: ownership transfers, org deletions, archive/restore, and suspicious access.
                    </p>
                </div>
            </CardContent>
        </Card>
    )
}
