<script setup lang="ts">
import type { NavigationMenuItem } from '@nuxt/ui'
import { computed } from 'vue'
import AppLayout from '@/layouts/app/Index.vue'
import { useAppLayout } from '@/composables/useAppLayout'

const props = defineProps<{
    title: string
    description: string
}>()

const { currentRoute } = useAppLayout()

const settingsPageTitle = computed(() => `Settings - ${props.title}`)

const items = computed<NavigationMenuItem[]>(() => [
    {
        label: 'Profile',
        icon: 'i-lucide-user-round',
        to: '/settings/profile',
        active: currentRoute.value === '/settings/profile'
    },
    {
        label: 'Password',
        icon: 'i-lucide-key-round',
        to: '/settings/password',
        active: currentRoute.value === '/settings/password'
    },
    {
        label: 'Appearance',
        icon: 'i-lucide-palette',
        to: '/settings/appearance',
        active: currentRoute.value === '/settings/appearance'
    }
])
</script>

<template>
    <AppLayout
        :title="settingsPageTitle"
        :description="props.description"
    >
        <UPageHeader
            :title="settingsPageTitle"
            :description="props.description"
        />

        <UPage>
            <template #left>
                <UPageAside>
                    <UNavigationMenu
                        :items="items"
                        orientation="vertical"
                        variant="pill"
                    />
                </UPageAside>
            </template>

            <div class="lg:hidden">
                <UNavigationMenu
                    :items="items"
                    orientation="vertical"
                    variant="pill"
                    class="w-full mt-6"
                />
            </div>

            <UPageBody>
                <slot />
            </UPageBody>
        </UPage>
    </AppLayout>
</template>
