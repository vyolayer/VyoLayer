import { useState } from 'react'
import { useOrganization } from '@/contexts/OrganizationContext'
import { useOrganizationMembers } from '@/features/organizations/hooks/useOrganizationMembers'
import { useDeleteOrganization, useTransferOwnership } from '@/features/organizations/hooks/useOrganizationMutations'
import { ConfirmDialog } from '@/components/organization/ConfirmDialog'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { AlertCircle, Trash2, ArrowRightLeft } from 'lucide-react'
import { useNavigate } from '@tanstack/react-router'
import type { OrganizationMember } from '../types'

export function OrganizationDangerZone() {
    const { org, me, isOwner } = useOrganization()
    const navigate = useNavigate()
    const { data: members = org.members } = useOrganizationMembers(org.id)

    const [transferMemberId, setTransferMemberId] = useState('')
    const [showTransferDialog, setShowTransferDialog] = useState(false)
    const [showDeleteDialog, setShowDeleteDialog] = useState(false)

    const deleteMutation = useDeleteOrganization(org.id)
    const transferMutation = useTransferOwnership(org.id)

    if (!isOwner) {
        return (
            <div className="py-20 text-center space-y-4">
                <AlertCircle className="mx-auto h-10 w-10 text-muted-foreground" />
                <h2 className="text-lg font-semibold">Access Denied</h2>
                <p className="text-sm text-muted-foreground">Only the organization owner can access the danger zone.</p>
            </div>
        )
    }

    const otherMembers = members.filter((m) => m.userId !== me?.userId)
    const selectedMember = otherMembers.find((m) => m.id === transferMemberId)

    return (
        <div className="space-y-8">
            <div>
                <h2 className="text-lg font-semibold text-destructive">Danger Zone</h2>
                <p className="text-sm text-muted-foreground">
                    Irreversible actions. Please be certain before proceeding.
                </p>
            </div>

            {/* ── Transfer Ownership ── */}
            <OrganizationTransferOwnership
                orgId={org.id}
                members={otherMembers}
                transferMemberId={transferMemberId}
                setTransferMemberId={setTransferMemberId}
                setShowTransferDialog={setShowTransferDialog}
            />

            <Separator />

            {/* ── Delete Organization ── */}
            <Card className="border-destructive/40">
                <CardHeader>
                    <CardTitle className="flex items-center gap-2 text-base text-destructive">
                        <Trash2 className="w-4 h-4" />
                        Delete Organization
                    </CardTitle>
                    <CardDescription>
                        This will permanently delete the organization, all members, and all associated data. This action
                        is irreversible.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <Button variant="destructive" onClick={() => setShowDeleteDialog(true)}>
                        <Trash2 className="w-4 h-4 mr-2" />
                        Delete this organization
                    </Button>
                </CardContent>
            </Card>

            {/* Transfer Dialog */}
            <ConfirmDialog
                open={showTransferDialog}
                onOpenChange={setShowTransferDialog}
                title="Transfer Ownership"
                description={`Are you sure you want to transfer ownership to ${selectedMember?.fullName ?? ''}? You will be downgraded to a Member role.`}
                confirmLabel="Transfer Ownership"
                confirmVariant="destructive"
                isLoading={transferMutation.isPending}
                onConfirm={() =>
                    transferMutation.mutate(
                        { newOwnerMemberId: transferMemberId },
                        { onSuccess: () => setShowTransferDialog(false) },
                    )
                }
            />

            {/* Delete Dialog */}
            <ConfirmDialog
                open={showDeleteDialog}
                onOpenChange={setShowDeleteDialog}
                title="Delete Organization"
                description="This will permanently delete the organization and all its data. This cannot be undone."
                confirmLabel="Delete Organization"
                confirmVariant="destructive"
                confirmText={org.name}
                confirmPlaceholder="Type the organization name"
                isLoading={deleteMutation.isPending}
                onConfirm={() =>
                    deleteMutation.mutate(
                        { confirmName: org.name },
                        {
                            onSuccess: () => {
                                setShowDeleteDialog(false)
                                navigate({ to: '/' })
                            },
                        },
                    )
                }
            />
        </div>
    )
}

const OrganizationTransferOwnership: React.FC<{
    orgId: string
    members: OrganizationMember[]
    transferMemberId: string
    setTransferMemberId: (id: string) => void
    setShowTransferDialog: (show: boolean) => void
}> = ({ members, transferMemberId, setTransferMemberId, setShowTransferDialog, orgId }) => {
    const transferMutation = useTransferOwnership(orgId)

    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2 text-base">
                    <ArrowRightLeft className="w-4 h-4" />
                    Transfer Ownership
                </CardTitle>
                <CardDescription>
                    Transfer ownership to another member. You will be downgraded to a regular Member.
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                {members.length === 0 ? (
                    <p className="text-sm text-muted-foreground">
                        No other members to transfer ownership to. Invite members first.
                    </p>
                ) : (
                    <>
                        <div className="space-y-2">
                            <Label>Select new owner</Label>
                            <Select value={transferMemberId} onValueChange={setTransferMemberId}>
                                <SelectTrigger>
                                    <SelectValue placeholder="Choose a member..." />
                                </SelectTrigger>
                                <SelectContent>
                                    {members.map((m) => (
                                        <SelectItem key={m.id} value={m.id}>
                                            {m.fullName} ({m.email})
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                        <Button
                            variant="outline"
                            onClick={() => setShowTransferDialog(true)}
                            disabled={!transferMemberId}
                        >
                            Transfer Ownership
                        </Button>
                    </>
                )}

                {transferMutation.isError && (
                    <p className="text-sm text-destructive">
                        {transferMutation.error instanceof Error ? transferMutation.error.message : 'Transfer failed'}
                    </p>
                )}
                {transferMutation.isSuccess && (
                    <p className="text-sm text-emerald-600">Ownership transferred successfully.</p>
                )}
            </CardContent>
        </Card>
    )
}
