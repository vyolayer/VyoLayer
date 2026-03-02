import { useQuery } from '@tanstack/react-query'
import { getMyMembershipInOrg } from '../services/organization.service'

function useMyMembershipInOrgQuery(orgId?: string) {
    return useQuery({
        queryKey: ['organizations', orgId, 'members', 'me'],
        queryFn: () => getMyMembershipInOrg(orgId || ''),
        enabled: !!orgId,
        retry: 1,
        staleTime: 1000 * 60 * 5,
    })
}

export const useMyMembershipInOrg = (orgId?: string) => {
    const { data, isLoading, error } = useMyMembershipInOrgQuery(orgId)
    return { data, isLoading, error }
}
