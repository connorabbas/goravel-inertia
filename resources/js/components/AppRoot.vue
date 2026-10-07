<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { router, usePage } from '@inertiajs/vue3'
import { useToast } from '@nuxt/ui/composables'
import { errorToast, flashToasts } from '../utils/notifications'

const toast = useToast()
const page = usePage()

function showFlash(flash: unknown) {
    flashToasts(flash).forEach((notification) => toast.add(notification))
}

// The Inertia child initializes usePage during mount. Only fresh visits should
// show flash messages; browser-history navigation must not replay old toasts.
onMounted(() => showFlash(page.props?.flash))

const unsubscribe = [
    router.on('success', (event) => showFlash(event.detail.page.props.flash)),
    router.on('httpException', (event) => {
        const notification = errorToast(event.detail.response.data, event.detail.response.status)
        if (!notification) return

        event.preventDefault()
        toast.add(notification)
    }),
    router.on('networkError', (event) => {
        event.preventDefault()
        toast.add({
            color: 'error',
            title: 'Connection Error',
            description: 'An unexpected error occurred while loading this page. Please try again.',
            icon: 'i-lucide-wifi-off',
        })
    }),
]

onUnmounted(() => unsubscribe.forEach((remove) => remove()))
</script>

<template>
    <UApp>
        <slot />
    </UApp>
</template>
