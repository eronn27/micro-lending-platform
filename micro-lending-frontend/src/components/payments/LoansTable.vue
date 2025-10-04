<template>
  <div class="loans-table">
    <div class="table-responsive">
      <table class="table">
        <thead>
          <tr>
            <th class="select-column">
              <input
                type="checkbox"
                :checked="allSelected"
                @change="$emit('select-all')"
                class="select-checkbox"
              >
            </th>
            <th v-for="column in columns" :key="column.key" :class="column.class">
              {{ column.label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="loan in loans"
            :key="loan.id"
            :class="[
              'loan-row',
              loan.status,
              { selected: isSelected(loan) }
            ]"
          >
            <td class="select-column">
              <input
                type="checkbox"
                :checked="isSelected(loan)"
                @change="$emit('select-loan', loan)"
                class="select-checkbox"
              >
            </td>
            <td class="control-number">
              <strong>{{ loan.control_number }}</strong>
            </td>
            <td class="client-name">
              <div class="name">{{ loan.client.name }}</div>
              <div class="contact">{{ loan.client.contact_number }}</div>
            </td>
            <td class="date-release">
              {{ formatDate(loan.date_of_release) }}
            </td>
            <td class="total-amount">
              {{ formatCurrency(loan.total_amount) }}
            </td>
            <td class="ammortization">
              {{ formatCurrency(loan.ammortization) }}
            </td>
            <td class="terms">
              {{ loan.terms }} month{{ loan.terms > 1 ? 's' : '' }}
            </td>
            <td class="payment-progress">
              <PaymentProgress :progress="loan.payment_progress" />
            </td>
            <td class="status">
              <StatusBadge :status="loan.status" />
            </td>
            <td class="actions">
              <PaymentButton
                :loan="loan"
                :loading="loading"
                @mark-paid="$emit('mark-paid', loan)"
              />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import PaymentProgress from './PaymentProgress.vue'
import StatusBadge from './StatusBadge.vue'
import PaymentButton from './PaymentButton.vue'

defineProps({
  loans: Array,
  selectedLoans: Array,
  loading: Boolean
})

defineEmits(['select-loan', 'select-all', 'mark-paid'])

const columns = [
  { key: 'control_number', label: 'Control #', class: 'control-number' },
  { key: 'name', label: 'Name', class: 'client-name' },
  { key: 'date_release', label: 'Date of Release', class: 'date-release' },
  { key: 'total_amount', label: 'Total Amount', class: 'total-amount' },
  { key: 'ammortization', label: 'Amortization', class: 'ammortization' },
  { key: 'terms', label: 'Terms', class: 'terms' },
  { key: 'progress', label: 'Payment Progress', class: 'payment-progress' },
  { key: 'status', label: 'Status', class: 'status' },
  { key: 'actions', label: 'Action', class: 'actions' }
]

const isSelected = (loan) => {
  return props.selectedLoans.some(selected => selected.id === loan.id)
}

const allSelected = computed(() => {
  return props.loans.length > 0 && props.selectedLoans.length === props.loans.length
})

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-PH')
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-PH', {
    style: 'currency',
    currency: 'PHP'
  }).format(amount)
}
</script>

<style scoped>
.loans-table {
  width: 100%;
}

.table-responsive {
  overflow-x: auto;
}

.table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.table th {
  background: #f8f9fa;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #495057;
  border-bottom: 2px solid #e9ecef;
  white-space: nowrap;
}

.table td {
  padding: 1rem;
  border-bottom: 1px solid #e9ecef;
  vertical-align: middle;
}

.loan-row {
  transition: background-color 0.2s ease;
}

.loan-row:hover {
  background: #f8f9fa;
}

.loan-row.selected {
  background: #e3f2fd;
}

.loan-row.due_today {
  background: #fff3cd;
  border-left: 4px solid #ffc107;
}

.loan-row.overdue {
  background: #f8d7da;
  border-left: 4px solid #dc3545;
  animation: pulse-overdue 2s infinite;
}

@keyframes pulse-overdue {
  0%, 100% { background-color: #f8d7da; }
  50% { background-color: #f5c6cb; }
}

.select-column {
  width: 40px;
  text-align: center;
}

.select-checkbox {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.control-number {
  font-family: 'Courier New', monospace;
  font-weight: 600;
}

.client-name .name {
  font-weight: 600;
  color: #2c3e50;
}

.client-name .contact {
  font-size: 0.75rem;
  color: #6c757d;
  margin-top: 0.25rem;
}

.total-amount,
.ammortization {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  text-align: right;
}

.terms {
  text-align: center;
}

.payment-progress {
  min-width: 120px;
}

.status {
  min-width: 100px;
}

.actions {
  min-width: 120px;
  text-align: center;
}

@media (max-width: 1024px) {
  .table {
    font-size: 0.8rem;
  }

  .table th,
  .table td {
    padding: 0.75rem 0.5rem;
  }
}

@media (max-width: 768px) {
  .client-name .contact {
    display: none;
  }

  .terms {
    display: none;
  }
}
</style>
