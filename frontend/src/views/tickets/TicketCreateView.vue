<template>
  <!-- VULN-03: Form submits ticket text to API with length-only style validation server-side; XSS if rendered as HTML elsewhere. -->
  <div class="mx-auto max-w-2xl space-y-6">
    <div>
      <RouterLink
        :to="paths.dashboard.tickets"
        class="inline-flex items-center gap-1 text-xs font-semibold text-[var(--brand-green-dark)] hover:underline sm:text-sm"
      >
        ← Back to tickets
      </RouterLink>
    </div>

    <header>
      <h1 class="text-xl font-semibold tracking-tight text-[var(--text-primary)] sm:text-2xl">
        New ticket
      </h1>
      <p class="mt-1 text-sm text-[var(--text-secondary)]">
        Describe your issue so SecWeb support can help.
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

      <div>
        <label
          for="tk-title"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Title</label>
        <input
          id="tk-title"
          v-model.trim="title"
          type="text"
          required
          minlength="3"
          maxlength="180"
          class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm text-[var(--text-primary)] shadow-sm outline-none ring-[var(--brand-green)] focus:border-transparent focus:ring-2"
          placeholder="Short summary of the problem"
        >
      </div>

      <div>
        <label
          for="tk-category"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Category</label>
        <input
          id="tk-category"
          v-model.trim="category"
          type="text"
          required
          minlength="2"
          maxlength="80"
          class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm text-[var(--text-primary)] shadow-sm outline-none ring-[var(--brand-green)] focus:border-transparent focus:ring-2"
          placeholder="e.g. Billing, Access, Bug"
        >
      </div>

      <div>
        <label
          for="tk-desc"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Description</label>
        <textarea
          id="tk-desc"
          v-model="description"
          required
          minlength="5"
          maxlength="5000"
          rows="8"
          class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm text-[var(--text-primary)] shadow-sm outline-none ring-[var(--brand-green)] focus:border-transparent focus:ring-2"
          placeholder="Steps to reproduce, what you expected, what happened instead…"
        />
        <p class="mt-1 text-xs text-[var(--text-muted)]">
          {{ description.length }} / 5000 characters (minimum 5)
        </p>
      </div>

      <div class="flex flex-wrap gap-3 pt-1">
        <button
          type="submit"
          :disabled="submitting"
          class="rounded-full bg-[var(--brand-green)] px-6 py-3 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ submitting ? 'Creating…' : 'Create ticket' }}
        </button>
        <RouterLink
          :to="paths.dashboard.tickets"
          class="inline-flex items-center rounded-full border border-[var(--border-subtle)] px-6 py-3 text-sm font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)]"
        >
          Cancel
        </RouterLink>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
// VULN-03: createTicket() sends raw title/description/category — server persists with weak sanitization.
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createTicket } from '@/api/tickets'
import { paths } from '@/constants/routes'
import { getSessionCsrfToken } from '@/stores/auth-session'

const router = useRouter()

const title = ref('')
const category = ref('')
const description = ref('')
const submitting = ref(false)
const errorMessage = ref('')

async function onSubmit() {
  errorMessage.value = ''
  if (title.value.length < 3) {
    errorMessage.value = 'Title must be at least 3 characters.'
    return
  }
  if (category.value.length < 2) {
    errorMessage.value = 'Category must be at least 2 characters.'
    return
  }
  if (description.value.trim().length < 5) {
    errorMessage.value = 'Description must be at least 5 characters.'
    return
  }

  const csrf = getSessionCsrfToken()
  if (!csrf) {
    errorMessage.value = 'Your session expired. Sign in again.'
    return
  }

  submitting.value = true
  try {
    const result = await createTicket(
      {
        title: title.value,
        category: category.value,
        description: description.value,
      },
      csrf,
    )
    if (!result.ok) {
      errorMessage.value = result.message
      return
    }
    await router.replace(paths.dashboard.ticketDetail(String(result.data.ticket_id)))
  }
  catch (e) {
    const msg = e instanceof Error ? e.message : String(e)
    errorMessage.value = msg.includes('fetch')
      ? 'Could not reach the server. Check your connection and API URL.'
      : 'Something went wrong. Please try again.'
  }
  finally {
    submitting.value = false
  }
}
</script>
