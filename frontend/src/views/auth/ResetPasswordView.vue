<template>
  <div>
    <p
      class="inline-flex rounded-full border border-[var(--border-subtle)] bg-[var(--surface-main)] px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)] lg:hidden"
    >
      New password
    </p>
    <h1 class="mt-4 text-2xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:mt-0 lg:text-3xl">
      Set a new
      <span class="text-[var(--brand-green-dark)]">password</span>
    </h1>
    <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
      Paste the 64-character token from the reset link in server logs, or open a URL with <code class="rounded bg-[var(--surface-muted)] px-1">?token=…</code>.
    </p>

    <div
      v-if="!tokenFromQuery"
      class="mt-6 rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-950"
      role="status"
    >
      No <code class="rounded bg-white px-1">token</code> in the URL — paste the token into the field below.
    </div>

    <form
      class="mt-8 space-y-5"
      @submit.prevent="onSubmit"
    >
      <!-- VULN-04: Public CSRF + POST. -->
      <div
        v-if="errorMessage"
        class="rounded-2xl border border-red-200/80 bg-red-50/90 px-4 py-3 text-sm text-red-900"
        role="alert"
      >
        {{ errorMessage }}
      </div>

      <div>
        <label
          for="rp-token"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Reset token</label>
        <input
          id="rp-token"
          v-model.trim="tokenField"
          type="text"
          maxlength="64"
          required
          class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 font-mono text-xs text-[var(--text-primary)] shadow-sm outline-none focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
        >
      </div>

      <div>
        <label
          for="rp-password"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >New password</label>
        <input
          id="rp-password"
          v-model="password"
          type="password"
          autocomplete="new-password"
          required
          minlength="8"
          maxlength="128"
          class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
        >
      </div>

      <div>
        <label
          for="rp-password2"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Confirm new password</label>
        <input
          id="rp-password2"
          v-model="password2"
          type="password"
          autocomplete="new-password"
          required
          minlength="8"
          maxlength="128"
          class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
        >
      </div>

      <button
        type="submit"
        :disabled="submitting"
        class="auth-submit mt-2 w-full rounded-full bg-[var(--brand-green)] py-3.5 text-sm font-semibold text-[var(--text-on-green)] shadow-md transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {{ submitting ? 'Updating…' : 'Update password' }}
      </button>
    </form>

    <p class="mt-8 border-t border-[var(--border-subtle)] pt-8 text-center text-sm text-[var(--text-secondary)]">
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
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchPublicCsrfToken, resetPasswordRequest } from '@/api/auth'
import { paths } from '@/constants/routes'

const route = useRoute()
const router = useRouter()

const tokenFromQuery = computed(() => {
  const q = route.query.token
  return typeof q === 'string' && q.trim().length === 64 ? q.trim() : ''
})

const tokenField = ref('')
watch(
  tokenFromQuery,
  (t) => {
    if (t)
      tokenField.value = t
  },
  { immediate: true },
)

const password = ref('')
const password2 = ref('')
const submitting = ref(false)
const errorMessage = ref('')

async function onSubmit() {
  errorMessage.value = ''
  if (password.value !== password2.value) {
    errorMessage.value = 'Passwords do not match.'
    return
  }
  if (password.value.length < 8) {
    errorMessage.value = 'Password must be at least 8 characters.'
    return
  }
  const tok = tokenField.value.trim()
  if (tok.length !== 64) {
    errorMessage.value = 'Token must be exactly 64 characters.'
    return
  }

  submitting.value = true
  try {
    const csrf = await fetchPublicCsrfToken()
    if (!csrf.ok) {
      errorMessage.value = csrf.message
      return
    }
    const result = await resetPasswordRequest(tok, password.value, csrf.token)
    if (!result.ok) {
      errorMessage.value = result.message
      return
    }
    const target = result.data.redirect_to.startsWith('/') ? result.data.redirect_to : paths.login
    await router.replace(target)
  }
  finally {
    submitting.value = false
  }
}
</script>
