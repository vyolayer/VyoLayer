import * as React from 'react'
import type { FieldErrors, UseFormHandleSubmit, UseFormRegister } from 'react-hook-form'
import type { CreateOrganizationFormValues } from '@/lib/schemas/organization.schema'

import { FormField } from '@/components/auth/FormField'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { LoadingButtonChildren } from '@/components/auth/AuthFormUtils'

export const CreateOrganizationForm: React.FC<{
    register: UseFormRegister<CreateOrganizationFormValues>
    errors: FieldErrors<CreateOrganizationFormValues>
    isSubmitting: boolean
    onSubmit: (data: CreateOrganizationFormValues) => void
    handleSubmit: UseFormHandleSubmit<CreateOrganizationFormValues>
    submitLabel?: string
}> = ({ register, errors, isSubmitting, onSubmit, handleSubmit, submitLabel = 'Create Organization' }) => {
    return (
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
            <FormField
                id="org-name"
                label="Organization name"
                error={errors.name?.message}
                hint="This name will appear across your workspace."
            >
                <Input
                    id="org-name"
                    type="text"
                    autoComplete="organization"
                    placeholder="Acme Labs"
                    {...register('name')}
                />
            </FormField>

            <FormField
                id="org-description"
                label="Description"
                error={errors.description?.message}
                hint="Optional — briefly describe what this organization is for."
            >
                <Textarea
                    id="org-description"
                    placeholder="Building amazing products…"
                    rows={3}
                    className="resize-none"
                    {...register('description')}
                />
            </FormField>

            <Button type="submit" disabled={isSubmitting} className="w-full mt-1" size="lg">
                <LoadingButtonChildren isLoading={isSubmitting} loadingText="Creating…">
                    {submitLabel}
                </LoadingButtonChildren>
            </Button>
        </form>
    )
}
