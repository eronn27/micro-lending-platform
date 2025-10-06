<template>
  <div class="client-form">
    <header class="form-header">
      <h1>Create New Client</h1>
      <button @click="goBack" class="back-btn">← Back to Dashboard</button>
    </header>

    <div class="form-layout">
      <!-- Left Column: Loan Calculator -->
      <div class="calculator-column">
        <div class="calculator-card">
          <h2>Loan Calculator</h2>
          
          <!-- Loan Amount -->
          <div class="form-group">
            <label for="loanAmount">Loan Amount (₱) *</label>
            <input
              type="number"
              id="loanAmount"
              v-model.number="loanCalculator.loanAmount"
              @input="calculateLoan"
              placeholder="0.00"
              step="0.01"
              min="0"
              required
            >
          </div>

          <!-- Term Selection -->
          <div class="form-group">
            <label for="term">Loan Term (Months) *</label>
            <select 
              id="term"
              v-model.number="loanCalculator.term" 
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

          <!-- Deduction Type -->
          <div class="form-group">
            <label>Deduction Type *</label>
            <div class="radio-group">
              <label class="radio-option">
                <input 
                  type="radio" 
                  value="Without" 
                  v-model="loanCalculator.deductionType"
                  @change="calculateLoan"
                >
                <span class="radio-label">Without Deduction</span>
              </label>
              <label class="radio-option">
                <input 
                  type="radio" 
                  value="With" 
                  v-model="loanCalculator.deductionType"
                  @change="calculateLoan"
                >
                <span class="radio-label">With Deduction</span>
              </label>
            </div>
          </div>

          <!-- Calculation Results -->
          <div v-if="loanCalculationResult" class="calculation-results">
            <h3>Loan Calculation Results</h3>
            
            <div class="result-grid">
              <div class="result-item">
                <span class="result-label">Total Amount:</span>
                <span class="result-value">₱{{ formatCurrency(loanCalculationResult.totalAmount) }}</span>
              </div>
              
              <div class="result-item">
                <span class="result-label">Ammortization:</span>
                <span class="result-value">₱{{ formatCurrency(loanCalculationResult.ammortization) }}</span>
              </div>
              
              <div class="result-item">
                <span class="result-label">Weekly Payment:</span>
                <span class="result-value">₱{{ formatCurrency(loanCalculationResult.weeklyPayment) }}</span>
              </div>
              
              <div class="result-item">
                <span class="result-label">Number of Weeks:</span>
                <span class="result-value">{{ loanCalculationResult.numberOfWeeks }}</span>
              </div>
              
              <div class="result-item">
                <span class="result-label">Total Interest:</span>
                <span class="result-value">₱{{ formatCurrency(loanCalculationResult.totalInterest) }}</span>
              </div>
              
              <div class="result-item highlight">
                <span class="result-label">Outstanding Balance:</span>
                <span class="result-value">₱{{ formatCurrency(loanCalculationResult.outstandingBalance) }}</span>
              </div>
            </div>

            <div v-if="loanCalculationApplied" class="calculation-note">
              ✅ Loan details applied automatically
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Client Form -->
      <div class="form-column">
        <form @submit.prevent="submitForm" class="form-container">
          <!-- Personal Information Section -->
          <div class="form-section">
            <h2>Personal Information</h2>
            
            <div class="form-row">
              <div class="form-group">
                <label for="firstName">First Name *</label>
                <input
                  type="text"
                  id="firstName"
                  v-model="client.first_name"
                  required
                  @blur="checkDuplicate"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="middleName">Middle Name</label>
                <input
                  type="text"
                  id="middleName"
                  v-model="client.middle_name"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="lastName">Last Name *</label>
                <input
                  type="text"
                  id="lastName"
                  v-model="client.last_name"
                  required
                  @blur="checkDuplicate"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="nickname">Nickname</label>
                <input
                  type="text"
                  id="nickname"
                  v-model="client.nickname"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="dateOfBirth">Date of Birth</label>
                <input
                  type="date"
                  id="dateOfBirth"
                  v-model="client.date_of_birth"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="gender">Gender</label>
                <select id="gender" v-model="client.gender" :disabled="submitting">
                  <option value="">Select Gender</option>
                  <option value="Male">Male</option>
                  <option value="Female">Female</option>
                  <option value="Other">Other</option>
                </select>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="civilStatus">Civil Status</label>
                <select id="civilStatus" v-model="client.civil_status" :disabled="submitting">
                  <option value="">Select Status</option>
                  <option value="Single">Single</option>
                  <option value="Married">Married</option>
                  <option value="Divorced">Divorced</option>
                  <option value="Widowed">Widowed</option>
                </select>
              </div>
              
              <div class="form-group">
                <label for="religion">Religion</label>
                <input
                  type="text"
                  id="religion"
                  v-model="client.religion"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-group">
              <label for="homeAddress">Home Address *</label>
              <textarea
                id="homeAddress"
                v-model="client.home_address"
                required
                rows="3"
                :disabled="submitting"
              ></textarea>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="yearsOfResidence">Years of Residence</label>
                <input
                  type="number"
                  id="yearsOfResidence"
                  v-model.number="client.years_of_residence"
                  min="0"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="contactNumber">Contact Number *</label>
                <input
                  type="tel"
                  id="contactNumber"
                  v-model="client.contact_number"
                  required
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="facebookAccount">Facebook Account</label>
                <input
                  type="text"
                  id="facebookAccount"
                  v-model="client.facebook_account"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="age">Age</label>
                <input
                  type="number"
                  id="age"
                  v-model.number="client.age"
                  min="0"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div v-if="duplicateWarning" class="warning-banner">
              ⚠️ Warning: A client with a similar name already exists.
              <br>
              <small>Similar records found: {{ duplicateCount }}</small>
            </div>
          </div>

          <!-- Income Information Section -->
          <div class="form-section">
            <h2>Income Information</h2>
            
            <div class="form-group">
              <label>Family Income</label>
              <div class="form-row">
                <div class="form-group">
                  <label for="familyIncomeDaily">Daily</label>
                  <input
                    type="number"
                    id="familyIncomeDaily"
                    v-model.number="income.family_income_daily"
                    step="0.01"
                    min="0"
                    :disabled="submitting"
                  >
                </div>
                <div class="form-group">
                  <label for="familyIncomeMonthly">Monthly</label>
                  <input
                    type="number"
                    id="familyIncomeMonthly"
                    v-model.number="income.family_income_monthly"
                    step="0.01"
                    min="0"
                    :disabled="submitting"
                  >
                </div>
              </div>
            </div>

            <div class="form-group">
              <label>Total Cost</label>
              <div class="form-row">
                <div class="form-group">
                  <label for="totalCostDaily">Daily</label>
                  <input
                    type="number"
                    id="totalCostDaily"
                    v-model.number="income.total_cost_daily"
                    step="0.01"
                    min="0"
                    :disabled="submitting"
                  >
                </div>
                <div class="form-group">
                  <label for="totalCostMonthly">Monthly</label>
                  <input
                    type="number"
                    id="totalCostMonthly"
                    v-model.number="income.total_cost_monthly"
                    step="0.01"
                    min="0"
                    :disabled="submitting"
                  >
                </div>
              </div>
            </div>

            <div class="form-group">
              <label>Net Income</label>
              <div class="form-row">
                <div class="form-group">
                  <label for="netIncomeDaily">Daily</label>
                  <input
                    type="number"
                    id="netIncomeDaily"
                    v-model.number="income.net_income_daily"
                    step="0.01"
                    min="0"
                    :disabled="submitting"
                  >
                </div>
                <div class="form-group">
                  <label for="netIncomeMonthly">Monthly</label>
                  <input
                    type="number"
                    id="netIncomeMonthly"
                    v-model.number="income.net_income_monthly"
                    step="0.01"
                    min="0"
                    :disabled="submitting"
                  >
                </div>
              </div>
            </div>
          </div>

          <!-- Loan Information Section (Auto-filled from calculator) -->
          <div class="form-section">
            <h2>Loan Information</h2>
            
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
                <label for="mode">Mode</label>
                <input
                  type="text"
                  id="mode"
                  v-model="loan.mode"
                  value="Weekly"
                  disabled
                  class="disabled-field"
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
                <label for="methodOfPayment">Method of Payment</label>
                <select id="methodOfPayment" v-model="loan.method_of_payment" :disabled="submitting">
                  <option value="">Select Method</option>
                  <option value="Cash">Cash</option>
                  <option value="Bank Transfer">Bank Transfer</option>
                  <option value="GCash">GCash</option>
                  <option value="Other">Other</option>
                </select>
              </div>
            </div>

            <div class="form-row">
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

            <div class="calculation-note" v-if="loanCalculationApplied">
              ✅ Loan details calculated automatically
            </div>
          </div>

          <!-- Co-maker Information Section -->
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

          <!-- Family Information Section -->
          <div class="form-section">
            <h2>Family Information</h2>
            
            <div class="form-row">
              <div class="form-group">
                <label for="fatherName">Father's Name</label>
                <input
                  type="text"
                  id="fatherName"
                  v-model="family.father_name"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="motherName">Mother's Name</label>
                <input
                  type="text"
                  id="motherName"
                  v-model="family.mother_name"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-group">
              <label for="familyAddress">Family Address</label>
              <textarea
                id="familyAddress"
                v-model="family.address"
                rows="2"
                :disabled="submitting"
              ></textarea>
            </div>

            <!-- Siblings -->
            <div class="siblings-section">
              <h3>Siblings</h3>
              <div v-for="(sibling, index) in siblings" :key="index" class="sibling-item">
                <div class="form-row">
                  <div class="form-group">
                    <label :for="`siblingName${index}`">Name</label>
                    <input
                      type="text"
                      :id="`siblingName${index}`"
                      v-model="sibling.name"
                      :disabled="submitting"
                    >
                  </div>
                  
                  <div class="form-group">
                    <label :for="`siblingAge${index}`">Age</label>
                    <input
                      type="number"
                      :id="`siblingAge${index}`"
                      v-model.number="sibling.age"
                      min="0"
                      :disabled="submitting"
                    >
                  </div>
                </div>
                
                <div class="form-group">
                  <label :for="`siblingAddress${index}`">Address</label>
                  <input
                    type="text"
                    :id="`siblingAddress${index}`"
                    v-model="sibling.address"
                    :disabled="submitting"
                  >
                </div>

                <button
                  v-if="siblings.length > 1"
                  type="button"
                  @click="removeSibling(index)"
                  class="btn-remove"
                  :disabled="submitting"
                >
                  Remove Sibling
                </button>
              </div>

              <button
                type="button"
                @click="addSibling"
                class="btn-add"
                :disabled="submitting"
              >
                + Add Sibling
              </button>
            </div>
          </div>

          <!-- Spouse Information Section -->
          <div class="form-section">
            <h2>Spouse Information (If Applicable)</h2>
            
            <div class="form-row">
              <div class="form-group">
                <label for="spouseName">Spouse Name</label>
                <input
                  type="text"
                  id="spouseName"
                  v-model="spouse.name"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="spouseAge">Age</label>
                <input
                  type="number"
                  id="spouseAge"
                  v-model.number="spouse.age"
                  min="0"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="spouseNickname">Nickname</label>
                <input
                  type="text"
                  id="spouseNickname"
                  v-model="spouse.nickname"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="spouseBirthday">Birthday</label>
                <input
                  type="date"
                  id="spouseBirthday"
                  v-model="spouse.birthday"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="spouseWork">Work</label>
                <input
                  type="text"
                  id="spouseWork"
                  v-model="spouse.work"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="spouseContactNumber">Contact Number</label>
                <input
                  type="tel"
                  id="spouseContactNumber"
                  v-model="spouse.contact_number"
                  :disabled="submitting"
                >
              </div>
            </div>

            <!-- Dependents -->
            <div class="dependents-section">
              <h3>Dependents</h3>
              <div v-for="(dependent, index) in dependents" :key="index" class="dependent-item">
                <div class="form-row">
                  <div class="form-group">
                    <label :for="`dependentName${index}`">Name</label>
                    <input
                      type="text"
                      :id="`dependentName${index}`"
                      v-model="dependent.name"
                      :disabled="submitting"
                    >
                  </div>
                  
                  <div class="form-group">
                    <label :for="`dependentSurname${index}`">Surname</label>
                    <input
                      type="text"
                      :id="`dependentSurname${index}`"
                      v-model="dependent.surname"
                      :disabled="submitting"
                    >
                  </div>
                </div>
                
                <div class="form-row">
                  <div class="form-group">
                    <label :for="`dependentAge${index}`">Age</label>
                    <input
                      type="number"
                      :id="`dependentAge${index}`"
                      v-model.number="dependent.age"
                      min="0"
                      :disabled="submitting"
                    >
                  </div>
                  
                  <div class="form-group">
                    <label :for="`dependentSpouse${index}`">Spouse</label>
                    <input
                      type="text"
                      :id="`dependentSpouse${index}`"
                      v-model="dependent.spouse"
                      :disabled="submitting"
                    >
                  </div>
                </div>
                
                <div class="form-group">
                  <label :for="`dependentAddress${index}`">Address</label>
                  <input
                    type="text"
                    :id="`dependentAddress${index}`"
                    v-model="dependent.address"
                    :disabled="submitting"
                  >
                </div>

                <button
                  v-if="dependents.length > 1"
                  type="button"
                  @click="removeDependent(index)"
                  class="btn-remove"
                  :disabled="submitting"
                >
                  Remove Dependent
                </button>
              </div>

              <button
                type="button"
                @click="addDependent"
                class="btn-add"
                :disabled="submitting"
              >
                + Add Dependent
              </button>
            </div>
          </div>

          <!-- Company/Approval Section -->
          <div class="form-section">
            <h2>Company Information</h2>
            
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
                <label for="loanCycle">Loan Cycle</label>
                <input
                  type="number"
                  id="loanCycle"
                  v-model.number="loan.loan_cycle"
                  min="0"
                  :disabled="submitting"
                >
              </div>
              
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
            </div>

            <div class="form-row">
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
              
              <div class="form-group">
                <label for="checkedBy">Checked By</label>
                <input
                  type="text"
                  id="checkedBy"
                  v-model="loan.checked_by"
                  :disabled="submitting"
                >
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label for="nameCI">Name CI</label>
                <input
                  type="text"
                  id="nameCI"
                  v-model="loan.name_ci"
                  :disabled="submitting"
                >
              </div>
              
              <div class="form-group">
                <label for="notedBy">Noted By</label>
                <input
                  type="text"
                  id="notedBy"
                  v-model="loan.noted_by"
                  :disabled="submitting"
                >
              </div>
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

          <!-- Form Actions -->
          <div class="form-actions">
            <button type="button" @click="goBack" class="btn-secondary" :disabled="submitting">
              Cancel
            </button>
            <button type="submit" :disabled="submitting || !isFormValid" class="btn-primary">
              {{ submitting ? 'Creating Client...' : 'Create Client' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../services/api'

const router = useRouter()

// Loan Calculator Logic
const useLoanCalculator = () => {
  const calculateLoan = (loanAmount, term, deductionType) => {
    if (!loanAmount || !term) return null;

    // Constants based on your formula
    const INTEREST_RATE = 0.035; // 3.5% weekly interest
    const DEDUCTION_FEE = 1200; // Deduction fee
    
    // Calculate base values
    const numberOfWeeks = term * 4; // 4 weeks per month
    const weeklyInterest = loanAmount * INTEREST_RATE;
    const totalInterest = weeklyInterest * term;
    const baseTotalAmount = totalInterest + loanAmount;
    
    let finalTotalAmount = 0;
    let ammortization = 0;
    let amountRelease = loanAmount;
    
    // Apply deduction type
    if (deductionType === 'With') {
      // With Deduction: Deduct fee from release amount
      finalTotalAmount = baseTotalAmount;
      // amountRelease = loanAmount - DEDUCTION_FEE; // Commented out as requested
    } else if (deductionType === 'Without') {
      // Without Deduction: Add fee to total amount
      finalTotalAmount = baseTotalAmount + DEDUCTION_FEE;
      amountRelease = loanAmount;
    }
    
    // Calculate weekly payment (ammortization)
    ammortization = Math.round(finalTotalAmount / numberOfWeeks);
    
    return {
      loanAmount: parseFloat(loanAmount),
      term: parseInt(term),
      deductionType,
      numberOfWeeks,
      totalInterest: parseFloat(totalInterest.toFixed(2)),
      totalAmount: parseFloat(finalTotalAmount.toFixed(2)),
      ammortization: parseFloat(ammortization.toFixed(2)),
      weeklyPayment: parseFloat(ammortization.toFixed(2)),
      outstandingBalance: parseFloat(finalTotalAmount.toFixed(2)),
      amountRelease: parseFloat(amountRelease.toFixed(2))
    };
  };

  return { calculateLoan };
};

// Main reactive data
const { calculateLoan } = useLoanCalculator();

// Loan Calculator State
const loanCalculator = reactive({
  loanAmount: 0,
  term: '',
  deductionType: 'Without'
});

const loanCalculationResult = ref(null);
const loanCalculationApplied = ref(false);

// Client Form Data
const client = reactive({
  first_name: '',
  middle_name: '',
  last_name: '',
  nickname: '',
  date_of_birth: '',
  gender: '',
  religion: '',
  civil_status: '',
  home_address: '',
  years_of_residence: 0,
  facebook_account: '',
  age: 0,
  contact_number: ''
});

const income = reactive({
  family_income_daily: 0,
  family_income_monthly: 0,
  total_cost_daily: 0,
  total_cost_monthly: 0,
  net_income_daily: 0,
  net_income_monthly: 0
});

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
  credit_history: '',
  method_of_payment: '',
  recommended_by: '',
  approved_by: '',
  loan_cycle: 0,
  recommended_loan_amount: 0,
  approved_loan_amount: 0,
  checked_by: '',
  name_ci: '',
  noted_by: '',
  application_date: ''
});

const comakers = ref([{ name: '', address: '', business: '' }]);
const family = reactive({ father_name: '', mother_name: '', address: '' });
const siblings = ref([{ name: '', age: 0, address: '' }]);
const spouse = reactive({ 
  name: '', 
  age: 0, 
  nickname: '', 
  birthday: '', 
  work: '', 
  contact_number: '' 
});
const dependents = ref([{ name: '', surname: '', age: 0, spouse: '', address: '' }]);

// UI State
const duplicateWarning = ref(false);
const duplicateCount = ref(0);
const submitting = ref(false);
const duplicateCheckTimeout = ref(null);

// Computed Properties
const isFormValid = computed(() => {
  return client.first_name && client.last_name && client.contact_number && 
         client.home_address && loan.total_amount > 0 && loan.terms > 0 && 
         loan.ammortization > 0 && loan.amount_release >= 0;
});

const hasRequiredCalculatorFields = computed(() => {
  return loanCalculator.loanAmount > 0 && loanCalculator.term && loanCalculator.deductionType;
});

// Methods
const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-PH', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount);
};

const calculateLoanDetails = () => {
  if (!hasRequiredCalculatorFields.value) {
    loanCalculationResult.value = null;
    loanCalculationApplied.value = false;
    return;
  }
  
  loanCalculationResult.value = calculateLoan(
    parseFloat(loanCalculator.loanAmount),
    parseInt(loanCalculator.term),
    loanCalculator.deductionType
  );
  
  // Automatically apply the calculation to the form
  applyLoanCalculation();
};

const applyLoanCalculation = () => {
  if (!loanCalculationResult.value) return;
  
  const result = loanCalculationResult.value;
  
  // Apply calculation results to loan form
  loan.total_amount = result.totalAmount;
  loan.terms = result.term;
  loan.ammortization = result.ammortization;
  loan.outstanding_balance = result.outstandingBalance;
  loan.deductions = result.deductionType;
  loan.amount_release = result.amountRelease;
  loan.mode = 'Weekly';
  loan.payment_period_weeks = result.numberOfWeeks; // ADD THIS LINE 

  loanCalculationApplied.value = true;
};

// Methods for dynamic form sections
const addComaker = () => {
  comakers.value.push({ name: '', address: '', business: '' });
};

const removeComaker = (index) => {
  if (comakers.value.length > 1) {
    comakers.value.splice(index, 1);
  }
};

const addSibling = () => {
  siblings.value.push({ name: '', age: 0, address: '' });
};

const removeSibling = (index) => {
  if (siblings.value.length > 1) {
    siblings.value.splice(index, 1);
  }
};

const addDependent = () => {
  dependents.value.push({ name: '', surname: '', age: 0, spouse: '', address: '' });
};

const removeDependent = (index) => {
  if (dependents.value.length > 1) {
    dependents.value.splice(index, 1);
  }
};

// Duplicate checking with debounce
const checkDuplicate = async () => {
  if (!client.first_name && !client.last_name) return;
  
  clearTimeout(duplicateCheckTimeout.value);
  duplicateCheckTimeout.value = setTimeout(async () => {
    const fullName = `${client.first_name} ${client.last_name}`.trim();
    if (!fullName) return;
    
    try {
      const response = await api.get(`/clients/check-duplicate?name=${encodeURIComponent(fullName)}`);
      duplicateWarning.value = response.data.is_duplicate;
      duplicateCount.value = response.data.similar_count;
    } catch (error) {
      console.error('Error checking duplicate:', error);
      duplicateWarning.value = false;
      duplicateCount.value = 0;
    }
  }, 500);
};

// Form submission
const submitForm = async () => {
  if (submitting.value) return;
  
  submitting.value = true;
  try {
    // Prepare the complete data payload
    const formData = {
      client: { ...client },
      income: { ...income },
      loan: { 
        ...loan,
        status: 'Active'
      },
      comakers: comakers.value.filter(comaker => comaker.name.trim() !== ''),
      family: { ...family },
      siblings: siblings.value.filter(sibling => sibling.name.trim() !== ''),
      spouse: spouse.name.trim() !== '' ? { ...spouse } : null,
      dependents: dependents.value.filter(dependent => dependent.name.trim() !== '')
    };

    const response = await api.post('/clients', formData);
    alert('Client created successfully!');
    router.push('/dashboard');
  } catch (error) {
    console.error('Error creating client:', error);
    const errorMessage = error.response?.data?.error || 'Error creating client. Please try again.';
    alert(errorMessage);
  } finally {
    submitting.value = false;
  }
};

const goBack = () => {
  router.push('/dashboard');
};

// Watch for changes in calculator inputs
watch(
  () => [loanCalculator.loanAmount, loanCalculator.term, loanCalculator.deductionType],
  () => {
    calculateLoanDetails();
  },
  { immediate: true }
);

// Watch for name changes to trigger duplicate checking
watch([() => client.first_name, () => client.last_name], () => {
  if (client.first_name || client.last_name) {
    checkDuplicate();
  }
});
</script>

<style scoped>
.client-form {
  padding: 1rem;
  max-width: 1400px;
  margin: 0 auto;
  background-color: #f8f9fa;
  min-height: 100vh;
}

.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.form-header h1 {
  margin: 0;
  color: #2c3e50;
  font-size: 1.5rem;
}

.back-btn {
  padding: 0.5rem 1rem;
  background-color: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.back-btn:hover:not(:disabled) {
  background-color: #5a6268;
}

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

/* Calculator Specific Styles */
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
  align-items: center;
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

.disabled-field {
  background-color: #e9ecef;
  color: #6c757d;
}

/* Form Sections */
.form-section {
  margin-bottom: 2.5rem;
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

.form-section h3 {
  margin: 1rem 0 0.75rem 0;
  color: #495057;
  font-size: 1.1rem;
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
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
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
  opacity: 1;
  cursor: not-allowed;
}

.warning-banner {
  background-color: #fff3cd;
  border: 1px solid #ffeaa7;
  color: #856404;
  padding: 1rem;
  border-radius: 4px;
  margin-top: 1rem;
  font-size: 0.9rem;
}

.warning-banner small {
  font-size: 0.8rem;
  opacity: 0.8;
}

.comaker-section,
.sibling-item,
.dependent-item {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 6px;
  margin-bottom: 1rem;
  border-left: 4px solid #3498db;
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

.btn-add:disabled,
.btn-remove:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

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

/* Responsive Design */
@media (max-width: 1024px) {
  .form-layout {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  .calculator-column {
    position: static;
  }
}

@media (max-width: 768px) {
  .client-form {
    padding: 0.5rem;
  }
  
  .form-header {
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
