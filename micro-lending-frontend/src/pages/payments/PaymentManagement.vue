<template>
  <div class="payment-management">
    <!-- Header Section -->
    <header class="page-header">
      <div class="header-content">
        <h1>Payment Management</h1>
        <div class="due-today-badge" v-if="dueTodayCount > 0">
          <span class="badge-count">{{ dueTodayCount }}</span>
          Due Today
        </div>
      </div>
      
      <div class="header-actions">
        <div class="bulk-actions" v-if="selectedClients.length > 0">
          <button 
            @click="processBulkPayments" 
            class="btn-bulk-pay"
            :disabled="processingBulk"
          >
            <span v-if="processingBulk">Processing...</span>
            <span v-else>Mark {{ selectedClients.length }} Selected as Paid</span>
          </button>
          <button @click="clearSelection" class="btn-clear">Clear Selection</button>
        </div>
        
        <button @click="goBack" class="back-btn">
          ← Back to Dashboard
        </button>
      </div>
    </header>

    <!-- Filters & Search -->
    <div class="filters-section">
      <div class="search-box">
        <input
          type="text"
          v-model="searchQuery"
          placeholder="Search by name or control number..."
          class="search-input"
          @input="handleSearch"
        />
        <span class="search-icon">🔍</span>
      </div>

      <div class="filter-buttons">
        <button
          v-for="filter in statusFilters"
          :key="filter.value"
          @click="setStatusFilter(filter.value)"
          :class="['filter-btn', { active: statusFilter === filter.value }]"
        >
          {{ filter.label }}
          <span class="filter-count" v-if="filter.count > 0">
            {{ filter.count }}
          </span>
        </button>
      </div>
    </div>

    <!-- Clients Table -->
    <div class="table-container">
      <table class="clients-table">
        <thead>
          <tr>
            <th class="checkbox-column">
              <input
                type="checkbox"
                v-model="selectAll"
                @change="toggleSelectAll"
              />
            </th>
            <th>Control #</th>
            <th>Name</th>
            <th>Date of Release</th>
            <th>Outstanding Balance</th>
            <th>Amortization</th>
            <th>Terms</th>
            <th>Mode</th>
            <th>Due Date</th>
            <th>Payment Progress</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="client in filteredClients"
            :key="client.id"
            :class="getRowClass(client)"
          >
            <!-- Checkbox -->
            <td class="checkbox-column">
              <input
                type="checkbox"
                :value="client.id"
                v-model="selectedClients"
                :disabled="getClientStatus(client) === 'paid'"
              />
            </td>

            <!-- Control Number -->
            <td class="control-number">
              {{ client.control_number }}
            </td>

            <!-- Name -->
            <td class="client-name">
              {{ client.first_name }} {{ client.last_name }}
            </td>

            <!-- Date of Release -->
            <td class="date-release">
              {{ formatDate(getLoanField(client, 'date_of_release')) }}
            </td>

            <!-- Outstanding Balance -->
            <td class="amount">
              ₱{{ formatCurrency(getLoanField(client, 'outstanding_balance')) }}
            </td>

            <!-- Amortization -->
            <td class="amortization">
              ₱{{ formatCurrency(getLoanField(client, 'ammortization')) }}
            </td>

            <!-- Terms (in months) -->
            <td class="terms">
              {{ getLoanField(client, 'terms') || 0 }} months
            </td>

            <!-- Mode -->
            <td class="mode">
              {{ getLoanField(client, 'mode') || 'Weekly' }}
            </td>

            <!-- Due Date -->
            <td class="due-date">
              {{ getLoanField(client, 'due_date') || '-' }}
            </td>

            <!-- Payment Progress -->
            <td class="progress-cell">
              <div class="progress-container">
                <div class="progress-info">
                  {{ getPaidWeeks(client) }}/{{ getLoanField(client, 'payment_period_weeks') || 0 }} weeks
                </div>
                <div class="progress-bar">
                  <div
                    class="progress-fill"
                    :style="{
                      width: `${getProgressPercentage(client)}%`,
                      backgroundColor: getProgressColor(client)
                    }"
                  ></div>
                </div>
              </div>
            </td>

            <!-- Status -->
            <td class="status-cell">
              <span :class="['status-badge', getStatusClass(client)]">
                {{ getStatusText(client) }}
              </span>
            </td>

            <!-- Action -->
            <td class="action-cell">
              <button
                v-if="getClientStatus(client) !== 'paid'"
                @click="processPayment(client)"
                :disabled="isProcessingPayment(client.id) || isPaidToday(client)"
                :class="['btn-pay', { processing: isProcessingPayment(client.id) }]"
              >
                <span v-if="isProcessingPayment(client.id)">Processing...</span>
                <span v-else>Mark Paid</span>
              </button>
              <span v-else class="paid-text">🎉 Paid</span>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Empty State -->
      <div v-if="filteredClients.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">💸</div>
        <h3>No clients found</h3>
        <p>Try adjusting your search or filters</p>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading payments...</p>
      </div>
    </div>

    <!-- Notifications -->
    <div v-if="notification.show" :class="['notification', notification.type]">
      {{ notification.message }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../services/api'

const router = useRouter()
const goBack = () => {
  router.push('/dashboard')
}


// Reactive state
const clients = ref([])
const selectedClients = ref([])
const selectAll = ref(false)
const searchQuery = ref('')
const statusFilter = ref('all')
const processingPayments = ref(new Set())
const processingBulk = ref(false)
const loading = ref(false)
const notification = ref({ show: false, message: '', type: 'success' })

// Helper function to get loan data from client
const getLoanField = (client, field) => {
  if (client.loans && client.loans.length > 0) {
    return client.loans[0][field]
  }
  return client.loan ? client.loan[field] : null
}

// Status filters with counts
const statusFilters = computed(() => [
  { label: 'All', value: 'all', count: clients.value.length },
  { label: 'Due Today', value: 'due_today', count: dueTodayCount.value },
  { label: 'Overdue', value: 'overdue', count: overdueCount.value },
  { label: 'Active', value: 'active', count: activeCount.value },
  { label: 'Paid-in-Full', value: 'paid', count: paidCount.value }
])

// Computed properties
const filteredClients = computed(() => {
  let filtered = clients.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(client =>
      client.control_number?.toLowerCase().includes(query) ||
      client.first_name?.toLowerCase().includes(query) ||
      client.last_name?.toLowerCase().includes(query)
    )
  }

  if (statusFilter.value !== 'all') {
    filtered = filtered.filter(client => getClientStatus(client) === statusFilter.value)
  }

  return filtered.sort((a, b) => {
    const statusOrder = { due_today: 0, overdue: 1, active: 2, paid: 3 }
    return statusOrder[getClientStatus(a)] - statusOrder[getClientStatus(b)]
  })
})

const dueTodayCount = computed(() => 
  clients.value.filter(client => getClientStatus(client) === 'due_today').length
)

const overdueCount = computed(() => 
  clients.value.filter(client => getClientStatus(client) === 'overdue').length
)

const activeCount = computed(() => 
  clients.value.filter(client => getClientStatus(client) === 'active').length
)

const paidCount = computed(() => 
  clients.value.filter(client => getClientStatus(client) === 'paid').length
)

// Methods
const fetchClients = async () => {
  loading.value = true
  try {
    const response = await api.get('/clients')
    
    clients.value = response.data.clients?.map(client => {
      // Use the paid_weeks from database directly
      const paidWeeks = getLoanField(client, 'paid_weeks') || 0
      const status = calculateClientStatus(client, paidWeeks)
      
      return {
        ...client,
        paid_weeks: paidWeeks,
        status: status
      }
    }) || []
    
  } catch (error) {
    showNotification('Failed to load clients', 'error')
    console.error('Error fetching clients:', error)
  } finally {
    loading.value = false
  }
}

// SIMPLIFIED - Get from database field
const calculatePaidWeeks = (client) => {
  return getLoanField(client, 'paid_weeks') || 0
}

const calculateClientStatus = (client, paidWeeks) => {
  const paymentPeriodWeeks = getLoanField(client, 'payment_period_weeks') || 0
  const loanStatus = getLoanField(client, 'status')
  const outstandingBalance = getLoanField(client, 'outstanding_balance') || 0
  
  if (loanStatus === 'Paid' || outstandingBalance <= 0 || paidWeeks >= paymentPeriodWeeks) {
    return 'paid'
  }
  
  const dateOfRelease = getLoanField(client, 'date_of_release')
  if (!dateOfRelease) return 'active'
  
  const releaseDate = new Date(dateOfRelease)
  const today = new Date()
  const daysSinceRelease = Math.floor((today - releaseDate) / (1000 * 60 * 60 * 24))
  const currentWeek = Math.floor(daysSinceRelease / 7) + 1
  
  if (currentWeek > paidWeeks && currentWeek <= paymentPeriodWeeks) {
    const daysInCurrentWeek = daysSinceRelease % 7
    const dueDate = getLoanField(client, 'due_date')
    
    if (dueDate && isDueToday(dueDate)) {
      return 'due_today'
    }
    
    if (daysInCurrentWeek >= 3 && daysInCurrentWeek <= 6) {
      return 'due_today'
    }
    return 'active'
  }
  
  if (currentWeek > paidWeeks && currentWeek > paymentPeriodWeeks) {
    return 'overdue'
  }
  
  return 'active'
}

const isDueToday = (dueDate) => {
  const daysOfWeek = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
  const today = new Date()
  const todayDay = daysOfWeek[today.getDay()]
  return dueDate.toLowerCase() === todayDay.toLowerCase()
}

const getClientStatus = (client) => {
  return client.status || calculateClientStatus(client, client.paid_weeks || 0)
}

// Always get from the database field
const getPaidWeeks = (client) => {
  return getLoanField(client, 'paid_weeks') || 0
}

const handleSearch = () => {
  // Filtering handled by computed property
}

const processPayment = async (client) => {
  processingPayments.value.add(client.id)
  
  try {
    const loanId = getLoanField(client, 'id')
    const nextWeek = getPaidWeeks(client) + 1
    const paymentAmount = getLoanField(client, 'ammortization')
    const currentBalance = getLoanField(client, 'outstanding_balance') || 0
    const newBalance = Math.max(0, currentBalance - paymentAmount)
    
    // Create payment record
    const paymentPayload = {
      loan_id: loanId,
      week_number: nextWeek,
      amount_due: paymentAmount,
      amount_paid: paymentAmount,
      payment_date: new Date().toISOString().split('T')[0],
      status: 'Paid',
      payment_method: 'Cash'
    }

    await api.post('/payments', paymentPayload)
    
    // Refresh data from database to get updated paid_weeks
    await fetchClients()
    
    showNotification(`Payment recorded for ${client.first_name} ${client.last_name} - Week ${nextWeek}`)
    
    const index = selectedClients.value.indexOf(client.id)
    if (index > -1) {
      selectedClients.value.splice(index, 1)
    }
  } catch (error) {
    showNotification('Failed to process payment', 'error')
    console.error('Error processing payment:', error)
  } finally {
    processingPayments.value.delete(client.id)
  }
}

const processBulkPayments = async () => {
  processingBulk.value = true
  
  try {
    const paymentPromises = selectedClients.value.map(clientId => {
      const client = clients.value.find(c => c.id === clientId)
      if (client && getClientStatus(client) !== 'paid') {
        return processPayment(client)
      }
      return Promise.resolve()
    })

    await Promise.all(paymentPromises)
    showNotification(`Processed payments for ${selectedClients.value.length} clients`)
    selectedClients.value = []
  } catch (error) {
    showNotification('Some payments failed to process', 'error')
  } finally {
    processingBulk.value = false
  }
}

const toggleSelectAll = () => {
  if (selectAll.value) {
    selectedClients.value = filteredClients.value
      .filter(client => getClientStatus(client) !== 'paid')
      .map(client => client.id)
  } else {
    selectedClients.value = []
  }
}

const clearSelection = () => {
  selectedClients.value = []
  selectAll.value = false
}

const setStatusFilter = (filter) => {
  statusFilter.value = filter
}

const getRowClass = (client) => {
  const classes = []
  const status = getClientStatus(client)
  if (status === 'due_today') classes.push('due-today-row')
  if (status === 'overdue') classes.push('overdue-row')
  if (selectedClients.value.includes(client.id)) classes.push('selected-row')
  return classes
}

const getStatusClass = (client) => {
  return `status-${getClientStatus(client)}`
}

const getStatusText = (client) => {
  const statusMap = {
    due_today: '⚠️ Due Today',
    overdue: '🚨 Overdue',
    active: 'Active',
    paid: '🎉 Paid'
  }
  return statusMap[getClientStatus(client)] || 'Active'
}

const getProgressPercentage = (client) => {
  const total = getLoanField(client, 'payment_period_weeks') || 1
  const paid = getPaidWeeks(client) || 0
  return (paid / total) * 100
}

const getProgressColor = (client) => {
  const percentage = getProgressPercentage(client)
  if (percentage >= 100) return '#10b981'
  if (percentage >= 75) return '#3b82f6'
  if (percentage >= 50) return '#f59e0b'
  return '#ef4444'
}

const isProcessingPayment = (clientId) => {
  return processingPayments.value.has(clientId)
}

const isPaidToday = (client) => {
  return false
}

const formatCurrency = (amount) => {
  if (!amount && amount !== 0) return '0.00'
  return new Intl.NumberFormat('en-PH', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount || 0)
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  try {
    return new Date(dateString).toLocaleDateString('en-PH')
  } catch (error) {
    return '-'
  }
}

const showNotification = (message, type = 'success') => {
  notification.value = { show: true, message, type }
  setTimeout(() => {
    notification.value.show = false
  }, 3000)
}

watch(filteredClients, (newFiltered) => {
  const selectableClients = newFiltered.filter(client => getClientStatus(client) !== 'paid')
  if (selectableClients.length > 0) {
    selectAll.value = selectedClients.value.length === selectableClients.length
  } else {
    selectAll.value = false
  }
})

onMounted(() => {
  fetchClients()
})
</script>


<style scoped>
.payment-management {
  padding: 1.5rem;
  max-width: 1600px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding: 1.5rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-content h1 {
  margin: 0;
  color: #1f2937;
  font-size: 1.75rem;
}

.due-today-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: white;
  border-radius: 20px;
  font-weight: 600;
  font-size: 0.875rem;
}

.badge-count {
  background: rgba(255, 255, 255, 0.2);
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
}

.bulk-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.btn-bulk-pay {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-bulk-pay:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

.btn-bulk-pay:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-clear {
  background: #6b7280;
  color: white;
  border: none;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-clear:hover {
  background: #4b5563;
}

.filters-section {
  display: flex;
  gap: 2rem;
  margin-bottom: 1.5rem;
  align-items: center;
}

.search-box {
  position: relative;
  flex: 1;
  max-width: 400px;
}

.search-input {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 2.5rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 0.875rem;
  transition: border-color 0.2s;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  color: #6b7280;
}

.filter-buttons {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 20px;
  background: white;
  color: #6b7280;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-btn.active {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.filter-btn:hover:not(.active) {
  background: #f3f4f6;
}

.filter-count {
  background: rgba(255, 255, 255, 0.2);
  padding: 0.125rem 0.5rem;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 600;
}

.table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow-x: auto;
}

.clients-table {
  width: 100%;
  border-collapse: collapse;
  min-width: 1200px;
}

.clients-table th {
  background: #f8fafc;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #374151;
  font-size: 0.875rem;
  border-bottom: 1px solid #e5e7eb;
  white-space: nowrap;
}

.clients-table td {
  padding: 1rem;
  border-bottom: 1px solid #f3f4f6;
  font-size: 0.875rem;
}

.clients-table tr:last-child td {
  border-bottom: none;
}

.due-today-row {
  background: #fffbeb !important;
  border-left: 4px solid #f59e0b;
}

.overdue-row {
  background: #fef2f2 !important;
  border-left: 4px solid #ef4444;
  animation: pulse-overdue 2s infinite;
}

.selected-row {
  background: #eff6ff !important;
}

@keyframes pulse-overdue {
  0%, 100% { background-color: #fef2f2; }
  50% { background-color: #fee2e2; }
}

.checkbox-column {
  width: 40px;
  text-align: center;
}

.control-number {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  color: #1f2937;
}

.client-name {
  font-weight: 500;
  color: #1f2937;
}

.amount, .amortization {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  color: #059669;
}

.mode {
  font-weight: 500;
  color: #6b7280;
  text-align: center;
}

.due-date {
  font-weight: 500;
  color: #1f2937;
  text-align: center;
}

.progress-container {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  min-width: 120px;
}

.progress-info {
  font-size: 0.75rem;
  color: #6b7280;
  text-align: center;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: #e5e7eb;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.3s ease;
}

.status-badge {
  padding: 0.375rem 0.75rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  text-align: center;
  display: inline-block;
  white-space: nowrap;
}

.status-due_today {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #f59e0b;
}

.status-overdue {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #ef4444;
  animation: pulse-badge 2s infinite;
}

.status-active {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #3b82f6;
}

.status-paid {
  background: #f0fdf4;
  color: #059669;
  border: 1px solid #10b981;
}

@keyframes pulse-badge {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.btn-pay {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 90px;
  white-space: nowrap;
}

.btn-pay:hover:not(:disabled):not(.processing) {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(59, 130, 246, 0.3);
}

.btn-pay:disabled {
  background: #9ca3af;
  cursor: not-allowed;
  transform: none;
}

.btn-pay.processing {
  background: #6b7280;
  cursor: not-allowed;
}

.paid-text {
  color: #059669;
  font-weight: 600;
  font-size: 0.875rem;
  white-space: nowrap;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
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

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #6b7280;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.empty-state h3 {
  margin: 0 0 0.5rem 0;
  color: #374151;
}

.empty-state p {
  margin: 0;
  font-size: 0.875rem;
}

.notification {
  position: fixed;
  top: 2rem;
  right: 2rem;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  color: white;
  font-weight: 500;
  z-index: 1000;
  animation: slideIn 0.3s ease;
}

.notification.success {
  background: #10b981;
}

.notification.error {
  background: #ef4444;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@media (max-width: 1400px) {
  .payment-management {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .filters-section {
    flex-direction: column;
    gap: 1rem;
  }

  .search-box {
    max-width: none;
  }
}

@media (max-width: 768px) {
  .filter-buttons {
    justify-content: center;
  }

  .clients-table th,
  .clients-table td {
    padding: 0.75rem 0.5rem;
    font-size: 0.8rem;
  }
}
.header-actions {
  display: flex;
  justify-content: flex-end; 
  align-items: center;
  gap: 1rem;
  flex: 1;
  margin-left: 2rem;
}

.back-btn {
  padding: 0.75rem 1.5rem;
  background-color: #6c757d;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
 
}

.bulk-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin-right: auto; 
}
.back-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(156, 151, 151, 0.3);
}
</style>
