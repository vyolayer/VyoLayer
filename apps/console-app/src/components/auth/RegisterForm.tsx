import * as React from 'react'
import { type FieldErrors, type UseFormHandleSubmit, type UseFormRegister } from 'react-hook-form'
import type { RegisterFormValues } from '@/lib/schemas/auth.schema'

import { FormField } from '@/components/auth/FormField'
import { PasswordInput } from '@/components/auth/PasswordInput'
import { LoadingButtonChildren } from '@/components/auth/AuthFormUtils'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

export const RegisterForm: React.FC<{
    register: UseFormRegister<RegisterFormValues>
    errors: FieldErrors<RegisterFormValues>
    isSubmitting: boolean
    onSubmit: (data: RegisterFormValues) => void
    handleSubmit: UseFormHandleSubmit<RegisterFormValues>
}> = ({ register, errors, isSubmitting, onSubmit, handleSubmit }) => {
    return (
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
            <FormField id="register-fullName" label="Full name" error={errors.fullName?.message}>
                <Input
                    id="register-fullName"
                    type="text"
                    autoComplete="name"
                    placeholder="John Doe"
                    {...register('fullName')}
                />
            </FormField>

            <FormField id="register-email" label="Email address" error={errors.email?.message}>
                <Input
                    id="register-email"
                    type="email"
                    autoComplete="email"
                    placeholder="you@company.com"
                    {...register('email')}
                />
            </FormField>

            <FormField
                id="register-password"
                label="Password"
                error={errors.password?.message}
                hint="8–20 characters with at least one special character (!@#$%^&*)"
            >
                <PasswordInput
                    id="register-password"
                    autoComplete="new-password"
                    placeholder="••••••••"
                    {...register('password')}
                />
            </FormField>

            <Button type="submit" disabled={isSubmitting} className="w-full mt-1" size="lg">
                <LoadingButtonChildren isLoading={isSubmitting} loadingText="Creating account…">
                    Create account
                </LoadingButtonChildren>
            </Button>
        </form>
    )
}
