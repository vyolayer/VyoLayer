import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import type React from 'react'
import type { NavItem } from './OrganizationHeader'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'

export const OrganizationBreadcrumb: React.FC<{
    navItems: NavItem[]
    currentPath: string
    orgSlug: string
    orgName: string
}> = ({ navItems, currentPath, orgName, orgSlug }) => {
    const currentLabel = navItems.find((i) => i.to === currentPath)?.label

    return (
        <Breadcrumb>
            <BreadcrumbList>
                <BreadcrumbItem>
                    <BreadcrumbLink to="/dashboard">Organizations</BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator />
                <BreadcrumbItem>
                    <BreadcrumbLink to={`/org/${orgSlug}`}>{orgName}</BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator />
                <DropdownMenu>
                    <DropdownMenuTrigger>
                        <BreadcrumbPage>{currentLabel}</BreadcrumbPage>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent>
                        {navItems.map((item) => (
                            <DropdownMenuItem key={item.to} asChild>
                                <BreadcrumbLink to={item.to}>{item.label}</BreadcrumbLink>
                            </DropdownMenuItem>
                        ))}
                    </DropdownMenuContent>
                </DropdownMenu>
                <BreadcrumbItem></BreadcrumbItem>
            </BreadcrumbList>
        </Breadcrumb>
    )
}
