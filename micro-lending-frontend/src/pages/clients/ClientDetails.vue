<template>
  <div class="client-details">
    <!-- Header -->
    <header class="page-header">
      <div class="header-content">
        <button @click="goBack" class="back-btn">← Back to Search</button>
        <h1>Client Details</h1>
        <div class="header-actions">
          <button @click="editClient" class="edit-btn">Edit Client</button>
        </div>
      </div>
    </header>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading client details...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Error Loading Client</h3>
      <p>{{ error }}</p>
      <button @click="loadClient" class="retry-btn">Retry</button>
    </div>

    <!-- Client Details -->
    <div v-else-if="client" class="client-content">
      <!-- Basic Information Card -->
      <div class="detail-card">
        <div class="card-header">
          <h2>Basic Information</h2>
          <span class="control-number">{{ client.control_number }}</span>
        </div>
        <div class="card-content">
          <div class="info-grid">
            <div class="info-item">
              <label>Name:</label>
              <span class="info-value">{{ client.first_name }} {{ client.middle_name }} {{ client.last_name }}</span>
            </div>
            <div class="info-item">
              <label>Contact:</label>
              <span class="info-value">{{ client.contact_number }}</span>
            </div>
            <div class="info-item">
              <label>Age:</label>
              <span class="info-value">{{ client.age }} years old</span>
            </div>
            <div class="info-item">
              <label>Gender:</label>
              <span class="info-value">{{ client.gender || 'N/A' }}</span>
            </div>
            <div class="info-item">
              <label>Civil Status:</label>
              <span class="info-value">{{ client.civil_status || 'N/A' }}</span>
            </div>
            <div class="info-item">
              <label>Religion:</label>
              <span class="info-value">{{ client.religion || 'N/A' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Address Information -->
      <div class="detail-card">
        <h2>Address Information</h2>
        <div class="card-content">
          <div class="info-item">
            <label>Home Address:</label>
            <span class="info-value address">{{ client.home_address }}</span>
          </div>
          <div class="info-item">
            <label>Years of Residence:</label>
            <span class="info-value">{{ client.years_of_residence || 'N/A' }}</span>
          </div>
          <div class="info-item">
            <label>Facebook:</label>
            <span class="info-value">{{ client.facebook_account || 'N/A' }}</span>
          </div>
        </div>
      </div>

      <!-- Loan Information -->
      <div v-if="client.loans && client.loans.length > 0" class="detail-card">
        <h2>Loan Information</h2>
        <div class="card-content">
          <div class="info-grid">
            <div class="info-item">
              <label>Total Amount:</label>
              <span class="info-value">{{ formatCurrency(client.loans[0].total_amount) }}</span>
            </div>
            <div class="info-item">
              <label>Outstanding Balance:</label>
              <span class="info-value">{{ formatCurrency(client.loans[0].outstanding_balance) }}</span>
            </div>
            <div class="info-item">
              <label>Ammortization:</label>
              <span class="info-value">{{ formatCurrency(client.loans[0].ammortization) }}</span>
            </div>
            <div class="info-item">
              <label>Terms:</label>
              <span class="info-value">{{ client.loans[0].terms }} months</span>
            </div>
            <div class="info-item">
              <label>Status:</label>
              <span :class="['status-badge', getStatusClass(client.loans[0].status)]">
                {{ client.loans[0].status }}
              </span>
            </div>
            <div class="info-item">
              <label>Due Date:</label>
              <span class="info-value">{{ client.loans[0].due_date || 'N/A' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- No Loan State -->
      <div v-else class="detail-card">
        <h2>Loan Information</h2>
        <div class="card-content">
          <p class="no-data">No active loans for this client.</p>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="action-buttons">
        <button @click="goBack" class="btn-secondary">Back to List</button>
        <button @click="editClient" class="btn-primary">Edit Client</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../../services/api'

const route = useRoute()
const router = useRouter()

const client = ref(null)
const loading = ref(false)
const error = ref('')

const loadClient = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await api.get(`/clients/${route.params.id}`)
    client.value = response.data
  } catch (err) {
    console.error('Error loading client:', err)
    error.value = 'Failed to load client details. Please try again.'
  } finally {
    loading.value = false
  }
}

const formatCurrency = (amount) => {
  if (!amount) return '₱0.00'
  return new Intl.NumberFormat('en-PH', {
    style: 'currency',
    currency: 'PHP',
    minimumFractionDigits: 2
  }).format(amount)
}

const getStatusClass = (status) => {
  const statusMap = {
    'Active': 'status-active',
    'Paid': 'status-paid',
    'Overdue': 'status-overdue',
    'Default': 'status-default'
  }
  return statusMap[status] || 'status-unknown'
}

const goBack = () => {
  router.push('/clients/search')
}

const editClient = (client) => {
  
  router.push(`/clients/${route.params.id}/edit`)

}



onMounted(() => {
  loadClient()
})
</script>

<style scoped>
.client-details {
  padding: 1rem;
  max-width: 800px;
  margin: 0 auto;
  background-color: #f8f9fa;
  min-height: 100vh;
}

/* Header */
.page-header {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  margin-bottom: 2rem;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h1 {
  margin: 0;
  color: #2c3e50;
  flex: 1;
  text-align: center;
}

.back-btn {
  padding: 0.5rem 1rem;
  background-color: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.edit-btn {
  padding: 0.5rem 1rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.back-btn:hover {
  background-color: #5a6268;
}

.edit-btn:hover {
  background-color: #0056b3;
}

/* Detail Cards */
.detail-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  margin-bottom: 1.5rem;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.card-header h2 {
  margin: 0;
  color: #2c3e50;
}

.control-number {
  background: #007bff;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-weight: 600;
  font-size: 0.9rem;
}

.card-content {
  padding: 1.5rem;
}

/* Info Grid */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.info-item label {
  font-weight: 600;
  color: #495057;
  font-size: 0.9rem;
}

.info-value {
  color: #2c3e50;
  font-size: 1rem;
}

.address {
  line-height: 1.4;
}

/* Status Badges */
.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  display: inline-block;
}

.status-active {
  background-color: #d4edda;
  color: #155724;
}

.status-paid {
  background-color: #d1ecf1;
  color: #0c5460;
}

.status-overdue {
  background-color: #f8d7da;
  color: #721c24;
}

.status-default {
  background-color: #fff3cd;
  color: #856404;
}

.status-unknown {
  background-color: #e2e3e5;
  color: #383d41;
}

/* No Data State */
.no-data {
  color: #6c757d;
  font-style: italic;
  text-align: center;
  margin: 0;
}

/* Loading and Error States */
.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.retry-btn {
  padding: 0.75rem 1.5rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  margin-top: 1rem;
}

.retry-btn:hover {
  background-color: #2980b9;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  padding: 2rem 0;
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-primary:hover {
  background-color: #0056b3;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

/* Responsive Design */
@media (max-width: 768px) {
  .client-details {
    padding: 0.5rem;
  }
  
  .header-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>
