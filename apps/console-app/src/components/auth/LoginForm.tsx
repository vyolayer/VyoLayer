import type { LoginFormValues } from '@/lib/schemas/auth.schema'
import * as React from 'react'
import type { FieldErrors, UseFormHandleSubmit, UseFormRegister } from 'react-hook-form'
import { Button } from '../ui/button'
import { LoadingButtonChildren } from './AuthFormUtils'
import { Input } from '../ui/input'
import { FormField } from './FormField'
import { PasswordInput } from './PasswordInput'

export const LoginForm: React.FC<{
    onSubmit: (data: LoginFormValues) => void
    errors: FieldErrors<LoginFormValues>
    isSubmitting: boolean
    register: UseFormRegister<LoginFormValues>
    handleSubmit: UseFormHandleSubmit<LoginFormValues>
}> = ({ onSubmit, errors, isSubmitting, register, handleSubmit }) => {
    return (
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
            <FormField id="login-email" label="Email address" error={errors.email?.message}>
                <Input
                    id="login-email"
                    type="email"
                    autoComplete="email"
                    placeholder="you@company.com"
                    {...register('email')}
                />
            </FormField>

            <FormField id="login-password" label="Password" error={errors.password?.message}>
                <PasswordInput
                    id="login-password"
                    autoComplete="current-password"
                    placeholder="••••••••"
                    {...register('password')}
                />
            </FormField>

            <Button type="submit" disabled={isSubmitting} className="w-full mt-1" size="lg">
                <LoadingButtonChildren isLoading={isSubmitting} loadingText="Signing in…">
                    Sign in
                </LoadingButtonChildren>
            </Button>
        </form>
    )
}
