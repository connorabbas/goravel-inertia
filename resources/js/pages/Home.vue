<script setup lang="ts">
import { Head, Link, usePage } from '@inertiajs/vue3'

defineProps<{ message: string }>()
const page = usePage()
</script>

<template>

    <Head title="Home" />
    <main>
        <h1>{{ page.props.appName }}</h1>
        <p>{{ message }}</p>
        <nav aria-label="Main navigation">
            <Link href="/about">About this test</Link>
        </nav>
    </main>
</template>
