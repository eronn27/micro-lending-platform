// Add these functions to your existing api.js file or create a new loanService.js

import { api } from './api'

/**
 * Loan API Service
 * Handles all loan-related API operations
 */

export const loanService = {
  /**
   * Create a new loan for an existing client
   * @param {Object} loanData - Loan data including client_id, loan details, and comakers
   * @returns {Promise} API response with created loan
   */
  createLoan: async (loanData) => {
    try {
      const response = await api.post('/loans', loanData)
      return response.data
    } catch (error) {
      console.error('Error creating loan:', error)
      throw error
    }
  },

  /**
   * Get all loans for a specific client
   * @param {number} clientId - Client ID
   * @returns {Promise} API response with loan array
   */
  getClientLoans: async (clientId) => {
    try {
      const response = await api.get(`/loans/client/${clientId}`)
      return response.data
    } catch (error) {
      console.error('Error fetching client loans:', error)
      throw error
    }
  },

  /**
   * Get a specific loan by ID
   * @param {number} loanId - Loan ID
   * @returns {Promise} API response with loan details
   */
  getLoanById: async (loanId) => {
    try {
      const response = await api.get(`/loans/${loanId}`)
      return response.data
    } catch (error) {
      console.error('Error fetching loan:', error)
      throw error
    }
  },

  /**
   * Update an existing loan
   * @param {number} loanId - Loan ID
   * @param {Object} loanData - Updated loan data
   * @returns {Promise} API response with updated loan
   */
  updateLoan: async (loanId, loanData) => {
    try {
      const response = await api.put(`/loans/${loanId}`, loanData)
      return response.data
    } catch (error) {
      console.error('Error updating loan:', error)
      throw error
    }
  },

  /**
   * Delete a loan
   * @param {number} loanId - Loan ID
   * @returns {Promise} API response
   */
  deleteLoan: async (loanId) => {
    try {
      const response = await api.delete(`/loans/${loanId}`)
      return response.data
    } catch (error) {
      console.error('Error deleting loan:', error)
      throw error
    }
  },

  /**
   * Get all loans with optional filters
   * @param {Object} params - Query parameters (status, offset, limit)
   * @returns {Promise} API response with loans array
   */
  getAllLoans: async (params = {}) => {
    try {
      const response = await api.get('/loans', { params })
      return response.data
    } catch (error) {
      console.error('Error fetching loans:', error)
      throw error
    }
  },

  /**
   * Get loan statistics
   * @returns {Promise} API response with loan statistics
   */
  getLoanStats: async () => {
    try {
      const response = await api.get('/loans/stats')
      return response.data
    } catch (error) {
      console.error('Error fetching loan stats:', error)
      throw error
    }
  },

  /**
   * Calculate loan cycle for a client
   * @param {number} clientId - Client ID
   * @returns {Promise} Calculated loan cycle number
   */
  calculateLoanCycle: async (clientId) => {
    try {
      const loans = await loanService.getClientLoans(clientId)
      return loans.length + 1
    } catch (error) {
      console.error('Error calculating loan cycle:', error)
      return 1 // Default to 1 if error occurs
    }
  }
}

// Export as default for easy importing
export default loanService
