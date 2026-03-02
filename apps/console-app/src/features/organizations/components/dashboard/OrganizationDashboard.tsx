import { Link, useRouter } from '@tanstack/react-router'
import { Building2Icon, FolderOpen, Plus, Users } from 'lucide-react'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { formatDate } from '@/lib/utils'
import type { Organization } from '@/features/organizations/types'

export function OrganizationDashboard({ organizations }: { organizations: Organization[] }) {
    return (
        <>
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-semibold tracking-tight">Your Organizations</h1>
                    <p className="text-muted-foreground text-sm mt-1">Select a workspace to enter.</p>
                </div>

                <Button variant="outline" asChild>
                    <Link to="/onboarding">
                        <Plus className="h-4 w-4 mr-2" />
                        New Organization
                    </Link>
                </Button>
            </div>

            <Separator />

            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
                {organizations.map((org) => (
                    <OrganizationCard key={org.id} org={org} />
                ))}
            </div>
        </>
    )
}

function OrganizationCard({ org }: { org: Organization }) {
    const router = useRouter()

    return (
        <Card
            key={org.id}
            className="group cursor-pointer transition-all hover:border-primary hover:shadow-md"
            onClick={() => router.navigate({ to: `/org/$slug`, params: { slug: org.slug } })}
        >
            <CardHeader className="pb-3">
                <div className="flex items-start justify-between gap-2">
                    <CardTitle className="text-base leading-snug line-clamp-1">{org.name}</CardTitle>
                    {!org.isActive && (
                        <Badge variant="secondary" className="shrink-0 text-xs">
                            Inactive
                        </Badge>
                    )}
                </div>

                {org.description && (
                    <CardDescription className="line-clamp-2 text-xs mt-1">{org.description}</CardDescription>
                )}
            </CardHeader>

            <CardContent className="pt-0 space-y-3">
                <div className="flex gap-4 text-xs text-muted-foreground">
                    <span className="flex items-center gap-1">
                        <Users className="h-3.5 w-3.5" />
                        {org.memberCount}/{org.maxMembers} members
                    </span>
                    <span className="flex items-center gap-1">
                        <FolderOpen className="h-3.5 w-3.5" />
                        {org.maxProjects} projects
                    </span>
                </div>

                <p className="text-xs text-muted-foreground">Created {formatDate(org.createdAt)}</p>

                <Button
                    variant="secondary"
                    size="sm"
                    className="w-full group-hover:bg-primary group-hover:text-primary-foreground transition-colors"
                >
                    Enter workspace
                </Button>
            </CardContent>
        </Card>
    )
}

export function EmptyOrganization() {
    return (
        <div className="mx-auto max-w-2xl py-20 px-4 text-center space-y-6">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full border border-border bg-card mb-2">
                <Building2Icon className="w-7 h-7 text-muted-foreground" />
            </div>

            <div className="space-y-2">
                <h1 className="text-2xl font-semibold tracking-tight">Welcome to Worklayer</h1>
                <p className="text-muted-foreground">
                    Create your first organization to start managing projects, API keys and workflows.
                </p>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Create Organization</CardTitle>
                    <CardDescription>An organization is your workspace container.</CardDescription>
                </CardHeader>
                <CardContent>
                    <Button asChild className="w-full">
                        <Link to="/onboarding">
                            <Plus className="h-4 w-4 mr-2" />
                            Create First Organization
                        </Link>
                    </Button>
                </CardContent>
            </Card>
        </div>
    )
}
