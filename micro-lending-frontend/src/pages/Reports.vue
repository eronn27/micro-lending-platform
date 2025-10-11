<template>
  <div class="reports-container">
    <!-- Header -->
    <header class="reports-header">
      <!-- Left Section: Title -->
      <div class="header-left">
        <div class="header-content">
          <h1>Reports & Analytics</h1>
          <p class="subtitle">Weekly performance metrics and key statistics</p>
        </div>
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

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading reports...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Failed to load reports</h3>
      <p>{{ error }}</p>
      <button @click="refreshData" class="btn-retry">Try Again</button>
    </div>

    <!-- Main Content -->
    <div v-if="!loading && !error" class="reports-content">
      <!-- Key Metrics Grid -->
      <div class="metrics-grid">
        <!-- Weekly Payments Card -->
        <div class="metric-card payments">
          <div class="metric-icon">💰</div>
          <div class="metric-content">
            <h3>Weekly Payments</h3>
            <div class="metric-value">₱{{ formatCurrency(reportData.weekly_payment_total) }}</div>
            <p class="metric-description">Total payments received this week</p>
          </div>
          <div class="metric-trend positive">
            <span>↑ Active</span>
          </div>
        </div>

        <!-- Weekly Releases Card -->
        <div class="metric-card releases">
          <div class="metric-icon">📤</div>
          <div class="metric-content">
            <h3>Weekly Releases</h3>
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

        <!-- Client Health Card -->
        <div class="metric-card health">
          <div class="metric-icon">📊</div>
          <div class="metric-content">
            <h3>Client Health</h3>
            <div class="metric-value">{{ getClientHealthPercentage() }}%</div>
            <p class="metric-description">Healthy vs total clients ratio</p>
          </div>
          <div class="metric-trend positive">
            <span>{{ reportData.active_clients }} Active</span>
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
                <span class="stat-label">Total This Week:</span>
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
                <span class="stat-label">Total This Week:</span>
                <span class="stat-value">₱{{ formatCurrency(reportData.weekly_release_total) }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Average per Loan:</span>
                <span class="stat-value">
                  ₱{{ formatCurrency(getAverageReleasePerLoan()) }}
                </span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Total Payment vs Release:</span>
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

          <!-- Health Indicators -->
          <div class="analytics-card">
            <h3>Health Indicators</h3>
            <div class="analytics-content">
              <div class="progress-item">
                <label>Active Rate</label>
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: getActivePercentage() + '%', backgroundColor: '#10b981' }"></div>
                </div>
                <span class="progress-label">{{ getActivePercentage() }}%</span>
              </div>
              <div class="progress-item">
                <label>Overdue Rate</label>
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: getOverduePercentage() + '%', backgroundColor: '#ef4444' }"></div>
                </div>
                <span class="progress-label">{{ getOverduePercentage() }}%</span>
              </div>
              <div class="progress-item">
                <label>Health Score</label>
                <div class="progress-bar">
                  <div class="progress-fill" :style="{ width: getClientHealthPercentage() + '%', backgroundColor: '#3b82f6' }"></div>
                </div>
                <span class="progress-label">{{ getClientHealthPercentage() }}%</span>
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
    const response = await reportService.getWeeklyReport()
    reportData.value = response.data
    lastUpdated.value = new Date().toLocaleTimeString('en-PH')
  } catch (err) {
    error.value = err.message || 'Failed to load reports'
    console.error('Error fetching reports:', err)
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  fetchReports()
  showNotification('Reports refreshed successfully')
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
    timestamp: new Date().toISOString(),
    weeklyPayments: reportData.value.weekly_payment_total,
    weeklyReleases: reportData.value.weekly_release_total,
    totalClients: reportData.value.total_clients,
    activeClients: reportData.value.active_clients,
    overdueClients: reportData.value.overdue_clients,
    clientHealth: getClientHealthPercentage()
  }

  const dataStr = JSON.stringify(reportContent, null, 2)
  const dataBlob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(dataBlob)
  const link = document.createElement('a')
  link.href = url
  link.download = `report-${new Date().toISOString().split('T')[0]}.json`
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
.reports-container {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
}

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
}
</style>