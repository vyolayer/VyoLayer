import { LandingPage } from '@/features/landing/pages/Landing'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_landing/')({
    component: LandingPage,
})
