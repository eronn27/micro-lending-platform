<template>
  <button
    :class="[
      'payment-button',
      buttonClass,
      { loading: loading }
    ]"
    :disabled="isDisabled"
    @click="$emit('mark-paid')"
  >
    <span v-if="loading" class="button-loading">⏳</span>
    <span v-else class="button-icon">{{ buttonIcon }}</span>
    <span class="button-text">{{ buttonText }}</span>
    <span v-if="loan.next_due_week" class="week-number">Week {{ loan.next_due_week }}</span>
  </button>
</template>

<script setup>
const props = defineProps({
  loan: Object,
  loading: Boolean
})

defineEmits(['mark-paid'])

const isDisabled = computed(() => {
  return props.loan.status === 'paid' || props.loading
})

const buttonConfig = {
  due_today: { icon: '💳', text: 'Mark Paid', class: 'due-today' },
  overdue: { icon: '⚠️', text: 'Mark Paid', class: 'overdue' },
  active: { icon: '✅', text: 'Mark Paid', class: 'active' },
  paid: { icon: '🎉', text: 'Paid', class: 'paid' }
}

const buttonClass = computed(() => buttonConfig[props.loan.status]?.class || 'active')
const buttonIcon = computed(() => buttonConfig[props.loan.status]?.icon || '✅')
const buttonText = computed(() => buttonConfig[props.loan.status]?.text || 'Mark Paid')
</script>

<style scoped>
.payment-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 80px;
  min-height: 60px;
}

.payment-button:not(:disabled):hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.payment-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.payment-button.due-today {
  background: #28a745;
  color: white;
}

.payment-button.overdue {
  background: #dc3545;
  color: white;
}

.payment-button.active {
  background: #007bff;
  color: white;
}

.payment-button.paid {
  background: #6c757d;
  color: white;
}

.payment-button.loading {
  opacity: 0.7;
  cursor: wait;
}

.button-loading,
.button-icon {
  font-size: 1.1rem;
}

.button-text {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.week-number {
  font-size: 0.6rem;
  opacity: 0.9;
  font-weight: normal;
}
</style>
