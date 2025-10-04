<template>
  <span :class="['status-badge', status]">
    <span class="badge-icon">{{ statusIcon }}</span>
    <span class="badge-text">{{ statusText }}</span>
  </span>
</template>

<script setup>
const props = defineProps({
  status: String
})

const statusConfig = {
  due_today: { icon: '⏰', text: 'Due Today', class: 'due-today' },
  overdue: { icon: '⚠️', text: 'Overdue', class: 'overdue' },
  active: { icon: '🟢', text: 'Active', class: 'active' },
  paid: { icon: '🎉', text: 'Paid', class: 'paid' }
}

const statusIcon = computed(() => statusConfig[props.status]?.icon || '❓')
const statusText = computed(() => statusConfig[props.status]?.text || 'Unknown')
</script>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.75rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-badge.due-today {
  background: #fff3cd;
  color: #856404;
  border: 1px solid #ffeaa7;
}

.status-badge.overdue {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
  animation: pulse 2s infinite;
}

.status-badge.active {
  background: #d1ecf1;
  color: #0c5460;
  border: 1px solid #bee5eb;
}

.status-badge.paid {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.badge-icon {
  font-size: 0.875rem;
}

.badge-text {
  font-size: 0.7rem;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}
</style>
