import { useQuery } from '@tanstack/react-query'
import { getOrganizations } from '../services/organization.service'

function useGetOrganizations() {
    return useQuery({
        queryKey: ['organizations'],
        queryFn: getOrganizations,
        retry: 1,
        staleTime: 1000 * 60 * 5,
    })
}

export const useOrganizations = () => {
    const { data, isLoading, error } = useGetOrganizations()
    return { data, isLoading, error }
}
