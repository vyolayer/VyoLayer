import { OrganizationStats } from '../components/overview/OrganizationStats'
import { OrganizationAdminQuickAction } from '../components/overview/OrganizationAdminQuickAction'
import { Separator } from '@/components/ui/separator'
import { useProjects } from '@/features/projects/hooks/useProjects'
import { useOrganization } from '@/contexts/OrganizationContext'
import { Card, CardContent } from '@/components/ui/card'
import { FolderOpen } from 'lucide-react'
import { Link } from '@tanstack/react-router'

export function OrganizationOverviewPage() {
    const { org } = useOrganization()
    const { data: projects, isLoading } = useProjects(org.id)
    return (
        <div className="space-y-8">
            {/* ── Stats ── */}
            <OrganizationStats />

            {/* --- Projects --- */}
            <Separator />

            <div className="space-y-4">
                <div className="flex items-center justify-between">
                    <h2 className="text-lg font-semibold">Projects</h2>
                    <Link
                        to="/org/$slug/projects"
                        params={{ slug: org.slug }}
                        className="text-sm text-primary hover:underline"
                    >
                        View all &rarr;
                    </Link>
                </div>

                {isLoading ? (
                    <div className="h-40 w-full animate-pulse bg-muted rounded-xl" />
                ) : !projects?.length ? (
                    <p className="text-sm text-muted-foreground">No projects yet.</p>
                ) : (
                    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                        {projects.slice(0, 6).map((project) => (
                            <Link
                                key={project.id}
                                to="/org/$slug/projects/$projectId"
                                params={{ slug: org.slug, projectId: project.id }}
                            >
                                <Card className="transition-all hover:shadow-md hover:border-primary/30">
                                    <CardContent className="flex items-center gap-4 pt-6">
                                        <div className="flex items-center justify-center w-10 h-10 rounded-lg bg-primary/10">
                                            <FolderOpen className="w-5 h-5 text-primary" />
                                        </div>
                                        <div>
                                            <p className="text-2xl font-bold">{project.name}</p>
                                            <p className="text-xs text-muted-foreground">Project</p>
                                        </div>
                                    </CardContent>
                                </Card>
                            </Link>
                        ))}
                    </div>
                )}
            </div>

            {/* ── Quick Actions (admins only) ── */}
            <OrganizationAdminQuickAction />
        </div>
    )
}
