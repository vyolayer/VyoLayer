import { createFileRoute } from '@tanstack/react-router'
import { ProjectSettingsPage } from '@/features/projects/pages/ProjectSettingsPage'

export const Route = createFileRoute(
    '/_authenticated/org/$slug/_layout/projects/$projectId/_layout/settings',
)({
    component: () => {
        const { projectId } = Route.useParams()
        return <ProjectSettingsPage projectId={projectId} />
    },
})
