import { createFileRoute } from '@tanstack/react-router'
import { ProjectMembersPage } from '@/features/projects/pages/ProjectMembersPage'

export const Route = createFileRoute(
    '/_authenticated/org/$slug/_layout/projects/$projectId/_layout/members',
)({
    component: () => {
        const { projectId } = Route.useParams()
        return <ProjectMembersPage projectId={projectId} />
    },
})
