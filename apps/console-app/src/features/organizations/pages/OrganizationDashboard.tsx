import { Container } from '@/components/Container'
import {
    OrganizationDashboardError,
    OrganizationDashboardLoading,
    EmptyOrganization,
    OrganizationDashboard,
} from '@/org/components/dashboard'
import { useOrganizations } from '@/org/hooks/useOrganizations'

export function OrganizationDashboardPage() {
    const { data: organizations, isLoading, error } = useOrganizations()

    if (isLoading) return <OrganizationDashboardLoading />

    if (error) return <OrganizationDashboardError error={error} />

    // ── First-time user (no orgs) ─────────────────────────────────
    if (!organizations || organizations.length === 0) return <EmptyOrganization />

    // ── Organizations list ────────────────────────────────────────
    return (
        <Container>
            <OrganizationDashboard organizations={organizations} />
        </Container>
    )
}
