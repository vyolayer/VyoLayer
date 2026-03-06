import { Link, Outlet, useRouterState } from '@tanstack/react-router'
import { cn } from '@/lib/utils'
import { LayoutDashboard, Users, Key, Settings } from 'lucide-react'
import { useOrganization } from '@/contexts/OrganizationContext'

interface ProjectLayoutProps {
    projectId: string
}

const navItemDefs = [
    { label: 'Overview', segment: '', icon: LayoutDashboard },
    { label: 'Members', segment: '/members', icon: Users },
    { label: 'API Keys', segment: '/api-keys', icon: Key },
    { label: 'Settings', segment: '/settings', icon: Settings },
] as const

export function ProjectLayout({ projectId }: ProjectLayoutProps) {
    const { org } = useOrganization()
    const currentPath = useRouterState().location.pathname
    const basePath = `/org/${org.slug}/projects/${projectId}`

    return (
        <div className="space-y-6">
            {/* Project sub-navigation */}
            <div className="flex gap-1 bg-muted py-1 px-2 rounded-lg">
                {navItemDefs.map((item) => {
                    const to = `${basePath}${item.segment}`
                    const isActive =
                        item.segment === ''
                            ? currentPath === basePath || currentPath === `${basePath}/`
                            : currentPath.startsWith(to)
                    const Icon = item.icon

                    return (
                        <Link
                            key={item.label}
                            to={to}
                            className={cn(
                                'flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                                isActive
                                    ? 'bg-primary/10 text-primary'
                                    : 'text-muted-foreground hover:bg-muted hover:text-foreground',
                            )}
                        >
                            <Icon className="h-4 w-4 shrink-0" />
                            {item.label}
                        </Link>
                    )
                })}
            </div>

            <Outlet />
        </div>
    )
}
