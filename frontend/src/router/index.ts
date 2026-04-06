/** App routes — landing, auth, and dashboard. */
import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: () => import('@/views/LandingView.vue'),
    },
    {
      path: '/login',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'login',
          component: () => import('@/views/LoginView.vue'),
        },
      ],
    },
    {
      path: '/signup',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'signup',
          component: () => import('@/views/SignupView.vue'),
        },
      ],
    },
    {
      path: '/dashboard',
      component: AppLayout,
      children: [
        {
          path: 'tickets/:id',
          name: 'ticket-detail',
          meta: { title: 'Ticket detail' },
          component: () => import('@/views/tickets/TicketDetailView.vue'),
        },
        {
          path: 'tickets',
          name: 'dashboard-tickets',
          meta: { title: 'Tickets & search' },
          component: () => import('@/views/tickets/TicketsListView.vue'),
        },
        {
          path: 'profile',
          name: 'dashboard-profile',
          meta: { title: 'Your profile' },
          component: () => import('@/views/account/ProfileView.vue'),
        },
        {
          path: 'admin/users',
          name: 'admin-users',
          meta: { title: 'Users' },
          component: () => import('@/views/admin/AdminUsersListView.vue'),
        },
        {
          path: 'admin/staff/new',
          name: 'admin-staff-new',
          meta: { title: 'Create staff' },
          component: () => import('@/views/admin/AdminStaffView.vue'),
        },
        {
          path: 'admin/staff',
          name: 'admin-staff',
          meta: { title: 'Staff' },
          component: () => import('@/views/admin/AdminStaffListView.vue'),
        },
        {
          path: '',
          name: 'dashboard',
          meta: { title: 'Support home' },
          component: () => import('@/views/DashboardView.vue'),
        },
      ],
    },
  ],
})

export default router
