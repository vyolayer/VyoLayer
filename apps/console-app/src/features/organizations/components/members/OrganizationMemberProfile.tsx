import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Calendar, Mail, Shield, UserMinus } from 'lucide-react'
import * as React from 'react'
import { formatDate } from '@/lib/utils'
import { RoleBadge } from '../RoleBadge'
import type { OrganizationMember } from '../../types'
import { Button } from '@/components/ui/button'
import { ConfirmDialog } from '@/components/organization/ConfirmDialog'
import { useRemoveMember } from '../../hooks/useOrganizationMutations'

export const OrganizationMemberProfile: React.FC<{
    member: OrganizationMember
    isMemberOwner: boolean
    isSelf: boolean
}> = ({ member, isMemberOwner, isSelf }) => {
    return (
        <Card>
            <CardContent className="pt-6">
                <div className="flex items-start gap-4">
                    <Avatar className="h-16 w-16 shrink-0">
                        <AvatarFallback className="text-lg font-semibold bg-primary/10 text-primary">
                            {member.fullName?.charAt(0)?.toUpperCase() ?? '?'}
                        </AvatarFallback>
                    </Avatar>
                    <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 flex-wrap">
                            <h2 className="text-xl font-bold">{member.fullName}</h2>
                            {isMemberOwner && <RoleBadge name="Owner" size="md" />}
                            {isSelf && (
                                <span className="text-xs px-2 py-0.5 rounded-full bg-primary/10 text-primary">You</span>
                            )}
                        </div>

                        <div className="mt-2 space-y-1">
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <Mail className="h-4 w-4" />
                                {member.email}
                            </div>
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <Calendar className="h-4 w-4" />
                                Joined {formatDate(member.joinedAt)}
                            </div>
                            {member.invitedAt && (
                                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                    <Shield className="h-4 w-4" />
                                    Invited {formatDate(member.invitedAt)}
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}

export const OrganizationMemberRemove: React.FC<{
    orgId: string
    member: OrganizationMember
}> = ({ member, orgId }) => {
    const [showRemoveDialog, setShowRemoveDialog] = React.useState(false)
    const removeMutation = useRemoveMember(orgId)

    const handleRemove = () => {
        removeMutation.mutate(member.id, {
            onSuccess: () => {
                setShowRemoveDialog(false)
            },
        })
    }

    return (
        <>
            {/* --- Remove Member --- */}
            <Card className="border-destructive/30">
                <CardHeader>
                    <CardTitle className="text-base text-destructive">Remove Member</CardTitle>
                    <CardDescription>
                        Remove this member from the organization. They will lose all access.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <Button variant="destructive" size="sm" onClick={() => setShowRemoveDialog(true)}>
                        <UserMinus className="w-4 h-4 mr-2" />
                        Remove {member.fullName}
                    </Button>
                </CardContent>
            </Card>

            {/* --- Confirm Dialog --- */}
            <ConfirmDialog
                open={showRemoveDialog}
                onOpenChange={setShowRemoveDialog}
                title="Remove Member"
                description={`Are you sure you want to remove ${member.fullName} (${member.email}) from this organization?`}
                confirmLabel="Remove Member"
                isLoading={removeMutation.isPending}
                onConfirm={handleRemove}
            />
        </>
    )
}
