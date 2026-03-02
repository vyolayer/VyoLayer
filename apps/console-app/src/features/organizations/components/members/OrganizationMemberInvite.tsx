import { useState } from 'react'

import { Loader2, Send } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

import { useInviteMember } from '@/org/hooks/useOrganizationMutations'
import { useOrganizationRoles } from '@/org/hooks/useOrganizationRoles'

export const OrganizationMembersInvite: React.FC<{
    orgId: string
    isAdmin: boolean
}> = ({ orgId, isAdmin }) => {
    const inviteMutation = useInviteMember(orgId)
    const { data: roles = [] } = useOrganizationRoles(isAdmin ? orgId : undefined)

    const [inviteEmail, setInviteEmail] = useState('')
    const [inviteRoleId, setInviteRoleId] = useState<string>('')

    const handleInvite = () => {
        const payload: { email: string; roleIds?: string[] } = { email: inviteEmail }
        if (inviteRoleId) payload.roleIds = [inviteRoleId]
        inviteMutation.mutate(payload, {
            onSuccess: () => {
                setInviteEmail('')
                setInviteRoleId('')
            },
        })
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-base">Invite Member</CardTitle>
                <CardDescription>Send an email invitation to add someone to this organization.</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="flex flex-col sm:flex-row gap-3">
                    <Input
                        placeholder="colleague@example.com"
                        type="email"
                        value={inviteEmail}
                        onChange={(e) => setInviteEmail(e.target.value)}
                        className="flex-1"
                    />
                    {roles.length > 0 && (
                        <Select value={inviteRoleId} onValueChange={setInviteRoleId}>
                            <SelectTrigger className="w-full sm:w-40">
                                <SelectValue placeholder="Role (optional)" />
                            </SelectTrigger>
                            <SelectContent>
                                {roles
                                    .filter((r) => r.name !== 'Owner')
                                    .map((r) => (
                                        <SelectItem key={r.id} value={r.id}>
                                            {r.name}
                                        </SelectItem>
                                    ))}
                            </SelectContent>
                        </Select>
                    )}
                    <Button onClick={handleInvite} disabled={!inviteEmail || inviteMutation.isPending}>
                        {inviteMutation.isPending ? (
                            <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                        ) : (
                            <Send className="w-4 h-4 mr-2" />
                        )}
                        Invite
                    </Button>
                </div>
                {inviteMutation.isError && (
                    <p className="text-sm text-destructive mt-2">
                        {inviteMutation.error instanceof Error ? inviteMutation.error.message : 'Failed to invite'}
                    </p>
                )}
                {inviteMutation.isSuccess && <p className="text-sm text-emerald-600 mt-2">Invitation sent!</p>}
            </CardContent>
        </Card>
    )
}
