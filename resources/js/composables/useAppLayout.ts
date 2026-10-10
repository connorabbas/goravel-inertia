import type { DropdownMenuItem, NavigationMenuItem } from '@nuxt/ui'
import { router, usePage } from '@inertiajs/vue3'
import { computed } from 'vue'

export function useAppLayout() {
    const page = usePage()

    const currentRoute = computed(() => page.url.split('?')[0])

    const user = computed(() => page.props.auth.user)

    const navMenuItems = computed<NavigationMenuItem[][]>(() => {
        return [
            [
                {
                    label: 'Home',
                    icon: 'i-lucide-house',
                    to: '/',
                    active: currentRoute.value === '/'
                },
                {
                    label: 'Dashboard',
                    icon: 'i-lucide-layout-dashboard',
                    to: '/dashboard',
                    active: currentRoute.value === '/dashboard'
                }
            ],
            [
                {
                    label: 'Goravel Docs',
                    icon: 'i-lucide-book-open',
                    to: 'https://www.goravel.dev/',
                    target: '_blank'
                }
            ]
        ]
    })

    const userMenuItems = computed<DropdownMenuItem[][]>(() => {
        const items: DropdownMenuItem[][] = [[
            {
                label: 'Settings',
                icon: 'i-lucide-settings',
                to: '/settings/profile'
            }
        ]]

        items.push([
            {
                label: 'Log out',
                icon: 'i-lucide-log-out',
                onSelect: () => {
                    router.post('/logout')
                }
            }
        ])

        return items
    })

    return {
        currentRoute,
        navMenuItems,
        userMenuItems,
        user
    }
}
