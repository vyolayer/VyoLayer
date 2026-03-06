import { useState } from 'react'
import { useOrganization } from '@/contexts/OrganizationContext'
import { useApiKeys, useGenerateApiKey, useRevokeApiKey } from '../hooks/useApiKeys'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Key, Plus, Copy, Loader2, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import type { ApiKeyCreated, ApiKeyMode } from '../types'

interface ProjectApiKeysPageProps {
    projectId: string
}

export function ProjectApiKeysPage({ projectId }: ProjectApiKeysPageProps) {
    const { org } = useOrganization()
    const { data: apiKeys, isLoading } = useApiKeys(org.id, projectId)
    const generateKey = useGenerateApiKey(org.id, projectId)
    const revokeKey = useRevokeApiKey(org.id, projectId)

    const [newKey, setNewKey] = useState<ApiKeyCreated | null>(null)
    const [createOpen, setCreateOpen] = useState(false)
    const [keyName, setKeyName] = useState('')
    const [keyMode, setKeyMode] = useState<ApiKeyMode>('dev')

    const handleGenerate = () => {
        if (!keyName.trim()) return
        generateKey.mutate(
            { name: keyName, mode: keyMode },
            {
                onSuccess: (data) => {
                    setNewKey(data)
                    setCreateOpen(false)
                    setKeyName('')
                    toast.success('API key generated')
                },
                onError: (err) => toast.error(err instanceof Error ? err.message : 'Failed to generate key'),
            },
        )
    }

    const handleRevoke = (apiKeyId: string) => {
        revokeKey.mutate(apiKeyId, {
            onSuccess: () => toast.success('API key revoked'),
            onError: (err) => toast.error(err instanceof Error ? err.message : 'Failed to revoke key'),
        })
    }

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text)
        toast.success('Copied to clipboard')
    }

    if (isLoading) {
        return (
            <div className="space-y-3">
                {[1, 2].map((i) => (
                    <div key={i} className="h-20 animate-pulse rounded-lg bg-muted" />
                ))}
            </div>
        )
    }

    return (
        <div className="space-y-6">
            {/* Newly generated key banner */}
            {newKey && (
                <Card className="border-green-500/50 bg-green-500/5">
                    <CardContent className="pt-6 space-y-3">
                        <div className="flex items-center gap-2">
                            <Key className="h-5 w-5 text-green-600" />
                            <p className="font-semibold text-green-700 dark:text-green-400">
                                New key generated — copy it now!
                            </p>
                        </div>
                        <p className="text-xs text-muted-foreground">
                            This key will not be shown again. Store it securely.
                        </p>
                        <div className="flex items-center gap-2">
                            <code className="flex-1 rounded bg-muted px-3 py-2 text-sm font-mono break-all">
                                {newKey.rawKey}
                            </code>
                            <Button variant="outline" size="icon" onClick={() => copyToClipboard(newKey.rawKey)}>
                                <Copy className="h-4 w-4" />
                            </Button>
                        </div>
                        <Button variant="ghost" size="sm" onClick={() => setNewKey(null)}>
                            Dismiss
                        </Button>
                    </CardContent>
                </Card>
            )}

            <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold">API Keys</h3>

                <Dialog open={createOpen} onOpenChange={setCreateOpen}>
                    <DialogTrigger asChild>
                        <Button size="sm">
                            <Plus className="h-4 w-4 mr-1" />
                            Generate Key
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Generate API Key</DialogTitle>
                            <DialogDescription>Create a new API key for this project.</DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4">
                            <div className="space-y-2">
                                <Label>Name</Label>
                                <Input
                                    placeholder="e.g. Production Key"
                                    value={keyName}
                                    onChange={(e) => setKeyName(e.target.value)}
                                />
                            </div>
                            <div className="space-y-2">
                                <Label>Mode</Label>
                                <Select value={keyMode} onValueChange={(v) => setKeyMode(v as ApiKeyMode)}>
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="dev">Development</SelectItem>
                                        <SelectItem value="live">Live</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="flex justify-end gap-2">
                                <Button variant="outline" onClick={() => setCreateOpen(false)}>
                                    Cancel
                                </Button>
                                <Button onClick={handleGenerate} disabled={generateKey.isPending || !keyName.trim()}>
                                    {generateKey.isPending && <Loader2 className="h-4 w-4 mr-1 animate-spin" />}
                                    Generate
                                </Button>
                            </div>
                        </div>
                    </DialogContent>
                </Dialog>
            </div>

            {!apiKeys?.length ? (
                <div className="flex flex-col items-center justify-center py-16 text-center">
                    <Key className="h-12 w-12 text-muted-foreground/50 mb-4" />
                    <h3 className="text-lg font-semibold">No API keys</h3>
                    <p className="text-sm text-muted-foreground max-w-sm">
                        Generate an API key to start making authenticated requests.
                    </p>
                </div>
            ) : (
                <div className="space-y-3">
                    {apiKeys.map((key) => (
                        <Card key={key.id}>
                            <CardContent className="flex items-center justify-between py-4">
                                <div className="space-y-1">
                                    <div className="flex items-center gap-2">
                                        <p className="font-medium text-sm">{key.name}</p>
                                        <Badge variant={key.mode === 'live' ? 'default' : 'secondary'}>
                                            {key.mode}
                                        </Badge>
                                        {key.isRevoked && <Badge variant="destructive">Revoked</Badge>}
                                    </div>
                                    <p className="text-xs text-muted-foreground font-mono">{key.keyPrefix}•••</p>
                                    <p className="text-xs text-muted-foreground">
                                        Created {new Date(key.createdAt).toLocaleDateString()}
                                        {key.lastUsedAt &&
                                            ` · Last used ${new Date(key.lastUsedAt).toLocaleDateString()}`}
                                    </p>
                                </div>

                                {!key.isRevoked && (
                                    <AlertDialog>
                                        <AlertDialogTrigger asChild>
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                className="text-destructive hover:bg-destructive/10"
                                            >
                                                <Trash2 className="h-4 w-4" />
                                            </Button>
                                        </AlertDialogTrigger>
                                        <AlertDialogContent>
                                            <AlertDialogHeader>
                                                <AlertDialogTitle>Revoke API Key</AlertDialogTitle>
                                                <AlertDialogDescription>
                                                    This action cannot be undone. All requests using this key will be
                                                    rejected.
                                                </AlertDialogDescription>
                                            </AlertDialogHeader>
                                            <AlertDialogFooter>
                                                <AlertDialogCancel>Cancel</AlertDialogCancel>
                                                <AlertDialogAction
                                                    onClick={() => handleRevoke(key.id)}
                                                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                                                >
                                                    Revoke
                                                </AlertDialogAction>
                                            </AlertDialogFooter>
                                        </AlertDialogContent>
                                    </AlertDialog>
                                )}
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    )
}
