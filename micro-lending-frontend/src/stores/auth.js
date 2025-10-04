import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../services/api'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(null)
    const token = ref(localStorage.getItem('token'))
    
    const isAuthenticated = computed(() => !!token.value)
    const isAdmin = computed(() => user.value?.is_admin || false)
    
    function setToken(newToken) {
        token.value = newToken
        localStorage.setItem('token', newToken)
        // Configure axios to use the token
        api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
    }
    
    function setUser(userData) {
        user.value = userData
    }
    
    async function login(credentials) {
        try {
            const response = await api.post('/auth/login', credentials)
            setToken(response.data.token)
            setUser(response.data.user)
            return response
        } catch (error) {
            logout()
            throw error
        }
    }
    
    async function logout() {
        user.value = null
        token.value = null
        localStorage.removeItem('token')
        delete api.defaults.headers.common['Authorization']
    }
    
    async function fetchCurrentUser() {
        try {
            const response = await api.get('/auth/me')
            setUser(response.data)
            return response
        } catch (error) {
            logout()
            throw error
        }
    }
    
    // Initialize auth state if token exists
    if (token.value) {
        api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
        fetchCurrentUser().catch(() => logout())
    }
    
    return {
        user,
        token,
        isAuthenticated,
        isAdmin,
        login,
        logout,
        fetchCurrentUser
    }
})
