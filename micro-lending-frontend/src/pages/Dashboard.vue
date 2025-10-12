<template>
    <div class="dashboard">
        <!-- Header -->
        <header class="dashboard-header">
            <div class="header-content">
                <div class="logo-title-container">
                    <img src="/Screenshot 2025-01-14 131433.png" alt="Company Logo" class="company-logo" />
                    <h1 class="dashboard-title">Admin Dashboard</h1>
                </div>
            </div>
            <button @click="handleLogout" class="logout-btn">
                <span class="logout-icon"> </span>
                Logout
            </button>
        </header>
        
        <!-- Main Content -->
        <main class="dashboard-main">
            <div class="welcome-section">
                <h2>Welcome Back</h2>
                <p>What would you like to do today?</p>
            </div>

            <div class="action-grid">
                <!-- Client Management -->
                <div class="action-card" @click="navigateToClientForm">
                    <div class="card-header">
                        <div class="card-icon">👤</div>
                        <h3>Create New Client</h3>
                    </div>
                    <p class="card-description">Add a new client to the system and set up their loan details</p>
                    <div class="card-arrow">→</div>
                </div>
                
                <div class="action-card" @click="navigateToSearch">
                    <div class="card-header">
                        <div class="card-icon">🔍</div>
                        <h3>Search Clients</h3>
                    </div>
                    <p class="card-description">Find and manage existing client records and information</p>
                    <div class="card-arrow">→</div>
                </div>

                <!-- Payment Management -->
                <div class="action-card" @click="navigateToPayments">
                    <div class="card-header">
                        <div class="card-icon">💰</div>
                        <h3>Payment Management</h3>
                    </div>
                    <p class="card-description">Process weekly payments and track client payment progress</p>
                    <div class="card-arrow">→</div>
                </div>
                
                <!-- Reports & Analytics -->
                <div class="action-card" @click="navigateToReports">
                    <div class="card-header">
                        <div class="card-icon">📊</div>
                        <h3>Reports & Analytics</h3>
                    </div>
                    <p class="card-description">View business insights, performance metrics and analytics</p>
                    <div class="card-arrow">→</div>
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

const navigateToSearch = () => {
    router.push('/clients/search')
}

const handleLogout = () => {
    authStore.logout()
    router.push('/login')
}

const navigateToPayments = () => {
  router.push('/payments')
}

const navigateToReports = () => {
  router.push('/reports')
}
</script>

<style scoped>
.dashboard {
    min-height: 100vh;
    background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
    padding: 2rem;
}

.dashboard-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 3rem;
}

.header-content {
    color: #2c3e50;
}

.dashboard-title {
    font-size: 2.5rem;
    font-weight: 700;
    margin: 0 0 0.5rem 0;
    color: #2c3e50;
}

.dashboard-subtitle {
    color: #7f8c8d;
    font-size: 1.1rem;
    margin: 0;
}

.logout-btn {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.75rem 1.8rem;
    background: #e74c3c;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 500;
    transition: all 0.3s ease;
}

.logout-btn:hover {
    background: #c0392b;
    transform: translateY(-2px);
}

.logout-icon {
    font-size: 1.2rem;
}

.dashboard-main {
    max-width: 1200px;
    margin: 0 auto;
}

.welcome-section {
    text-align: center;
    margin-bottom: 3rem;
    color: #2c3e50;
}

.welcome-section h2 {
    font-size: 2rem;
    margin: 0 0 0.5rem 0;
    font-weight: 600;
}

.welcome-section p {
    font-size: 1.1rem;
    color: #7f8c8d;
    margin: 0;
}

.action-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.5rem;
    margin-bottom: 3rem;
}

.action-card {
    background: white;
    padding: 2rem;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    border: 1px solid #e1e8ed;
    position: relative;
    overflow: hidden;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
}

.action-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: #3498db;
}

.action-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.card-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1rem;
}

.card-icon {
    font-size: 2rem;
    width: 60px;
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #3498db;
    color: white;
    border-radius: 10px;
}

.action-card h3 {
    margin: 0;
    color: #2c3e50;
    font-size: 1.25rem;
    font-weight: 600;
}

.card-description {
    color: #7f8c8d;
    line-height: 1.6;
    margin: 0 0 1.5rem 0;
}

.card-arrow {
    font-size: 1.5rem;
    color: #3498db;
    font-weight: bold;
    transition: transform 0.3s ease;
}

.action-card:hover .card-arrow {
    transform: translateX(5px);
}

.stats-section {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1.5rem;
}

.stat-card {
    background: white;
    padding: 1.5rem;
    border-radius: 10px;
    border: 1px solid #e1e8ed;
    color: #2c3e50;
    display: flex;
    align-items: center;
    gap: 1rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.stat-icon {
    font-size: 2rem;
}

.stat-content h4 {
    margin: 0 0 0.25rem 0;
    font-size: 0.9rem;
    font-weight: 500;
    color: #7f8c8d;
}

.stat-number {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 700;
}

.dashboard-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 3rem;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 1rem;
}

.logo-container {
    display: flex;
    align-items: center;
    gap: 1rem;
}

.company-logo {
    height: 50px;
    width: 50px;
    object-fit: contain;
    padding: 8px;
    background: white;
    border-radius: 10px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    transition: all 0.3s ease;
}

.company-logo:hover {
    box-shadow: 0 6px 20px rgba(52, 152, 219, 0.3);
}

.brand-text {
    display: flex;
    flex-direction: column;
}

.company-name {
    font-size: 1.75rem;
    font-weight: 700;
    background: linear-gradient(45deg, #2c3e50, #3498db);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin: 0;
}


/* Responsive Design */
@media (max-width: 768px) {
    .dashboard {
        padding: 1rem;
    }
    
    .dashboard-header {
        flex-direction: column;
        gap: 1rem;
        align-items: stretch;
    }
    
    .dashboard-title {
        font-size: 2rem;
    }
    
    .action-grid {
        grid-template-columns: 1fr;
    }
    
    .stats-section {
        grid-template-columns: 1fr;
    }
    
    .card-header {
        flex-direction: column;
        text-align: center;
        gap: 0.5rem;
    }
}

@media (max-width: 480px) {
    .dashboard-title {
        font-size: 1.75rem;
    }
    
    .action-card {
        padding: 1.5rem;
    }
}
</style>