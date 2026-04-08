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
          v-if="errorMessage"
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-900"
          role="alert"
        >
          {{ errorMessage }}
        </div>

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
              Sends <code class="rounded bg-black/5 px-1 py-0.5 text-[11px]">role: admin</code> in the create request.
              Only a <strong>super admin</strong> may create an administrator; a normal admin receives an error if this is checked.
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
          :disabled="submitting"
          class="w-full rounded-full bg-[var(--brand-green)] py-3.5 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ submitting ? 'Creating…' : 'Create staff' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import AdminSubnav from '@/components/admin/AdminSubnav.vue'
import { createStaffUser } from '@/api/admin'
import { getSessionCsrfToken } from '@/stores/auth-session'

const form = reactive({
  email: '',
  firstName: '',
  lastName: '',
  isAdmin: false,
  password: '',
})

const banner = ref('')
const errorMessage = ref('')
const submitting = ref(false)

async function onSubmit() {
  errorMessage.value = ''
  banner.value = ''

  const csrf = getSessionCsrfToken()
  if (!csrf) {
    errorMessage.value = 'Your session is missing a security token. Sign out and sign in again, then retry.'
    return
  }

  submitting.value = true
  try {
    const res = await createStaffUser(
      {
        email: form.email.trim(),
        password: form.password,
        first_name: form.firstName.trim(),
        last_name: form.lastName.trim(),
        role: form.isAdmin ? 'admin' : 'staff',
      },
      csrf,
    )
    if (!res.ok) {
      errorMessage.value = res.message
      return
    }
    const note =
      res.data.role === 'admin'
        ? ' Account role is administrator.'
        : ''
    banner.value = `Account created for ${res.data.email} (id ${res.data.user_id}, role ${res.data.role}).${note}`
    form.email = ''
    form.firstName = ''
    form.lastName = ''
    form.password = ''
    form.isAdmin = false
    window.setTimeout(() => {
      banner.value = ''
    }, 8000)
  }
  finally {
    submitting.value = false
  }
}
</script>
