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

    <!-- Payment Management Table -->
    <div class="table-section">
      <h2 class="table-title">Payment Management</h2>
      <p class="table-subtitle">Clients requiring payment attention</p>
      
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

              <!-- Amortization (UPDATED) -->
              <td class="amortization">
                <div class="amortization-display">
                  <div class="full-amount">₱{{ formatCurrency(getLoanField(client, 'ammortization')) }}</div>
                  <div 
                    v-if="getRemainingBalance(client) > 0" 
                    :class="['remaining-balance', getRemainingBalanceClass(client)]"
                  >
                    ₱{{ formatCurrency(getRemainingBalance(client)) }} remaining
                  </div>
                  <div 
                    v-else-if="hasPartialPayments(client)"
                    class="balance-paid"
                  >
                    Week completed with partial payments
                  </div>
                </div>
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
                    <span v-if="hasPartialPayments(client)" class="partial-indicator">
                      (+{{ getCurrentWeekPartialPayments(client).length }} partial)
                    </span>
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

              <!-- Action (UPDATED) -->
              <td class="action-cell">
                <div class="action-buttons">
                  <button
                    v-if="getClientStatus(client) !== 'paid' && getRemainingBalance(client) === 0"
                    @click="processPayment(client)"
                    :disabled="isProcessingPayment(client.id) || isPaidToday(client)"
                    :class="['btn-pay', { processing: isProcessingPayment(client.id) }]"
                  >
                    <span v-if="isProcessingPayment(client.id)">Processing...</span>
                    <span v-else>Complete Week</span>
                  </button>
                  
                  <button
                    v-if="getClientStatus(client) !== 'paid' && getRemainingBalance(client) > 0"
                    @click="openPartialPaymentModal(client)"
                    :disabled="isProcessingPayment(client.id)"
                    class="btn-partial-pay"
                  >
                    Partial Payment
                  </button>
                  
                  <button
                    v-if="getClientStatus(client) !== 'paid' && getRemainingBalance(client) > 0"
                    @click="processFullPayment(client)"
                    :disabled="isProcessingPayment(client.id)"
                    class="btn-pay-full"
                  >
                    Pay Full ₱{{ formatCurrency(getRemainingBalance(client)) }}
                  </button>
                  
                  <span v-else-if="getClientStatus(client) === 'paid'" class="paid-text">🎉 Paid</span>
                </div>
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
    </div>

    <!-- Partial Payment Modal -->
    <div v-if="showPartialPaymentModal" class="modal-overlay" @click="closePartialPaymentModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Partial Payment</h3>
          <button class="modal-close" @click="closePartialPaymentModal">×</button>
        </div>
        
        <div class="modal-body">
          <div class="client-info">
            <h4>{{ selectedClient?.first_name }} {{ selectedClient?.last_name }}</h4>
            <p>Control #: {{ selectedClient?.control_number }}</p>
            <p>Current Week: {{ getCurrentWeek(selectedClient) }}</p>
          </div>

          <div class="payment-summary">
            <div class="summary-item">
              <span>Full Amortization:</span>
              <span class="amount">₱{{ formatCurrency(getLoanField(selectedClient, 'ammortization')) }}</span>
            </div>
            <div class="summary-item">
              <span>Remaining Balance:</span>
              <span :class="['amount', getRemainingBalanceClass(selectedClient)]">
                ₱{{ formatCurrency(getRemainingBalance(selectedClient)) }}
              </span>
            </div>
            <div class="summary-item" v-if="getCurrentWeekPartialPayments(selectedClient).length > 0">
              <span>Previous Partial Payments:</span>
              <span class="amount">₱{{ formatCurrency(getCurrentWeekPartialTotal(selectedClient)) }}</span>
            </div>
          </div>

          <form @submit.prevent="processPartialPayment" class="payment-form">
            <div class="form-group">
              <label for="paymentAmount">Payment Amount</label>
              <input
                id="paymentAmount"
                type="number"
                v-model="partialPaymentAmount"
                :max="getRemainingBalance(selectedClient)"
                :min="1"
                step="0.01"
                required
                class="amount-input"
                placeholder="Enter amount"
              />
              <div class="input-hint">
                Maximum: ₱{{ formatCurrency(getRemainingBalance(selectedClient)) }}
              </div>
            </div>

            <div class="form-group">
              <label for="paymentMethod">Payment Method</label>
              <select id="paymentMethod" v-model="partialPaymentMethod" required>
                <option value="Cash">Cash</option>
                <option value="Bank Transfer">Bank Transfer</option>
                <option value="GCash">GCash</option>
                <option value="Maya">Maya</option>
                <option value="Check">Check</option>
              </select>
            </div>

            <div class="payment-preview">
              <div class="preview-item">
                <span>Payment Amount:</span>
                <span class="amount">₱{{ formatCurrency(partialPaymentAmount || 0) }}</span>
              </div>
              <div class="preview-item">
                <span>New Remaining Balance:</span>
                <span class="amount">
                  ₱{{ formatCurrency(getRemainingBalance(selectedClient) - (partialPaymentAmount || 0)) }}
                </span>
              </div>
              <div class="preview-item" v-if="(getRemainingBalance(selectedClient) - (partialPaymentAmount || 0)) <= 0">
                <span class="completion-notice">✅ This payment will complete the week!</span>
              </div>
            </div>

            <div class="modal-actions">
              <button
                type="button"
                @click="closePartialPaymentModal"
                class="btn-cancel"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="!partialPaymentAmount || partialPaymentAmount <= 0 || processingPartialPayment"
                class="btn-confirm-partial"
              >
                <span v-if="processingPartialPayment">Processing...</span>
                <span v-else>Confirm Partial Payment</span>
              </button>
            </div>
          </form>
        </div>
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

// Reactive state
const clients = ref([])
const selectedClients = ref([])
const selectAll = ref(false)
const searchQuery = ref('')
const statusFilter = ref('all') // Default to show all loans
const processingPayments = ref(new Set())
const processingBulk = ref(false)
const loading = ref(false)
const notification = ref({ show: false, message: '', type: 'success' })

// Partial payment modal state
const showPartialPaymentModal = ref(false)
const selectedClient = ref(null)
const partialPaymentAmount = ref(0)
const partialPaymentMethod = ref('Cash')
const processingPartialPayment = ref(false)

// Payment progress cache
const paymentProgressCache = ref(new Map())

// Helper function to get loan data from client
const getLoanField = (client, field) => {
  if (client.loans && client.loans.length > 0) {
    return client.loans[0][field]
  }
  return client.loan ? client.loan[field] : null
}

// Get current week number
const getCurrentWeek = (client) => {
  return (getLoanField(client, 'paid_weeks') || 0) + 1
}

// Fetch payment progress for a client
const fetchPaymentProgress = async (client) => {
  const loanId = getLoanField(client, 'id')
  if (!loanId) return null

  try {
    const response = await api.get(`/payments/loan/${loanId}`)
    const payments = response.data.payments || []
    
    // Calculate current week partial payments
    const currentWeek = getCurrentWeek(client)
    const currentWeekPayments = payments.filter(payment => 
      payment.week_number === currentWeek && payment.is_partial
    )
    
    const progress = {
      partial_payments: currentWeekPayments,
      remaining_balance: calculateRemainingBalance(client, currentWeekPayments)
    }
    
    paymentProgressCache.value.set(loanId, progress)
    return progress
  } catch (error) {
    console.error('Failed to fetch payment progress:', error)
    return null
  }
}

// Calculate remaining balance
const calculateRemainingBalance = (client, partialPayments = null) => {
  const amortization = getLoanField(client, 'ammortization') || 0
  const payments = partialPayments || getCurrentWeekPartialPayments(client)
  const partialTotal = payments.reduce((total, payment) => total + payment.amount_paid, 0)
  return Math.max(0, amortization - partialTotal)
}

// Get current week partial payments
const getCurrentWeekPartialPayments = (client) => {
  const loanId = getLoanField(client, 'id')
  const progress = paymentProgressCache.value.get(loanId)
  return progress ? progress.partial_payments : []
}

// Calculate total of current week partial payments
const getCurrentWeekPartialTotal = (client) => {
  return getCurrentWeekPartialPayments(client).reduce((total, payment) => total + payment.amount_paid, 0)
}

// Get remaining balance
const getRemainingBalance = (client) => {
  const loanId = getLoanField(client, 'id')
  const progress = paymentProgressCache.value.get(loanId)
  if (progress) {
    return progress.remaining_balance
  }
  
  // Fallback calculation
  return calculateRemainingBalance(client)
}

// Check if client has any partial payments
const hasPartialPayments = (client) => {
  return getCurrentWeekPartialPayments(client).length > 0
}

// Get CSS class for remaining balance display
const getRemainingBalanceClass = (client) => {
  const remaining = getRemainingBalance(client)
  const amortization = getLoanField(client, 'ammortization') || 0
  
  if (remaining === 0) return 'balance-paid'
  if (remaining <= amortization * 0.3) return 'balance-low'
  if (remaining <= amortization * 0.7) return 'balance-medium'
  return 'balance-high'
}

// Status filters with counts - includes all statuses
const statusFilters = computed(() => [
  { label: 'All Loans', value: 'all', count: clients.value.length },
  { label: 'Active', value: 'active', count: activeCount.value },
  { label: 'Due Today', value: 'due_today', count: dueTodayCount.value },
  { label: 'Overdue', value: 'overdue', count: overdueCount.value },
  { label: 'Paid-in-Full', value: 'paid', count: paidCount.value }
])

// Computed properties
const activeClients = computed(() => {
  return filteredClients.value.filter(client => getClientStatus(client) !== 'paid')
})

const paidClients = computed(() => {
  return filteredClients.value.filter(client => getClientStatus(client) === 'paid')
})

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
    // Fetch all clients with their loans (both active and paid)
    const response = await api.get('/clients')
    
    clients.value = response.data.clients?.map(client => {
      const paidWeeks = getLoanField(client, 'paid_weeks') || 0
      const status = calculateClientStatus(client, paidWeeks)

      return {
        ...client,
        paid_weeks: paidWeeks,
        status: status
      }
    }) || []

    // Fetch payment progress for all clients (both active and paid)
    await Promise.all(clients.value.map(client => fetchPaymentProgress(client)))
    
  } catch (error) {
    showNotification('Failed to load clients', 'error')
    console.error('Error fetching clients:', error)
  } finally {
    loading.value = false
  }
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

const getPaidWeeks = (client) => {
  return getLoanField(client, 'paid_weeks') || 0
}

// New method to get completion date for paid clients
const getCompletionDate = (client) => {
  const lastPaymentDate = getLoanField(client, 'last_payment_date') || getLoanField(client, 'updated_at')
  return lastPaymentDate ? formatDate(lastPaymentDate) : 'Completed'
}

const handleSearch = () => {
  // Filtering handled by computed property
}

// Partial Payment Modal Methods
const openPartialPaymentModal = async (client) => {
  selectedClient.value = client
  // Ensure we have latest payment progress
  await fetchPaymentProgress(client)
  partialPaymentAmount.value = getRemainingBalance(client)
  partialPaymentMethod.value = 'Cash'
  showPartialPaymentModal.value = true
}

const closePartialPaymentModal = () => {
  showPartialPaymentModal.value = false
  selectedClient.value = null
  partialPaymentAmount.value = 0
  processingPartialPayment.value = false
}

const processPartialPayment = async () => {
  if (!selectedClient.value || !partialPaymentAmount.value) return
  
  processingPartialPayment.value = true
  
  try {
    const loanId = getLoanField(selectedClient.value, 'id')
    const currentWeek = getCurrentWeek(selectedClient.value)
    const paymentAmount = parseFloat(partialPaymentAmount.value)
    const amortization = getLoanField(selectedClient.value, 'ammortization')
    
    // Calculate if this payment completes the week
    const currentRemaining = getRemainingBalance(selectedClient.value)
    const completesWeek = (currentRemaining - paymentAmount) <= 0
    
    // Create partial payment record
    const paymentPayload = {
      loan_id: loanId,
      week_number: currentWeek,
      amount_due: amortization,
      amount_paid: paymentAmount,
      payment_date: new Date().toISOString().split('T')[0],
      status: 'Partial',
      payment_method: partialPaymentMethod.value,
      is_partial: true,
      completes_week: completesWeek
    }

    await api.post('/payments', paymentPayload)
    
    // Refresh client data and payment progress
    await fetchClients()
    
    let message = `Partial payment of ₱${formatCurrency(paymentAmount)} recorded for ${selectedClient.value.first_name} ${selectedClient.value.last_name}`
    if (completesWeek) {
      message += ' - Week completed!'
    }
    
    showNotification(message)
    
    closePartialPaymentModal()
  } catch (error) {
    showNotification('Failed to process partial payment', 'error')
    console.error('Error processing partial payment:', error)
  } finally {
    processingPartialPayment.value = false
  }
}

// Process full payment for the remaining balance
const processFullPayment = async (client) => {
  processingPayments.value.add(client.id)
  
  try {
    const loanId = getLoanField(client, 'id')
    const currentWeek = getCurrentWeek(client)
    const paymentAmount = getRemainingBalance(client)
    
    if (paymentAmount <= 0) {
      showNotification('No remaining balance to pay', 'error')
      return
    }

    // Create payment record to complete the week
    const paymentPayload = {
      loan_id: loanId,
      week_number: currentWeek,
      amount_due: getLoanField(client, 'ammortization'),
      amount_paid: paymentAmount,
      payment_date: new Date().toISOString().split('T')[0],
      status: 'Paid',
      payment_method: 'Cash',
      is_partial: false,
      completes_week: true
    }

    await api.post('/payments', paymentPayload)
    
    // Refresh client data
    await fetchClients()
    
    showNotification(`Full payment of ₱${formatCurrency(paymentAmount)} recorded for ${client.first_name} ${client.last_name} - Week ${currentWeek} completed!`)
    
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

// Process payment to complete the week (when remaining balance is 0)
const processPayment = async (client) => {
  processingPayments.value.add(client.id)
  
  try {
    const loanId = getLoanField(client, 'id')
    const currentWeek = getCurrentWeek(client)
    const amortization = getLoanField(client, 'ammortization')
    
    // Create payment record to mark week as completed
    const paymentPayload = {
      loan_id: loanId,
      week_number: currentWeek,
      amount_due: amortization,
      amount_paid: 0, // No payment needed since partial payments already covered it
      payment_date: new Date().toISOString().split('T')[0],
      status: 'Paid',
      payment_method: 'Cash',
      is_partial: false,
      completes_week: true
    }

    await api.post('/payments', paymentPayload)
    
    // Refresh client data
    await fetchClients()
    
    showNotification(`Week ${currentWeek} completed for ${client.first_name} ${client.last_name}`)
    
    const index = selectedClients.value.indexOf(client.id)
    if (index > -1) {
      selectedClients.value.splice(index, 1)
    }
  } catch (error) {
    showNotification('Failed to complete week', 'error')
    console.error('Error completing week:', error)
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
        if (getRemainingBalance(client) === 0) {
          return processPayment(client)
        } else {
          return processFullPayment(client)
        }
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
    // Only select active clients (not paid ones)
    selectedClients.value = activeClients.value.map(client => client.id)
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
  if (status === 'paid') classes.push('paid-row')
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

// Watch for changes in active clients to update select all checkbox
watch(activeClients, (newActive) => {
  if (newActive.length > 0) {
    selectAll.value = selectedClients.value.length === newActive.length
  } else {
    selectAll.value = false
  }
})

onMounted(() => {
  fetchClients()
})
</script>

<style scoped>
/* Your existing CSS styles remain the same, just adding the new button style */

.btn-pay-full {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
  border: none;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-pay-full:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(59, 130, 246, 0.3);
}

.btn-pay-full:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.completion-notice {
  color: #059669;
  font-weight: 600;
  font-size: 0.875rem;
  text-align: center;
  width: 100%;
  display: block;
  margin-top: 0.5rem;
}

/* Rest of your existing CSS styles remain unchanged */
.payment-management {
  padding: 1.5rem;
  max-width: 1600px;
  margin: 0 auto;
}

/* Amortization Display Styles */
.amortization-display {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.full-amount {
  font-weight: 600;
  color: #059669;
}

.remaining-balance {
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  border-radius: 8px;
  font-weight: 500;
}

.balance-paid {
  background: #f0fdf4;
  color: #059669;
  border: 1px solid #bbf7d0;
}

.balance-low {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #fcd34d;
}

.balance-medium {
  background: #fef3c7;
  color: #b45309;
  border: 1px solid #f59e0b;
}

.balance-high {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fca5a5;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  min-width: 120px;
}

.btn-partial-pay {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: white;
  border: none;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-partial-pay:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(245, 158, 11, 0.3);
}

.btn-partial-pay:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

/* Partial Payment Indicator */
.partial-indicator {
  font-size: 0.7rem;
  color: #f59e0b;
  font-weight: 500;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  margin: 0;
  color: #1f2937;
  font-size: 1.25rem;
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6b7280;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.modal-close:hover {
  background: #f3f4f6;
  color: #374151;
}

.modal-body {
  padding: 1.5rem;
}

.client-info {
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: #f8fafc;
  border-radius: 8px;
}

.client-info h4 {
  margin: 0 0 0.5rem 0;
  color: #1f2937;
}

.client-info p {
  margin: 0;
  color: #6b7280;
  font-size: 0.875rem;
}

.payment-summary {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.summary-item:last-child {
  margin-bottom: 0;
}

.summary-item span:first-child {
  color: #374151;
  font-weight: 500;
}

.summary-item .amount {
  font-weight: 600;
  color: #059669;
  font-family: 'Courier New', monospace;
}

.payment-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 500;
  color: #374151;
  font-size: 0.875rem;
}

.amount-input {
  padding: 0.75rem;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  transition: border-color 0.2s;
  font-family: 'Courier New', monospace;
}

.amount-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.input-hint {
  font-size: 0.75rem;
  color: #6b7280;
}

.form-group select {
  padding: 0.75rem;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 0.875rem;
  background: white;
  cursor: pointer;
  transition: border-color 0.2s;
}

.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.payment-preview {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  padding: 1rem;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.preview-item:last-child {
  margin-bottom: 0;
}

.preview-item span:first-child {
  color: #374151;
  font-weight: 500;
}

.preview-item .amount {
  font-weight: 600;
  color: #1d4ed8;
  font-family: 'Courier New', monospace;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1rem;
}

.btn-cancel {
  padding: 0.75rem 1.5rem;
  border: 1px solid #d1d5db;
  background: white;
  color: #374151;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-cancel:hover {
  background: #f9fafb;
  border-color: #9ca3af;
}

.btn-confirm-partial {
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-confirm-partial:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(245, 158, 11, 0.3);
}

.btn-confirm-partial:disabled {
  background: #9ca3af;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* Responsive Design */
@media (max-width: 768px) {
  .action-buttons {
    flex-direction: row;
    flex-wrap: wrap;
  }
  
  .btn-partial-pay {
    font-size: 0.65rem;
    padding: 0.375rem 0.5rem;
  }
  
  .modal-content {
    margin: 1rem;
    width: calc(100% - 2rem);
  }
  
  .modal-actions {
    flex-direction: column;
  }
  
  .modal-actions button {
    width: 100%;
  }
}

/* Existing table styles remain unchanged */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.due-today-badge {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.badge-count {
  background: rgba(255, 255, 255, 0.2);
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.875rem;
}

.filters-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.search-box {
  position: relative;
  flex: 1;
  min-width: 300px;
}

.search-input {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 2.5rem;
  border: 2px solid #e5e7eb;
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
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  background: white;
  color: #6b7280;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-btn:hover {
  border-color: #9ca3af;
}

.filter-btn.active {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.filter-count {
  background: rgba(255, 255, 255, 0.2);
  padding: 0.125rem 0.375rem;
  border-radius: 8px;
  font-size: 0.75rem;
}

.table-section {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
  overflow: hidden;
}

.table-title {
  padding: 1.5rem 1.5rem 0.5rem;
  margin: 0;
  color: #1f2937;
  font-size: 1.25rem;
}

.table-subtitle {
  padding: 0 1.5rem 1.5rem;
  margin: 0;
  color: #6b7280;
  font-size: 0.875rem;
}

.table-container {
  overflow-x: auto;
}

.clients-table {
  width: 100%;
  border-collapse: collapse;
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

.checkbox-column {
  width: 40px;
  text-align: center;
}

.control-number {
  font-weight: 600;
  color: #1f2937;
}

.client-name {
  font-weight: 500;
  color: #1f2937;
}

.amount {
  font-weight: 600;
  color: #059669;
  font-family: 'Courier New', monospace;
}

.progress-cell {
  min-width: 150px;
}

.progress-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.progress-info {
  font-size: 0.75rem;
  color: #6b7280;
  display: flex;
  align-items: center;
  gap: 0.25rem;
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

.status-cell {
  white-space: nowrap;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  display: inline-block;
}

.status-due_today {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #fcd34d;
}

.status-overdue {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fca5a5;
}

.status-active {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.status-paid {
  background: #f0fdf4;
  color: #059669;
  border: 1px solid #bbf7d0;
}

.action-cell {
  white-space: nowrap;
}

.btn-pay {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-pay:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(16, 185, 129, 0.3);
}

.btn-pay:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.btn-pay.processing {
  background: #6b7280;
}

.paid-text {
  color: #059669;
  font-weight: 600;
  font-size: 0.75rem;
}

.due-today-row {
  background: #fffbeb !important;
  border-left: 4px solid #f59e0b;
}

.overdue-row {
  background: #fef2f2 !important;
  border-left: 4px solid #ef4444;
}

.selected-row {
  background: #eff6ff !important;
}

.empty-state, .loading-state {
  padding: 3rem;
  text-align: center;
  color: #6b7280;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.spinner {
  border: 3px solid #f3f4f6;
  border-top: 3px solid #3b82f6;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.notification {
  position: fixed;
  top: 1rem;
  right: 1rem;
  padding: 1rem 1.5rem;
  border-radius: 8px;
  color: white;
  font-weight: 500;
  z-index: 1001;
  animation: slideIn 0.3s ease;
}

.notification.success {
  background: linear-gradient(135deg, #10b981, #059669);
}

.notification.error {
  background: linear-gradient(135deg, #ef4444, #dc2626);
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

@media (max-width: 1024px) {
  .clients-table {
    min-width: 1200px;
  }
}

@media (max-width: 768px) {
  .payment-management {
    padding: 1rem;
  }
  
  .page-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .header-content {
    justify-content: space-between;
  }
  
  .filters-section {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-box {
    min-width: auto;
  }
  
  .filter-buttons {
    justify-content: center;
  }
}
</style>
