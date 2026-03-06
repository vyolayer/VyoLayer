import { useOrganization } from '@/contexts/OrganizationContext'
import { useProjects } from '../hooks/useProjects'
import { CreateProjectDialog } from '../components/CreateProjectDialog'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { FolderOpen, Users, Key } from 'lucide-react'
import { Link } from '@tanstack/react-router'

export function ProjectsOverviewPage() {
    const { org } = useOrganization()
    const { data: projects, isLoading } = useProjects(org.id)

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-2xl font-bold tracking-tight">Projects</h2>
                    <p className="text-sm text-muted-foreground">Manage all projects in this organization.</p>
                </div>
                <CreateProjectDialog orgId={org.id} />
            </div>

            {isLoading ? (
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="h-36 animate-pulse rounded-xl bg-muted" />
                    ))}
                </div>
            ) : !projects?.length ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <FolderOpen className="h-12 w-12 text-muted-foreground/50 mb-4" />
                    <h3 className="text-lg font-semibold">No projects yet</h3>
                    <p className="text-sm text-muted-foreground max-w-sm mb-4">
                        Create your first project to start managing API keys and members.
                    </p>
                    <CreateProjectDialog orgId={org.id} />
                </div>
            ) : (
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    {projects.map((project) => (
                        <Link
                            key={project.id}
                            to="/org/$slug/projects/$projectId"
                            params={{ slug: org.slug, projectId: project.id }}
                            className="group"
                        >
                            <Card className="transition-all hover:shadow-md hover:border-primary/30 group-hover:scale-[1.01]">
                                <CardContent className="pt-6 space-y-3">
                                    <div className="flex items-start justify-between">
                                        <div className="flex items-center gap-3">
                                            <div className="flex items-center justify-center w-10 h-10 rounded-lg bg-primary/10">
                                                <FolderOpen className="w-5 h-5 text-primary" />
                                            </div>
                                            <div>
                                                <p className="font-semibold">{project.name}</p>
                                                <p className="text-xs text-muted-foreground">{project.slug}</p>
                                            </div>
                                        </div>
                                        {!project.isActive && <Badge variant="secondary">Archived</Badge>}
                                    </div>

                                    {project.description && (
                                        <p className="text-sm text-muted-foreground line-clamp-2">
                                            {project.description}
                                        </p>
                                    )}

                                    <div className="flex gap-4 text-xs text-muted-foreground pt-1">
                                        <span className="flex items-center gap-1">
                                            <Users className="h-3 w-3" />
                                            {project.memberCount}/{project.maxMembers} members
                                        </span>
                                        <span className="flex items-center gap-1">
                                            <Key className="h-3 w-3" />
                                            {project.maxApiKeys} keys max
                                        </span>
                                    </div>
                                </CardContent>
                            </Card>
                        </Link>
                    ))}
                </div>
            )}
        </div>
    )
}
