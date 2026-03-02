import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { changePasswordSchema, type ChangePasswordFormValues } from '@/lib/schemas/auth.schema'
import { changeUserPassword } from '@/lib/services/auth.service'
import { ApiError } from '@/lib/api'

export const useChangePassword = () => {
    const [serverError, setServerError] = useState<string | null>(null)
    const [isSuccess, setIsSuccess] = useState(false)

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors, isSubmitting },
    } = useForm<ChangePasswordFormValues>({
        resolver: zodResolver(changePasswordSchema),
        defaultValues: { currentPassword: '', newPassword: '', confirmPassword: '' },
    })

    const onSubmit = async (data: ChangePasswordFormValues) => {
        setServerError(null)
        try {
            await changeUserPassword({
                currentPassword: data.currentPassword,
                newPassword: data.newPassword,
            })
            reset()
            setIsSuccess(true)
        } catch (err: unknown) {
            setServerError(err instanceof ApiError ? err.message : 'An unexpected error occurred. Please try again.')
        }
    }

    return {
        register,
        handleSubmit,
        reset,
        formState: { errors, isSubmitting },
        serverError,
        isSuccess,
        onSubmit,
    }
}
