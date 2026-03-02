import { z } from 'zod'

export const createOrganizationSchema = z.object({
    name: z.string().min(3, 'Name must be at least 3 characters').max(100, 'Name must be at most 100 characters'),
    description: z.string().max(500, 'Description must be at most 500 characters').optional(),
})

export type CreateOrganizationFormValues = z.infer<typeof createOrganizationSchema>
