import { useOrganization } from '@/contexts/OrganizationContext'
import { useProject } from '../hooks/useProjects'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Users, Key, Calendar, FolderOpen } from 'lucide-react'

interface ProjectDashboardPageProps {
    projectId: string
}

export function ProjectDashboardPage({ projectId }: ProjectDashboardPageProps) {
    const { org } = useOrganization()
    const { data: project, isLoading } = useProject(org.id, projectId)

    if (isLoading) {
        return (
            <div className="space-y-4">
                <div className="h-8 w-48 animate-pulse rounded bg-muted" />
                <div className="grid gap-4 sm:grid-cols-3">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="h-28 animate-pulse rounded-xl bg-muted" />
                    ))}
                </div>
            </div>
        )
    }

    if (!project) return null

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center gap-4">
                <div className="flex items-center justify-center w-12 h-12 rounded-xl bg-primary/10">
                    <FolderOpen className="w-6 h-6 text-primary" />
                </div>
                <div>
                    <h2 className="text-2xl font-bold tracking-tight">{project.name}</h2>
                    <p className="text-sm text-muted-foreground">{project.description || 'No description'}</p>
                </div>
                {!project.isActive && <Badge variant="secondary">Archived</Badge>}
            </div>

            {/* Stats */}
            <div className="grid gap-4 sm:grid-cols-3">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                        <CardTitle className="text-sm font-medium">Members</CardTitle>
                        <Users className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <p className="text-2xl font-bold">{project.memberCount}</p>
                        <p className="text-xs text-muted-foreground">of {project.maxMembers} max</p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                        <CardTitle className="text-sm font-medium">API Keys</CardTitle>
                        <Key className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <p className="text-2xl font-bold">{project.maxApiKeys}</p>
                        <p className="text-xs text-muted-foreground">max allowed</p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                        <CardTitle className="text-sm font-medium">Created</CardTitle>
                        <Calendar className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <p className="text-2xl font-bold">{new Date(project.createdAt).toLocaleDateString()}</p>
                        <p className="text-xs text-muted-foreground">project created</p>
                    </CardContent>
                </Card>
            </div>

            {/* Members preview */}
            {project.members && project.members.length > 0 && (
                <Card>
                    <CardHeader>
                        <CardTitle className="text-base">Recent Members</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-3">
                            {project.members.slice(0, 5).map((m) => (
                                <div key={m.id} className="flex items-center justify-between">
                                    <div>
                                        <p className="text-sm font-medium">{m.fullName}</p>
                                        <p className="text-xs text-muted-foreground">{m.email}</p>
                                    </div>
                                    <Badge variant="outline">{m.role}</Badge>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    )
}
