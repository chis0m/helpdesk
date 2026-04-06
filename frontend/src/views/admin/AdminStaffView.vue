<template>
  <div class="space-y-6">
    <AdminSubnav />

    <div class="mx-auto max-w-lg space-y-6">
      <header>
      <p class="text-xs font-semibold uppercase tracking-wider text-[var(--brand-green-dark)]">
        Admin
      </p>
      <h2 class="mt-2 text-lg font-semibold text-[var(--text-primary)]">
        Create staff account
      </h2>
      <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
        Invite someone to the internal team. They appear on the Staff list once created.
      </p>
    </header>

    <form
      class="space-y-5 rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-sm"
      @submit.prevent="onSubmit"
    >
      <div
        v-if="banner"
        class="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-900"
        role="status"
      >
        {{ banner }}
      </div>

      <div>
        <label
          for="st-email"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Work email</label>
        <input
          id="st-email"
          v-model="form.email"
          type="email"
          required
          autocomplete="email"
          placeholder="new.staff@secweb.internal"
          class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm outline-none ring-[var(--brand-green)] focus:border-transparent focus:ring-2"
        >
      </div>

      <div class="grid gap-4 sm:grid-cols-2">
        <div>
          <label
            for="st-first"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >First name</label>
          <input
            id="st-first"
            v-model="form.firstName"
            type="text"
            required
            class="w-full rounded-2xl border border-[var(--border-subtle)] px-4 py-3 text-sm"
          >
        </div>
        <div>
          <label
            for="st-last"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Last name</label>
          <input
            id="st-last"
            v-model="form.lastName"
            type="text"
            required
            class="w-full rounded-2xl border border-[var(--border-subtle)] px-4 py-3 text-sm"
          >
        </div>
      </div>

      <div class="flex items-start gap-3 rounded-2xl border border-[var(--border-subtle)] bg-[var(--surface-muted)]/40 px-4 py-3">
        <input
          id="st-admin"
          v-model="form.isAdmin"
          type="checkbox"
          class="mt-1 h-4 w-4 rounded border-[var(--border-strong)] text-[var(--brand-green)]"
        >
        <div>
          <label
            for="st-admin"
            class="text-sm font-medium text-[var(--text-primary)]"
          >Administrator</label>
          <p class="mt-0.5 text-xs text-[var(--text-muted)]">
            Full access to admin areas and sensitive operations.
          </p>
        </div>
      </div>

      <div>
        <label
          for="st-temp"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Temporary password</label>
        <input
          id="st-temp"
          v-model="form.password"
          type="password"
          autocomplete="new-password"
          placeholder="At least 8 characters"
          class="w-full rounded-2xl border border-[var(--border-subtle)] px-4 py-3 text-sm"
        >
      </div>

      <button
        type="submit"
        class="w-full rounded-full bg-[var(--brand-green)] py-3.5 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95"
      >
        Create staff
      </button>
    </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import AdminSubnav from '@/components/admin/AdminSubnav.vue'

const form = reactive({
  email: '',
  firstName: '',
  lastName: '',
  isAdmin: false,
  password: '',
})

const banner = ref('')

function onSubmit() {
  const role = form.isAdmin ? 'Administrator' : 'Staff'
  banner.value = `Invitation prepared for ${form.email} (${role}).`
  window.setTimeout(() => {
    banner.value = ''
  }, 4000)
}
</script>
