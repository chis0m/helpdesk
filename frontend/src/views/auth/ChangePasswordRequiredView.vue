<!-- Forced password change before dashboard access. -->
<template>
  <div>
    <p
      class="inline-flex rounded-full border border-[var(--border-subtle)] bg-[var(--surface-main)] px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)] lg:hidden"
    >
      Security
    </p>
    <h1 class="mt-4 text-2xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:mt-0 lg:text-3xl">
      Choose a new password
    </h1>
    <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
      Your account requires a password update before you can use the helpdesk.
    </p>

    <form
      class="mt-8 space-y-5"
      @submit.prevent="onSubmit"
    >
      <div
        v-if="errorMessage"
        class="rounded-2xl border border-red-200/80 bg-red-50/90 px-4 py-3 text-sm text-red-900"
        role="alert"
      >
        {{ errorMessage }}
      </div>

      <div>
        <label
          for="req-current-password"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Current password</label>
        <div class="relative">
          <input
            id="req-current-password"
            v-model="currentPassword"
            :type="showCurrentPw ? 'text' : 'password'"
            name="current-password"
            autocomplete="current-password"
            required
            class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white py-3.5 pl-4 pr-14 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
          >
          <button
            type="button"
            class="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg px-2.5 py-1 text-xs font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] hover:text-[var(--text-primary)]"
            :aria-pressed="showCurrentPw"
            :aria-label="showCurrentPw ? 'Hide current password' : 'Show current password'"
            @click="showCurrentPw = !showCurrentPw"
          >
            {{ showCurrentPw ? 'Hide' : 'Show' }}
          </button>
        </div>
      </div>

      <div>
        <label
          for="req-new-password"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >New password</label>
        <div class="relative">
          <input
            id="req-new-password"
            v-model="newPassword"
            :type="showNewPw ? 'text' : 'password'"
            name="new-password"
            autocomplete="new-password"
            required
            minlength="8"
            maxlength="128"
            class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white py-3.5 pl-4 pr-14 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
          >
          <button
            type="button"
            class="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg px-2.5 py-1 text-xs font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] hover:text-[var(--text-primary)]"
            :aria-pressed="showNewPw"
            :aria-label="showNewPw ? 'Hide new password' : 'Show new password'"
            @click="showNewPw = !showNewPw"
          >
            {{ showNewPw ? 'Hide' : 'Show' }}
          </button>
        </div>
      </div>

      <div>
        <label
          for="req-new-password-2"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Confirm new password</label>
        <div class="relative">
          <input
            id="req-new-password-2"
            v-model="newPassword2"
            :type="showNewPw2 ? 'text' : 'password'"
            name="new-password-confirm"
            autocomplete="new-password"
            required
            minlength="8"
            maxlength="128"
            class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white py-3.5 pl-4 pr-14 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
          >
          <button
            type="button"
            class="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg px-2.5 py-1 text-xs font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] hover:text-[var(--text-primary)]"
            :aria-pressed="showNewPw2"
            :aria-label="showNewPw2 ? 'Hide confirm password' : 'Show confirm password'"
            @click="showNewPw2 = !showNewPw2"
          >
            {{ showNewPw2 ? 'Hide' : 'Show' }}
          </button>
        </div>
      </div>

      <button
        type="submit"
        :disabled="submitting"
        class="auth-submit mt-2 w-full rounded-full bg-[var(--brand-green)] py-3.5 text-sm font-semibold text-[var(--text-on-green)] shadow-md transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {{ submitting ? 'Saving…' : 'Update password' }}
      </button>
    </form>

    <p class="mt-8 border-t border-[var(--border-subtle)] pt-8 text-center text-sm text-[var(--text-secondary)]">
      <button
        type="button"
        class="font-semibold text-[var(--brand-green-dark)] transition hover:underline"
        @click="onSignOut"
      >
        Sign out
      </button>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { changePasswordRequest, fetchMe } from '@/api/auth'
import { paths } from '@/constants/routes'
import { getSessionCsrfToken, mergeAuthUserFromMe } from '@/stores/auth-session'
import { performLogout } from '@/utils/perform-logout'

const router = useRouter()

const currentPassword = ref('')
const newPassword = ref('')
const newPassword2 = ref('')
const showCurrentPw = ref(false)
const showNewPw = ref(false)
const showNewPw2 = ref(false)
const submitting = ref(false)
const errorMessage = ref('')

async function onSubmit() {
  errorMessage.value = ''
  if (newPassword.value !== newPassword2.value) {
    errorMessage.value = 'New passwords do not match.'
    return
  }
  if (newPassword.value.length < 8) {
    errorMessage.value = 'New password must be at least 8 characters.'
    return
  }
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    errorMessage.value = 'Your session expired. Sign in again.'
    return
  }

  submitting.value = true
  try {
    const result = await changePasswordRequest(
      {
        current_password: currentPassword.value,
        new_password: newPassword.value,
      },
      csrf,
    )
    if (!result.ok) {
      errorMessage.value = result.message
      return
    }

    const me = await fetchMe()
    if (!me.ok) {
      errorMessage.value = me.message
      return
    }
    if (me.data.must_change_password) {
      errorMessage.value = 'Password was updated but the server still requires a change. Try again or contact support.'
      return
    }
    mergeAuthUserFromMe(me.data)
    await router.replace(paths.dashboard.home)
  }
  catch (e) {
    const msg = e instanceof Error ? e.message : String(e)
    errorMessage.value = msg.includes('fetch')
      ? 'Could not reach the server.'
      : 'Something went wrong. Please try again.'
  }
  finally {
    submitting.value = false
  }
}

function onSignOut() {
  void performLogout(router)
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
</style>
