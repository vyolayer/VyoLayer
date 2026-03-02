import * as React from 'react'
import { cn } from '@/lib/utils'

import { Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { RoleBadge } from '@/features/organizations/components/RoleBadge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useChangeMemberRole } from '../../hooks/useOrganizationMutations'
import type { OrganizationMember, OrganizationRole } from '@/features/organizations/types'

export const OrganizationMemberProfileRoleManagement: React.FC<{
    member: OrganizationMember
    roles: OrganizationRole[]
    isOwner: boolean
    orgId: string
}> = ({ member, roles, isOwner, orgId }) => {
    const [selectedRoleId, setSelectedRoleId] = React.useState<string>('')
    const changeRoleMutation = useChangeMemberRole(orgId)

    const handleUpdateRole = () => {
        changeRoleMutation.mutate({
            memberId: member.id,
            payload: { roleId: selectedRoleId },
        })
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-base">Change Role</CardTitle>
                <CardDescription>
                    Update this member's role. This affects their permissions within the organization.
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                <div className="grid gap-3">
                    {roles
                        .filter((r) => r.name !== 'Owner' || isOwner)
                        .map((role) => (
                            <button
                                key={role.id}
                                onClick={() => setSelectedRoleId(role.id)}
                                className={cn(
                                    'flex items-start gap-3 p-3 rounded-lg border text-left transition-colors',
                                    selectedRoleId === role.id
                                        ? 'border-primary bg-primary/5'
                                        : 'border-border hover:border-muted-foreground/30',
                                )}
                            >
                                <div className="pt-0.5">
                                    <RoleBadge name={role.name} />
                                </div>
                                <div>
                                    <p className="text-sm font-medium">{role.name}</p>
                                    <p className="text-xs text-muted-foreground">{role.description}</p>
                                </div>
                            </button>
                        ))}
                </div>
                <Button onClick={handleUpdateRole} disabled={!selectedRoleId || changeRoleMutation.isPending}>
                    {changeRoleMutation.isPending && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
                    Update Role
                </Button>
                {changeRoleMutation.isError && (
                    <p className="text-sm text-destructive">
                        {changeRoleMutation.error instanceof Error
                            ? changeRoleMutation.error.message
                            : 'Failed to change role'}
                    </p>
                )}
                {changeRoleMutation.isSuccess && <p className="text-sm text-emerald-600">Role updated successfully.</p>}
            </CardContent>
        </Card>
    )
}
