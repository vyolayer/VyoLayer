import { Card, CardContent } from '@/components/ui/card'
import { FolderOpen, Users, Calendar } from 'lucide-react'
import { formatDate } from '@/lib/utils'
import { useOrganization } from '@/contexts/OrganizationContext'

export const OrganizationStats = () => {
    const { org } = useOrganization()

    return (
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
    )
}
