import * as React from 'react'
import { type FieldErrors, type UseFormHandleSubmit, type UseFormRegister } from 'react-hook-form'
import type { ChangePasswordFormValues } from '@/lib/schemas/auth.schema'

import { FormField } from '@/components/auth/FormField'
import { PasswordInput } from '@/components/auth/PasswordInput'
import { LoadingButtonChildren } from '@/components/auth/AuthFormUtils'
import { Button } from '@/components/ui/button'

export const ChangePasswordForm: React.FC<{
    register: UseFormRegister<ChangePasswordFormValues>
    errors: FieldErrors<ChangePasswordFormValues>
    isSubmitting: boolean
    onSubmit: (data: ChangePasswordFormValues) => void
    handleSubmit: UseFormHandleSubmit<ChangePasswordFormValues>
}> = ({ register, errors, isSubmitting, onSubmit, handleSubmit }) => {
    return (
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
            <FormField id="cp-current" label="Current password" error={errors.currentPassword?.message}>
                <PasswordInput
                    id="cp-current"
                    autoComplete="current-password"
                    placeholder="••••••••"
                    {...register('currentPassword')}
                />
            </FormField>

            <FormField
                id="cp-new"
                label="New password"
                error={errors.newPassword?.message}
                hint="8–20 characters with at least one special character (!@#$%^&*)"
            >
                <PasswordInput
                    id="cp-new"
                    autoComplete="new-password"
                    placeholder="••••••••"
                    {...register('newPassword')}
                />
            </FormField>

            <FormField id="cp-confirm" label="Confirm new password" error={errors.confirmPassword?.message}>
                <PasswordInput
                    id="cp-confirm"
                    autoComplete="new-password"
                    placeholder="••••••••"
                    {...register('confirmPassword')}
                />
            </FormField>

            <Button type="submit" disabled={isSubmitting} className="w-full mt-1" size="lg">
                <LoadingButtonChildren isLoading={isSubmitting} loadingText="Updating password…">
                    Update password
                </LoadingButtonChildren>
            </Button>
        </form>
    )
}
