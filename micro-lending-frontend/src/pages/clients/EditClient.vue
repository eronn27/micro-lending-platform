<template>
  <div class="edit-client">
    <!-- Header -->
    <header class="page-header">
      <button @click="goBack" class="back-btn">← Cancel</button>
      <h1>Edit Client</h1>
      <button @click="saveClient" class="save-btn" :disabled="saving">
        {{ saving ? 'Saving...' : 'Save Changes' }}
      </button>
    </header>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading client data...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Error Loading Client</h3>
      <p>{{ error }}</p>
      <button @click="loadClient" class="retry-btn">Retry</button>
    </div>

    <!-- Edit Form -->
    <div v-else-if="client" class="edit-form">
      <div class="form-section">
        <h2>Basic Information</h2>
        <div class="form-grid">
          <div class="form-group">
            <label>First Name *</label>
            <input type="text" v-model="client.first_name" required>
          </div>
          <div class="form-group">
            <label>Middle Name</label>
            <input type="text" v-model="client.middle_name">
          </div>
          <div class="form-group">
            <label>Last Name *</label>
            <input type="text" v-model="client.last_name" required>
          </div>
          <div class="form-group">
            <label>Contact Number *</label>
            <input type="tel" v-model="client.contact_number" required>
          </div>
          <div class="form-group">
            <label>Age *</label>
            <input type="number" v-model="client.age" required min="18">
          </div>
          <div class="form-group">
            <label>Gender</label>
            <select v-model="client.gender">
              <option value="">Select Gender</option>
              <option value="Male">Male</option>
              <option value="Female">Female</option>
              <option value="Other">Other</option>
            </select>
          </div>
        </div>
      </div>

      <div class="form-section">
        <h2>Address Information</h2>
        <div class="form-group">
          <label>Home Address *</label>
          <textarea v-model="client.home_address" rows="3" required></textarea>
        </div>
        <div class="form-grid">
          <div class="form-group">
            <label>Years of Residence</label>
            <input type="number" v-model="client.years_of_residence" min="0">
          </div>
          <div class="form-group">
            <label>Facebook Account</label>
            <input type="text" v-model="client.facebook_account">
          </div>
        </div>
      </div>

      <div class="form-section">
        <h2>Additional Information</h2>
        <div class="form-grid">
          <div class="form-group">
            <label>Civil Status</label>
            <select v-model="client.civil_status">
              <option value="">Select Status</option>
              <option value="Single">Single</option>
              <option value="Married">Married</option>
              <option value="Divorced">Divorced</option>
              <option value="Widowed">Widowed</option>
            </select>
          </div>
          <div class="form-group">
            <label>Religion</label>
            <input type="text" v-model="client.religion">
          </div>
          <div class="form-group">
            <label>Nickname</label>
            <input type="text" v-model="client.nickname">
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="action-buttons">
        <button @click="goBack" class="btn-secondary">Cancel</button>
        <button @click="saveClient" class="btn-primary" :disabled="saving">
          {{ saving ? 'Saving...' : 'Save Changes' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../../services/api'

const route = useRoute()
const router = useRouter()

const client = ref(null)
const loading = ref(false)
const saving = ref(false)
const error = ref('')

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

const saveClient = async () => {
  if (!client.value) return

  // Basic validation
  if (!client.value.first_name || !client.value.last_name || !client.value.contact_number || !client.value.home_address || !client.value.age) {
    alert('Please fill in all required fields (marked with *)')
    return
  }

  saving.value = true
  
  try {
    await api.put(`/clients/${route.params.id}`, client.value)
    
    // Redirect back to client details after successful save
    router.push(`/clients/${route.params.id}`)
  } catch (err) {
    console.error('Error updating client:', err)
    alert('Failed to update client. Please try again.')
  } finally {
    saving.value = false
  }
}

const goBack = () => {
  router.push(`/clients/${route.params.id}`)
}

onMounted(() => {
  loadClient()
})
</script>

<style scoped>
.edit-client {
  padding: 1rem;
  max-width: 800px;
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
}

.back-btn {
  padding: 0.75rem 1.5rem;
  background-color: #6c757d;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.save-btn {
  padding: 0.75rem 1.5rem;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
}

.save-btn:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.back-btn:hover {
  background-color: #5a6268;
}

.save-btn:hover:not(:disabled) {
  background-color: #218838;
}

/* Form Styles */
.edit-form {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  padding: 2rem;
}

.form-section {
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #e9ecef;
}

.form-section:last-of-type {
  border-bottom: none;
  margin-bottom: 0;
}

.form-section h2 {
  margin: 0 0 1.5rem 0;
  color: #2c3e50;
  font-size: 1.3rem;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 600;
  color: #495057;
  font-size: 0.9rem;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 0.75rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.2s ease;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
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

/* Action Buttons */
.action-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  padding-top: 2rem;
  margin-top: 2rem;
  border-top: 1px solid #e9ecef;
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  min-width: 140px;
}

.btn-primary {
  background-color: #28a745;
  color: white;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #218838;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

.btn-primary:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

/* Responsive Design */
@media (max-width: 768px) {
  .edit-client {
    padding: 0.5rem;
  }
  
  .page-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>
