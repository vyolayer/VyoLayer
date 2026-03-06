import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useOrganization } from '@/contexts/OrganizationContext'
import {
    useProject,
    useUpdateProject,
    useArchiveProject,
    useRestoreProject,
    useDeleteProject,
} from '../hooks/useProjects'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
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
import { Loader2, Archive, RotateCcw, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import { useNavigate } from '@tanstack/react-router'

const settingsSchema = z.object({
    name: z.string().min(3).max(100),
    description: z.string().max(500).optional(),
})

type FormValues = z.infer<typeof settingsSchema>

interface ProjectSettingsPageProps {
    projectId: string
}

export function ProjectSettingsPage({ projectId }: ProjectSettingsPageProps) {
    const { org } = useOrganization()
    const { data: project, isLoading } = useProject(org.id, projectId)
    const updateProject = useUpdateProject(org.id, projectId)
    const archiveProject = useArchiveProject(org.id, projectId)
    const restoreProject = useRestoreProject(org.id, projectId)
    const deleteProject = useDeleteProject(org.id, projectId)
    const navigate = useNavigate()
    const [confirmName, setConfirmName] = useState('')

    const {
        register,
        handleSubmit,
        formState: { errors, isDirty },
    } = useForm<FormValues>({
        resolver: zodResolver(settingsSchema),
        values: project ? { name: project.name, description: project.description || '' } : undefined,
    })

    const onSubmit = (data: FormValues) => {
        updateProject.mutate(data, {
            onSuccess: () => toast.success('Project updated'),
            onError: (err) => toast.error(err instanceof Error ? err.message : 'Failed to update'),
        })
    }

    const handleArchive = () => {
        archiveProject.mutate(undefined, {
            onSuccess: () => toast.success('Project archived'),
            onError: (err) => toast.error(err instanceof Error ? err.message : 'Failed to archive'),
        })
    }

    const handleRestore = () => {
        restoreProject.mutate(undefined, {
            onSuccess: () => toast.success('Project restored'),
            onError: (err) => toast.error(err instanceof Error ? err.message : 'Failed to restore'),
        })
    }

    const handleDelete = () => {
        deleteProject.mutate(undefined, {
            onSuccess: () => {
                toast.success('Project deleted')
                navigate({ to: `/org/${org.slug}/projects` })
            },
            onError: (err) => toast.error(err instanceof Error ? err.message : 'Failed to delete'),
        })
    }

    if (isLoading || !project) {
        return <div className="h-48 animate-pulse rounded-xl bg-muted" />
    }

    return (
        <div className="space-y-6">
            {/* General settings */}
            <Card>
                <CardHeader>
                    <CardTitle>General</CardTitle>
                    <CardDescription>Update your project details.</CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="settings-name">Name</Label>
                            <Input id="settings-name" {...register('name')} />
                            {errors.name && <p className="text-xs text-destructive">{errors.name.message}</p>}
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="settings-desc">Description</Label>
                            <Textarea id="settings-desc" {...register('description')} />
                            {errors.description && (
                                <p className="text-xs text-destructive">{errors.description.message}</p>
                            )}
                        </div>
                        <Button type="submit" disabled={!isDirty || updateProject.isPending}>
                            {updateProject.isPending && <Loader2 className="h-4 w-4 mr-1 animate-spin" />}
                            Save Changes
                        </Button>
                    </form>
                </CardContent>
            </Card>

            <Separator />

            {/* Archive / Restore */}
            <Card>
                <CardHeader>
                    <CardTitle>{project.isActive ? 'Archive Project' : 'Restore Project'}</CardTitle>
                    <CardDescription>
                        {project.isActive
                            ? 'Archiving a project hides it and disables API keys.'
                            : 'Restore the project to make it active again.'}
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    {project.isActive ? (
                        <Button variant="outline" onClick={handleArchive} disabled={archiveProject.isPending}>
                            {archiveProject.isPending ? (
                                <Loader2 className="h-4 w-4 mr-1 animate-spin" />
                            ) : (
                                <Archive className="h-4 w-4 mr-1" />
                            )}
                            Archive
                        </Button>
                    ) : (
                        <Button variant="outline" onClick={handleRestore} disabled={restoreProject.isPending}>
                            {restoreProject.isPending ? (
                                <Loader2 className="h-4 w-4 mr-1 animate-spin" />
                            ) : (
                                <RotateCcw className="h-4 w-4 mr-1" />
                            )}
                            Restore
                        </Button>
                    )}
                </CardContent>
            </Card>

            <Separator />

            {/* Danger zone — Delete */}
            <Card className="border-destructive/50">
                <CardHeader>
                    <CardTitle className="text-destructive">Delete Project</CardTitle>
                    <CardDescription>
                        Permanently delete this project and all associated data. This action cannot be undone.
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <AlertDialog>
                        <AlertDialogTrigger asChild>
                            <Button variant="destructive">
                                <Trash2 className="h-4 w-4 mr-1" />
                                Delete Project
                            </Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                            <AlertDialogHeader>
                                <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                                <AlertDialogDescription>
                                    Type <strong>{project.name}</strong> to confirm deletion.
                                </AlertDialogDescription>
                            </AlertDialogHeader>
                            <Input
                                placeholder={project.name}
                                value={confirmName}
                                onChange={(e) => setConfirmName(e.target.value)}
                            />
                            <AlertDialogFooter>
                                <AlertDialogCancel>Cancel</AlertDialogCancel>
                                <AlertDialogAction
                                    disabled={confirmName !== project.name || deleteProject.isPending}
                                    onClick={handleDelete}
                                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                                >
                                    {deleteProject.isPending && <Loader2 className="h-4 w-4 mr-1 animate-spin" />}
                                    Delete
                                </AlertDialogAction>
                            </AlertDialogFooter>
                        </AlertDialogContent>
                    </AlertDialog>
                </CardContent>
            </Card>
        </div>
    )
}
