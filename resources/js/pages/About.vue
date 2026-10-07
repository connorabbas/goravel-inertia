<script setup lang="ts">
import { Head, Link } from '@inertiajs/vue3'

defineProps<{ message: string }>()
</script>

<template>

    <Head title="About" />
    <main>
        <h1>About</h1>
        <p>{{ message }}</p>
        <nav aria-label="Main navigation">
            <Link href="/">Back home</Link>
        </nav>
    </main>
</template>
