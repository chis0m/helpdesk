<template>
  <div class="flex min-h-screen bg-[var(--surface-page)] text-[var(--text-primary)]">
    <AppSidebar />
    <div class="flex min-w-0 flex-1 flex-col">
      <AppHeader :title="pageTitle" />
      <main class="relative flex-1 overflow-auto bg-[var(--surface-main)]">
        <div
          class="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_90%_60%_at_50%_-10%,rgba(34,197,94,0.09),transparent_55%)]"
          aria-hidden="true"
        />
        <div class="relative mx-auto max-w-6xl px-5 py-6 lg:px-8 lg:py-8">
          <RouterView />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import AppHeader from '@/components/layout/AppHeader.vue'

const route = useRoute()

const pageTitle = computed(() => {
  const m = route.meta.title
  return typeof m === 'string' ? m : 'Support home'
})
</script>
