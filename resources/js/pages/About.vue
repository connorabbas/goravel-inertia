<script setup lang="ts">
import { Head } from '@inertiajs/vue3'

defineProps<{ message: string }>()
</script>

<template>

    <Head title="About" />
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
</template>
