import { Link, useRouterState } from '@tanstack/react-router'
import { useOrganization } from '@/contexts/OrganizationContext'
import { Building2, Users, UserCircle, Settings, Archive, AlertTriangle, ScrollText } from 'lucide-react'
import { cn } from '@/lib/utils'
import { OrganizationBreadcrumb } from './OrganizationBreadcrumb'

export interface NavItem {
    label: string
    to: string
    icon: typeof Building2
    requireAdmin?: boolean
    requireOwner?: boolean
}

export function OrganizationHeader() {
    const { org, isAdmin, isOwner } = useOrganization()
    const routerState = useRouterState()
    const currentPath = routerState.location.pathname

    const basePath = `/org/${org.slug}`

    const navItems: NavItem[] = [
        { label: 'Overview', to: basePath, icon: Building2 },
        { label: 'Members', to: `${basePath}/members`, icon: Users },
        { label: 'My Membership', to: `${basePath}/me`, icon: UserCircle },
        { label: 'Settings', to: `${basePath}/settings`, icon: Settings, requireAdmin: true },
        { label: 'Archive', to: `${basePath}/archive`, icon: Archive, requireAdmin: true },
        { label: 'Audit Log', to: `${basePath}/audit`, icon: ScrollText, requireAdmin: true },
        { label: 'Danger Zone', to: `${basePath}/danger`, icon: AlertTriangle, requireOwner: true },
    ]

    const filteredItems = navItems.filter((item) => {
        if (item.requireOwner) return isOwner
        if (item.requireAdmin) return isAdmin
        return true
    })

    return (
        <>
            <OrganizationBreadcrumb
                orgName={org.name}
                orgSlug={org.slug}
                navItems={filteredItems}
                currentPath={currentPath}
            />

            <div className="flex gap-1 bg-muted py-1 px-2 rounded-lg">
                {filteredItems.map((item, index) => (
                    <OrgHeaderItem key={index} {...item} basePath={basePath} currentPath={currentPath} />
                ))}
            </div>
        </>
    )
}

function OrgHeaderItem({
    icon,
    label,
    to,
    basePath,
    currentPath,
}: NavItem & { basePath: string; currentPath: string }) {
    const Icon = icon
    const isActive = to === basePath ? currentPath === basePath || currentPath === basePath + '/' : currentPath === to

    return (
        <Link
            key={to}
            to={to}
            className={cn(
                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                isActive ? 'bg-primary/10 text-primary ' : 'text-muted-foreground hover:bg-muted hover:text-foreground',
                label == 'Danger Zone' && 'text-destructive',
                isActive && label == 'Danger Zone' && 'bg-destructive/10 text-destructive ring ring-destructive',
            )}
        >
            <Icon className="h-4 w-4 shrink-0" />
            {label}
        </Link>
    )
}
