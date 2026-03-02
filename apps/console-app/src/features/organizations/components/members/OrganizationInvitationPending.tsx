import { Button } from '@/components/ui/button'
import { Mail, RotateCcw, XCircle } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { formatDate } from '@/lib/utils'

import { useOrganizationInvitations } from '@/org/hooks/useOrganizationInvitations'
import { useCancelInvitation, useResendInvitation } from '@/org/hooks/useOrganizationMutations'

export const OrganizationInvitationPending: React.FC<{
    orgId: string
    isAdmin: boolean
}> = ({ orgId, isAdmin }) => {
    const { data: invitations = [] } = useOrganizationInvitations(isAdmin ? orgId : undefined)

    const resendMutation = useResendInvitation(orgId)
    const cancelMutation = useCancelInvitation(orgId)

    const pendingInvitations = invitations.filter((i) => i.isPending)

    if (pendingInvitations.length === 0) return null

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-base">Pending Invitations</CardTitle>
                <CardDescription>{pendingInvitations.length} outstanding</CardDescription>
            </CardHeader>
            <CardContent className="p-0">
                <ul className="divide-y divide-border">
                    {pendingInvitations.map((inv) => (
                        <li key={inv.id} className="flex items-center justify-between px-4 py-3">
                            <div className="flex items-center gap-3">
                                <Mail className="w-4 h-4 text-muted-foreground" />
                                <div>
                                    <p className="text-sm font-medium">{inv.email}</p>
                                    <p className="text-xs text-muted-foreground">Expires {formatDate(inv.expiredAt)}</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-1">
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => resendMutation.mutate(inv.id)}
                                    disabled={resendMutation.isPending}
                                >
                                    <RotateCcw className="w-3.5 h-3.5 mr-1" />
                                    Resend
                                </Button>
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => cancelMutation.mutate(inv.id)}
                                    disabled={cancelMutation.isPending}
                                    className="text-destructive hover:text-destructive"
                                >
                                    <XCircle className="w-3.5 h-3.5 mr-1" />
                                    Cancel
                                </Button>
                            </div>
                        </li>
                    ))}
                </ul>
            </CardContent>
        </Card>
    )
}
