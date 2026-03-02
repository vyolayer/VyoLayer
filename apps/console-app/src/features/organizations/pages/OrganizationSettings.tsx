import { useState } from 'react'
import { Save, Loader2, AlertCircle } from 'lucide-react'

import { useOrganization } from '@/contexts/OrganizationContext'
import { useUpdateOrganization } from '@/org/hooks/useOrganizationMutations'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export function OrganizationSettingsPage() {
    const { org, isAdmin } = useOrganization()

    const [name, setName] = useState(org.name)
    const [description, setDescription] = useState(org.description || '')

    const updateMutation = useUpdateOrganization(org.id)

    // If the user is not an admin, they don't have permission to edit settings
    if (!isAdmin) {
        return (
            <div className="py-20 text-center space-y-4">
                <AlertCircle className="mx-auto h-10 w-10 text-muted-foreground" />
                <p className="text-muted-foreground">You don't have permission to edit settings.</p>
            </div>
        )
    }

    const hasChanges = name !== org.name || description !== (org.description || '')

    const handleSave = () => {
        updateMutation.mutate({
            name: name !== org.name ? name : undefined,
            description: description !== (org.description || '') ? description : undefined,
        })
    }

    return (
        <div className="space-y-6">
            <div>
                <h2 className="text-lg font-semibold">Organization Settings</h2>
                <p className="text-sm text-muted-foreground">Manage your organization's profile information.</p>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">General Information</CardTitle>
                    <CardDescription>Update the organization's name, slug, and description.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-5">
                    <div className="space-y-2">
                        <Label htmlFor="org-name">Organization Name</Label>
                        <Input
                            id="org-name"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="My Organization"
                            maxLength={100}
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="org-slug">Slug</Label>
                        <div className="flex items-center gap-2">
                            <span className="text-sm text-muted-foreground whitespace-nowrap">/org/</span>
                            <Input id="org-slug" value={org.slug} disabled />
                        </div>
                        <p className="text-xs text-muted-foreground">
                            URL-safe identifier. Only lowercase letters, numbers, and hyphens.
                        </p>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="org-description">Description</Label>
                        <Textarea
                            id="org-description"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            placeholder="Tell people what this organization does..."
                            rows={3}
                            maxLength={500}
                        />
                        <p className="text-xs text-muted-foreground text-right">{description.length}/500</p>
                    </div>

                    <div className="flex items-center gap-3 justify-end">
                        {updateMutation.isSuccess && <p className="text-sm text-emerald-600">Saved successfully.</p>}
                        {updateMutation.isError && (
                            <p className="text-sm text-destructive">
                                {updateMutation.error instanceof Error ? updateMutation.error.message : 'Save failed'}
                            </p>
                        )}
                        <Button onClick={handleSave} disabled={!hasChanges || updateMutation.isPending}>
                            {updateMutation.isPending ? (
                                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                            ) : (
                                <Save className="w-4 h-4 mr-2" />
                            )}
                            Save Changes
                        </Button>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
