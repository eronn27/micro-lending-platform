<template>
  <div class="client-details">
    <!-- Header -->
    <header class="page-header">
      <div class="header-content">
        <button @click="goBack" class="back-btn">← Back to Search</button>
        <h1>Client Details</h1>
        <div class="header-actions">
          <button @click="createNewLoan" class="create-loan-btn">+ Create New Loan</button>
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

      <!-- Loan History Section -->
      <div class="detail-card loan-history-card">
        <div class="card-header">
          <h2>Loan History</h2>
          <div class="loan-stats">
            <span class="stat-badge">Total Loans: {{ loanCount }}</span>
            <span class="stat-badge active">Active: {{ activeLoanCount }}</span>
          </div>
        </div>
        
        <div v-if="client.loans && client.loans.length > 0" class="card-content">
          <!-- Current/Active Loans -->
          <div v-if="activeLoans.length > 0" class="loan-section">
            <h3 class="section-title">Current Loans</h3>
            <div class="loans-grid">
              <div 
                v-for="loan in activeLoans" 
                :key="loan.id" 
                class="loan-card active-loan"
              >
                <div class="loan-header">
                  <span class="loan-number">Loan #{{ loan.loan_cycle || loan.id }}</span>
                  <span :class="['status-badge', getStatusClass(loan.status)]">
                    {{ loan.status }}
                  </span>
                </div>
                
                <div class="loan-details">
                  <div class="loan-detail-row">
                    <span class="detail-label">Release Date:</span>
                    <span class="detail-value">{{ formatDate(loan.date_of_release) }}</span>
                  </div>
                  <div class="loan-detail-row">
                    <span class="detail-label">Total Amount:</span>
                    <span class="detail-value amount">{{ formatCurrency(loan.total_amount) }}</span>
                  </div>
                  <div class="loan-detail-row">
                    <span class="detail-label">Outstanding:</span>
                    <span class="detail-value amount-outstanding">{{ formatCurrency(loan.outstanding_balance) }}</span>
                  </div>
                  <div class="loan-detail-row">
                    <span class="detail-label">Weekly Payment:</span>
                    <span class="detail-value">{{ formatCurrency(loan.ammortization) }}</span>
                  </div>
                  <div class="loan-detail-row">
                    <span class="detail-label">Terms:</span>
                    <span class="detail-value">{{ loan.terms }} months</span>
                  </div>
                  <div class="loan-detail-row">
                    <span class="detail-label">Due Date:</span>
                    <span class="detail-value">{{ loan.due_date || 'N/A' }}</span>
                  </div>
                  <div class="loan-detail-row">
                    <span class="detail-label">Progress:</span>
                    <span class="detail-value">
                      {{ loan.paid_weeks || 0 }} / {{ loan.payment_period_weeks || 0 }} weeks
                    </span>
                  </div>
                </div>

                <div class="loan-progress-bar">
                  <div 
                    class="progress-fill" 
                    :style="{ width: calculateProgress(loan) + '%' }"
                  ></div>
                </div>
              </div>
            </div>
          </div>

          <!-- Completed/Past Loans -->
          <div v-if="completedLoans.length > 0" class="loan-section">
            <h3 class="section-title">Loan History</h3>
            <div class="loans-grid">
              <div 
                v-for="loan in completedLoans" 
                :key="loan.id" 
                class="loan-card completed-loan"
              >
                <div class="loan-header">
                  <span class="loan-number">Loan #{{ loan.loan_cycle || loan.id }}</span>
                  <span :class="['status-badge', getStatusClass(loan.status)]">
                    {{ loan.status }}
                  </span>
                </div>
                
                <div class="loan-details">
                  <div class="loan-detail-row">
                    <span class="detail-label">Release Date:</span>
                    <span class="detail-value">{{ formatDate(loan.date_of_release) }}</span>
                  </div>
                  <div class="loan-detail-row">
                    <span class="detail-label">Total Amount:</span>
                    <span class="detail-value amount">{{ formatCurrency(loan.total_amount) }}</span>
                  </div>
                  <div class="loan-detail-row">
                    <span class="detail-label">Terms:</span>
                    <span class="detail-value">{{ loan.terms }} months</span>
                  </div>
                  <div class="loan-detail-row">
                    <span class="detail-label">Method:</span>
                    <span class="detail-value">{{ loan.method_of_payment || 'N/A' }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- No Loans State -->
        <div v-else class="card-content">
          <div class="no-loans-state">
            <div class="no-loans-icon">📋</div>
            <p class="no-loans-text">No loans found for this client.</p>
            <button @click="createNewLoan" class="btn-create-first-loan">
              Create First Loan
            </button>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="action-buttons">
        <button @click="goBack" class="btn-secondary">Back to List</button>
        <button @click="createNewLoan" class="btn-primary">Create New Loan</button>
        <button @click="editClient" class="btn-primary">Edit Client</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../../services/api'

const route = useRoute()
const router = useRouter()

const client = ref(null)
const loading = ref(false)
const error = ref('')

// Computed properties for loan statistics
const loanCount = computed(() => {
  return client.value?.loans?.length || 0
})

const activeLoanCount = computed(() => {
  return activeLoans.value.length
})

const activeLoans = computed(() => {
  if (!client.value?.loans) return []
  return client.value.loans.filter(loan => 
    loan.status === 'Active' || loan.status === 'Overdue'
  ).sort((a, b) => new Date(b.date_of_release) - new Date(a.date_of_release))
})

const completedLoans = computed(() => {
  if (!client.value?.loans) return []
  return client.value.loans.filter(loan => 
    loan.status === 'Paid' || loan.status === 'Default'
  ).sort((a, b) => new Date(b.date_of_release) - new Date(a.date_of_release))
})

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

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-PH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
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

const calculateProgress = (loan) => {
  if (!loan.payment_period_weeks || loan.payment_period_weeks === 0) return 0
  const progress = ((loan.paid_weeks || 0) / loan.payment_period_weeks) * 100
  return Math.min(Math.round(progress), 100)
}

const goBack = () => {
  router.push('/clients/search')
}

const editClient = () => {
  router.push(`/clients/${route.params.id}/edit`)
}

const createNewLoan = () => {
  router.push(`/clients/${route.params.id}/new-loan`)
}

onMounted(() => {
  loadClient()
})
</script>

<style scoped>
.client-details {
  padding: 1rem;
  max-width: 1200px;
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

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.back-btn, .edit-btn, .create-loan-btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s ease;
}

.back-btn {
  background-color: #6c757d;
  color: white;
}

.back-btn:hover {
  background-color: #5a6268;
}

.edit-btn {
  background-color: #007bff;
  color: white;
}

.edit-btn:hover {
  background-color: #0056b3;
}

.create-loan-btn {
  background-color: #28a745;
  color: white;
  font-size: 0.95rem;
}

.create-loan-btn:hover {
  background-color: #218838;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(40, 167, 69, 0.3);
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

/* Loan History Styles */
.loan-history-card {
  border: 2px solid #e9ecef;
}

.loan-stats {
  display: flex;
  gap: 0.75rem;
}

.stat-badge {
  background: #e9ecef;
  color: #495057;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.85rem;
  font-weight: 600;
}

.stat-badge.active {
  background: #d4edda;
  color: #155724;
}

.loan-section {
  margin-bottom: 2rem;
}

.loan-section:last-child {
  margin-bottom: 0;
}

.section-title {
  margin: 0 0 1rem 0;
  color: #495057;
  font-size: 1.1rem;
  font-weight: 600;
  border-bottom: 2px solid #e9ecef;
  padding-bottom: 0.5rem;
}

.loans-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.loan-card {
  background: #fff;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  padding: 1.25rem;
  transition: all 0.2s ease;
}

.loan-card:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  transform: translateY(-2px);
}

.active-loan {
  border-left: 4px solid #28a745;
}

.completed-loan {
  border-left: 4px solid #6c757d;
  opacity: 0.85;
}

.loan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #e9ecef;
}

.loan-number {
  font-weight: 700;
  color: #2c3e50;
  font-size: 1.1rem;
}

.loan-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.loan-detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-label {
  font-size: 0.85rem;
  color: #6c757d;
  font-weight: 500;
}

.detail-value {
  font-size: 0.9rem;
  color: #2c3e50;
  font-weight: 600;
}

.detail-value.amount {
  color: #28a745;
}

.detail-value.amount-outstanding {
  color: #dc3545;
}

.loan-progress-bar {
  margin-top: 1rem;
  height: 8px;
  background: #e9ecef;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #28a745, #20c997);
  transition: width 0.3s ease;
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

/* No Loans State */
.no-loans-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 2rem;
  text-align: center;
}

.no-loans-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.no-loans-text {
  color: #6c757d;
  font-size: 1.1rem;
  margin: 0 0 1.5rem 0;
}

.btn-create-first-loan {
  padding: 0.75rem 2rem;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  transition: all 0.2s ease;
}

.btn-create-first-loan:hover {
  background-color: #218838;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(40, 167, 69, 0.3);
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

.btn-primary:hover {
  background-color: #0056b3;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
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

  .header-actions {
    width: 100%;
    justify-content: center;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }

  .loans-grid {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>