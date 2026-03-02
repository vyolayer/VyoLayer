import { useOrganization } from '@/contexts/OrganizationContext'
import { useArchiveOrganization, useRestoreOrganization } from '@/features/organizations/hooks/useOrganizationMutations'
import { ConfirmDialog } from '@/components/organization/ConfirmDialog'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Archive, RotateCcw, AlertCircle, CheckCircle2, XCircle } from 'lucide-react'
import { useState } from 'react'
import { formatDate } from '@/lib/utils'

export function OrganizationArchivePage() {
    const { org, isAdmin } = useOrganization()
    const [showArchiveDialog, setShowArchiveDialog] = useState(false)
    const [showRestoreDialog, setShowRestoreDialog] = useState(false)

    const archiveMutation = useArchiveOrganization(org.id)
    const restoreMutation = useRestoreOrganization(org.id)

    if (!isAdmin) {
        return (
            <div className="py-20 text-center space-y-4">
                <AlertCircle className="mx-auto h-10 w-10 text-muted-foreground" />
                <p className="text-muted-foreground">You don't have permission to manage archive status.</p>
            </div>
        )
    }

    return (
        <div className="space-y-6">
            <div>
                <h2 className="text-lg font-semibold">Archive Management</h2>
                <p className="text-sm text-muted-foreground">Control the active status of your organization.</p>
            </div>

            <Card>
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <CardTitle className="text-base">Organization Status</CardTitle>
                        <Badge
                            variant="outline"
                            className={
                                org.isActive
                                    ? 'text-emerald-600 border-emerald-500/30 bg-emerald-500/10'
                                    : 'text-destructive border-destructive/30 bg-destructive/10'
                            }
                        >
                            {org.isActive ? (
                                <>
                                    <CheckCircle2 className="h-3 w-3 mr-1" />
                                    Active
                                </>
                            ) : (
                                <>
                                    <XCircle className="h-3 w-3 mr-1" />
                                    Archived
                                </>
                            )}
                        </Badge>
                    </div>
                    <CardDescription>
                        {org.isActive
                            ? 'This organization is currently active. All members can access it.'
                            : 'This organization is archived. Members cannot access it until restored.'}
                    </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    {org.isActive ? (
                        <>
                            <div className="rounded-lg border border-amber-500/30 bg-amber-500/5 p-4 space-y-2">
                                <p className="text-sm font-medium text-amber-700 dark:text-amber-400">
                                    What happens when you archive?
                                </p>
                                <ul className="text-sm text-muted-foreground space-y-1 list-disc pl-5">
                                    <li>All members will lose access</li>
                                    <li>No new actions can be performed</li>
                                    <li>Data is preserved and can be restored</li>
                                </ul>
                            </div>
                            <Button variant="outline" onClick={() => setShowArchiveDialog(true)}>
                                <Archive className="w-4 h-4 mr-2" />
                                Archive Organization
                            </Button>
                        </>
                    ) : (
                        <>
                            <div className="rounded-lg border border-emerald-500/30 bg-emerald-500/5 p-4 space-y-2">
                                <p className="text-sm font-medium text-emerald-700 dark:text-emerald-400">
                                    Restoring will:
                                </p>
                                <ul className="text-sm text-muted-foreground space-y-1 list-disc pl-5">
                                    <li>Re-enable access for all members</li>
                                    <li>Restore all previous settings</li>
                                    <li>Resume normal operations</li>
                                </ul>
                            </div>
                            <Button onClick={() => setShowRestoreDialog(true)}>
                                <RotateCcw className="w-4 h-4 mr-2" />
                                Restore Organization
                            </Button>
                        </>
                    )}

                    {org.deactivatedAt && (
                        <p className="text-xs text-muted-foreground">Archived on {formatDate(org.deactivatedAt)}</p>
                    )}
                </CardContent>
            </Card>

            <ConfirmDialog
                open={showArchiveDialog}
                onOpenChange={setShowArchiveDialog}
                title="Archive Organization"
                description={`Are you sure you want to archive "${org.name}"? All members will lose access until it is restored.`}
                confirmLabel="Archive"
                confirmVariant="destructive"
                isLoading={archiveMutation.isPending}
                onConfirm={() =>
                    archiveMutation.mutate(undefined, {
                        onSuccess: () => setShowArchiveDialog(false),
                    })
                }
            />

            <ConfirmDialog
                open={showRestoreDialog}
                onOpenChange={setShowRestoreDialog}
                title="Restore Organization"
                description={`Restore "${org.name}" to active status? All members will regain access.`}
                confirmLabel="Restore"
                confirmVariant="default"
                isLoading={restoreMutation.isPending}
                onConfirm={() =>
                    restoreMutation.mutate(undefined, {
                        onSuccess: () => setShowRestoreDialog(false),
                    })
                }
            />
        </div>
    )
}
