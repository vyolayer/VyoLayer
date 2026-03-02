import { useQuery } from '@tanstack/react-query'
import { getMembers } from '../services/organization.service'

export function useOrganizationMembers(orgId?: string) {
    return useQuery({
        queryKey: ['org-members', orgId],
        queryFn: () => getMembers(orgId!),
        enabled: !!orgId,
        staleTime: 1000 * 60 * 2,
    })
}
