<template>
  <div>
    <p
      class="inline-flex rounded-full border border-[var(--border-subtle)] bg-[var(--surface-main)] px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)] lg:hidden"
    >
      Create account
    </p>
    <h1 class="mt-4 text-2xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:mt-0 lg:text-3xl">
      Create your
      <span class="text-[var(--brand-green-dark)]">account</span>
    </h1>
    <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
      For customers and users of SecWeb products — a few fields so you can submit tickets when something needs fixing.
    </p>

    <form
      class="mt-8 space-y-6"
      @submit.prevent="onSubmit"
    >
      <div
        v-if="errorMessage"
        class="rounded-2xl border border-red-200/80 bg-red-50/90 px-4 py-3 text-sm text-red-900"
        role="alert"
      >
        {{ errorMessage }}
      </div>

      <fieldset class="space-y-4">
        <legend class="mb-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-muted)]">
          Your name
        </legend>
        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label
              for="su-first"
              class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
            >First name</label>
            <input
              id="su-first"
              v-model.trim="firstName"
              type="text"
              name="given-name"
              autocomplete="given-name"
              required
              minlength="2"
              maxlength="100"
              class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
              placeholder="Jane"
            >
          </div>
          <div>
            <label
              for="su-last"
              class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
            >Last name</label>
            <input
              id="su-last"
              v-model.trim="lastName"
              type="text"
              name="family-name"
              autocomplete="family-name"
              required
              minlength="2"
              maxlength="100"
              class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
              placeholder="Doe"
            >
          </div>
        </div>
        <div>
          <label
            for="su-middle"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Middle name <span class="font-normal text-[var(--text-muted)]">(optional)</span></label>
          <input
            id="su-middle"
            v-model.trim="middleName"
            type="text"
            maxlength="100"
            class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
            placeholder="Optional"
          >
        </div>
      </fieldset>

      <fieldset class="space-y-4">
        <legend class="mb-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-muted)]">
          Account
        </legend>
        <div>
          <label
            for="su-email"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Email</label>
          <input
            id="su-email"
            v-model.trim="email"
            type="email"
            name="email"
            autocomplete="email"
            required
            class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
            placeholder="you@company.com"
          >
        </div>

        <div>
          <label
            for="su-password"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Password</label>
          <input
            id="su-password"
            v-model="password"
            type="password"
            name="new-password"
            autocomplete="new-password"
            required
            minlength="8"
            maxlength="128"
            class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
            placeholder="Choose a password"
          >
          <p class="mt-1.5 text-xs text-[var(--text-muted)]">
            At least 8 characters (required by the server).
          </p>
        </div>

        <div>
          <label
            for="su-password2"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Confirm password</label>
          <input
            id="su-password2"
            v-model="passwordConfirm"
            type="password"
            required
            minlength="8"
            maxlength="128"
            class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
            placeholder="Repeat password"
          >
          <p
            v-if="passwordConfirm.length > 0"
            class="mt-1.5 text-xs font-medium"
            :class="passwordsMatch ? 'text-[var(--brand-green-dark)]' : 'text-amber-800'"
          >
            {{ passwordsMatch ? 'Passwords match' : 'Passwords do not match' }}
          </p>
        </div>
      </fieldset>

      <p class="text-center text-xs leading-relaxed text-[var(--text-muted)]">
        By creating a SecWeb Helpdesk account you agree to our Terms and Privacy Policy (placeholder copy).
      </p>

      <button
        type="submit"
        :disabled="submitting"
        class="auth-submit w-full rounded-full bg-[var(--brand-green)] py-3.5 text-sm font-semibold text-[var(--text-on-green)] shadow-md transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {{ submitting ? 'Creating account…' : 'Create account' }}
      </button>
    </form>

    <p class="mt-8 border-t border-[var(--border-subtle)] pt-8 text-center text-sm text-[var(--text-secondary)]">
      Already registered?
      <RouterLink
        to="/login"
        class="font-semibold text-[var(--brand-green-dark)] transition hover:underline"
      >
        Sign in
      </RouterLink>
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { fetchPublicCsrfToken, signupRequest } from '@/api/auth'

const router = useRouter()

const firstName = ref('')
const lastName = ref('')
const middleName = ref('')
const email = ref('')
const password = ref('')
const passwordConfirm = ref('')
const errorMessage = ref('')
const submitting = ref(false)

const passwordsMatch = computed(
  () => password.value.length > 0 && password.value === passwordConfirm.value,
)

async function onSubmit() {
  errorMessage.value = ''
  if (firstName.value.trim().length < 2 || lastName.value.trim().length < 2) {
    errorMessage.value = 'First and last name must be at least 2 characters.'
    return
  }
  if (password.value.length < 8) {
    errorMessage.value = 'Password must be at least 8 characters.'
    return
  }
  if (password.value !== passwordConfirm.value) {
    errorMessage.value = 'Passwords do not match'
    return
  }

  submitting.value = true
  try {
    const csrf = await fetchPublicCsrfToken()
    if (!csrf.ok) {
      errorMessage.value = csrf.message
      return
    }

    const result = await signupRequest(
      {
        email: email.value,
        password: password.value,
        first_name: firstName.value,
        last_name: lastName.value,
        middle_name: middleName.value.trim() || undefined,
      },
      csrf.token,
    )
    if (!result.ok) {
      errorMessage.value = result.message
      return
    }

    const redirect = result.data.redirect_to
    if (redirect.startsWith('/')) {
      await router.replace({ path: redirect, query: { registered: '1' } })
    }
    else {
      await router.replace({ name: 'login', query: { registered: '1' } })
    }
  }
  catch (e) {
    const msg = e instanceof Error ? e.message : String(e)
    errorMessage.value = msg.includes('fetch')
      ? 'Could not reach the server. Check that the API is running and VITE_API_BASE_URL matches it (and CORS FRONTEND_URL matches this app).'
      : 'Something went wrong. Please try again.'
  }
  finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.auth-submit {
  transform: translateZ(0);
  transition:
    transform 0.2s ease,
    filter 0.2s ease;
}

.auth-submit:hover {
  transform: scale(1.02);
}

.auth-submit:active {
  transform: scale(0.98);
}

fieldset {
  border: none;
  padding: 0;
  margin: 0;
}

legend {
  padding: 0;
}
</style>
