<template>
  <div class="mx-auto max-w-xl space-y-6">
    <header>
      <h2 class="text-lg font-semibold text-[var(--text-primary)]">
        Your profile
      </h2>
      <p class="mt-1 text-sm text-[var(--text-secondary)]">
        Update the details on your SecWeb Helpdesk account.
      </p>
    </header>

    <form
      class="space-y-5 rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-sm"
      @submit.prevent="onSave"
    >
      <div
        v-if="savedBanner"
        class="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-900"
        role="status"
      >
        Your profile was updated.
      </div>

      <div>
        <label
          for="pf-email"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Email</label>
        <input
          id="pf-email"
          v-model="form.email"
          type="email"
          autocomplete="email"
          class="w-full rounded-2xl border border-[var(--border-subtle)] bg-[var(--surface-main)] px-4 py-3 text-sm text-[var(--text-primary)]"
        >
      </div>

      <div class="grid gap-4 sm:grid-cols-2">
        <div>
          <label
            for="pf-first"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >First name</label>
          <input
            id="pf-first"
            v-model="form.firstName"
            type="text"
            class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm"
          >
        </div>
        <div>
          <label
            for="pf-last"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Last name</label>
          <input
            id="pf-last"
            v-model="form.lastName"
            type="text"
            class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm"
          >
        </div>
      </div>

      <button
        type="submit"
        class="rounded-full bg-[var(--brand-green)] px-6 py-3 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95"
      >
        Save changes
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'

const form = reactive({
  email: 'jordan.lee@customer.example',
  firstName: 'Jordan',
  lastName: 'Lee',
})

const savedBanner = ref(false)

function onSave() {
  savedBanner.value = true
  window.setTimeout(() => {
    savedBanner.value = false
  }, 3000)
}
</script>
