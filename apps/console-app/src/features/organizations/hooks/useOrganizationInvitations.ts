import { useQuery } from '@tanstack/react-query'
import { getInvitations } from '../services/organization.service'

export function useOrganizationInvitations(orgId?: string) {
    return useQuery({
        queryKey: ['org-invitations', orgId],
        queryFn: () => getInvitations(orgId!),
        enabled: !!orgId,
        staleTime: 1000 * 60 * 2,
    })
}
