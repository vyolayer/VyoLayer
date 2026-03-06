import { createFileRoute } from '@tanstack/react-router'
import { ProjectLayout } from '@/features/projects/components/ProjectLayout'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/projects/$projectId/_layout')({
    component: () => {
        const { projectId } = Route.useParams()
        return <ProjectLayout projectId={projectId} />
    },
})
