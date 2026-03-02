import { z } from 'zod'

export const loginSchema = z.object({
    email: z
        .string()
        .min(1, 'Email is required')
        .email('Please enter a valid email')
        .max(255, 'Email must be at most 255 characters'),
    password: z.string().min(1, 'Password is required'),
})

export const registerSchema = z.object({
    fullName: z.string().min(3, 'Name must be at least 3 characters').max(100, 'Name must be at most 100 characters'),
    email: z
        .string()
        .min(1, 'Email is required')
        .email('Please enter a valid email')
        .max(255, 'Email must be at most 255 characters'),
    password: z
        .string()
        .min(8, 'Password must be at least 8 characters')
        .max(20, 'Password must be at most 20 characters')
        .regex(/[!@#$%^&*]/, 'Password must contain at least one special character (!@#$%^&*)'),
})

export const changePasswordSchema = z
    .object({
        currentPassword: z.string().min(1, 'Current password is required'),
        newPassword: z
            .string()
            .min(8, 'Password must be at least 8 characters')
            .max(20, 'Password must be at most 20 characters')
            .regex(/[!@#$%^&*]/, 'Password must contain at least one special character (!@#$%^&*)'),
        confirmPassword: z.string().min(1, 'Please confirm your new password'),
    })
    .refine((data) => data.newPassword === data.confirmPassword, {
        message: 'Passwords do not match',
        path: ['confirmPassword'],
    })

export type LoginFormValues = z.infer<typeof loginSchema>
export type RegisterFormValues = z.infer<typeof registerSchema>
export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>
