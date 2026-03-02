import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { createOrganizationSchema, type CreateOrganizationFormValues } from '@/lib/schemas/organization.schema'
import { ApiError } from '@/lib/api'
import type { OrganizationWithMembers } from '../types'
import { onboardOrganization } from '../services/organization.service'

export const useCreateOrganization = () => {
    const [serverError, setServerError] = useState<string | null>(null)
    const [createdOrg, setCreatedOrg] = useState<OrganizationWithMembers | null>(null)

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors, isSubmitting },
    } = useForm<CreateOrganizationFormValues>({
        resolver: zodResolver(createOrganizationSchema),
        defaultValues: { name: '', description: '' },
    })

    const onSubmit = async (data: CreateOrganizationFormValues) => {
        setServerError(null)
        try {
            const org = await onboardOrganization({
                name: data.name,
                description: data.description || undefined,
            })
            reset()
            setCreatedOrg(org)
        } catch (err: unknown) {
            setServerError(err instanceof ApiError ? err.message : 'An unexpected error occurred. Please try again.')
        }
    }

    return {
        register,
        handleSubmit,
        errors,
        isSubmitting,
        serverError,
        createdOrg,
        onSubmit,
    }
}
