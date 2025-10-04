<template>
  <div class="loan-filters">
    <div class="filters-row">
      <!-- Search -->
      <div class="search-box">
        <div class="search-icon">🔍</div>
        <input
          type="text"
          placeholder="Search by name or control number..."
          :value="searchQuery"
          @input="$emit('update:search', $event.target.value)"
          class="search-input"
        >
      </div>

      <!-- Status Filters -->
      <div class="status-filters">
        <button
          v-for="filter in statusFilters"
          :key="filter.value"
          :class="['status-filter', { active: statusFilter === filter.value }]"
          @click="$emit('update:status', filter.value)"
        >
          <span class="filter-badge" :class="filter.value">{{ filter.count }}</span>
          {{ filter.label }}
        </button>
      </div>
    </div>

    <!-- Bulk Actions -->
    <div v-if="selectedCount > 0" class="bulk-actions">
      <div class="bulk-info">
        <strong>{{ selectedCount }}</strong> client{{ selectedCount > 1 ? 's' : '' }} selected
      </div>
      <button @click="$emit('bulk-pay')" class="bulk-pay-btn">
        💳 Mark {{ selectedCount }} Selected as Paid
      </button>
      <button @click="$emit('clear-selection')" class="clear-selection-btn">
        ✕ Clear
      </button>
    </div>

    <!-- Summary -->
    <div class="summary">
      Showing {{ totalCount }} loan{{ totalCount !== 1 ? 's' : '' }}
      <span v-if="searchQuery">matching "{{ searchQuery }}"</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  selectedCount: Number,
  totalCount: Number,
  searchQuery: String,
  statusFilter: String
})

defineEmits(['update:search', 'update:status', 'bulk-pay', 'clear-selection'])

const statusFilters = [
  { value: 'all', label: 'All Loans', count: '' },
  { value: 'due_today', label: 'Due Today', count: '' },
  { value: 'overdue', label: 'Overdue', count: '' },
  { value: 'active', label: 'Active', count: '' },
  { value: 'paid', label: 'Paid-in-Full', count: '' }
]
</script>

<style scoped>
.loan-filters {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  margin-bottom: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.filters-row {
  display: flex;
  gap: 1.5rem;
  align-items: center;
  margin-bottom: 1rem;
}

.search-box {
  position: relative;
  flex: 1;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  font-size: 1.1rem;
  opacity: 0.6;
}

.search-input {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 3rem;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 0.9rem;
  transition: border-color 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.1);
}

.status-filters {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.status-filter {
  padding: 0.5rem 1rem;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  background: white;
  cursor: pointer;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.2s ease;
}

.status-filter:hover {
  border-color: #3498db;
  background: #f8f9fa;
}

.status-filter.active {
  background: #3498db;
  border-color: #3498db;
  color: white;
}

.filter-badge {
  padding: 0.125rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.filter-badge.due_today {
  background: #fff3cd;
  color: #856404;
}

.filter-badge.overdue {
  background: #f8d7da;
  color: #721c24;
}

.filter-badge.active {
  background: #d1ecf1;
  color: #0c5460;
}

.filter-badge.paid {
  background: #d4edda;
  color: #155724;
}

.bulk-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: #e3f2fd;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.bulk-info {
  font-size: 0.9rem;
  color: #1976d2;
}

.bulk-pay-btn {
  background: #4caf50;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: background-color 0.2s ease;
}

.bulk-pay-btn:hover {
  background: #45a049;
}

.clear-selection-btn {
  background: transparent;
  border: 1px solid #dc3545;
  color: #dc3545;
  padding: 0.75rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.clear-selection-btn:hover {
  background: #dc3545;
  color: white;
}

.summary {
  font-size: 0.875rem;
  color: #6c757d;
  text-align: center;
}

@media (max-width: 768px) {
  .filters-row {
    flex-direction: column;
    align-items: stretch;
  }

  .search-box {
    max-width: none;
  }

  .status-filters {
    justify-content: center;
  }

  .bulk-actions {
    flex-direction: column;
    align-items: stretch;
    text-align: center;
  }
}
</style>
