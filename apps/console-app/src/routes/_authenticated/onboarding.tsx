import { createFileRoute } from '@tanstack/react-router'
import { OrganizationOnboardingPage } from '@/features/organizations/pages'

export const Route = createFileRoute('/_authenticated/onboarding')({
    component: OrganizationOnboardingPage,
})
