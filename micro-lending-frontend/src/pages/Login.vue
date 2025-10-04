<template>
    <div class="login-container">
        <form @submit.prevent="handleLogin" class="login-form">
            <h2>Micro-Lending Platform</h2>
            <div class="form-group">
                <label for="username">Username</label>
                <input 
                    type="text" 
                    id="username"
                    v-model="credentials.username"
                    required
                    :disabled="loading"
                >
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input 
                    type="password" 
                    id="password"
                    v-model="credentials.password"
                    required
                    :disabled="loading"
                >
            </div>
            <button type="submit" :disabled="loading">
                {{ loading ? 'Logging in...' : 'Login' }}
            </button>
            <div v-if="error" class="error-message">{{ error }}</div>
        </form>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const credentials = ref({
    username: '',
    password: ''
})
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
    loading.value = true
    error.value = ''
    
    try {
        await authStore.login(credentials.value)
        router.push('/dashboard')
    } catch (err) {
        error.value = err.response?.data?.error || 'Login failed'
    } finally {
        loading.value = false
    }
}
</script>

<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background-color: #f5f5f5;
}

.login-form {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    width: 100%;
    max-width: 400px;
}

.form-group {
    margin-bottom: 1rem;
}

label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
}

input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    box-sizing: border-box;
}

button {
    width: 100%;
    padding: 0.75rem;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}

button:disabled {
    background-color: #6c757d;
    cursor: not-allowed;
}

.error-message {
    color: #dc3545;
    margin-top: 1rem;
    text-align: center;
}
</style>
