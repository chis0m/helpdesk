<template>
  <div>
    <p
      class="inline-flex rounded-full border border-[var(--border-subtle)] bg-[var(--surface-main)] px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)] lg:hidden"
    >
      Reset access
    </p>
    <h1 class="mt-4 text-2xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:mt-0 lg:text-3xl">
      Forgot
      <span class="text-[var(--brand-green-dark)]">password</span>
    </h1>
    <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
      Enter your email. If it’s registered, a reset link is written to the server logs (CA: no email delivery).
    </p>

    <div
      v-if="done"
      class="mt-8 rounded-2xl border border-emerald-200/80 bg-emerald-50/90 px-4 py-3 text-sm text-emerald-900"
      role="status"
    >
      If that email is registered, a password reset link was written to the server logs. Check your backend logs for the URL.
    </div>

    <form
      v-else
      class="mt-8 space-y-5"
      @submit.prevent="onSubmit"
    >
      <!-- VULN-04: Public CSRF + POST (same weak middleware as login). -->
      <div
        v-if="errorMessage"
        class="rounded-2xl border border-red-200/80 bg-red-50/90 px-4 py-3 text-sm text-red-900"
        role="alert"
      >
        {{ errorMessage }}
      </div>

      <div>
        <label
          for="fp-email"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Email</label>
        <input
          id="fp-email"
          v-model.trim="email"
          type="email"
          name="email"
          autocomplete="email"
          required
          class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
        >
      </div>

      <button
        type="submit"
        :disabled="submitting"
        class="auth-submit mt-2 w-full rounded-full bg-[var(--brand-green)] py-3.5 text-sm font-semibold text-[var(--text-on-green)] shadow-md transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {{ submitting ? 'Sending…' : 'Request reset link' }}
      </button>
    </form>

    <p class="mt-8 border-t border-[var(--border-subtle)] pt-8 text-center text-sm text-[var(--text-secondary)]">
      <RouterLink
        :to="paths.resetPassword"
        class="font-semibold text-[var(--brand-green-dark)] transition hover:underline"
      >
        Already have a reset token?
      </RouterLink>
      <span class="mx-2 text-[var(--text-muted)]">·</span>
      <RouterLink
        :to="paths.login"
        class="font-semibold text-[var(--brand-green-dark)] transition hover:underline"
      >
        Back to sign in
      </RouterLink>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { fetchPublicCsrfToken, forgotPasswordRequest } from '@/api/auth'
import { paths } from '@/constants/routes'

const email = ref('')
const submitting = ref(false)
const errorMessage = ref('')
const done = ref(false)

async function onSubmit() {
  errorMessage.value = ''
  submitting.value = true
  try {
    const csrf = await fetchPublicCsrfToken()
    if (!csrf.ok) {
      errorMessage.value = csrf.message
      return
    }
    const result = await forgotPasswordRequest(email.value, csrf.token)
    if (!result.ok) {
      errorMessage.value = result.message
      return
    }
    done.value = true
  }
  finally {
    submitting.value = false
  }
}
</script>
