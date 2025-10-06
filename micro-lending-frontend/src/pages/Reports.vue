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

    <!-- Key Metrics Dashboard -->
    <div class="metrics-grid">
      <!-- Total Collections -->
      <div class="metric-card primary">
        <div class="metric-icon">💰</div>
        <div class="metric-content">
          <h3>Total Collections</h3>
          <div class="metric-value">₱{{ formatCurrency(summary.totalCollections) }}</div>
          <div class="metric-trend" :class="getTrendClass(summary.collectionTrend)">
            {{ getTrendIcon(summary.collectionTrend) }} {{ Math.abs(summary.collectionTrend) }}%
          </div>
        </div>
      </div>

      <!-- Total Loan Releases -->
      <div class="metric-card secondary">
        <div class="metric-icon">🚀</div>
        <div class="metric-content">
          <h3>Total Loan Releases</h3>
          <div class="metric-value">₱{{ formatCurrency(summary.totalReleases) }}</div>
          <div class="metric-trend" :class="getTrendClass(summary.releaseTrend)">
            {{ getTrendIcon(summary.releaseTrend) }} {{ Math.abs(summary.releaseTrend) }}%
          </div>
        </div>
      </div>

      <!-- Active Loans -->
      <div class="metric-card info">
        <div class="metric-icon">📈</div>
        <div class="metric-content">
          <h3>Active Loans</h3>
          <div class="metric-value">{{ summary.activeLoans }}</div>
          <div class="metric-subtitle">{{ summary.newLoansThisPeriod }} new this period</div>
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
      <div class="metric-card" :class="getEfficiencyClass(summary.collectionEfficiency)">
        <div class="metric-icon">📊</div>
        <div class="metric-content">
          <h3>Collection Efficiency</h3>
          <div class="metric-value">{{ summary.collectionEfficiency }}%</div>
          <div class="metric-subtitle">
            {{ summary.paymentsReceived }}/{{ summary.paymentsExpected }} payments
          </div>
        </div>
      </div>

      <!-- Average Loan Size -->
      <div class="metric-card success">
        <div class="metric-icon">📦</div>
        <div class="metric-content">
          <h3>Average Loan Size</h3>
          <div class="metric-value">₱{{ formatCurrency(summary.averageLoanSize) }}</div>
          <div class="metric-subtitle">Per client</div>
        </div>
      </div>
    </div>

    <!-- Charts Section -->
    <div class="charts-grid">
      <!-- Collections vs Releases Chart -->
      <div class="chart-card">
        <div class="chart-header">
          <h3>Collections vs Loan Releases</h3>
          <div class="chart-legend">
            <div class="legend-item">
              <span class="legend-color collections"></span>
              Collections
            </div>
            <div class="legend-item">
              <span class="legend-color releases"></span>
              Releases
            </div>
          </div>
        </div>
        <div class="chart-container">
          <canvas ref="collectionsChart"></canvas>
        </div>
      </div>

      <!-- Loan Status Distribution -->
      <div class="chart-card">
        <div class="chart-header">
          <h3>Loan Portfolio Distribution</h3>
        </div>
        <div class="chart-container">
          <canvas ref="loanStatusChart"></canvas>
        </div>
      </div>

      <!-- Weekly Collection Trend -->
      <div class="chart-card full-width">
        <div class="chart-header">
          <h3>Weekly Collection Trend</h3>
        </div>
        <div class="chart-container">
          <canvas ref="weeklyTrendChart"></canvas>
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
            @click="generateDetailedReport(report.id)"
            :class="['btn-report', { active: activeReport === report.id }]"
            :disabled="loading"
          >
            {{ report.label }}
          </button>
        </div>
      </div>

      <!-- Report Content -->
      <div class="report-content">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Generating report...</p>
        </div>

        <!-- Collection Summary Report -->
        <div v-if="activeReport === 'collections' && !loading" class="report-table">
          <h4>Collection Summary</h4>
          <table>
            <thead>
              <tr>
                <th>Week</th>
                <th>Date Range</th>
                <th>Expected</th>
                <th>Collected</th>
                <th>Efficiency</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="week in weeklyReports" :key="week.weekNumber">
                <td>Week {{ week.weekNumber }}</td>
                <td>{{ week.dateRange }}</td>
                <td>₱{{ formatCurrency(week.expectedAmount) }}</td>
                <td>₱{{ formatCurrency(week.collectedAmount) }}</td>
                <td>
                  <span :class="getEfficiencyBadgeClass(week.efficiency)">
                    {{ week.efficiency }}%
                  </span>
                </td>
                <td>
                  <span :class="getStatusBadgeClass(week.status)">
                    {{ week.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Loan Releases Report -->
        <div v-if="activeReport === 'releases' && !loading" class="report-table">
          <h4>Loan Releases Summary</h4>
          <table>
            <thead>
              <tr>
                <th>Client</th>
                <th>Control #</th>
                <th>Release Date</th>
                <th>Loan Amount</th>
                <th>Amount Released</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="release in loanReleases" :key="release.id">
                <td>{{ release.clientName }}</td>
                <td>{{ release.controlNumber }}</td>
                <td>{{ formatDate(release.releaseDate) }}</td>
                <td>₱{{ formatCurrency(release.loanAmount) }}</td>
                <td>₱{{ formatCurrency(release.amountReleased) }}</td>
                <td>
                  <span :class="getLoanStatusBadgeClass(release.status)">
                    {{ release.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Portfolio at Risk Report -->
        <div v-if="activeReport === 'risk' && !loading" class="report-table">
          <h4>Portfolio at Risk Analysis</h4>
          <table>
            <thead>
              <tr>
                <th>Risk Category</th>
                <th>Loan Count</th>
                <th>Amount at Risk</th>
                <th>Percentage</th>
                <th>Average Days Overdue</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="risk in riskAnalysis" :key="risk.category">
                <td>
                  <span :class="getRiskCategoryClass(risk.category)">
                    {{ risk.category }}
                  </span>
                </td>
                <td>{{ risk.loanCount }}</td>
                <td>₱{{ formatCurrency(risk.amountAtRisk) }}</td>
                <td>{{ risk.percentage }}%</td>
                <td>{{ risk.averageDaysOverdue }}</td>
              </tr>
            </tbody>
          </table>
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
        <span class="stat-value">{{ new Date().toLocaleDateString() }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">Total Clients:</span>
        <span class="stat-value">{{ summary.totalClients }}</span>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import Chart from 'chart.js/auto'

const router = useRouter()

// Reactive state
const loading = ref(false)
const dateRange = ref({
  start: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0],
  end: new Date().toISOString().split('T')[0]
})
const selectedPeriod = ref('month')
const activeReport = ref('collections')

// Chart references
const collectionsChart = ref(null)
const loanStatusChart = ref(null)
const weeklyTrendChart = ref(null)

// Report data
const summary = ref({
  totalCollections: 0,
  totalReleases: 0,
  activeLoans: 0,
  portfolioAtRisk: 0,
  collectionEfficiency: 0,
  averageLoanSize: 0,
  collectionTrend: 0,
  releaseTrend: 0,
  newLoansThisPeriod: 0,
  atRiskAmount: 0,
  paymentsReceived: 0,
  paymentsExpected: 0,
  totalClients: 0
})

const weeklyReports = ref([])
const loanReleases = ref([])
const riskAnalysis = ref([])

const availableReports = [
  { id: 'collections', label: 'Collection Summary' },
  { id: 'releases', label: 'Loan Releases' },
  { id: 'risk', label: 'Portfolio at Risk' }
]

// Methods
const fetchReportData = async () => {
  loading.value = true
  try {
    // Simulate API calls - replace with actual API endpoints
    await Promise.all([
      fetchSummaryData(),
      fetchWeeklyReports(),
      fetchLoanReleases(),
      fetchRiskAnalysis()
    ])
    
    // Initialize charts after data is loaded
    setTimeout(() => {
      initializeCharts()
    }, 100)
    
  } catch (error) {
    console.error('Error fetching report data:', error)
  } finally {
    loading.value = false
  }
}

const fetchSummaryData = async () => {
  // Mock data - replace with actual API call
  summary.value = {
    totalCollections: 125000,
    totalReleases: 180000,
    activeLoans: 45,
    portfolioAtRisk: 12.5,
    collectionEfficiency: 87.3,
    averageLoanSize: 7500,
    collectionTrend: 15.2,
    releaseTrend: 8.7,
    newLoansThisPeriod: 12,
    atRiskAmount: 22500,
    paymentsReceived: 183,
    paymentsExpected: 210,
    totalClients: 67
  }
}

const fetchWeeklyReports = async () => {
  // Mock data - replace with actual API call
  weeklyReports.value = [
    { weekNumber: 1, dateRange: 'Jan 1-7', expectedAmount: 25000, collectedAmount: 23000, efficiency: 92, status: 'Good' },
    { weekNumber: 2, dateRange: 'Jan 8-14', expectedAmount: 25500, collectedAmount: 24500, efficiency: 96, status: 'Excellent' },
    { weekNumber: 3, dateRange: 'Jan 15-21', expectedAmount: 26000, collectedAmount: 22000, efficiency: 85, status: 'Fair' },
    { weekNumber: 4, dateRange: 'Jan 22-28', expectedAmount: 26500, collectedAmount: 25000, efficiency: 94, status: 'Good' }
  ]
}

const fetchLoanReleases = async () => {
  // Mock data - replace with actual API call
  loanReleases.value = [
    { id: 1, clientName: 'Juan Dela Cruz', controlNumber: 'MLP-2024-001', releaseDate: '2024-01-15', loanAmount: 10000, amountReleased: 9500, status: 'Active' },
    { id: 2, clientName: 'Maria Santos', controlNumber: 'MLP-2024-002', releaseDate: '2024-01-18', loanAmount: 8000, amountReleased: 7600, status: 'Active' },
    { id: 3, clientName: 'Pedro Reyes', controlNumber: 'MLP-2024-003', releaseDate: '2024-01-22', loanAmount: 12000, amountReleased: 11400, status: 'Active' }
  ]
}

const fetchRiskAnalysis = async () => {
  // Mock data - replace with actual API call
  riskAnalysis.value = [
    { category: 'Current (0-7 days)', loanCount: 35, amountAtRisk: 0, percentage: 0, averageDaysOverdue: 0 },
    { category: 'Watch (8-30 days)', loanCount: 6, amountAtRisk: 8500, percentage: 4.7, averageDaysOverdue: 18 },
    { category: 'Substandard (31-90 days)', loanCount: 3, amountAtRisk: 12000, percentage: 6.7, averageDaysOverdue: 45 },
    { category: 'Doubtful (91-180 days)', loanCount: 1, amountAtRisk: 2000, percentage: 1.1, averageDaysOverdue: 120 }
  ]
}

const initializeCharts = () => {
  // Collections vs Releases Chart
  if (collectionsChart.value) {
    new Chart(collectionsChart.value, {
      type: 'bar',
      data: {
        labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
        datasets: [
          {
            label: 'Collections',
            data: [23000, 24500, 22000, 25000],
            backgroundColor: '#10b981',
            borderColor: '#059669',
            borderWidth: 1
          },
          {
            label: 'Releases',
            data: [30000, 25000, 35000, 28000],
            backgroundColor: '#3b82f6',
            borderColor: '#1d4ed8',
            borderWidth: 1
          }
        ]
      },
      options: {
        responsive: true,
        plugins: {
          legend: {
            position: 'top',
          },
          title: {
            display: true,
            text: 'Weekly Performance'
          }
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              callback: function(value) {
                return '₱' + value.toLocaleString()
              }
            }
          }
        }
      }
    })
  }

  // Loan Status Distribution Chart
  if (loanStatusChart.value) {
    new Chart(loanStatusChart.value, {
      type: 'doughnut',
      data: {
        labels: ['Active', 'Overdue', 'Paid', 'Pending'],
        datasets: [{
          data: [45, 8, 32, 5],
          backgroundColor: [
            '#3b82f6',
            '#ef4444',
            '#10b981',
            '#f59e0b'
          ],
          borderWidth: 2,
          borderColor: '#ffffff'
        }]
      },
      options: {
        responsive: true,
        plugins: {
          legend: {
            position: 'bottom'
          }
        }
      }
    })
  }

  // Weekly Collection Trend Chart
  if (weeklyTrendChart.value) {
    new Chart(weeklyTrendChart.value, {
      type: 'line',
      data: {
        labels: ['Jan W1', 'Jan W2', 'Jan W3', 'Jan W4', 'Feb W1', 'Feb W2'],
        datasets: [{
          label: 'Collection Trend',
          data: [22000, 24500, 23000, 25000, 26500, 28000],
          borderColor: '#10b981',
          backgroundColor: 'rgba(16, 185, 129, 0.1)',
          borderWidth: 3,
          fill: true,
          tension: 0.4
        }]
      },
      options: {
        responsive: true,
        plugins: {
          legend: {
            display: false
          }
        },
        scales: {
          y: {
            beginAtZero: false,
            ticks: {
              callback: function(value) {
                return '₱' + value.toLocaleString()
              }
            }
          }
        }
      }
    })
  }
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

const generateDetailedReport = (reportId) => {
  activeReport.value = reportId
  // In a real app, you would fetch specific report data here
}

const exportToPDF = () => {
  // Implement PDF export logic
  alert('PDF export feature will be implemented soon!')
}

const exportToCSV = () => {
  // Implement CSV export logic
  alert('CSV export feature will be implemented soon!')
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-PH', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount || 0)
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString('en-PH')
}

const getTrendClass = (trend) => {
  return trend >= 0 ? 'trend-up' : 'trend-down'
}

const getTrendIcon = (trend) => {
  return trend >= 0 ? '📈' : '📉'
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

const getEfficiencyBadgeClass = (efficiency) => {
  if (efficiency >= 90) return 'badge-success'
  if (efficiency >= 80) return 'badge-warning'
  return 'badge-danger'
}

const getStatusBadgeClass = (status) => {
  const statusMap = {
    'Excellent': 'badge-success',
    'Good': 'badge-info',
    'Fair': 'badge-warning',
    'Poor': 'badge-danger'
  }
  return statusMap[status] || 'badge-info'
}

const getLoanStatusBadgeClass = (status) => {
  const statusMap = {
    'Active': 'badge-info',
    'Overdue': 'badge-danger',
    'Paid': 'badge-success',
    'Pending': 'badge-warning'
  }
  return statusMap[status] || 'badge-info'
}

const getRiskCategoryClass = (category) => {
  const categoryMap = {
    'Current (0-7 days)': 'risk-current',
    'Watch (8-30 days)': 'risk-watch',
    'Substandard (31-90 days)': 'risk-substandard',
    'Doubtful (91-180 days)': 'risk-doubtful',
    'Loss (>180 days)': 'risk-loss'
  }
  return categoryMap[category] || ''
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

.metric-trend {
  font-size: 0.75rem;
  font-weight: 600;
}

.trend-up {
  color: #10b981;
}

.trend-down {
  color: #ef4444;
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.chart-header h3 {
  margin: 0;
  font-size: 1.125rem;
  color: #1f2937;
}

.chart-legend {
  display: flex;
  gap: 1rem;
  font-size: 0.75rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.legend-color.collections {
  background: #10b981;
}

.legend-color.releases {
  background: #3b82f6;
}

.chart-container {
  height: 300px;
  position: relative;
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

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  color: #6b7280;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.report-table {
  overflow-x: auto;
}

.report-table h4 {
  margin: 0 0 1rem 0;
  color: #1f2937;
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
}

.report-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #f3f4f6;
}

.report-table tr:last-child td {
  border-bottom: none;
}

.badge-success {
  background: #d1fae5;
  color: #065f46;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge-warning {
  background: #fef3c7;
  color: #92400e;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge-danger {
  background: #fee2e2;
  color: #991b1b;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge-info {
  background: #dbeafe;
  color: #1e40af;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.risk-current { color: #10b981; font-weight: 600; }
.risk-watch { color: #f59e0b; font-weight: 600; }
.risk-substandard { color: #ef4444; font-weight: 600; }
.risk-doubtful { color: #dc2626; font-weight: 600; }
.risk-loss { color: #7f1d1d; font-weight: 600; }

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
}
</style>
