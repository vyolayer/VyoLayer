import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent } from '@/components/ui/card'
import { Activity } from 'lucide-react'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/audit/')({
    component: AuditActivityPage,
})

function AuditActivityPage() {
    return (
        <Card>
            <CardContent className="py-16 text-center space-y-3">
                <Activity className="h-10 w-10 text-muted-foreground mx-auto" />
                <div>
                    <p className="text-sm font-medium">Activity Log</p>
                    <p className="text-sm text-muted-foreground mt-1">
                        Organization activity events will appear here as they occur.
                    </p>
                    <p className="text-xs text-muted-foreground mt-3">
                        Tracks: member joins, role changes, settings updates, and more.
                    </p>
                </div>
            </CardContent>
        </Card>
    )
}
