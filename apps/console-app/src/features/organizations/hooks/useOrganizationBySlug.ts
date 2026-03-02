import { useQuery } from '@tanstack/react-query'
import { getOrganizationBySlug } from '../services/organization.service'

export const organizationQueries = {
    bySlug: (slug: string) => ({
        queryKey: ['organizations', 'slug', slug],
        queryFn: () => getOrganizationBySlug(slug),
        retry: 1,
        staleTime: 1000 * 60 * 5,
    }),
}

export function useOrganizationBySlug(slug: string) {
    return useQuery(organizationQueries.bySlug(slug))
}
