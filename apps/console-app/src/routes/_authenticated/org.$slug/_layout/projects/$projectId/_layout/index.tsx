import { createFileRoute } from '@tanstack/react-router'
import { ProjectDashboardPage } from '@/features/projects/pages/ProjectDashboardPage'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/projects/$projectId/_layout/')({
    component: () => {
        const { projectId } = Route.useParams()
        return <ProjectDashboardPage projectId={projectId} />
    },
})
