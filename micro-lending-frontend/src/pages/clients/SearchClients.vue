<template>
  <div class="search-clients">
    <!-- Header with Logo and Title -->
    <header class="page-header">
      <div class="logo-section">
        <div class="logo-placeholder">
          <span class="logo-text">MLP</span>
        </div>
        <h1>Client Search</h1>
      </div>
      <div class="header-actions">
        <button @click="goBack" class="back-btn">
          ← Back to Dashboard
        </button>
      </div>
    </header>

    <!-- Search Section -->
    <div class="search-section">
      <div class="search-container">
        <div class="search-input-group">
          <div class="search-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="m21 21-4.3-4.3"></path>
            </svg>
          </div>
          <input
            type="text"
            v-model="searchQuery"
            @input="handleSearch"
            placeholder="Search clients by name, control number, contact number, or address..."
            class="search-input"
          />
          <button 
            @click="clearSearch" 
            class="clear-btn"
            v-if="searchQuery"
          >
            Clear
          </button>
        </div>

        <!-- Search Filters -->
        <div class="search-filters">
          <div class="filter-group">
            <label>Status:</label>
            <select v-model="filters.status" @change="applyFilters">
              <option value="">All Status</option>
              <option value="Active">Active</option>
              <option value="Paid">Paid</option>
              <option value="Overdue">Overdue</option>
              <option value="Default">Default</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Mode:</label>
            <select v-model="filters.mode" @change="applyFilters">
              <option value="">All Modes</option>
              <option value="Weekly">Weekly</option>
              <option value="Monthly">Monthly</option>
            </select>
          </div>

          <div class="filter-group">
            <label>Due Date:</label>
            <select v-model="filters.dueDate" @change="applyFilters">
              <option value="">All Days</option>
              <option value="Monday">Monday</option>
              <option value="Tuesday">Tuesday</option>
              <option value="Wednesday">Wednesday</option>
              <option value="Thursday">Thursday</option>
              <option value="Friday">Friday</option>
              <option value="Saturday">Saturday</option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <!-- Results Section -->
    <div class="results-section">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Searching clients...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-state">
        <div class="error-icon">⚠️</div>
        <p>{{ error }}</p>
        <button @click="loadClients" class="retry-btn">Retry</button>
      </div>

      <!-- Results -->
      <div v-else class="results-container">
        <!-- Results Header -->
        <div class="results-header">
          <div class="results-info">
            <h3>Client Results</h3>
            <span class="results-count">
              {{ totalClients }} client{{ totalClients !== 1 ? 's' : '' }} found
            </span>
          </div>
          
          <div class="results-actions">
            <button 
              @click="exportClients" 
              class="export-btn"
              :disabled="clients.length === 0"
            >
              Export CSV
            </button>
          </div>
        </div>

        <!-- Clients Table -->
        <div class="table-container">
          <table class="clients-table">
            <thead>
              <tr>
                <th>Control #</th>
                <th>Name</th>
                <th>Date of Release</th>
                <th>Total Amount</th>
                <th>Ammortization</th>
                <th>Terms</th>
                <th>Mode</th>
                <th>Outstanding Balance</th>
                <th>Status</th>
                <th>Due Date</th>
                <th>Deductions</th>
                <th>Number</th>
                <th>Address</th>
                <th>Amount Release</th>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="client in clients" 
                :key="client.id"
                @click="viewClientDetails(client)"
                class="client-row"
              >
                <td class="control-number">
                  {{ client.control_number || 'N/A' }}
                </td>
                <td class="client-name">
                  {{ client.first_name }} {{ client.last_name }}
                </td>
                <td class="date-release">
                  {{ formatDate(client.loans?.[0]?.date_of_release) }}
                </td>
                <td class="total-amount">
                  {{ formatCurrency(client.loans?.[0]?.total_amount) }}
                </td>
                <td class="ammortization">
                  {{ formatCurrency(client.loans?.[0]?.ammortization) }}
                </td>
                <td class="terms">
                  {{ client.loans?.[0]?.terms || 'N/A' }}
                </td>
                <td class="mode">
                  <span class="mode-badge">{{ client.loans?.[0]?.mode || 'N/A' }}</span>
                </td>
                <td class="outstanding-balance">
                  {{ formatCurrency(client.loans?.[0]?.outstanding_balance) }}
                </td>
                <td class="status">
                  <span :class="['status-badge', getStatusClass(client.loans?.[0]?.status)]">
                    {{ client.loans?.[0]?.status || 'N/A' }}
                  </span>
                </td>
                <td class="due-date">
                  {{ client.loans?.[0]?.due_date || 'N/A' }}
                </td>
                <td class="deductions">
                  {{ client.loans?.[0]?.deductions || 'N/A' }}
                </td>
                <td class="contact-number">
                  {{ client.contact_number || 'N/A' }}
                </td>
                <td class="address">
                  {{ truncateAddress(client.home_address) }}
                </td>
                <td class="amount-release">
                  {{ formatCurrency(client.loans?.[0]?.amount_release) }}
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Empty State -->
          <div v-if="clients.length === 0 && !loading" class="empty-state">
            <div class="empty-icon">📋</div>
            <h3>No clients found</h3>
            <p v-if="searchQuery || hasActiveFilters">
              Try adjusting your search criteria or filters
            </p>
            <p v-else>
              No clients have been added yet
            </p>
            <button @click="goToCreateClient" class="create-client-btn">
              Create New Client
            </button>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="totalClients > 0" class="pagination">
          <div class="pagination-info">
            Showing {{ startItem }}-{{ endItem }} of {{ totalClients }} clients
          </div>
          <div class="pagination-controls">
            <button 
              @click="previousPage" 
              :disabled="currentPage === 1"
              class="pagination-btn"
            >
              Previous
            </button>
            <span class="page-info">
              Page {{ currentPage }} of {{ totalPages }}
            </span>
            <button 
              @click="nextPage" 
              :disabled="currentPage === totalPages"
              class="pagination-btn"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../services/api'

const router = useRouter()

// Reactive data
const searchQuery = ref('')
const clients = ref([])
const loading = ref(false)
const error = ref('')
const currentPage = ref(1)
const totalClients = ref(0)
const limit = ref(20)

const filters = ref({
  status: '',
  mode: '',
  dueDate: ''
})

// Computed properties
const totalPages = computed(() => Math.ceil(totalClients.value / limit.value))
const startItem = computed(() => (currentPage.value - 1) * limit.value + 1)
const endItem = computed(() => Math.min(currentPage.value * limit.value, totalClients.value))
const hasActiveFilters = computed(() => {
  return filters.value.status || filters.value.mode || filters.value.dueDate
})

// Methods
const handleSearch = () => {
  currentPage.value = 1
  loadClients()
}

const applyFilters = () => {
  currentPage.value = 1
  loadClients()
}

const clearSearch = () => {
  searchQuery.value = ''
  currentPage.value = 1
  loadClients()
}

const loadClients = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await api.get('/clients', {
      params: {
        page: currentPage.value,
        limit: limit.value,
        search: searchQuery.value,
        status: filters.value.status,
        mode: filters.value.mode,
        dueDate: filters.value.dueDate
      }
    })
    
    clients.value = response.data.clients
    totalClients.value = response.data.pagination.total
  } catch (err) {
    console.error('Error loading clients:', err)
    error.value = 'Failed to load clients. Please try again.'
    clients.value = []
    totalClients.value = 0
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
  return date.toLocaleDateString('en-PH')
}

const truncateAddress = (address) => {
  if (!address) return 'N/A'
  if (address.length <= 30) return address
  return address.substring(0, 30) + '...'
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

const viewClientDetails = (client) => {
  router.push(`/clients/${client.id}/details`)
}

const goToCreateClient = () => {
  router.push('/clients/new')
}

const goBack = () => {
  router.push('/dashboard')
}

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadClients()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadClients()
  }
}

const exportClients = async () => {
  try {
    const response = await api.get('/clients/export')
    // In a real implementation, you would download the CSV file
    console.log('Export data:', response.data)
    alert('Export functionality will be implemented soon!')
  } catch (err) {
    console.error('Error exporting clients:', err)
    alert('Failed to export clients')
  }
}

// Lifecycle
onMounted(() => {
  loadClients()
})
</script>

<style scoped>
.search-clients {
  padding: 1rem;
  max-width: 100%;
  margin: 0 auto;
  background-color: #f8f9fa;
  min-height: 100vh;
}

/* Header Styles */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding: 1rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logo-placeholder {
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  font-size: 1.2rem;
}

.logo-section h1 {
  margin: 0;
  color: #2c3e50;
  font-size: 1.5rem;
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

.back-btn:hover {
  background-color: #5a6268;
}

/* Search Section */
.search-section {
  margin-bottom: 2rem;
}

.search-container {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.search-input-group {
  position: relative;
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}

.search-icon {
  position: absolute;
  left: 1rem;
  color: #6c757d;
}

.search-input {
  width: 100%;
  padding: 1rem 1rem 1rem 3rem;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 1rem;
  transition: border-color 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.1);
}

.clear-btn {
  position: absolute;
  right: 1rem;
  background: none;
  border: none;
  color: #6c757d;
  cursor: pointer;
  font-size: 0.8rem;
}

.clear-btn:hover {
  color: #dc3545;
}

.search-filters {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-group label {
  font-weight: 600;
  color: #495057;
  font-size: 0.9rem;
}

.filter-group select {
  padding: 0.5rem;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  background: white;
}

/* Results Section */
.results-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  overflow: hidden;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e9ecef;
}

.results-info h3 {
  margin: 0 0 0.25rem 0;
  color: #2c3e50;
}

.results-count {
  color: #6c757d;
  font-size: 0.9rem;
}

.export-btn {
  padding: 0.75rem 1.5rem;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
}

.export-btn:hover:not(:disabled) {
  background-color: #218838;
}

.export-btn:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

/* Table Styles */
.table-container {
  overflow-x: auto;
}

.clients-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}

.clients-table th {
  background-color: #f8f9fa;
  padding: 1rem 0.75rem;
  text-align: left;
  font-weight: 600;
  color: #495057;
  border-bottom: 1px solid #e9ecef;
  white-space: nowrap;
}

.clients-table td {
  padding: 1rem 0.75rem;
  border-bottom: 1px solid #e9ecef;
  white-space: nowrap;
}

.client-row {
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.client-row:hover {
  background-color: #f8f9fa;
}

/* Status Badges */
.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
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

/* Mode Badge */
.mode-badge {
  padding: 0.25rem 0.5rem;
  background-color: #e3f2fd;
  color: #1565c0;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
}

/* Loading and Error States */
.loading-state,
.error-state,
.empty-state {
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

.error-icon,
.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.retry-btn,
.create-client-btn {
  padding: 0.75rem 1.5rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  margin-top: 1rem;
}

.retry-btn:hover,
.create-client-btn:hover {
  background-color: #2980b9;
}

/* Pagination */
.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-top: 1px solid #e9ecef;
}

.pagination-info {
  color: #6c757d;
  font-size: 0.9rem;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.pagination-btn {
  padding: 0.5rem 1rem;
  background-color: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.pagination-btn:hover:not(:disabled) {
  background-color: #5a6268;
}

.pagination-btn:disabled {
  background-color: #e9ecef;
  color: #6c757d;
  cursor: not-allowed;
}

.page-info {
  color: #495057;
  font-size: 0.9rem;
}

/* Responsive Design */
@media (max-width: 768px) {
  .search-clients {
    padding: 0.5rem;
  }
  
  .page-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .search-filters {
    flex-direction: column;
  }
  
  .filter-group {
    justify-content: space-between;
  }
  
  .results-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .pagination {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .clients-table {
    font-size: 0.75rem;
  }
  
  .clients-table th,
  .clients-table td {
    padding: 0.5rem 0.25rem;
  }
}

/* Column-specific styles for better readability */
.control-number {
  font-weight: 600;
  color: #2c3e50;
}

.client-name {
  font-weight: 600;
  color: #2c3e50;
}

.total-amount,
.ammortization,
.outstanding-balance,
.amount-release {
  font-weight: 600;
  text-align: right;
}

.address {
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
