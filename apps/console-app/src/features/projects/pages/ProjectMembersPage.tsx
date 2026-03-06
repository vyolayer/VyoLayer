import { useOrganization } from '@/contexts/OrganizationContext'
import { useProjectMembers, useChangeProjectMemberRole, useRemoveProjectMember } from '../hooks/useProjectMembers'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Trash2, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import type { ProjectRole } from '../types'

interface ProjectMembersPageProps {
    projectId: string
}

export function ProjectMembersPage({ projectId }: ProjectMembersPageProps) {
    const { org } = useOrganization()
    const { data: members, isLoading } = useProjectMembers(org.id, projectId)
    const changeRole = useChangeProjectMemberRole(org.id, projectId)
    const removeMember = useRemoveProjectMember(org.id, projectId)

    const handleRoleChange = (memberId: string, role: ProjectRole) => {
        changeRole.mutate(
            { memberId, payload: { role } },
            {
                onSuccess: () => toast.success('Role updated'),
                onError: (err) => toast.error(err instanceof Error ? err.message : 'Failed to update role'),
            },
        )
    }

    const handleRemove = (memberId: string) => {
        removeMember.mutate(memberId, {
            onSuccess: () => toast.success('Member removed'),
            onError: (err) => toast.error(err instanceof Error ? err.message : 'Failed to remove member'),
        })
    }

    if (isLoading) {
        return (
            <div className="space-y-3">
                {[1, 2, 3].map((i) => (
                    <div key={i} className="h-16 animate-pulse rounded-lg bg-muted" />
                ))}
            </div>
        )
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>Project Members</CardTitle>
            </CardHeader>
            <CardContent>
                {!members?.length ? (
                    <p className="text-sm text-muted-foreground py-8 text-center">No members found.</p>
                ) : (
                    <div className="divide-y">
                        {members.map((member) => (
                            <div
                                key={member.id}
                                className="flex items-center justify-between py-3 first:pt-0 last:pb-0"
                            >
                                <div className="flex items-center gap-3">
                                    <div>
                                        <p className="text-sm font-medium">{member.fullName}</p>
                                        <p className="text-xs text-muted-foreground">{member.email}</p>
                                    </div>
                                    {!member.isActive && (
                                        <Badge variant="destructive" className="text-xs">
                                            Removed
                                        </Badge>
                                    )}
                                </div>

                                <div className="flex items-center gap-2">
                                    <Select
                                        defaultValue={member.role}
                                        onValueChange={(val) => handleRoleChange(member.id, val as ProjectRole)}
                                        disabled={changeRole.isPending}
                                    >
                                        <SelectTrigger className="w-28 h-8 text-xs">
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="admin">Admin</SelectItem>
                                            <SelectItem value="member">Member</SelectItem>
                                            <SelectItem value="viewer">Viewer</SelectItem>
                                        </SelectContent>
                                    </Select>

                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        className="h-8 w-8 text-destructive hover:bg-destructive/10"
                                        onClick={() => handleRemove(member.id)}
                                        disabled={removeMember.isPending}
                                    >
                                        {removeMember.isPending ? (
                                            <Loader2 className="h-4 w-4 animate-spin" />
                                        ) : (
                                            <Trash2 className="h-4 w-4" />
                                        )}
                                    </Button>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </CardContent>
        </Card>
    )
}
