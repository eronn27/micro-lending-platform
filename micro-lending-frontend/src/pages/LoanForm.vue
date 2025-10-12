<template>
  <div class="loan-form">
    <header class="page-header">
      <button @click="goBack" class="back-btn">← Back to Client</button>
      <h1>Create New Loan</h1>
      <div class="client-info-badge">
        <span class="label">Client:</span>
        <span class="value">{{ clientName }}</span>
      </div>
    </header>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading client information...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Error Loading Client</h3>
      <p>{{ error }}</p>
      <button @click="loadClient" class="retry-btn">Retry</button>
    </div>

    <!-- Loan Form -->
    <div v-else class="form-layout">
      <!-- Left Column: Loan Calculator -->
      <div class="calculator-column">
        <div class="calculator-card">
          <h2>Loan Calculator</h2>
          
          <div class="form-group">
            <label for="loanAmount">Loan Amount (₱) *</label>
            <input
              type="number"
              id="loanAmount"
              v-model.number="calculator.loanAmount"
              @input="calculateLoan"
              placeholder="0.00"
              step="0.01"
              min="0"
              required
            >
          </div>

          <div class="form-group">
            <label for="term">Loan Term (Months) *</label>
            <select 
              id="term"
              v-model.number="calculator.term" 
              @change="calculateLoan"
              required
            >
              <option value="">Select Term</option>
              <option value="2">2 Months</option>
              <option value="4">4 Months</option>
              <option value="5">5 Months</option>
              <option value="6">6 Months</option>
            </select>
          </div>

          <div class="form-group">
            <label>Deduction Type *</label>
            <div class="radio-group">
              <label class="radio-option">
                <input 
                  type="radio" 
                  value="Without" 
                  v-model="calculator.deductionType"
                  @change="calculateLoan"
                >
                <span class="radio-label">Without Deduction</span>
              </label>
              <label class="radio-option">
                <input 
                  type="radio" 
                  value="With" 
                  v-model="calculator.deductionType"
                  @change="calculateLoan"
                >
                <span class="radio-label">With Deduction</span>
              </label>
            </div>
          </div>

          <div v-if="calculationResult" class="calculation-results">
            <h3>Calculation Results</h3>
            
            <div class="result-grid">
              <div class="result-item">
                <span class="result-label">Total Amount:</span>
                <span class="result-value">₱{{ formatCurrency(calculationResult.totalAmount) }}</span>
              </div>
              
              <div class="result-item">
                <span class="result-label">Ammortization:</span>
                <span class="result-value">₱{{ formatCurrency(calculationResult.ammortization) }}</span>
              </div>
              
              <div class="result-item">
                <span class="result-label">Weekly Payment:</span>
                <span class="result-value">₱{{ formatCurrency(calculationResult.weeklyPayment) }}</span>
              </div>
              
              <div class="result-item">
                <span class="result-label">Number of Weeks:</span>
                <span class="result-value">{{ calculationResult.numberOfWeeks }}</span>
              </div>
              
              <div class="result-item">
                <span class="result-label">Total Interest:</span>
                <span class="result-value">₱{{ formatCurrency(calculationResult.totalInterest) }}</span>
              </div>
              
              <div class="result-item highlight">
                <span class="result-label">Outstanding Balance:</span>
                <span class="result-value">₱{{ formatCurrency(calculationResult.outstandingBalance) }}</span>
              </div>
            </div>

            <div v-if="calculationApplied" class="calculation-note">
              ✅ Loan details applied automatically
            </div>
          </div>

          <!-- Loan Cycle Info -->
          <div class="loan-cycle-info">
            <div class="cycle-badge">
              <span class="cycle-label">Loan Cycle:</span>
              <span class="cycle-number">{{ loanCycle }}</span>
            </div>
            <p class="cycle-description">
              This will be loan #{{ loanCycle }} for this client
            </p>
          </div>
        </div>
      </div>

      <!-- Right Column: Loan Details Form -->
      <div class="form-column">
        <form @submit.prevent="submitForm" class="form-container">
          <!-- Loan Details Section -->
          <div class="form-section">
            <h2>Loan Details</h2>
            
            <div class="form-row">
              <div class="form-group">
                <label for="totalAmount">Total Amount *</label>
                <input
                  type="number"
                  id="totalAmount"
                  v-model.number="loan.total_amount"
                  required
                  step="0.01"
                  min="0"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="terms">Terms (Months) *</label>
                <input
                  type="number"
                  id="terms"
                  v-model.number="loan.terms"
                  required
                  min="1"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="ammortization">Ammortization *</label>
                <input
                  type="number"
                  id="ammortization"
                  v-model.number="loan.ammortization"
                  required
                  step="0.01"
                  min="0"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="outstandingBalance">Outstanding Balance *</label>
                <input
                  type="number"
                  id="outstandingBalance"
                  v-model.number="loan.outstanding_balance"
                  required
                  step="0.01"
                  min="0"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="amountRelease">Amount Release *</label>
                <input
                  type="number"
                  id="amountRelease"
                  v-model.number="loan.amount_release"
                  required
                  step="0.01"
                  min="0"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="deductions">Deductions</label>
                <input
                  type="text"
                  id="deductions"
                  v-model="loan.deductions"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="dateOfRelease">Date of Release</label>
                <input
                  type="date"
                  id="dateOfRelease"
                  v-model="loan.date_of_release"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="dueDate">Due Date</label>
                <select id="dueDate" v-model="loan.due_date" :disabled="submitting">
                  <option value="">Select Due Date</option>
                  <option value="Monday">Monday</option>
                  <option value="Tuesday">Tuesday</option>
                  <option value="Wednesday">Wednesday</option>
                  <option value="Thursday">Thursday</option>
                  <option value="Friday">Friday</option>
                  <option value="Saturday">Saturday</option>
                </select>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="methodOfPayment">Method of Payment</label>
                <select id="methodOfPayment" v-model="loan.method_of_payment" :disabled="submitting">
                  <option value="">Select Method</option>
                  <option value="Cash">Cash</option>
                  <option value="Bank Transfer">Bank Transfer</option>
                  <option value="GCash">GCash</option>
                  <option value="Other">Other</option>
                </select>
              </div>
              
              <div class="form-group">
                <label for="creditHistory">Credit History</label>
                <select id="creditHistory" v-model="loan.credit_history" :disabled="submitting">
                  <option value="">Select History</option>
                  <option value="New">New</option>
                  <option value="Existing">Existing</option>
                  <option value="New W/ Existing">New W/ Existing</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Co-maker Section -->
          <div class="form-section">
            <h2>Co-maker Information</h2>
            
            <div v-for="(comaker, index) in comakers" :key="index" class="comaker-section">
              <h3>Co-maker {{ index + 1 }}</h3>
              
              <div class="form-row">
                <div class="form-group">
                  <label :for="`comakerName${index}`">Name</label>
                  <input
                    type="text"
                    :id="`comakerName${index}`"
                    v-model="comaker.name"
                    :disabled="submitting"
                  >
                </div>
                
                <div class="form-group">
                  <label :for="`comakerBusiness${index}`">Business</label>
                  <input
                    type="text"
                    :id="`comakerBusiness${index}`"
                    v-model="comaker.business"
                    :disabled="submitting"
                  >
                </div>
              </div>

              <div class="form-group">
                <label :for="`comakerAddress${index}`">Address</label>
                <textarea
                  :id="`comakerAddress${index}`"
                  v-model="comaker.address"
                  rows="2"
                  :disabled="submitting"
                ></textarea>
              </div>

              <button
                v-if="comakers.length > 1"
                type="button"
                @click="removeComaker(index)"
                class="btn-remove"
                :disabled="submitting"
              >
                Remove Co-maker
              </button>
            </div>

            <button
              type="button"
              @click="addComaker"
              class="btn-add"
              :disabled="submitting"
            >
              + Add Another Co-maker
            </button>
          </div>

          <!-- Approval Section -->
          <div class="form-section">
            <h2>Approval Information</h2>
            
            <div class="form-row">
              <div class="form-group">
                <label for="recommendedBy">Recommended By</label>
                <input
                  type="text"
                  id="recommendedBy"
                  v-model="loan.recommended_by"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="approvedBy">Approved By</label>
                <input
                  type="text"
                  id="approvedBy"
                  v-model="loan.approved_by"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="recommendedLoanAmount">Recommended Loan Amount</label>
                <input
                  type="number"
                  id="recommendedLoanAmount"
                  v-model.number="loan.recommended_loan_amount"
                  step="0.01"
                  min="0"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="approvedLoanAmount">Approved Loan Amount</label>
                <input
                  type="number"
                  id="approvedLoanAmount"
                  v-model.number="loan.approved_loan_amount"
                  step="0.01"
                  min="0"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="checkedBy">Checked By</label>
                <input
                  type="text"
                  id="checkedBy"
                  v-model="loan.checked_by"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="nameCI">Name CI</label>
                <input
                  type="text"
                  id="nameCI"
                  v-model="loan.name_ci"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="notedBy">Noted By</label>
                <input
                  type="text"
                  id="notedBy"
                  v-model="loan.noted_by"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="applicationDate">Application Date</label>
                <input
                  type="date"
                  id="applicationDate"
                  v-model="loan.application_date"
                  :disabled="submitting"
                >
              </div>
            </div>
          </div>

          <!-- Form Actions -->
          <div class="form-actions">
            <button type="button" @click="goBack" class="btn-secondary" :disabled="submitting">
              Cancel
            </button>
            <button type="submit" :disabled="submitting || !isFormValid" class="btn-primary">
              {{ submitting ? 'Creating Loan...' : 'Create Loan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../services/api'

const route = useRoute()
const router = useRouter()

// State
const client = ref(null)
const loading = ref(false)
const submitting = ref(false)
const error = ref('')
const calculationApplied = ref(false)

// Calculator
const calculator = reactive({
  loanAmount: 0,
  term: '',
  deductionType: 'Without'
})

const calculationResult = ref(null)

// Loan Form
const loan = reactive({
  total_amount: 0,
  terms: 0,
  ammortization: 0,
  outstanding_balance: 0,
  mode: 'Weekly',
  payment_period_weeks: 0,
  deductions: '',
  due_date: '',
  amount_release: 0,
  date_of_release: '',
  credit_history: 'Existing',
  method_of_payment: '',
  recommended_by: '',
  approved_by: '',
  recommended_loan_amount: 0,
  approved_loan_amount: 0,
  checked_by: '',
  name_ci: '',
  noted_by: '',
  application_date: ''
})

const comakers = ref([{ name: '', address: '', business: '' }])

// Computed
const clientName = computed(() => {
  if (!client.value) return ''
  return `${client.value.first_name} ${client.value.last_name}`
})

const loanCycle = computed(() => {
  if (!client.value || !client.value.loans) return 1
  return client.value.loans.length + 1
})

const isFormValid = computed(() => {
  return loan.total_amount > 0 && 
         loan.terms > 0 && 
         loan.ammortization > 0 && 
         loan.amount_release >= 0
})

// Methods
const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-PH', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

const calculateLoan = () => {
  if (!calculator.loanAmount || !calculator.term) {
    calculationResult.value = null
    calculationApplied.value = false
    return
  }

  const INTEREST_RATE = 0.035
  const DEDUCTION_FEE = 1200
  
  const numberOfWeeks = calculator.term * 4
  const weeklyInterest = calculator.loanAmount * INTEREST_RATE
  const totalInterest = weeklyInterest * calculator.term
  const baseTotalAmount = totalInterest + calculator.loanAmount
  
  let finalTotalAmount = 0
  let amountRelease = calculator.loanAmount
  
  if (calculator.deductionType === 'With') {
    finalTotalAmount = baseTotalAmount
  } else {
    finalTotalAmount = baseTotalAmount + DEDUCTION_FEE
  }
  
  const ammortization = Math.round(finalTotalAmount / numberOfWeeks)
  
  calculationResult.value = {
    loanAmount: parseFloat(calculator.loanAmount),
    term: parseInt(calculator.term),
    deductionType: calculator.deductionType,
    numberOfWeeks,
    totalInterest: parseFloat(totalInterest.toFixed(2)),
    totalAmount: parseFloat(finalTotalAmount.toFixed(2)),
    ammortization: parseFloat(ammortization.toFixed(2)),
    weeklyPayment: parseFloat(ammortization.toFixed(2)),
    outstandingBalance: parseFloat(finalTotalAmount.toFixed(2)),
    amountRelease: parseFloat(amountRelease.toFixed(2))
  }

  applyCalculation()
}

const applyCalculation = () => {
  if (!calculationResult.value) return
  
  const result = calculationResult.value
  
  loan.total_amount = result.totalAmount
  loan.terms = result.term
  loan.ammortization = result.ammortization
  loan.outstanding_balance = result.outstandingBalance
  loan.deductions = result.deductionType
  loan.amount_release = result.amountRelease
  loan.payment_period_weeks = result.numberOfWeeks
  
  calculationApplied.value = true
}

const addComaker = () => {
  comakers.value.push({ name: '', address: '', business: '' })
}

const removeComaker = (index) => {
  if (comakers.value.length > 1) {
    comakers.value.splice(index, 1)
  }
}

const loadClient = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await api.get(`/clients/${route.params.id}`)
    client.value = response.data
  } catch (err) {
    console.error('Error loading client:', err)
    error.value = 'Failed to load client data. Please try again.'
  } finally {
    loading.value = false
  }
}

const submitForm = async () => {
  if (submitting.value || !isFormValid.value) return
  
  submitting.value = true
  
  try {
    // Debug: Check if client exists first
    console.log('Client ID:', route.params.id)
    const clientCheck = await api.get(`/clients/${route.params.id}`)
    console.log('Client found:', clientCheck.data)

    const formData = {
      client_id: parseInt(route.params.id),
      loan: {
        ...loan,
        status: 'Active',
        loan_cycle: loanCycle.value
      },
      comakers: comakers.value.filter(c => c.name.trim() !== '')
    }

    console.log('Sending data:', formData)
    
    await api.post('/loans', formData)
    alert('Loan created successfully!')
    router.push(`/clients/${route.params.id}`)
  } catch (err) {
    console.error('Error creating loan:', err)
    console.error('Error response:', err.response)
    const errorMessage = err.response?.data?.error || 'Error creating loan. Please try again.'
    alert(errorMessage)
  } finally {
    submitting.value = false
  }
}
const goBack = () => {
  router.push(`/clients/${route.params.id}`)
}

// Watchers
watch(
  () => [calculator.loanAmount, calculator.term, calculator.deductionType],
  () => calculateLoan(),
  { immediate: true }
)

onMounted(() => {
  loadClient()
})
</script>

<style scoped>
/* Base Styles */
.loan-form {
  padding: 1rem;
  max-width: 1400px;
  margin: 0 auto;
  background-color: #f8f9fa;
  min-height: 100vh;
}

/* Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  margin-bottom: 2rem;
}

.page-header h1 {
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

.back-btn:hover {
  background-color: #5a6268;
}

.client-info-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: #e3f2fd;
  border-radius: 4px;
}

.client-info-badge .label {
  font-weight: 600;
  color: #495057;
}

.client-info-badge .value {
  color: #2c3e50;
  font-weight: 500;
}

/* Layout */
.form-layout {
  display: grid;
  grid-template-columns: 400px 1fr;
  gap: 2rem;
  align-items: start;
}

.calculator-column {
  position: sticky;
  top: 1rem;
}

.calculator-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  border-left: 4px solid #3498db;
}

.calculator-card h2 {
  margin-bottom: 1.5rem;
  color: #2c3e50;
  font-size: 1.3rem;
  text-align: center;
}

.form-column {
  min-width: 0;
}

.form-container {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

/* Loan Cycle Info */
.loan-cycle-info {
  margin-top: 1.5rem;
  padding: 1rem;
  background: #fff8e1;
  border-radius: 6px;
  border-left: 4px solid #ffc107;
}

.cycle-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.cycle-label {
  font-weight: 600;
  color: #495057;
}

.cycle-number {
  background: #ffc107;
  color: #000;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-weight: bold;
  font-size: 1.1rem;
}

.cycle-description {
  margin: 0;
  font-size: 0.9rem;
  color: #495057;
}

/* Form Sections */
.form-section {
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #e9ecef;
}

.form-section:last-of-type {
  border-bottom: none;
}

.form-section h2 {
  margin-bottom: 1.5rem;
  color: #2c3e50;
  font-size: 1.3rem;
  border-bottom: 2px solid #3498db;
  padding-bottom: 0.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #495057;
  font-size: 0.9rem;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 0.9rem;
  transition: border-color 0.15s ease-in-out;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.25);
}

.form-group input:disabled,
.form-group select:disabled,
.form-group textarea:disabled {
  background-color: #e9ecef;
  cursor: not-allowed;
}

/* Radio Group */
.radio-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.radio-option {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.radio-option:hover {
  border-color: #3498db;
  background-color: #f8f9fa;
}

.radio-option input {
  margin-right: 0.75rem;
}

.radio-label {
  font-weight: 500;
}

/* Calculation Results */
.calculation-results {
  margin-top: 1.5rem;
  padding: 1.5rem;
  background: #f8f9fa;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.calculation-results h3 {
  margin-bottom: 1rem;
  color: #2c3e50;
  font-size: 1.1rem;
  text-align: center;
}

.result-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.result-item {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem;
  background: white;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

.result-item.highlight {
  background: #e3f2fd;
  border-color: #3498db;
  font-weight: bold;
}

.result-label {
  font-size: 0.9rem;
  color: #495057;
}

.result-value {
  font-weight: 600;
  color: #2c3e50;
}

.calculation-note {
  background: #d4edda;
  color: #155724;
  padding: 0.75rem;
  border-radius: 4px;
  border: 1px solid #c3e6cb;
  font-size: 0.9rem;
  text-align: center;
  margin-top: 1rem;
}

/* Co-maker Section */
.comaker-section {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 6px;
  margin-bottom: 1rem;
  border-left: 4px solid #3498db;
}

.comaker-section h3 {
  margin: 0 0 1rem 0;
  color: #495057;
  font-size: 1.1rem;
}

.btn-add,
.btn-remove {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  margin-top: 0.5rem;
}

.btn-add {
  background-color: #28a745;
  color: white;
}

.btn-add:hover:not(:disabled) {
  background-color: #218838;
}

.btn-remove {
  background-color: #dc3545;
  color: white;
}

.btn-remove:hover:not(:disabled) {
  background-color: #c82333;
}

/* Form Actions */
.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e9ecef;
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  min-width: 140px;
}

.btn-primary {
  background-color: #3498db;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #2980b9;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #5a6268;
}

.btn-primary:disabled,
.btn-secondary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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

/* Responsive Design */
@media (max-width: 1024px) {
  .form-layout {
    grid-template-columns: 1fr;
  }
  
  .calculator-column {
    position: static;
  }
}

@media (max-width: 768px) {
  .loan-form {
    padding: 0.5rem;
  }
  
  .page-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .form-container {
    padding: 1rem;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn-primary,
  .btn-secondary {
    width: 100%;
  }
}
</style>
