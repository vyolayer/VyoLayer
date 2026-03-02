import { useQuery } from '@tanstack/react-query'
import { getRoles } from '../services/organization.service'

export function useOrganizationRoles(orgId?: string) {
    return useQuery({
        queryKey: ['org-roles', orgId],
        queryFn: () => getRoles(orgId!),
        enabled: !!orgId,
        staleTime: 1000 * 60 * 10,
    })
}
