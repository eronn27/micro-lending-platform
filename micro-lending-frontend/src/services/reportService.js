import { api } from './api'

export const reportService = {
  async getWeeklyReport() {
    try {
      const response = await api.get('/reports/weekly')
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.error || 'Failed to fetch reports')
    }
  },

  async getMonthlyReport() {
    try {
      const response = await api.get('/reports/monthly')
      return response.data
    } catch (error) {
      throw new Error(error.response?.data?.error || 'Failed to fetch monthly reports')
    }
  }
}
