<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <h1>Admin Dashboard</h1>
      <button @click="handleLogout" class="logout-btn">Logout</button>
    </header>
    
    <main class="dashboard-main">
      <div class="action-grid">
        <div class="action-card" @click="navigateToClientForm">
          <div class="action-icon">👤</div>
          <h3>Create New Client</h3>
          <p>Add a new client to the system</p>
        </div>
        
        <div class="action-card" @click="navigateToPaymentManagement">
          <div class="action-icon">💳</div>
          <h3>Payment Management</h3>
          <p>Manage client payments and collections</p>
        </div>
        
        <div class="action-card" @click="navigateToClientSearch">
          <div class="action-icon">🔍</div>
          <h3>Search Clients</h3>
          <p>Find and manage existing clients</p>
        </div>
        
        <div class="action-card" @click="navigateToReports">
          <div class="action-icon">📊</div>
          <h3>Reports & Analytics</h3>
          <p>View business insights and reports</p>
        </div>
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

const navigateToPaymentManagement = () => {
  router.push('/payments')
}

const navigateToClientSearch = () => {
  // Placeholder for future implementation
  alert('Client search feature coming soon!')
}

const navigateToReports = () => {
  // Placeholder for future implementation
  alert('Reports feature coming soon!')
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
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
  margin-bottom: 3rem;
  padding-bottom: 1.5rem;
  border-bottom: 2px solid #e5e7eb;
}

.dashboard-header h1 {
  margin: 0;
  color: #1f2937;
  font-size: 2rem;
  font-weight: 700;
}

.logout-btn {
  padding: 0.75rem 1.5rem;
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: background-color 0.2s;
}

.logout-btn:hover {
  background-color: #dc2626;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
}

.action-card {
  background: white;
  padding: 2.5rem 2rem;
  border-radius: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #f3f4f6;
  text-align: center;
}

.action-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  border-color: #3b82f6;
}

.action-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.action-card h3 {
  margin: 0 0 0.75rem 0;
  color: #1f2937;
  font-size: 1.25rem;
  font-weight: 600;
}

.action-card p {
  margin: 0;
  color: #6b7280;
  font-size: 0.875rem;
  line-height: 1.5;
}

/* Responsive Design */
@media (max-width: 768px) {
  .dashboard {
    padding: 1rem;
  }
  
  .dashboard-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .action-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
  
  .action-card {
    padding: 2rem 1.5rem;
  }
}
</style>
