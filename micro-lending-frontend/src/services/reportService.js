import { api } from './api'

export const reportService = {
  /**
   * Fetch weekly report data
   * @returns {Promise<Object>} Weekly report data
   */
  async getWeeklyReport() {
    try {
      const response = await api.get('/reports/weekly')
      return response.data
    } catch (error) {
      console.error('Error fetching weekly report:', error)
      throw new Error(
        error.response?.data?.error || 'Failed to fetch weekly reports'
      )
    }
  },

  /**
   * Fetch monthly report data
   * @returns {Promise<Object>} Monthly report data
   */
  async getMonthlyReport() {
    try {
      const response = await api.get('/reports/monthly')
      return response.data
    } catch (error) {
      console.error('Error fetching monthly report:', error)
      throw new Error(
        error.response?.data?.error || 'Failed to fetch monthly reports'
      )
    }
  },

  /**
   * Export report as JSON
   * @param {Object} reportData - Report data to export
   * @param {string} fileName - Name of the file to export
   */
  exportReportAsJSON(reportData, fileName = 'report') {
    const dataStr = JSON.stringify(reportData, null, 2)
    const dataBlob = new Blob([dataStr], { type: 'application/json' })
    const url = URL.createObjectURL(dataBlob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${fileName}-${new Date().toISOString().split('T')[0]}.json`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  },

  /**
   * Export report as CSV
   * @param {Object} reportData - Report data to export
   * @param {string} fileName - Name of the file to export
   */
  exportReportAsCSV(reportData, fileName = 'report') {
    const data = reportData.data
    const headers = [
      'Weekly Payments',
      'Weekly Releases',
      'Total Clients',
      'Active Clients',
      'Overdue Clients',
      'Active Payment Total',
      'Timestamp'
    ]
    const values = [
      data.weekly_payment_total,
      data.weekly_release_total,
      data.total_clients,
      data.active_clients,
      data.overdue_clients,
      data.active_payment_total,
      reportData.timestamp
    ]

    const csvContent = [
      headers.join(','),
      values.join(',')
    ].join('\n')

    const dataBlob = new Blob([csvContent], { type: 'text/csv' })
    const url = URL.createObjectURL(dataBlob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${fileName}-${new Date().toISOString().split('T')[0]}.csv`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }
}
