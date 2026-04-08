<template>
  <section
    class="rounded-2xl border border-[var(--border-subtle)] bg-white p-5 shadow-[var(--shadow-card)] ring-1 ring-black/[0.03]"
    aria-labelledby="ticket-workflow-heading"
  >
    <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
      <div>
        <p class="text-[10px] font-bold uppercase tracking-[0.14em] text-[var(--text-muted)]">
          Workflow
        </p>
        <h2
          id="ticket-workflow-heading"
          class="mt-1 text-lg font-semibold text-[var(--text-primary)]"
        >
          Ticket status
        </h2>
        <p class="mt-1 max-w-xl text-sm leading-snug text-[var(--text-secondary)]">
          Pick the next step, then <strong class="font-semibold text-[var(--text-primary)]">Save status</strong>. Cancel restores the last saved value.
        </p>
      </div>
      <div class="shrink-0 sm:text-right">
        <p class="text-[10px] font-medium uppercase tracking-wide text-[var(--text-muted)]">
          Saved on ticket
        </p>
        <div class="mt-1">
          <TicketStatusBadge
            :status="status"
            size="sm"
          />
        </div>
        <p
          v-if="dirty"
          class="mt-2 rounded-md border border-amber-200 bg-amber-50 px-2 py-1.5 text-[11px] font-medium text-amber-950 sm:text-xs"
          role="status"
        >
          Unsaved: {{ pendingLabel }}
        </p>
      </div>
    </div>

    <div
      class="mt-3 grid gap-2 sm:grid-cols-2"
      role="group"
      aria-label="Choose ticket status"
    >
      <button
        v-for="def in TICKET_STATUS_DEFINITIONS"
        :key="def.value"
        type="button"
        class="rounded-lg border-2 px-3 py-2.5 text-left transition focus:outline-none focus-visible:ring-2 focus-visible:ring-[var(--brand-green)] focus-visible:ring-offset-1"
        :class="selected === def.value
          ? 'border-[var(--brand-green)] bg-[var(--surface-mint)]/50 shadow-sm'
          : 'border-[var(--border-subtle)] bg-[var(--surface-main)] hover:border-[var(--border-strong)] hover:bg-[var(--surface-hover)]'"
        :aria-pressed="selected === def.value"
        @click="selected = def.value"
      >
        <span class="block text-sm font-semibold text-[var(--text-primary)]">{{ def.title }}</span>
        <span class="mt-0.5 block text-xs leading-snug text-[var(--text-secondary)]">{{ def.description }}</span>
      </button>
    </div>

    <div
      v-if="errorMessage"
      class="mt-3 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-900"
      role="alert"
    >
      {{ errorMessage }}
    </div>

    <div class="mt-3 flex flex-wrap items-center gap-2 border-t border-[var(--border-subtle)] pt-3">
      <button
        type="button"
        class="rounded-full bg-[var(--brand-green)] px-4 py-2 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-45"
        :disabled="!dirty || saving"
        @click="save"
      >
        {{ saving ? 'Saving…' : 'Save status' }}
      </button>
      <button
        type="button"
        class="rounded-full border border-[var(--border-strong)] bg-white px-4 py-2 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--surface-hover)] disabled:cursor-not-allowed disabled:opacity-45"
        :disabled="!dirty || saving"
        @click="cancel"
      >
        Cancel
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { TicketStatus } from '@/types/ticket'
import TicketStatusBadge from '@/components/tickets/TicketStatusBadge.vue'
import { TICKET_STATUS_DEFINITIONS } from '@/utils/ticket-ui'

const props = withDefaults(
  defineProps<{
    status: TicketStatus
    saving?: boolean
    errorMessage?: string
  }>(),
  {
    saving: false,
    errorMessage: '',
  },
)

const emit = defineEmits<{
  'update:status': [value: TicketStatus]
}>()

const selected = ref<TicketStatus>(props.status)

watch(
  () => props.status,
  (s) => {
    selected.value = s
  },
)

const dirty = computed(() => selected.value !== props.status)

const pendingLabel = computed(() => {
  return TICKET_STATUS_DEFINITIONS.find((d) => d.value === selected.value)?.title ?? selected.value
})

function save() {
  if (!dirty.value) return
  emit('update:status', selected.value)
}

function cancel() {
  selected.value = props.status
}
</script>
