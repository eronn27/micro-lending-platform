<template>
  <div class="reports-container">
    <!-- Header -->
    <header class="reports-header">
      <!-- Left Section: Title -->
      <div class="header-content">
        <h1>Reports & Analytics</h1>
        <p class="subtitle">{{ getPeriodDescription() }}</p>
      </div>

      <!-- Right Section: Actions -->
      <div class="header-right">
        <div class="last-updated">
          Last updated: {{ lastUpdated }}
        </div>
        <div class="header-buttons">
          <button @click="goBack" class="back-btn">
            ← Back
          </button>
          <button @click="refreshData" :disabled="loading" class="btn-refresh">
            <span v-if="!loading">🔄 Refresh</span>
            <span v-else>Loading...</span>
          </button>
        </div>
      </div>
    </header>

    <!-- Period Toggle Buttons -->
    <div class="period-toggle">
      <button
        @click="setPeriod('weekly')"
        :class="['period-btn', { active: selectedPeriod === 'weekly' }]"
      >
        📅 Weekly
      </button>
      <button
        @click="setPeriod('monthly')"
        :class="['period-btn', { active: selectedPeriod === 'monthly' }]"
      >
        📆 Monthly
      </button>
      <button
        @click="toggleHistory"
        :class="['period-btn', { active: showHistory }]"
      >
        📊 History
      </button>
    </div>

    <!-- History View -->
    <div v-if="showHistory" class="history-view">
      <div class="history-header">
        <h2>Historical Metrics</h2>
        <div class="history-controls">
          <select v-model="historyPeriod" @change="loadHistoricalData" class="history-select">
            <option value="weekly">Weekly History</option>
            <option value="monthly">Monthly History</option>
          </select>
          <select v-model="historyRange" @change="loadHistoricalData" class="history-select">
            <option value="4">Last 4 Periods</option>
            <option value="8">Last 8 Periods</option>
            <option value="12">Last 12 Periods</option>
          </select>
        </div>
      </div>

      <!-- Historical Data Table -->
      <div class="history-table-container">
        <table class="history-table">
          <thead>
            <tr>
              <th>Period</th>
              <th>Total Payments</th>
              <th>Total Releases</th>
              <th>Active Clients</th>
              <th>Overdue Clients</th>
              <th>Net Flow</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(record, index) in historicalData" :key="index">
              <td class="period-cell">{{ record.period }}</td>
              <td class="amount-cell">₱{{ formatCurrency(record.payments) }}</td>
              <td class="amount-cell">₱{{ formatCurrency(record.releases) }}</td>
              <td class="number-cell">{{ record.activeClients }}</td>
              <td class="number-cell">{{ record.overdueClients }}</td>
              <td :class="['amount-cell', record.netFlow >= 0 ? 'positive' : 'negative']">
                ₱{{ formatCurrency(Math.abs(record.netFlow)) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- History Summary Cards -->
      <div class="history-summary">
        <div class="summary-card">
          <h4>Average Payments</h4>
          <p class="summary-value">₱{{ formatCurrency(getHistoryAverage('payments')) }}</p>
        </div>
        <div class="summary-card">
          <h4>Average Releases</h4>
          <p class="summary-value">₱{{ formatCurrency(getHistoryAverage('releases')) }}</p>
        </div>
        <div class="summary-card">
          <h4>Trend</h4>
          <p :class="['summary-value', getTrendClass()]">
            {{ getTrendText() }}
          </p>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading && !showHistory" class="loading-state">
      <div class="spinner"></div>
      <p>Loading reports...</p>
    </div>

    <!-- Error State -->
    <div v-if="error && !showHistory" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Failed to load reports</h3>
      <p>{{ error }}</p>
      <button @click="refreshData" class="btn-retry">Try Again</button>
    </div>

    <!-- Main Content (Current Period) -->
    <div v-if="!loading && !error && !showHistory" class="reports-content">
      <!-- Key Metrics Grid -->
      <div class="metrics-grid">
        <!-- Weekly Payments Card -->
        <div class="metric-card payments">
          <div class="metric-icon">💰</div>
          <div class="metric-content">
            <h3>{{ selectedPeriod === 'weekly' ? 'Weekly' : 'Monthly' }} Payments</h3>
            <div class="metric-value">₱{{ formatCurrency(reportData.weekly_payment_total) }}</div>
            <p class="metric-description">Total payments received</p>
          </div>
          <div class="metric-trend positive">
            <span>↑ Active</span>
          </div>
        </div>

        <!-- Weekly Releases Card -->
        <div class="metric-card releases">
          <div class="metric-icon">📤</div>
          <div class="metric-content">
            <h3>{{ selectedPeriod === 'weekly' ? 'Weekly' : 'Monthly' }} Releases</h3>
            <div class="metric-value">₱{{ formatCurrency(reportData.weekly_release_total) }}</div>
            <p class="metric-description">Total loan amounts released</p>
          </div>
          <div class="metric-trend positive">
            <span>↑ Growing</span>
          </div>
        </div>

        <!-- Total Clients Card -->
        <div class="metric-card clients">
          <div class="metric-icon">👥</div>
          <div class="metric-content">
            <h3>Total Clients</h3>
            <div class="metric-value">{{ reportData.total_clients }}</div>
            <p class="metric-description">Registered clients in system</p>
          </div>
          <div class="metric-trend info">
            <span>ℹ All</span>
          </div>
        </div>

        <!-- Active Clients Card -->
        <div class="metric-card active">
          <div class="metric-icon">✅</div>
          <div class="metric-content">
            <h3>Active Clients</h3>
            <div class="metric-value">{{ reportData.active_clients }}</div>
            <p class="metric-description">Clients with active loans</p>
          </div>
          <div class="metric-trend positive">
            <span>{{ getActivePercentage() }}%</span>
          </div>
        </div>

        <!-- Overdue Clients Card -->
        <div class="metric-card overdue">
          <div class="metric-icon">🚨</div>
          <div class="metric-content">
            <h3>Overdue Clients</h3>
            <div class="metric-value">{{ reportData.overdue_clients }}</div>
            <p class="metric-description">Clients with overdue payments</p>
          </div>
          <div class="metric-trend negative">
            <span>⚠️ Attention</span>
          </div>
        </div>

      </div>

      <!-- Detailed Analytics Section -->
      <div class="analytics-section">
        <h2>Detailed Analytics</h2>
        
        <div class="analytics-grid">
          <!-- Payment Summary -->
          <div class="analytics-card">
            <h3>Payment Summary</h3>
            <div class="analytics-content">
              <div class="stat-row">
                <span class="stat-label">Total This {{ selectedPeriod === 'weekly' ? 'Week' : 'Month' }}:</span>
                <span class="stat-value">₱{{ formatCurrency(reportData.weekly_payment_total) }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Average per Client:</span>
                <span class="stat-value">
                  ₱{{ formatCurrency(getAveragePaymentPerClient()) }}
                </span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Active Payment:</span>
                <span class="stat-value">₱{{ formatCurrency(reportData.active_payment_total) }}</span>
              </div>
            </div>
          </div>

          <!-- Release Summary -->
          <div class="analytics-card">
            <h3>Release Summary</h3>
            <div class="analytics-content">
              <div class="stat-row">
                <span class="stat-label">Total This {{ selectedPeriod === 'weekly' ? 'Week' : 'Month' }}:</span>
                <span class="stat-value">₱{{ formatCurrency(reportData.weekly_release_total) }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Average per Loan:</span>
                <span class="stat-value">
                  ₱{{ formatCurrency(getAverageReleasePerLoan()) }}
                </span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Payment vs Release Ratio:</span>
                <span class="stat-value" :class="getPaymentReleaseRatio() > 1 ? 'positive' : 'neutral'">
                  {{ getPaymentReleaseRatio().toFixed(2) }}x
                </span>
              </div>
            </div>
          </div>

          <!-- Client Statistics -->
          <div class="analytics-card">
            <h3>Client Statistics</h3>
            <div class="analytics-content">
              <div class="stat-row">
                <span class="stat-label">Total Clients:</span>
                <span class="stat-value">{{ reportData.total_clients }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Active Clients:</span>
                <span class="stat-value positive">{{ reportData.active_clients }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Overdue Clients:</span>
                <span class="stat-value negative">{{ reportData.overdue_clients }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Inactive Clients:</span>
                <span class="stat-value neutral">
                  {{ getInactiveClients() }}
                </span>
              </div>
            </div>
          </div>

        </div>
      </div>

      <!-- Export Section -->
      <div class="export-section">
        <button @click="exportReport" class="btn-export">
          📥 Export Report
        </button>
      </div>
    </div>

    <!-- Notification -->
    <div v-if="notification.show" :class="['notification', notification.type]">
      {{ notification.message }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { reportService } from '../services/reportService'

const router = useRouter()
const loading = ref(false)
const error = ref(null)
const lastUpdated = ref('-')
const notification = ref({ show: false, message: '', type: 'success' })
const selectedPeriod = ref('weekly')
const showHistory = ref(false)
const historyPeriod = ref('weekly')
const historyRange = ref('4')
const historicalData = ref([])

const reportData = ref({
  weekly_payment_total: 0,
  weekly_release_total: 0,
  total_clients: 0,
  active_clients: 0,
  overdue_clients: 0,
  active_payment_total: 0,
  total_payment_this_week: 0
})

// Back button function
const goBack = () => {
  router.push('/dashboard')
}

const fetchReports = async () => {
  loading.value = true
  error.value = null
  try {
    const endpoint = selectedPeriod.value === 'weekly' ? 'getWeeklyReport' : 'getMonthlyReport'
    const response = await reportService[endpoint]()
    reportData.value = response.data
    lastUpdated.value = new Date().toLocaleTimeString('en-PH')
  } catch (err) {
    error.value = err.message || 'Failed to load reports'
    console.error('Error fetching reports:', err)
  } finally {
    loading.value = false
  }
}

const setPeriod = (period) => {
  selectedPeriod.value = period
  showHistory.value = false
  fetchReports()
}

const toggleHistory = () => {
  showHistory.value = !showHistory.value
  if (showHistory.value) {
    loadHistoricalData()
  }
}

const loadHistoricalData = async () => {
  loading.value = true
  try {
    const response = await reportService.getHistoricalReport(
      historyPeriod.value,
      parseInt(historyRange.value)
    )
    
    // Map API response to component format
    historicalData.value = response.data.map(record => ({
      period: record.period,
      payments: record.payments,
      releases: record.releases,
      activeClients: record.active_clients,
      overdueClients: record.overdue_clients,
      netFlow: record.net_flow
    }))
    
    showNotification('Historical data loaded successfully')
  } catch (error) {
    console.error('Error loading historical data:', error)
    showNotification('Failed to load historical data', 'error')
    // Fallback to empty array on error
    historicalData.value = []
  } finally {
    loading.value = false
  }
}

const getHistoryAverage = (field) => {
  if (historicalData.value.length === 0) return 0
  const sum = historicalData.value.reduce((acc, record) => acc + record[field], 0)
  return sum / historicalData.value.length
}

const getTrendText = () => {
  if (historicalData.value.length < 2) return 'N/A'
  
  const recent = historicalData.value[historicalData.value.length - 1].payments
  const previous = historicalData.value[historicalData.value.length - 2].payments
  const change = ((recent - previous) / previous * 100).toFixed(1)
  
  if (change > 0) return `↑ ${change}%`
  if (change < 0) return `↓ ${Math.abs(change)}%`
  return '→ No Change'
}

const getTrendClass = () => {
  if (historicalData.value.length < 2) return 'neutral'
  
  const recent = historicalData.value[historicalData.value.length - 1].payments
  const previous = historicalData.value[historicalData.value.length - 2].payments
  
  if (recent > previous) return 'positive'
  if (recent < previous) return 'negative'
  return 'neutral'
}

const formatDateShort = (date) => {
  return date.toLocaleDateString('en-PH', { month: 'short', day: 'numeric' })
}

const getPeriodDescription = () => {
  if (showHistory) return 'Historical performance metrics'
  return selectedPeriod.value === 'weekly' 
    ? 'Weekly performance metrics and key statistics'
    : 'Monthly performance metrics and key statistics'
}

const refreshData = () => {
  if (showHistory) {
    loadHistoricalData()
    showNotification('History refreshed successfully')
  } else {
    fetchReports()
    showNotification('Reports refreshed successfully')
  }
}

const formatCurrency = (amount) => {
  if (!amount && amount !== 0) return '0.00'
  return new Intl.NumberFormat('en-PH', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount || 0)
}

const getActivePercentage = () => {
  if (reportData.value.total_clients === 0) return 0
  return Math.round(
    (reportData.value.active_clients / reportData.value.total_clients) * 100
  )
}

const getOverduePercentage = () => {
  if (reportData.value.total_clients === 0) return 0
  return Math.round(
    (reportData.value.overdue_clients / reportData.value.total_clients) * 100
  )
}

const getInactiveClients = () => {
  return Math.max(
    0,
    reportData.value.total_clients - 
    reportData.value.active_clients - 
    reportData.value.overdue_clients
  )
}

const getClientHealthPercentage = () => {
  if (reportData.value.total_clients === 0) return 0
  const healthy = reportData.value.total_clients - reportData.value.overdue_clients
  return Math.round((healthy / reportData.value.total_clients) * 100)
}

const getAveragePaymentPerClient = () => {
  if (reportData.value.active_clients === 0) return 0
  return reportData.value.weekly_payment_total / reportData.value.active_clients
}

const getAverageReleasePerLoan = () => {
  if (reportData.value.total_clients === 0) return 0
  return reportData.value.weekly_release_total / reportData.value.total_clients
}

const getPaymentReleaseRatio = () => {
  if (reportData.value.weekly_release_total === 0) return 0
  return reportData.value.weekly_payment_total / reportData.value.weekly_release_total
}

const exportReport = () => {
  const reportContent = {
    period: selectedPeriod.value,
    timestamp: new Date().toISOString(),
    data: reportData.value,
    calculations: {
      activePercentage: getActivePercentage(),
      overduePercentage: getOverduePercentage(),
      clientHealth: getClientHealthPercentage()
    }
  }

  const dataStr = JSON.stringify(reportContent, null, 2)
  const dataBlob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(dataBlob)
  const link = document.createElement('a')
  link.href = url
  link.download = `report-${selectedPeriod.value}-${new Date().toISOString().split('T')[0]}.json`
  link.click()
  showNotification('Report exported successfully')
}

const showNotification = (message, type = 'success') => {
  notification.value = { show: true, message, type }
  setTimeout(() => {
    notification.value.show = false
  }, 3000)
}

onMounted(() => {
  fetchReports()
})
</script>

<style scoped>
/* Existing styles remain the same... */
.reports-container {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
}

/* Period Toggle Buttons */
.period-toggle {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  background: white;
  padding: 1rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
  justify-content: center;
}

.period-btn {
  padding: 0.75rem 2rem;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  background: white;
  color: #6b7280;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.period-btn:hover {
  border-color: #3b82f6;
  color: #3b82f6;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.2);
}

.period-btn.active {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
}

/* History View */
.history-view {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
  margin-bottom: 2rem;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #e5e7eb;
}

.history-header h2 {
  margin: 0;
  color: #1f2937;
  font-size: 1.5rem;
}

.history-controls {
  display: flex;
  gap: 1rem;
}

.history-select {
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: border-color 0.2s;
}

.history-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.history-table-container {
  overflow-x: auto;
  margin-bottom: 2rem;
}

.history-table {
  width: 100%;
  border-collapse: collapse;
}

.history-table thead {
  background: #f8fafc;
}

.history-table th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #374151;
  font-size: 0.875rem;
  border-bottom: 2px solid #e5e7eb;
}

.history-table td {
  padding: 1rem;
  border-bottom: 1px solid #f3f4f6;
}

.period-cell {
  font-weight: 600;
  color: #1f2937;
}

.amount-cell {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  color: #059669;
}

.amount-cell.positive {
  color: #059669;
}

.amount-cell.negative {
  color: #dc2626;
}

.number-cell {
  text-align: center;
  font-weight: 500;
  color: #1f2937;
}

.history-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.summary-card {
  background: #f8fafc;
  padding: 1.5rem;
  border-radius: 8px;
  text-align: center;
}

.summary-card h4 {
  margin: 0 0 0.5rem 0;
  color: #6b7280;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.summary-value {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  font-family: 'Courier New', monospace;
}

.summary-value.positive {
  color: #059669;
}

.summary-value.negative {
  color: #dc2626;
}

.summary-value.neutral {
  color: #6b7280;
}

/* Rest of existing styles... */
.reports-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
  gap: 2rem;
}

.header-left {
  display: flex;
  align-items: flex-start;
  flex: 1;
}

.header-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1rem;
}

.header-buttons {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.back-btn {
  padding: 0.75rem 1.25rem;
  background: #6b7280;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.back-btn:hover {
  background: #4b5563;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.header-content h1 {
  margin: 0 0 0.5rem 0;
  font-size: 2rem;
  color: #1f2937;
}

.subtitle {
  margin: 0;
  color: #6b7280;
  font-size: 1rem;
}

.btn-refresh {
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.btn-refresh:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(59, 130, 246, 0.3);
}

.btn-refresh:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.last-updated {
  font-size: 0.875rem;
  color: #6b7280;
  padding: 0.5rem 1rem;
  background: #f3f4f6;
  border-radius: 6px;
  white-space: nowrap;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 6rem 2rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
}

.spinner {
  width: 50px;
  height: 50px;
  border: 5px solid #f3f4f6;
  border-top: 5px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
  text-align: center;
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.error-state h3 {
  margin: 0 0 0.5rem 0;
  color: #ef4444;
  font-size: 1.25rem;
}

.error-state p {
  margin: 0 0 1.5rem 0;
  color: #6b7280;
}

.btn-retry {
  padding: 0.75rem 1.5rem;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-retry:hover {
  background: #dc2626;
  transform: translateY(-2px);
}

.reports-content {
  space-y: 2rem;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.metric-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  border-left: 5px solid #3b82f6;
}

.metric-card::before {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  width: 100px;
  height: 100px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  border-radius: 50%;
  transform: translate(30%, -30%);
}

.metric-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
}

.metric-card.payments {
  border-left-color: #10b981;
}

.metric-card.releases {
  border-left-color: #f59e0b;
}

.metric-card.clients {
  border-left-color: #3b82f6;
}

.metric-card.active {
  border-left-color: #8b5cf6;
}

.metric-card.overdue {
  border-left-color: #ef4444;
}

.metric-card.health {
  border-left-color: #06b6d4;
}

.metric-icon {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.metric-content h3 {
  margin: 0;
  color: #6b7280;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 600;
}

.metric-value {
  margin: 0.5rem 0;
  font-size: 1.875rem;
  font-weight: 700;
  color: #1f2937;
  font-family: 'Courier New', monospace;
}

.metric-description {
  margin: 0;
  color: #9ca3af;
  font-size: 0.875rem;
}

.metric-trend {
  margin-top: 1rem;
  display: inline-block;
  padding: 0.375rem 0.75rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}

.metric-trend.positive {
  background: #d1fae5;
  color: #065f46;
}

.metric-trend.negative {
  background: #fee2e2;
  color: #7f1d1d;
}

.metric-trend.info {
  background: #dbeafe;
  color: #1e3a8a;
}

.analytics-section {
  margin-top: 2rem;
}

.analytics-section h2 {
  margin: 0 0 1.5rem 0;
  color: #1f2937;
  font-size: 1.5rem;
}

.analytics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 1.5rem;
}

.analytics-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
}

.analytics-card h3 {
  margin: 0 0 1rem 0;
  color: #1f2937;
  font-size: 1.125rem;
  font-weight: 600;
  border-bottom: 2px solid #f3f4f6;
  padding-bottom: 0.75rem;
}

.analytics-content {
  space-y: 0.75rem;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 0;
  border-bottom: 1px solid #f3f4f6;
}

.stat-row:last-child {
  border-bottom: none;
}

.stat-label {
  color: #6b7280;
  font-size: 0.875rem;
  font-weight: 500;
}

.stat-value {
  font-weight: 700;
  color: #1f2937;
  font-family: 'Courier New', monospace;
}

.stat-value.positive {
  color: #059669;
}

.stat-value.negative {
  color: #dc2626;
}

.stat-value.neutral {
  color: #6b7280;
}

.progress-item {
  margin-bottom: 1.5rem;
}

.progress-item:last-child {
  margin-bottom: 0;
}

.progress-item label {
  display: block;
  margin-bottom: 0.5rem;
  color: #6b7280;
  font-size: 0.875rem;
  font-weight: 500;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-label {
  display: block;
  text-align: right;
  color: #1f2937;
  font-size: 0.875rem;
  font-weight: 600;
}

.export-section {
  margin-top: 2rem;
  display: flex;
  justify-content: flex-end;
}

.btn-export {
  padding: 0.875rem 1.75rem;
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-export:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(16, 185, 129, 0.3);
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

@media (max-width: 1024px) {
  .metrics-grid {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  }

  .analytics-grid {
    grid-template-columns: 1fr;
  }
  
  .period-toggle {
    flex-wrap: wrap;
  }
  
  .history-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .history-controls {
    width: 100%;
    flex-direction: column;
  }
  
  .history-select {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .reports-container {
    padding: 1rem;
  }

  .reports-header {
    flex-direction: column;
    padding: 1.5rem;
    gap: 1.5rem;
  }

  .header-left {
    width: 100%;
  }

  .header-right {
    width: 100%;
    align-items: stretch;
  }

  .header-buttons {
    justify-content: space-between;
  }

  .last-updated {
    text-align: center;
  }

  .metrics-grid {
    grid-template-columns: 1fr;
  }

  .metric-card {
    padding: 1.25rem;
  }

  .metric-value {
    font-size: 1.5rem;
  }
  
  .period-toggle {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .period-btn {
    width: 100%;
    justify-content: center;
  }
  
  .history-table {
    font-size: 0.75rem;
  }
  
  .history-table th,
  .history-table td {
    padding: 0.5rem;
  }
}
</style>