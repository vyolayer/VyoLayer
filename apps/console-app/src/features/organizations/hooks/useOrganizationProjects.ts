import { useEffect, useState } from 'react'

type OrganizationProject = {
    id: string
    name: string
    slug: string
    description: string
}

const dummyProjects: OrganizationProject[] = [
    {
        id: '1',
        name: 'Project 1',
        slug: 'project-1',
        description: 'Project 1 description',
    },
    {
        id: '2',
        name: 'Project 2',
        slug: 'project-2',
        description: 'Project 2 description',
    },
    {
        id: '3',
        name: 'Project 3',
        slug: 'project-3',
        description: 'Project 3 description',
    },
]

export const useOrganizationProjects = (orgId: string) => {
    const [isLoading, setIsLoading] = useState(true)
    // random number 0-3 to simulate loading
    const randomNumber = Math.floor(Math.random() * 100) % 4
    // dummy projects
    if (randomNumber === 0) {
        return {
            projects: dummyProjects,
            isLoading: false,
            error: null,
        }
    }

    const theProjects = dummyProjects.slice(0, randomNumber)

    useEffect(() => {
        setIsLoading(true)
        const randomTimeout = Math.floor(Math.random() * 1000) * 5
        const timer = setTimeout(() => setIsLoading(false), randomTimeout)

        return () => clearTimeout(timer)
    }, [])

    return {
        projects: theProjects,
        isLoading: isLoading,
        error: null,
    }
}
