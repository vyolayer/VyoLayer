import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import type { LoginFormValues } from '@/lib/schemas/auth.schema'

import { ApiError } from '@/lib/api'
import { loginSchema } from '@/lib/schemas/auth.schema'

import { useAuth } from '@/contexts/AuthContext'

export const useLogin = () => {
    const { login } = useAuth()
    const [serverError, setServerError] = useState<string | null>(null)

    const {
        register,
        handleSubmit,
        formState: { errors, isSubmitting },
    } = useForm<LoginFormValues>({
        resolver: zodResolver(loginSchema),
        defaultValues: { email: '', password: '' },
    })

    const onSubmit = async (data: LoginFormValues) => {
        setServerError(null)
        try {
            await login(data)
            window.location.href = '/'
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
        onSubmit,
    }
}
