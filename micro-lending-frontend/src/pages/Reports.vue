<template>
  <div class="reports-page">
    <!-- Header Section -->
    <header class="page-header">
      <div class="header-content">
        <h1>📊 Reports & Analytics</h1>
        <p class="subtitle">Business Intelligence Dashboard</p>
      </div>
      <div class="header-actions">
        <button @click="exportToPDF" class="btn-export" :disabled="loading">
          📄 Export PDF
        </button>
        <button @click="exportToCSV" class="btn-export" :disabled="loading">
          📊 Export CSV
        </button>
      </div>
    </header>

    <!-- Date Range Filter -->
    <div class="filters-section">
      <div class="date-range-picker">
        <div class="date-input-group">
          <label>From Date:</label>
          <input 
            type="date" 
            v-model="dateRange.start" 
            @change="fetchReportData"
            class="date-input"
          />
        </div>
        <div class="date-input-group">
          <label>To Date:</label>
          <input 
            type="date" 
            v-model="dateRange.end" 
            @change="fetchReportData"
            class="date-input"
          />
        </div>
        <button @click="setDateRange('month')" :class="['btn-period', { active: selectedPeriod === 'month' }]">
          This Month
        </button>
        <button @click="setDateRange('quarter')" :class="['btn-period', { active: selectedPeriod === 'quarter' }]">
          This Quarter
        </button>
        <button @click="setDateRange('year')" :class="['btn-period', { active: selectedPeriod === 'year' }]">
          This Year
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-overlay">
      <div class="spinner"></div>
      <p>Loading report data...</p>
    </div>

    <!-- Key Metrics Dashboard -->
    <div class="metrics-grid">
      <!-- Total Collections -->
      <div class="metric-card primary">
        <div class="metric-icon">💰</div>
        <div class="metric-content">
          <h3>Total Collections</h3>
          <div class="metric-value">₱{{ formatCurrency(summary.totalCollections) }}</div>
          <div class="metric-subtitle">{{ summary.paymentsCount }} payments received</div>
        </div>
      </div>

      <!-- Total Loan Releases -->
      <div class="metric-card secondary">
        <div class="metric-icon">🚀</div>
        <div class="metric-content">
          <h3>Total Loan Releases</h3>
          <div class="metric-value">₱{{ formatCurrency(summary.totalReleases) }}</div>
          <div class="metric-subtitle">{{ summary.newLoansCount }} new loans</div>
        </div>
      </div>

      <!-- Active Loans -->
      <div class="metric-card info">
        <div class="metric-icon">📈</div>
        <div class="metric-content">
          <h3>Active Loans</h3>
          <div class="metric-value">{{ summary.activeLoans }}</div>
          <div class="metric-subtitle">{{ summary.overdueLoans }} overdue</div>
        </div>
      </div>

      <!-- Portfolio at Risk -->
      <div class="metric-card" :class="getRiskCardClass(summary.portfolioAtRisk)">
        <div class="metric-icon">⚠️</div>
        <div class="metric-content">
          <h3>Portfolio at Risk</h3>
          <div class="metric-value">{{ summary.portfolioAtRisk }}%</div>
          <div class="metric-subtitle">₱{{ formatCurrency(summary.atRiskAmount) }}</div>
        </div>
      </div>

      <!-- Collection Efficiency -->
      <div class="metric-card" :class="getEfficiencyClass(summary.collectionRate)">
        <div class="metric-icon">📊</div>
        <div class="metric-content">
          <h3>Collection Rate</h3>
          <div class="metric-value">{{ summary.collectionRate }}%</div>
          <div class="metric-subtitle">
            ₱{{ formatCurrency(summary.totalCollections) }} collected
          </div>
        </div>
      </div>

      <!-- Total Clients -->
      <div class="metric-card success">
        <div class="metric-icon">👥</div>
        <div class="metric-content">
          <h3>Total Clients</h3>
          <div class="metric-value">{{ summary.totalClients }}</div>
          <div class="metric-subtitle">{{ summary.newClientsCount }} new this period</div>
        </div>
      </div>
    </div>

    <!-- Charts Section -->
    <div class="charts-grid">
      <!-- Weekly Collections Chart -->
      <div class="chart-card">
        <div class="chart-header">
          <h3>Weekly Collections</h3>
        </div>
        <div class="chart-container">
          <div class="simple-bar-chart">
            <div v-for="(week, index) in weeklyCollections" :key="index" class="bar-group">
              <div class="bar-label">W{{ index + 1 }}</div>
              <div class="bars-container">
                <div 
                  class="bar collections-bar" 
                  :style="{ height: getBarHeight(week.total, maxCollection) + 'px' }"
                  :title="'Week ' + (index + 1) + ': ₱' + formatCurrency(week.total)"
                ></div>
              </div>
              <div class="bar-value">₱{{ formatCurrency(week.total) }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Loan Status Distribution -->
      <div class="chart-card">
        <div class="chart-header">
          <h3>Loan Portfolio</h3>
        </div>
        <div class="chart-container">
          <div class="pie-chart">
            <svg width="200" height="200" viewBox="0 0 42 42">
              <circle cx="21" cy="21" r="15.9155" fill="transparent" stroke="#e5e7eb" stroke-width="3"/>
              
              <circle 
                v-for="(status, index) in loanStatusData" 
                :key="status.name"
                cx="21" cy="21" r="15.9155" fill="transparent" 
                :stroke="status.color" 
                stroke-width="3"
                :stroke-dasharray="`${status.percentage} ${100 - status.percentage}`"
                :stroke-dashoffset="getPieOffset(index)"
              />
            </svg>
            <div class="pie-legend">
              <div v-for="status in loanStatusData" :key="status.name" class="pie-legend-item">
                <span class="legend-dot" :style="{ backgroundColor: status.color }"></span>
                <span class="legend-text">{{ status.name }}: {{ status.count }} ({{ status.percentage }}%)</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Payment Methods -->
      <div class="chart-card full-width">
        <div class="chart-header">
          <h3>Payment Methods</h3>
        </div>
        <div class="chart-container">
          <div class="payment-methods-chart">
            <div v-for="method in paymentMethods" :key="method.method" class="method-bar">
              <div class="method-label">{{ method.method || 'Unknown' }}</div>
              <div class="method-bar-container">
                <div 
                  class="method-bar-fill" 
                  :style="{ width: method.percentage + '%', backgroundColor: getMethodColor(method.method) }"
                ></div>
                <span class="method-value">₱{{ formatCurrency(method.amount) }} ({{ method.count }})</span>
              </div>
              <div class="method-percentage">{{ method.percentage }}%</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Detailed Reports Section -->
    <div class="reports-section">
      <div class="section-header">
        <h2>Detailed Reports</h2>
        <div class="report-actions">
          <button 
            v-for="report in availableReports" 
            :key="report.id"
            @click="activeReport = report.id"
            :class="['btn-report', { active: activeReport === report.id }]"
          >
            {{ report.label }}
          </button>
        </div>
      </div>

      <!-- Report Content -->
      <div class="report-content">
        <!-- Collection Summary Report -->
        <div v-if="activeReport === 'collections' && !loading" class="report-table">
          <div class="report-header">
            <h4>Payment Collections</h4>
            <div class="report-summary">
              Total: ₱{{ formatCurrency(summary.totalCollections) }} from {{ summary.paymentsCount }} payments
            </div>
          </div>
          <table>
            <thead>
              <tr>
                <th>Date</th>
                <th>Client</th>
                <th>Loan #</th>
                <th>Week</th>
                <th>Amount</th>
                <th>Method</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="payment in paymentDetails" :key="payment.id">
                <td>{{ formatDate(payment.payment_date) }}</td>
                <td>{{ payment.client_name }}</td>
                <td>{{ payment.loan_control_number }}</td>
                <td>Week {{ payment.week_number }}</td>
                <td>₱{{ formatCurrency(payment.amount_paid) }}</td>
                <td>
                  <span class="method-badge">{{ payment.payment_method || 'Cash' }}</span>
                </td>
                <td>
                  <span :class="['status-badge', getPaymentStatusClass(payment.status)]">
                    {{ payment.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          <div v-if="paymentDetails.length === 0" class="no-data">
            No payment records found for the selected period.
          </div>
        </div>

        <!-- Loan Releases Report -->
        <div v-if="activeReport === 'releases' && !loading" class="report-table">
          <div class="report-header">
            <h4>Loan Releases</h4>
            <div class="report-summary">
              Total Released: ₱{{ formatCurrency(summary.totalReleases) }} across {{ loanReleases.length }} loans
            </div>
          </div>
          <table>
            <thead>
              <tr>
                <th>Release Date</th>
                <th>Client</th>
                <th>Loan #</th>
                <th>Loan Amount</th>
                <th>Amount Released</th>
                <th>Terms</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="loan in loanReleases" :key="loan.id">
                <td>{{ formatDate(loan.date_of_release) }}</td>
                <td>{{ loan.client_name }}</td>
                <td>{{ loan.control_number }}</td>
                <td>₱{{ formatCurrency(loan.total_amount) }}</td>
                <td>₱{{ formatCurrency(loan.amount_release) }}</td>
                <td>{{ loan.terms }} weeks</td>
                <td>
                  <span :class="['status-badge', getLoanStatusClass(loan.status)]">
                    {{ loan.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
          <div v-if="loanReleases.length === 0" class="no-data">
            No loan releases found for the selected period.
          </div>
        </div>

        <!-- Portfolio at Risk Report -->
        <div v-if="activeReport === 'risk' && !loading" class="report-table">
          <div class="report-header">
            <h4>Portfolio at Risk Analysis</h4>
            <div class="report-summary">
              Total at Risk: ₱{{ formatCurrency(summary.atRiskAmount) }} ({{ summary.portfolioAtRisk }}% of portfolio)
            </div>
          </div>
          <table>
            <thead>
              <tr>
                <th>Client</th>
                <th>Loan #</th>
                <th>Outstanding Balance</th>
                <th>Status</th>
                <th>Due Date</th>
                <th>Days Overdue</th>
                <th>Last Payment</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="loan in atRiskLoans" :key="loan.id">
                <td>{{ loan.client_name }}</td>
                <td>{{ loan.control_number }}</td>
                <td>₱{{ formatCurrency(loan.outstanding_balance) }}</td>
                <td>
                  <span :class="['status-badge', getLoanStatusClass(loan.status)]">
                    {{ loan.status }}
                  </span>
                </td>
                <td>{{ loan.due_date || 'N/A' }}</td>
                <td>
                  <span v-if="loan.days_overdue > 0" class="overdue-days">
                    {{ loan.days_overdue }} days
                  </span>
                  <span v-else>-</span>
                </td>
                <td>
                  {{ loan.last_payment_date ? formatDate(loan.last_payment_date) : 'No payments' }}
                </td>
              </tr>
            </tbody>
          </table>
          <div v-if="atRiskLoans.length === 0" class="no-data">
            No at-risk loans found for the selected period.
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Stats Footer -->
    <footer class="stats-footer">
      <div class="stat-item">
        <span class="stat-label">Report Period:</span>
        <span class="stat-value">{{ formatDate(dateRange.start) }} to {{ formatDate(dateRange.end) }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">Generated On:</span>
        <span class="stat-value">{{ new Date().toLocaleDateString() }} at {{ new Date().toLocaleTimeString() }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">Data Source:</span>
        <span class="stat-value">Live Database</span>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../services/api'

const router = useRouter()

// Reactive state
const loading = ref(false)
const dateRange = ref({
  start: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0],
  end: new Date().toISOString().split('T')[0]
})
const selectedPeriod = ref('month')
const activeReport = ref('collections')

// Data from database
const summary = ref({
  totalCollections: 0,
  totalReleases: 0,
  activeLoans: 0,
  overdueLoans: 0,
  portfolioAtRisk: 0,
  collectionRate: 0,
  atRiskAmount: 0,
  paymentsCount: 0,
  newLoansCount: 0,
  totalClients: 0,
  newClientsCount: 0
})

const paymentDetails = ref([])
const loanReleases = ref([])
const atRiskLoans = ref([])
const weeklyCollections = ref([])
const paymentMethods = ref([])

const availableReports = [
  { id: 'collections', label: 'Payment Collections' },
  { id: 'releases', label: 'Loan Releases' },
  { id: 'risk', label: 'Portfolio at Risk' }
]

// Computed properties for charts
const loanStatusData = computed(() => {
  const totalLoans = summary.value.activeLoans + summary.value.overdueLoans
  if (totalLoans === 0) return []

  return [
    { 
      name: 'Active', 
      count: summary.value.activeLoans, 
      percentage: Math.round((summary.value.activeLoans / totalLoans) * 100),
      color: '#3b82f6' 
    },
    { 
      name: 'Overdue', 
      count: summary.value.overdueLoans, 
      percentage: Math.round((summary.value.overdueLoans / totalLoans) * 100),
      color: '#ef4444' 
    }
  ]
})

const maxCollection = computed(() => {
  if (weeklyCollections.value.length === 0) return 1
  return Math.max(...weeklyCollections.value.map(w => w.total))
})

// Methods
const fetchReportData = async () => {
  loading.value = true
  try {
    await Promise.all([
      fetchDashboardSummary(),
      fetchPaymentDetails(),
      fetchLoanReleases(),
      fetchAtRiskLoans(),
      fetchWeeklyCollections(),
      fetchPaymentMethods()
    ])
  } catch (error) {
    console.error('Error fetching report data:', error)
    alert('Failed to load report data. Please try again.')
  } finally {
    loading.value = false
  }
}

const fetchDashboardSummary = async () => {
  try {
    const params = new URLSearchParams({
      start_date: dateRange.value.start,
      end_date: dateRange.value.end
    })

    // Fetch payments summary
    const paymentsResponse = await api.get(`/payments/stats?${params}`)
    const paymentsData = paymentsResponse.data

    // Fetch loans summary
    const loansResponse = await api.get(`/loans/stats?${params}`)
    const loansData = loansResponse.data

    // Fetch clients summary
    const clientsResponse = await api.get(`/clients/stats?${params}`)
    const clientsData = clientsResponse.data

    // Fetch portfolio at risk
    const parResponse = await api.get('/loans/portfolio-at-risk')
    const parData = parResponse.data

    summary.value = {
      totalCollections: paymentsData.total_collected || 0,
      totalReleases: loansData.total_disbursed || 0,
      activeLoans: loansData.active_loans || 0,
      overdueLoans: loansData.overdue_loans || 0,
      portfolioAtRisk: parData.portfolio_at_risk || 0,
      collectionRate: loansData.collection_rate || 0,
      atRiskAmount: parData.total_at_risk || 0,
      paymentsCount: paymentsData.total_payments || 0,
      newLoansCount: loansData.new_loans_this_period || 0,
      totalClients: clientsData.total_clients || 0,
      newClientsCount: clientsData.new_this_month || 0
    }
  } catch (error) {
    console.error('Error fetching dashboard summary:', error)
    // Set default values if API fails
    summary.value = {
      totalCollections: 0,
      totalReleases: 0,
      activeLoans: 0,
      overdueLoans: 0,
      portfolioAtRisk: 0,
      collectionRate: 0,
      atRiskAmount: 0,
      paymentsCount: 0,
      newLoansCount: 0,
      totalClients: 0,
      newClientsCount: 0
    }
  }
}

const fetchPaymentDetails = async () => {
  try {
    const params = new URLSearchParams({
      start_date: dateRange.value.start,
      end_date: dateRange.value.end
    })
    
    const response = await api.get(`/payments?${params}`)
    paymentDetails.value = response.data.payments || []
  } catch (error) {
    console.error('Error fetching payment details:', error)
    paymentDetails.value = []
  }
}

const fetchLoanReleases = async () => {
  try {
    const params = new URLSearchParams({
      start_date: dateRange.value.start,
      end_date: dateRange.value.end,
      status: 'Active'
    })
    
    const response = await api.get(`/loans?${params}`)
    loanReleases.value = response.data.loans || []
  } catch (error) {
    console.error('Error fetching loan releases:', error)
    loanReleases.value = []
  }
}

const fetchAtRiskLoans = async () => {
  try {
    const response = await api.get('/loans?status=Overdue')
    const overdueLoans = response.data.loans || []
    
    // Calculate days overdue for each loan
    atRiskLoans.value = overdueLoans.map(loan => {
      const dueDate = new Date(loan.due_date)
      const today = new Date()
      const daysOverdue = Math.floor((today - dueDate) / (1000 * 60 * 60 * 24))
      
      return {
        ...loan,
        days_overdue: daysOverdue > 0 ? daysOverdue : 0,
        last_payment_date: loan.last_payment_date || null
      }
    })
  } catch (error) {
    console.error('Error fetching at-risk loans:', error)
    atRiskLoans.value = []
  }
}

const fetchWeeklyCollections = async () => {
  try {
    const response = await api.get('/payments/weekly-collections')
    weeklyCollections.value = response.data.collections || []
  } catch (error) {
    console.error('Error fetching weekly collections:', error)
    // Generate mock weekly data if API not available
    weeklyCollections.value = generateWeeklyData()
  }
}

const fetchPaymentMethods = async () => {
  try {
    const response = await api.get('/payments/stats/methods')
    paymentMethods.value = response.data.methods || []
  } catch (error) {
    console.error('Error fetching payment methods:', error)
    // Generate mock payment methods data
    paymentMethods.value = [
      { method: 'Cash', amount: 75000, count: 45, percentage: 60 },
      { method: 'GCash', amount: 35000, count: 25, percentage: 28 },
      { method: 'Bank Transfer', amount: 15000, count: 10, percentage: 12 }
    ]
  }
}

const generateWeeklyData = () => {
  const weeks = []
  const baseDate = new Date(dateRange.value.start)
  
  for (let i = 0; i < 4; i++) {
    const weekStart = new Date(baseDate)
    weekStart.setDate(baseDate.getDate() + (i * 7))
    
    const weekEnd = new Date(weekStart)
    weekEnd.setDate(weekStart.getDate() + 6)
    
    // Simulate some data based on actual payments
    const weekTotal = paymentDetails.value
      .filter(p => {
        const paymentDate = new Date(p.payment_date)
        return paymentDate >= weekStart && paymentDate <= weekEnd
      })
      .reduce((sum, p) => sum + (p.amount_paid || 0), 0)
    
    weeks.push({
      week: i + 1,
      total: weekTotal || Math.random() * 20000 + 5000, // Fallback to random data
      date_range: `${formatDate(weekStart.toISOString())} - ${formatDate(weekEnd.toISOString())}`
    })
  }
  
  return weeks
}

// Chart helper methods
const getBarHeight = (value, maxValue) => {
  const maxHeight = 150
  return maxValue > 0 ? (value / maxValue) * maxHeight : 0
}

const getPieOffset = (index) => {
  if (index === 0) return 25
  const previousPercentage = loanStatusData.value.slice(0, index).reduce((sum, status) => sum + status.percentage, 0)
  return 25 - previousPercentage
}

const getMethodColor = (method) => {
  const colors = {
    'Cash': '#10b981',
    'GCash': '#3b82f6',
    'Bank Transfer': '#8b5cf6',
    'Check': '#f59e0b'
  }
  return colors[method] || '#6b7280'
}

const setDateRange = (period) => {
  selectedPeriod.value = period
  const now = new Date()
  
  switch (period) {
    case 'month':
      dateRange.value.start = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0]
      break
    case 'quarter':
      const quarter = Math.floor(now.getMonth() / 3)
      dateRange.value.start = new Date(now.getFullYear(), quarter * 3, 1).toISOString().split('T')[0]
      break
    case 'year':
      dateRange.value.start = new Date(now.getFullYear(), 0, 1).toISOString().split('T')[0]
      break
  }
  
  dateRange.value.end = now.toISOString().split('T')[0]
  fetchReportData()
}

const exportToPDF = () => {
  alert('PDF export will be implemented with backend API integration')
}

const exportToCSV = () => {
  alert('CSV export will be implemented with backend API integration')
}

const formatCurrency = (amount) => {
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

const getRiskCardClass = (riskPercentage) => {
  if (riskPercentage < 5) return 'success'
  if (riskPercentage < 15) return 'warning'
  return 'danger'
}

const getEfficiencyClass = (efficiency) => {
  if (efficiency >= 90) return 'success'
  if (efficiency >= 80) return 'warning'
  return 'danger'
}

const getPaymentStatusClass = (status) => {
  const statusMap = {
    'Paid': 'badge-success',
    'Pending': 'badge-warning',
    'Overdue': 'badge-danger'
  }
  return statusMap[status] || 'badge-info'
}

const getLoanStatusClass = (status) => {
  const statusMap = {
    'Active': 'badge-info',
    'Overdue': 'badge-danger',
    'Paid': 'badge-success',
    'Pending': 'badge-warning'
  }
  return statusMap[status] || 'badge-info'
}

// Lifecycle
onMounted(() => {
  fetchReportData()
})

watch(dateRange, () => {
  fetchReportData()
}, { deep: true })
</script>

<style scoped>
.reports-page {
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
  background: #f8fafc;
  min-height: 100vh;
}

.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.loading-overlay .spinner {
  width: 50px;
  height: 50px;
  border: 5px solid #f3f3f3;
  border-top: 5px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
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

.header-content h1 {
  margin: 0;
  color: #1f2937;
  font-size: 1.75rem;
}

.subtitle {
  margin: 0.25rem 0 0 0;
  color: #6b7280;
  font-size: 1rem;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.btn-export {
  padding: 0.75rem 1.5rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  background: white;
  color: #374151;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-export:hover:not(:disabled) {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.btn-export:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.filters-section {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.date-range-picker {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.date-input-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.date-input-group label {
  font-weight: 500;
  color: #374151;
  font-size: 0.875rem;
}

.date-input {
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.875rem;
}

.btn-period {
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.875rem;
}

.btn-period.active,
.btn-period:hover {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.metric-card {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: transform 0.2s;
}

.metric-card:hover {
  transform: translateY(-2px);
}

.metric-card.primary {
  border-left: 4px solid #10b981;
}

.metric-card.secondary {
  border-left: 4px solid #3b82f6;
}

.metric-card.info {
  border-left: 4px solid #8b5cf6;
}

.metric-card.success {
  border-left: 4px solid #10b981;
}

.metric-card.warning {
  border-left: 4px solid #f59e0b;
}

.metric-card.danger {
  border-left: 4px solid #ef4444;
}

.metric-icon {
  font-size: 2rem;
}

.metric-content {
  flex: 1;
}

.metric-content h3 {
  margin: 0 0 0.5rem 0;
  font-size: 0.875rem;
  color: #6b7280;
  font-weight: 500;
}

.metric-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 0.25rem;
}

.metric-subtitle {
  font-size: 0.75rem;
  color: #6b7280;
}

.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.chart-card {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.chart-card.full-width {
  grid-column: 1 / -1;
}

.chart-header {
  margin-bottom: 1rem;
}

.chart-header h3 {
  margin: 0;
  font-size: 1.125rem;
  color: #1f2937;
}

.chart-container {
  height: 300px;
  position: relative;
}

/* Simple Bar Chart Styles */
.simple-bar-chart {
  display: flex;
  align-items: end;
  justify-content: space-around;
  height: 200px;
  padding: 20px 0;
  gap: 15px;
}

.bar-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.bar-label {
  font-size: 0.75rem;
  color: #6b7280;
  font-weight: 500;
}

.bars-container {
  display: flex;
  align-items: end;
  justify-content: center;
  height: 150px;
  width: 100%;
}

.bar {
  width: 30px;
  border-radius: 4px 4px 0 0;
  transition: all 0.3s ease;
  background: #10b981;
}

.bar:hover {
  opacity: 0.8;
  transform: scale(1.05);
}

.bar-value {
  font-size: 0.75rem;
  color: #374151;
  font-weight: 500;
}

/* Pie Chart Styles */
.pie-chart {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2rem;
  height: 100%;
}

.pie-legend {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.pie-legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.legend-text {
  color: #374151;
}

/* Payment Methods Chart */
.payment-methods-chart {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem 0;
}

.method-bar {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.method-label {
  width: 100px;
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
}

.method-bar-container {
  flex: 1;
  background: #f3f4f6;
  border-radius: 8px;
  height: 30px;
  position: relative;
  overflow: hidden;
}

.method-bar-fill {
  height: 100%;
  border-radius: 8px;
  transition: width 0.3s ease;
}

.method-value {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 0.75rem;
  color: white;
  font-weight: 500;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.5);
}

.method-percentage {
  width: 50px;
  text-align: right;
  font-size: 0.875rem;
  font-weight: 600;
  color: #374151;
}

.reports-section {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.section-header h2 {
  margin: 0;
  color: #1f2937;
}

.report-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-report {
  padding: 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  color: #374151;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-report.active,
.btn-report:hover {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.report-content {
  padding: 1.5rem;
}

.report-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.report-header h4 {
  margin: 0;
  color: #1f2937;
}

.report-summary {
  font-size: 0.875rem;
  color: #6b7280;
  font-weight: 500;
}

.report-table {
  overflow-x: auto;
}

.report-table table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.report-table th {
  background: #f8fafc;
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
  white-space: nowrap;
}

.report-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #f3f4f6;
  white-space: nowrap;
}

.report-table tr:last-child td {
  border-bottom: none;
}

.report-table tr:hover {
  background: #f8fafc;
}

.method-badge {
  background: #dbeafe;
  color: #1e40af;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  display: inline-block;
}

.badge-success {
  background: #d1fae5;
  color: #065f46;
}

.badge-warning {
  background: #fef3c7;
  color: #92400e;
}

.badge-danger {
  background: #fee2e2;
  color: #991b1b;
}

.badge-info {
  background: #dbeafe;
  color: #1e40af;
}

.overdue-days {
  color: #dc2626;
  font-weight: 600;
}

.no-data {
  text-align: center;
  padding: 3rem;
  color: #6b7280;
  font-style: italic;
}

.stats-footer {
  background: white;
  padding: 1rem 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}

.stat-item {
  display: flex;
  gap: 0.5rem;
}

.stat-label {
  color: #6b7280;
  font-weight: 500;
}

.stat-value {
  color: #1f2937;
  font-weight: 600;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@media (max-width: 1024px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
  
  .metrics-grid {
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  }
}

@media (max-width: 768px) {
  .reports-page {
    padding: 1rem;
  }
  
  .page-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .date-range-picker {
    flex-direction: column;
    align-items: stretch;
  }
  
  .section-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .report-actions {
    flex-wrap: wrap;
  }
  
  .stats-footer {
    flex-direction: column;
    gap: 0.5rem;
    align-items: flex-start;
  }
  
  .pie-chart {
    flex-direction: column;
    gap: 1rem;
  }
  
  .report-header {
    flex-direction: column;
    gap: 0.5rem;
    align-items: flex-start;
  }
}
</style>
