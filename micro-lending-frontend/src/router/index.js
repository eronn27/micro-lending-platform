import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/login',
            name: 'Login',
            component: () => import('../pages/Login.vue'),
            meta: { requiresGuest: true }
        },
        {
            path: '/dashboard',
            name: 'Dashboard',
            component: () => import('../pages/Dashboard.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/clients/new',
            name: 'ClientForm',
            component: () => import('../pages/clients/ClientForm.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/clients/search',
            name: 'SearchClients',
            component: () => import('../pages/clients/SearchClients.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/clients/:id',
            name: 'ClientDetails',
            component: () => import('../pages/clients/ClientDetails.vue'), // Future feature
            meta: { requiresAuth: true }
        },
        {
            path: '/clients/:id/edit',
            name: 'EditClient',
            component: () => import('../pages/clients/EditClient.vue'), // Future feature
            meta: { requiresAuth: true }
        },
        // NEW: Payment Management Route
        {
            path: '/payments',
            name: 'PaymentManagement',
            component: () => import('../pages/payments/PaymentManagement.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/reports',
            name: 'Reports',
            component: () => import('../pages/Reports.vue'),
            meta: { requiresAuth: true }
        },

        {
            path: '/',
            redirect: '/dashboard'
        }
    ]
})

router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()
    
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        next('/login')
    } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
        next('/dashboard')
    } else {
        next()
    }
})

export default router
