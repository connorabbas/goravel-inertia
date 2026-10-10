<script setup lang="ts">
import { Link, usePage } from '@inertiajs/vue3'
import FlashAlerts from '@/components/FlashAlerts.vue'
import AppHead from '@/components/AppHead.vue'
import AppLogo from '@/components/AppLogo.vue'
import type { NavigationMenuItem } from '@nuxt/ui'

const props = defineProps<{
    title: string
    description?: string
}>()
const page = usePage()

const navMenuItems: NavigationMenuItem[] = [
    {
        label: 'Goravel Docs',
        icon: 'i-lucide-book-open',
        to: 'https://www.goravel.dev/',
        target: '_blank'
    },
    {
        label: 'Nuxt UI Docs',
        icon: 'i-lucide-book-open',
        to: 'https://ui.nuxt.com/',
        target: '_blank'
    }
]
</script>

<template>
    <div>
        <AppHead
            :title="props.title"
            :description="props.description"
        />

        <UHeader :ui="{ right: 'flex sm:gap-3' }">
            <template #left>
                <Link
                    href="/"
                    aria-label="Application logo"
                >
                    <AppLogo class="h-6 w-auto shrink-0" />
                </Link>
            </template>

            <UNavigationMenu
                :items="navMenuItems"
                variant="link"
            />

            <template #right>
                <UButton
                    v-if="!page.props.auth.user"
                    icon="i-lucide-log-in"
                    aria-label="Sign in"
                    color="neutral"
                    variant="ghost"
                    class="lg:hidden"
                    :to="'/login'"
                />

                <UButton
                    v-if="!page.props.auth.user"
                    label="Sign in"
                    color="neutral"
                    variant="outline"
                    class="hidden lg:inline-flex"
                    :to="'/login'"
                />

                <UButton
                    v-if="!page.props.auth.user"
                    label="Sign up"
                    color="neutral"
                    trailing-icon="i-lucide-arrow-right"
                    class="hidden lg:inline-flex"
                    :to="'/register'"
                />

                <UButton
                    v-if="page.props.auth.user"
                    to="/dashboard"
                    label="Dashboard"
                />

                <UColorModeButton />

                <UButton
                    to="https://github.com/goravel/goravel"
                    target="_blank"
                    icon="simple-icons:github"
                    aria-label="GitHub"
                    color="neutral"
                    variant="ghost"
                />
            </template>

            <template #body>
                <UNavigationMenu
                    :items="navMenuItems"
                    orientation="vertical"
                    class="-mx-2.5"
                />
            </template>
        </UHeader>

        <UMain>
            <UContainer class="pt-4">
                <FlashAlerts />
            </UContainer>

            <slot />
        </UMain>

        <USeparator icon="i-lucide-box" />

        <UFooter>
            <template #left>
                <p class="text-muted text-sm">
                    Built with Nuxt UI • © {{ new Date().getFullYear() }}
                </p>
            </template>

            <template #right>
                <UButton
                    to="https://github.com/goravel/goravel"
                    target="_blank"
                    icon="simple-icons:github"
                    aria-label="GitHub"
                    color="neutral"
                    variant="ghost"
                />
            </template>
        </UFooter>
    </div>
</template>
