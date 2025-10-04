<template>
    <div class="dashboard">
        <header class="dashboard-header">
            <h1>Admin Dashboard</h1>
            <button @click="handleLogout" class="logout-btn">Logout</button>
        </header>
        
        <main class="dashboard-main">
            <div class="action-grid">
                <div class="action-card" @click="navigateToClientForm">
                    <h3>Create New Client</h3>
                    <p>Add a new client to the system</p>
                </div>
                
                <div class="action-card" @click="navigateToSearch">
                    <h3>Search Clients</h3>
                    <p>Find and manage existing clients</p>
                </div>
            </div>
        <!-- NEW: Payment Management Card -->
        <div class="action-card" @click="navigateToPayments">
          <div class="card-icon">💰</div>
          <h3>Payment Management</h3>
          <p>Process weekly payments and track progress</p>
        </div>
        
        <div class="action-card" @click="navigateToReports">
          <div class="card-icon">📊</div>
          <h3>Reports & Analytics</h3>
          <p>View business insights and performance</p>
        </div>
        </main>
    </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const navigateToClientForm = () => {
    router.push('/clients/new')
}

const navigateToSearch = () => {
    router.push('/clients/search')
}

const handleLogout = () => {
    authStore.logout()
    router.push('/login')
}

// NEW: Navigation to Payment Management
const navigateToPayments = () => {
  router.push('/payments')
}

const navigateToReports = () => {
  alert('Reports feature coming soon!')
}
</script>

<style scoped>
.dashboard {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
}

.dashboard-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #ddd;
}

.logout-btn {
    padding: 0.5rem 1rem;
    background-color: #dc3545;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}

.action-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
}

.action-card {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s;
}

.action-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 20px rgba(0,0,0,0.15);
}

.action-card h3 {
    margin-bottom: 0.5rem;
    color: #007bff;
}
</style>
