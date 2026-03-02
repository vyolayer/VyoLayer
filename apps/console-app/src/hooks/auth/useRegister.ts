import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'

import { useAuth } from '@/contexts/AuthContext'
import { registerSchema, type RegisterFormValues } from '@/lib/schemas/auth.schema'
import { ApiError } from '@/lib/api'
import { useForm } from 'react-hook-form'

export const useRegister = () => {
    const { register: registerUser } = useAuth()
    const [serverError, setServerError] = useState<string | null>(null)
    const [isRegistered, setIsRegistered] = useState(false)

    const {
        register,
        handleSubmit,
        formState: { errors, isSubmitting },
    } = useForm<RegisterFormValues>({
        resolver: zodResolver(registerSchema),
        defaultValues: { fullName: '', email: '', password: '' },
    })

    const onSubmit = async (data: RegisterFormValues) => {
        setServerError(null)
        try {
            await registerUser(data)
            setIsRegistered(true)
        } catch (err: unknown) {
            setServerError(err instanceof ApiError ? err.message : 'An unexpected error occurred. Please try again.')
        }
    }

    return {
        register,
        handleSubmit,
        errors,
        isSubmitting,
        onSubmit,
        isRegistered,
        serverError,
    }
}
