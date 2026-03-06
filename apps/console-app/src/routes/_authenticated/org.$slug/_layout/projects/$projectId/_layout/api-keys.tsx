import { createFileRoute } from '@tanstack/react-router'
import { ProjectApiKeysPage } from '@/features/projects/pages/ProjectApiKeysPage'

export const Route = createFileRoute(
    '/_authenticated/org/$slug/_layout/projects/$projectId/_layout/api-keys',
)({
    component: () => {
        const { projectId } = Route.useParams()
        return <ProjectApiKeysPage projectId={projectId} />
    },
})
