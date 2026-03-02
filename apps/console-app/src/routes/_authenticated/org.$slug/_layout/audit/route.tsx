import { cn } from '@/lib/utils'
import { AlertCircle, ScrollText } from 'lucide-react'
import { useOrganization } from '@/contexts/OrganizationContext'
import { createFileRoute, Link, Outlet, useRouterState } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/audit')({
    component: AuditLayout,
})

function AuditLayout() {
    const { org, isAdmin } = useOrganization()
    const routerState = useRouterState()
    const currentPath = routerState.location.pathname

    if (!isAdmin) {
        return (
            <div className="py-20 text-center space-y-4">
                <AlertCircle className="mx-auto h-10 w-10 text-muted-foreground" />
                <p className="text-muted-foreground">You don't have permission to view audit logs.</p>
            </div>
        )
    }

    const basePath = `/org/${org.slug}/audit`
    const tabs = [
        { label: 'Activity', to: basePath },
        { label: 'Security', to: `${basePath}/security` },
        { label: 'Invitations', to: `${basePath}/invitations` },
    ]

    return (
        <div className="space-y-6">
            <div className="flex items-center gap-3">
                <ScrollText className="h-5 w-5 text-muted-foreground" />
                <div>
                    <h2 className="text-lg font-semibold">Audit Log</h2>
                    <p className="text-sm text-muted-foreground">Track organization actions and events.</p>
                </div>
            </div>

            {/* ── Sub-tabs ── */}
            <div className="flex gap-1 rounded-lg bg-muted p-1 w-fit">
                {tabs.map((tab) => {
                    const isActive =
                        tab.to === basePath
                            ? currentPath === basePath || currentPath === basePath + '/'
                            : currentPath.startsWith(tab.to)

                    return (
                        <Link
                            key={tab.to}
                            to={tab.to}
                            className={cn(
                                'px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
                                isActive
                                    ? 'bg-background text-foreground shadow-sm'
                                    : 'text-muted-foreground hover:text-foreground',
                            )}
                        >
                            {tab.label}
                        </Link>
                    )
                })}
            </div>

            <Outlet />
        </div>
    )
}
