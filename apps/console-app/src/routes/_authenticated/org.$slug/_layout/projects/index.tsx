import { createFileRoute } from '@tanstack/react-router'
import { ProjectsOverviewPage } from '@/features/projects/pages/ProjectsOverviewPage'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/projects/')({
    component: ProjectsOverviewPage,
})
