import { api } from '@/lib/api'
import type { ApiKey, ApiKeyCreated, ApiKeyMode } from '../types'

// ─── Payloads ─────────────────────────────────────────────────────────────────

export interface CreateApiKeyPayload {
    name: string
    mode: ApiKeyMode
}

// ─── API Key API ──────────────────────────────────────────────────────────────

const base = (orgId: string, projectId: string) => `/organizations/${orgId}/projects/${projectId}/api-keys`

export async function listApiKeys(orgId: string, projectId: string): Promise<ApiKey[]> {
    return api.get(base(orgId, projectId))
}

export async function getApiKeyById(orgId: string, projectId: string, apiKeyId: string): Promise<ApiKey> {
    return api.get(`${base(orgId, projectId)}/${apiKeyId}`)
}

export async function generateApiKey(
    orgId: string,
    projectId: string,
    payload: CreateApiKeyPayload,
): Promise<ApiKeyCreated> {
    return api.post(base(orgId, projectId), payload)
}

export async function revokeApiKey(orgId: string, projectId: string, apiKeyId: string): Promise<void> {
    return api.delete(`${base(orgId, projectId)}/${apiKeyId}`)
}
