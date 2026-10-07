<script setup lang="ts">
import { Head, usePage } from '@inertiajs/vue3'

defineProps<{ message: string }>()
const page = usePage()
</script>

<template>

    <Head title="Home" />
    <main class="mx-auto max-w-xl px-4 py-12">
        <UCard>
            <template #header>
                <h1 class="text-xl font-semibold text-highlighted">{{ page.props.appName }}</h1>
            </template>
            <p class="text-default">{{ message }}</p>
            <template #footer>
                <nav aria-label="Main navigation">
                    <UButton to="/about">About this test</UButton>
                </nav>
            </template>
        </UCard>
    </main>
</template>
