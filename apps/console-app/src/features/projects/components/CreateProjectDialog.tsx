import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Plus, Loader2 } from 'lucide-react'
import { useCreateProject } from '../hooks/useProjects'
import { toast } from 'sonner'

const createProjectSchema = z.object({
    name: z.string().min(3, 'Name must be at least 3 characters').max(100),
    description: z.string().max(500).optional(),
})

type FormValues = z.infer<typeof createProjectSchema>

interface CreateProjectDialogProps {
    orgId: string
}

export function CreateProjectDialog({ orgId }: CreateProjectDialogProps) {
    const [open, setOpen] = useState(false)
    const createProject = useCreateProject(orgId)

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<FormValues>({
        resolver: zodResolver(createProjectSchema),
        defaultValues: { name: '', description: '' },
    })

    const onSubmit = (data: FormValues) => {
        createProject.mutate(data, {
            onSuccess: () => {
                toast.success('Project created successfully')
                reset()
                setOpen(false)
            },
            onError: (err) => {
                toast.error(err instanceof Error ? err.message : 'Failed to create project')
            },
        })
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button size="sm">
                    <Plus className="h-4 w-4 mr-1" />
                    New Project
                </Button>
            </DialogTrigger>

            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Create Project</DialogTitle>
                    <DialogDescription>Add a new project to this organization.</DialogDescription>
                </DialogHeader>

                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="project-name">Name</Label>
                        <Input id="project-name" placeholder="My Project" {...register('name')} />
                        {errors.name && <p className="text-xs text-destructive">{errors.name.message}</p>}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="project-desc">Description (optional)</Label>
                        <Textarea id="project-desc" placeholder="A brief description…" {...register('description')} />
                        {errors.description && <p className="text-xs text-destructive">{errors.description.message}</p>}
                    </div>

                    <div className="flex justify-end gap-2">
                        <Button type="button" variant="outline" onClick={() => setOpen(false)}>
                            Cancel
                        </Button>
                        <Button type="submit" disabled={createProject.isPending}>
                            {createProject.isPending && <Loader2 className="h-4 w-4 mr-1 animate-spin" />}
                            Create
                        </Button>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    )
}
