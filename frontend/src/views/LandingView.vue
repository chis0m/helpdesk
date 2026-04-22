<template>
  <div class="relative min-h-screen overflow-hidden bg-[var(--surface-page)] text-[var(--text-primary)]">
    <!-- Ambient background -->
    <div
      class="pointer-events-none absolute inset-0 overflow-hidden"
      aria-hidden="true"
    >
      <div
        class="landing-blob absolute -right-32 -top-40 h-[520px] w-[520px] rounded-full bg-[var(--surface-mint)] opacity-70 blur-3xl"
      />
      <div
        class="landing-blob absolute -bottom-48 left-1/4 h-[420px] w-[420px] rounded-full bg-[var(--surface-muted)] opacity-90 blur-3xl"
      />
      <div
        class="landing-grid absolute inset-0 opacity-[0.35]"
      />
    </div>

    <header
      class="landing-nav sticky top-0 z-50 border-b border-[var(--border-subtle)]/80 bg-white/75 px-6 py-4 backdrop-blur-xl lg:px-10"
    >
      <div class="mx-auto flex max-w-6xl items-center justify-between">
        <RouterLink
          to="/"
          class="text-xl font-semibold tracking-tight text-[var(--brand-green-dark)] transition hover:opacity-90"
        >
          {{ appName }}
        </RouterLink>
        <div class="flex items-center gap-2 sm:gap-3">
          <RouterLink
            to="/login"
            class="landing-product-badge hidden rounded-full bg-[var(--brand-green)] px-3 py-1.5 text-xs font-semibold text-[var(--text-on-green)] shadow-sm sm:inline-flex"
          >
            Product support
          </RouterLink>
          <RouterLink
            to="/login"
            class="landing-auth-nav-link text-sm font-semibold text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
          >
            Sign in
          </RouterLink>
          <RouterLink
            to="/signup"
            class="landing-auth-nav-cta rounded-full bg-[var(--surface-mint)] px-4 py-2 text-sm font-semibold text-[var(--text-primary)] hover:bg-[var(--surface-mint-hover)]"
          >
            Create account
          </RouterLink>
        </div>
      </div>
    </header>

    <main class="relative mx-auto max-w-6xl px-6 pb-28 pt-10 lg:px-10 lg:pt-14">
      <!-- Hero -->
      <div
        class="grid items-center gap-12 lg:grid-cols-2 lg:gap-16"
      >
        <div class="land-in max-w-xl">
          <p
            class="inline-flex rounded-full border border-[var(--border-subtle)] bg-white/90 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)] shadow-sm"
          >
            {{ brandShort }} product support
          </p>
          <h1
            class="mt-5 text-4xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:text-[2.75rem] lg:leading-[1.12]"
          >
            Get help with
            <span class="text-[var(--brand-green-dark)]">{{ brandShort }}</span>
            products — report issues and track fixes.
          </h1>
          <p class="mt-5 text-lg leading-relaxed text-[var(--text-secondary)]">
            Create an account, sign in, and open a ticket when something breaks — the same idea as vendor portals like Sisense or other B2B product support sites.
          </p>
          <div class="mt-10 flex flex-wrap items-center gap-3">
            <RouterLink
              to="/signup"
              class="landing-cta-primary landing-auth-hero-pill inline-flex items-center justify-center rounded-full bg-[var(--brand-green)] px-6 py-3 text-sm font-semibold text-[var(--text-on-green)] shadow-md hover:brightness-95"
            >
              Get started
            </RouterLink>
            <RouterLink
              to="/login"
              class="landing-cta-secondary landing-auth-hero-pill inline-flex items-center justify-center rounded-full bg-[var(--surface-mint)] px-6 py-3 text-sm font-semibold text-[var(--text-primary)] hover:bg-[var(--surface-mint-hover)]"
            >
              Sign in
            </RouterLink>
            <RouterLink
              to="/dashboard"
              class="landing-cta-secondary inline-flex items-center justify-center rounded-full border border-[var(--border-strong)] bg-white/90 px-6 py-3 text-sm font-semibold text-[var(--text-primary)] shadow-sm transition hover:bg-[var(--surface-hover)]"
            >
              Support home
            </RouterLink>
          </div>
        </div>

        <!-- Product preview -->
        <div
          class="land-in land-delay-2 lg:justify-self-end"
        >
          <div
            class="w-full max-w-[440px] rounded-2xl border border-[var(--border-subtle)] bg-white shadow-[0_24px_80px_-20px_rgba(0,0,0,0.12)]"
          >
            <div
              class="flex items-center gap-2 border-b border-[var(--border-subtle)] bg-[var(--surface-muted)]/80 px-4 py-3"
            >
              <span class="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
              <span class="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
              <span class="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
              <div
                class="ml-3 flex-1 rounded-lg border border-[var(--border-subtle)] bg-white px-3 py-1.5 text-center text-xs text-[var(--text-muted)]"
              >
                {{ supportHostLabel }} · My requests
              </div>
            </div>
            <div class="space-y-4 p-5">
              <div class="flex items-end justify-between">
                <div>
                  <p class="text-xs font-medium text-[var(--text-muted)]">
                    API health
                  </p>
                  <p class="mt-0.5 text-2xl font-semibold tracking-tight tabular-nums text-[var(--text-primary)]">
                    <span v-if="healthLoading">…</span>
                    <span v-else-if="healthError">{{ healthError }}</span>
                    <span v-else>{{ healthStatus }}</span>
                  </p>
                </div>
                <span
                  class="rounded-full px-3 py-1 text-xs font-semibold"
                  :class="healthOk ? 'bg-[var(--brand-green)] text-[var(--text-on-green)]' : 'bg-amber-100 text-amber-900'"
                >{{ healthOk ? 'Live' : 'Check' }}</span>
              </div>
              <p class="text-sm leading-relaxed text-[var(--text-secondary)]">
                {{ healthMessage || 'Sign in after creating an account to see your real tickets and replies from support.' }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Bento -->
      <section class="mt-14">
        <h2 class="land-in land-delay-4 text-lg font-semibold text-[var(--text-primary)]">
          How {{ appName }} fits your workflow
        </h2>
        <p class="land-in land-delay-4 mt-1 text-sm text-[var(--text-secondary)]">
          Simple steps from account to ticket — staff triage and resolve on their side.
        </p>
        <div
          class="mt-8 grid gap-4 md:grid-cols-3"
        >
          <article
            class="land-in land-delay-5 landing-card group rounded-2xl border border-[var(--border-subtle)] bg-white p-7 shadow-sm transition duration-300 md:col-span-2"
          >
            <p class="text-base font-semibold text-[var(--text-primary)]">
              {{ featured.title }}
            </p>
            <p class="mt-3 max-w-xl text-sm leading-relaxed text-[var(--text-secondary)]">
              {{ featured.body }}
            </p>
            <div class="mt-6 flex flex-wrap gap-2">
              <span
                class="rounded-full bg-[var(--surface-mint)] px-3 py-1 text-xs font-medium text-[var(--brand-green-dark)]"
              >Your tickets</span>
              <span
                class="rounded-full bg-[var(--surface-muted)] px-3 py-1 text-xs font-medium text-[var(--text-secondary)]"
              >Staff triage</span>
              <span
                class="rounded-full bg-[var(--surface-muted)] px-3 py-1 text-xs font-medium text-[var(--text-secondary)]"
              >Product issues</span>
            </div>
          </article>
          <article
            class="land-in land-delay-6 landing-card group rounded-2xl border border-[var(--border-subtle)] bg-[var(--surface-muted)]/60 p-7 shadow-sm transition duration-300"
          >
            <p class="text-base font-semibold text-[var(--text-primary)]">
              {{ sideCard.title }}
            </p>
            <p class="mt-3 text-sm leading-relaxed text-[var(--text-secondary)]">
              {{ sideCard.body }}
            </p>
          </article>
          <article
            v-for="(card, i) in bottomCards"
            :key="card.title"
            class="landing-card group rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-sm transition duration-300"
            :class="landDelayClass(7 + i)"
          >
            <p class="text-sm font-semibold text-[var(--text-primary)]">
              {{ card.title }}
            </p>
            <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
              {{ card.body }}
            </p>
          </article>
        </div>
      </section>

      <!-- Testimonials -->
      <section
        class="land-in land-delay-10 mt-14"
        aria-labelledby="landing-testimonials-heading"
      >
        <h2
          id="landing-testimonials-heading"
          class="text-lg font-semibold text-[var(--text-primary)]"
        >
          Testimonials
        </h2>
        <p class="mt-1 text-sm text-[var(--text-secondary)]">
          What customers say about {{ brandShort }} support.
        </p>
        <div class="mt-6 grid gap-4 md:grid-cols-2">
          <figure
            v-for="t in testimonials"
            :key="t.name"
            class="landing-card rounded-2xl border border-[var(--border-subtle)] bg-white/90 p-6 shadow-sm"
          >
            <blockquote class="border-l-[3px] border-[var(--brand-green)] pl-4 text-sm leading-relaxed text-[var(--text-secondary)]">
              {{ t.quote }}
            </blockquote>
            <figcaption class="mt-4 flex items-center gap-3 border-t border-[var(--border-subtle)] pt-4">
              <div
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-[var(--surface-mint)] text-xs font-bold text-[var(--brand-green-dark)]"
                aria-hidden="true"
              >
                {{ t.initials }}
              </div>
              <div>
                <p class="text-sm font-semibold text-[var(--text-primary)]">
                  {{ t.name }}
                </p>
                <p class="text-xs text-[var(--text-muted)]">
                  {{ t.role }}
                </p>
                <p class="mt-0.5">
                  <a
                    :href="`mailto:${t.email}`"
                    class="text-xs text-[var(--text-secondary)] underline decoration-[var(--border-strong)] underline-offset-2 transition hover:text-[var(--brand-green-dark)]"
                  >{{ t.email }}</a>
                </p>
              </div>
            </figcaption>
          </figure>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { fetchHealth } from '@/api/health'
import { brandShortFromAppName, loadAppDetail } from '@/stores/app-detail'

const appName = ref('SecWeb HelpDesk')
const brandShort = ref('SecWeb')

const supportHostLabel = computed(() => {
  const t = brandShort.value.trim().toLowerCase()
  return `support.${t || 'helpdesk'}`
})

const healthLoading = ref(true)
const healthOk = ref(false)
const healthStatus = ref('')
const healthMessage = ref('')
const healthError = ref('')

onMounted(async () => {
  const [healthRes] = await Promise.all([
    fetchHealth(),
    loadAppDetail().then((d) => {
      appName.value = d.app_name
      brandShort.value = brandShortFromAppName(d.app_name)
    }),
  ])
  healthLoading.value = false
  if (healthRes.ok) {
    healthOk.value = true
    healthStatus.value = healthRes.status
    healthMessage.value = healthRes.message
  }
  else {
    healthError.value = healthRes.message
  }
})

const featured = computed(() => ({
  title: 'One place for your support requests',
  body: `You see the tickets you opened and their status; ${brandShort.value} staff see the queue on their side. Same calm layout whether you’re reporting a bug or asking for help.`,
}))

const sideCard = {
  title: 'Open a ticket in minutes',
  body: 'Describe the issue, pick a category if needed, and submit. No clutter — just what product support portals are meant to do.',
}

const bottomCards = computed(() => [
  {
    title: 'Account required',
    body: 'Every reporter signs in — so we know who to reply to and can keep your history in one place.',
  },
  {
    title: 'Staff triage',
    body: 'Internal roles assign and resolve tickets; you stay informed as status changes.',
  },
  {
    title: `Built for ${brandShort.value} CA`,
    body: 'Coursework baseline for Secure Web Development — authentication and tickets use the running API.',
  },
])

const testimonials = computed(() => [
  {
    name: 'Mark Anthony',
    role: 'Product Manager',
    email: 'mark.anthony@company-a.com',
    initials: 'MA',
    quote: `${brandShort.value} customer service has been outstanding — quick responses, clear updates, and they actually follow through until the issue is sorted.`,
  },
  {
    name: 'Jane Doe',
    role: 'IT Lead',
    email: 'jane.doe@company-b.com',
    initials: 'JD',
    quote: `Our team relies on timely support when integrations act up. ${brandShort.value}'s helpdesk keeps us moving; escalation feels professional every time.`,
  },
])

function landDelayClass(n: number) {
  return `land-in land-delay-${n}`
}
</script>

<style scoped>
.landing-grid {
  background-image: radial-gradient(circle at center, var(--border-subtle) 1px, transparent 1px);
  background-size: 24px 24px;
}

.landing-blob {
  animation: blob-drift 18s ease-in-out infinite alternate;
}

.landing-blob:nth-child(2) {
  animation-duration: 22s;
  animation-delay: -4s;
}

@keyframes blob-drift {
  0% {
    transform: translate(0, 0) scale(1);
  }
  100% {
    transform: translate(24px, 16px) scale(1.05);
  }
}

.land-in {
  animation: land-in 0.75s cubic-bezier(0.22, 1, 0.36, 1) both;
}

.land-delay-2 {
  animation-delay: 0.1s;
}
.land-delay-3 {
  animation-delay: 0.18s;
}
.land-delay-4 {
  animation-delay: 0.24s;
}
.land-delay-5 {
  animation-delay: 0.3s;
}
.land-delay-6 {
  animation-delay: 0.36s;
}
.land-delay-7 {
  animation-delay: 0.42s;
}
.land-delay-8 {
  animation-delay: 0.48s;
}
.land-delay-9 {
  animation-delay: 0.54s;
}
.land-delay-10 {
  animation-delay: 0.6s;
}

@keyframes land-in {
  from {
    opacity: 0;
    transform: translateY(18px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.landing-cta-primary,
.landing-cta-secondary {
  transform: translateZ(0);
  transition:
    transform 0.22s cubic-bezier(0.22, 1, 0.36, 1),
    box-shadow 0.22s ease,
    filter 0.22s ease,
    background-color 0.22s ease;
}

.landing-cta-primary:hover,
.landing-cta-secondary:hover {
  transform: scale(1.03) translateY(-2px);
  box-shadow: 0 10px 28px -12px rgba(0, 0, 0, 0.22);
}

.landing-cta-primary:active,
.landing-cta-secondary:active {
  transform: scale(0.97) translateY(0);
  transition-duration: 0.08s;
}

/* Header: Sign in + Create account — hover lift, click press */
.landing-auth-nav-link {
  display: inline-block;
  transform: translateZ(0);
  transition:
    transform 0.2s cubic-bezier(0.22, 1, 0.36, 1),
    color 0.2s ease;
}

.landing-auth-nav-link:hover {
  transform: translateY(-2px);
}

.landing-auth-nav-link:active {
  transform: translateY(0) scale(0.96);
  transition-duration: 0.08s;
}

.landing-auth-nav-cta {
  display: inline-block;
  transform: translateZ(0);
  transition:
    transform 0.2s cubic-bezier(0.22, 1, 0.36, 1),
    box-shadow 0.2s ease,
    background-color 0.2s ease;
}

.landing-auth-nav-cta:hover {
  transform: translateY(-2px) scale(1.03);
  box-shadow: 0 8px 20px -10px rgba(0, 0, 0, 0.2);
}

.landing-auth-nav-cta:active {
  transform: scale(0.96);
  transition-duration: 0.08s;
}

/* Product support badge → sign in */
.landing-product-badge {
  transform: translateZ(0);
  transition:
    transform 0.2s cubic-bezier(0.22, 1, 0.36, 1),
    box-shadow 0.2s ease,
    filter 0.2s ease;
}

.landing-product-badge:hover {
  transform: translateY(-2px) scale(1.04);
  box-shadow: 0 6px 16px -8px rgba(0, 0, 0, 0.25);
  filter: brightness(1.05);
}

.landing-product-badge:active {
  transform: scale(0.96);
  transition-duration: 0.08s;
}

/* Hero Get started / Sign in — slightly snappier motion (classes stack with .landing-cta-*) */
.landing-auth-hero-pill:hover {
  transform: scale(1.04) translateY(-3px);
}

.landing-auth-hero-pill:active {
  transform: scale(0.96) translateY(0);
}

.landing-card {
  transform: translateZ(0);
}

.landing-card:hover {
  border-color: color-mix(in srgb, var(--brand-green) 28%, var(--border-subtle));
  box-shadow:
    0 12px 40px -24px rgba(0, 0, 0, 0.18),
    0 0 0 1px color-mix(in srgb, var(--brand-green) 12%, transparent);
}
</style>
