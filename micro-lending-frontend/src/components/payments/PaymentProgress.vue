<template>
  <div class="payment-progress">
    <div class="progress-info">
      <span class="weeks-count">{{ progress.paid }}/{{ progress.total }}</span>
      <span class="weeks-label">weeks</span>
    </div>
    <div class="progress-bar">
      <div
        class="progress-fill"
        :style="{ width: `${progress.percentage}%` }"
        :class="progressClass"
      ></div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  progress: Object
})

const progressClass = computed(() => {
  if (props.progress.percentage >= 100) return 'complete'
  if (props.progress.percentage >= 75) return 'high'
  if (props.progress.percentage >= 50) return 'medium'
  return 'low'
})
</script>

<style scoped>
.payment-progress {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.progress-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 50px;
}

.weeks-count {
  font-weight: 600;
  font-size: 0.875rem;
  color: #2c3e50;
}

.weeks-label {
  font-size: 0.7rem;
  color: #6c757d;
  text-transform: uppercase;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: #e9ecef;
  border-radius: 4px;
  overflow: hidden;
  min-width: 80px;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-fill.low {
  background: #dc3545;
}

.progress-fill.medium {
  background: #fd7e14;
}

.progress-fill.high {
  background: #20c997;
}

.progress-fill.complete {
  background: #28a745;
}
</style>
