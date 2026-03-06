import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
    listApiKeys,
    getApiKeyById,
    generateApiKey,
    revokeApiKey,
    type CreateApiKeyPayload,
} from '../services/api-key.service'

// ─── Query Keys ──────────────────────────────────────────────────────────────

export const apiKeyKeys = {
    all: (orgId: string, projectId: string) => ['api-keys', orgId, projectId] as const,
    detail: (orgId: string, projectId: string, keyId: string) => ['api-keys', orgId, projectId, keyId] as const,
}

// ─── Queries ─────────────────────────────────────────────────────────────────

export function useApiKeys(orgId: string, projectId: string) {
    return useQuery({
        queryKey: apiKeyKeys.all(orgId, projectId),
        queryFn: () => listApiKeys(orgId, projectId),
        enabled: !!orgId && !!projectId,
    })
}

export function useApiKey(orgId: string, projectId: string, apiKeyId: string) {
    return useQuery({
        queryKey: apiKeyKeys.detail(orgId, projectId, apiKeyId),
        queryFn: () => getApiKeyById(orgId, projectId, apiKeyId),
        enabled: !!orgId && !!projectId && !!apiKeyId,
    })
}

// ─── Mutations ───────────────────────────────────────────────────────────────

export function useGenerateApiKey(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: CreateApiKeyPayload) => generateApiKey(orgId, projectId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: apiKeyKeys.all(orgId, projectId) })
        },
    })
}

export function useRevokeApiKey(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (apiKeyId: string) => revokeApiKey(orgId, projectId, apiKeyId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: apiKeyKeys.all(orgId, projectId) })
        },
    })
}
