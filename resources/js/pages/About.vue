<script setup lang="ts">
import DefaultLayout from '@/layouts/Default.vue'

defineProps<{ message: string }>()
</script>

<template>

    <DefaultLayout title="About">
    <main class="mx-auto max-w-xl px-4 py-12">
        <UCard>
            <template #header>
                <h1 class="text-xl font-semibold text-highlighted">About</h1>
            </template>
            <p class="text-default">{{ message }}</p>
            <template #footer>
                <nav aria-label="Main navigation">
                    <UButton to="/" color="neutral" variant="outline">Back home</UButton>
                </nav>
            </template>
        </UCard>
    </main>
    </DefaultLayout>
</template>
