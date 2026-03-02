import { Link } from '@tanstack/react-router'
import {
    Building2,
    Users,
    FolderOpen,
    Calendar,
    AlertCircle,
    Crown,
    CheckCircle2,
    XCircle,
    ArrowLeft,
    Settings,
} from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { useOrganizationBySlug } from '@/features/organizations/hooks/useOrganizationBySlug'
import { useAuth } from '@/contexts/AuthContext'
import { useMyMembershipInOrg } from '@/features/organizations/hooks/useOrganizationMember'

function formatDate(iso: string) {
    return new Date(iso).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
    })
}

interface OrganizationPageProps {
    slug: string
}

export default function OrganizationPage({ slug }: OrganizationPageProps) {
    const { user } = useAuth()
    const { data: org, isLoading: orgLoading, error: orgError } = useOrganizationBySlug(slug)
    const { data: member, isLoading: memberLoading, error: memberError } = useMyMembershipInOrg(org?.id)
    console.log(member)

    // ── Loading ───────────────────────────────────────────────────
    if (orgLoading) {
        return (
            <div className="mx-auto max-w-5xl py-10 px-4 space-y-6">
                <div className="h-8 w-64 animate-pulse bg-muted rounded" />
                <div className="h-4 w-48 animate-pulse bg-muted rounded" />
                <div className="grid gap-4 sm:grid-cols-3 mt-8">
                    {[1, 2, 3].map((n) => (
                        <div key={n} className="h-28 animate-pulse bg-muted rounded-xl" />
                    ))}
                </div>
            </div>
        )
    }

    // ── Error ─────────────────────────────────────────────────────
    if (orgError || !org) {
        return (
            <div className="mx-auto max-w-xl py-20 px-4 text-center space-y-4">
                <AlertCircle className="mx-auto h-10 w-10 text-destructive" />
                <h2 className="text-lg font-semibold">Organization not found</h2>
                <p className="text-sm text-muted-foreground">
                    {orgError instanceof Error ? orgError.message : `No organization found with slug "${slug}".`}
                </p>
                <Button variant="outline" asChild>
                    <Link to="/">
                        <ArrowLeft className="h-4 w-4 mr-2" />
                        Back to home
                    </Link>
                </Button>
            </div>
        )
    }

    const isOwner = user?.id === org.ownerId

    return (
        <div className="mx-auto max-w-7xl py-10 px-4 space-y-8">
            {/* ── Breadcrumb ── */}
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Link to="/" className="hover:text-foreground transition-colors">
                    Home
                </Link>
                <span>/</span>
                <span className="text-foreground font-medium">{org.name}</span>
            </div>

            {/* ── Header ── */}
            <div className="flex items-start justify-between gap-4 flex-wrap">
                <div className="flex items-center gap-4">
                    <div className="flex items-center justify-center w-14 h-14 rounded-2xl border border-border bg-card">
                        <Building2 className="w-7 h-7 text-muted-foreground" />
                    </div>

                    <div>
                        <div className="flex items-center gap-2 flex-wrap">
                            <h1 className="text-2xl font-bold tracking-tight">{org.name}</h1>
                            {isOwner && (
                                <Badge variant="secondary" className="gap-1">
                                    <Crown className="h-3 w-3" />
                                    Owner
                                </Badge>
                            )}
                            {org.isActive ? (
                                <Badge variant="success" className="gap-1">
                                    <CheckCircle2 className="h-3 w-3" />
                                    Active
                                </Badge>
                            ) : (
                                <Badge variant="destructive" className="gap-1">
                                    <XCircle className="h-3 w-3" />
                                    Inactive
                                </Badge>
                            )}
                        </div>

                        {org.description && (
                            <p className="text-muted-foreground text-sm mt-1 max-w-xl">{org.description}</p>
                        )}
                    </div>
                </div>

                {isOwner && (
                    <Button variant="outline" size="sm" asChild>
                        <Link to="/settings">
                            <Settings className="h-4 w-4 mr-2" />
                            Settings
                        </Link>
                    </Button>
                )}
            </div>

            {/* ── Stats ── */}
            <div className="grid gap-4 sm:grid-cols-3">
                <Card>
                    <CardContent className="flex items-center gap-4 pt-6">
                        <div className="flex items-center justify-center w-10 h-10 rounded-lg bg-primary/10">
                            <Users className="w-5 h-5 text-primary" />
                        </div>
                        <div>
                            <p className="text-2xl font-bold">
                                {org.memberCount}
                                <span className="text-sm font-normal text-muted-foreground">/{org.maxMembers}</span>
                            </p>
                            <p className="text-xs text-muted-foreground">Members</p>
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="flex items-center gap-4 pt-6">
                        <div className="flex items-center justify-center w-10 h-10 rounded-lg bg-primary/10">
                            <FolderOpen className="w-5 h-5 text-primary" />
                        </div>
                        <div>
                            <p className="text-2xl font-bold">{org.maxProjects}</p>
                            <p className="text-xs text-muted-foreground">Max Projects</p>
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="flex items-center gap-4 pt-6">
                        <div className="flex items-center justify-center w-10 h-10 rounded-lg bg-primary/10">
                            <Calendar className="w-5 h-5 text-primary" />
                        </div>
                        <div>
                            <p className="text-sm font-semibold">{formatDate(org.createdAt)}</p>
                            <p className="text-xs text-muted-foreground">Created</p>
                        </div>
                    </CardContent>
                </Card>
            </div>

            <Separator />

            {/* ── Members ── */}
            <div className="space-y-4">
                <div className="flex items-center justify-between">
                    <h2 className="text-lg font-semibold">Members</h2>
                    {isOwner && (
                        <Button size="sm" variant="outline" disabled>
                            Invite Member
                        </Button>
                    )}
                </div>

                {org.members.length === 0 ? (
                    <Card>
                        <CardContent className="py-10 text-center text-muted-foreground text-sm">
                            No members yet.
                        </CardContent>
                    </Card>
                ) : (
                    <Card>
                        <CardContent className="p-0">
                            <ul className="divide-y divide-border">
                                {org.members.map((member) => (
                                    <li key={member.id} className="flex items-center justify-between px-4 py-3 gap-3">
                                        <div className="flex items-center gap-3 min-w-0">
                                            <Avatar className="h-8 w-8 shrink-0">
                                                <AvatarFallback className="text-xs font-medium bg-primary/10 text-primary">
                                                    {member.fullName.charAt(0).toUpperCase()}
                                                </AvatarFallback>
                                            </Avatar>
                                            <div className="min-w-0">
                                                <p className="text-sm font-medium truncate">{member.fullName}</p>
                                                <p className="text-xs text-muted-foreground truncate">{member.email}</p>
                                            </div>
                                        </div>

                                        <div className="flex items-center gap-2 shrink-0">
                                            {member.userId === org.ownerId && (
                                                <Badge variant="secondary" className="gap-1 text-xs">
                                                    <Crown className="h-2.5 w-2.5" />
                                                    Owner
                                                </Badge>
                                            )}
                                            {!member.isActive && (
                                                <Badge variant="outline" className="text-xs text-muted-foreground">
                                                    Inactive
                                                </Badge>
                                            )}
                                            <p className="text-xs text-muted-foreground hidden sm:block">
                                                Joined {formatDate(member.joinedAt)}
                                            </p>
                                        </div>
                                    </li>
                                ))}
                            </ul>
                        </CardContent>
                    </Card>
                )}
            </div>

            {/* ── Danger zone (owner only) ── */}
            {isOwner && (
                <>
                    <Separator />
                    <Card className="border-destructive/30">
                        <CardHeader>
                            <CardTitle className="text-base text-destructive">Danger Zone</CardTitle>
                            <CardDescription>Destructive actions that cannot be undone.</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <Button variant="destructive" size="sm" disabled>
                                Deactivate Organization
                            </Button>
                        </CardContent>
                    </Card>
                </>
            )}
        </div>
    )
}
