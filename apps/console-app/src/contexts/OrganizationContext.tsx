import { createContext, useContext, type ReactNode } from 'react'
import { useOrganizationBySlug } from '@/org/hooks/useOrganizationBySlug'
import { useMyMembershipInOrg } from '@/org/hooks/useOrganizationMember'
import { Loader2, AlertCircle } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'

import type { OrganizationMemberWithRBAC, OrganizationWithMembers } from '@/features/organizations/types'

interface OrgContextType {
    org: OrganizationWithMembers
    me: OrganizationMemberWithRBAC
    isOwner: boolean
    isAdmin: boolean
    isMember: boolean
}

const OrganizationContext = createContext<OrgContextType | undefined>(undefined)

export function useOrganization() {
    const ctx = useContext(OrganizationContext)
    if (!ctx) throw new Error('useOrganization must be used within OrganizationProvider')
    return ctx
}

interface ProviderProps {
    slug: string
    children: ReactNode
}

export function OrganizationProvider({ slug, children }: ProviderProps) {
    const { data: org, isLoading: orgLoading, error: orgError } = useOrganizationBySlug(slug)
    const { data: me, isLoading: meLoading } = useMyMembershipInOrg(org?.id)

    if (orgLoading || meLoading) {
        return (
            <div className="flex items-center justify-center min-h-[60vh]">
                <div className="flex flex-col items-center gap-3">
                    <Loader2 className="w-8 h-8 animate-spin text-primary" />
                    <p className="text-sm text-muted-foreground">Loading organization…</p>
                </div>
            </div>
        )
    }

    if (orgError || !org) {
        return (
            <div className="flex items-center justify-center min-h-[60vh]">
                <div className="flex flex-col items-center gap-4 text-center">
                    <AlertCircle className="w-10 h-10 text-destructive" />
                    <h2 className="text-lg font-semibold">Organization not found</h2>
                    <p className="text-sm text-muted-foreground max-w-sm">
                        {orgError instanceof Error ? orgError.message : `No organization found with slug "${slug}".`}
                    </p>
                    <Button variant="outline" asChild>
                        <Link to="/">Back to home</Link>
                    </Button>
                </div>
            </div>
        )
    }

    const roles = me?.roles ?? []
    const isOwner = roles.some((r: { name: string }) => r.name === 'Owner') ?? false
    const isAdmin = roles.some((r: { name: string }) => r.name === 'Admin' || r.name === 'Owner') ?? false
    const isMember = !!me

    return (
        <OrganizationContext.Provider value={{ org, me: me!, isOwner, isAdmin, isMember }}>
            {children}
        </OrganizationContext.Provider>
    )
}
