// backend/internal/services/report_service.go
package services

import (
	"fmt"
	"micro-lending-platform/backend/internal/models"
	"micro-lending-platform/backend/internal/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// GetWeeklyReport returns weekly report data for the current week
func (s *ReportService) GetWeeklyReport() (*models.WeeklyReportData, error) {
	now := time.Now().UTC()
	
	// Get start of current week (Sunday)
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	
	// Get end of current week (Saturday)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)
	endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 0, endOfWeek.Location())

	// Get payment total for this week
	payments, err := s.repo.GetPaymentTotalForPeriod(startOfWeek, endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment total: %w", err)
	}

	// Get release total for this week
	releases, err := s.repo.GetReleaseTotalForPeriod(startOfWeek, endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to get release total: %w", err)
	}

	// Get active clients
	activeClients, err := s.repo.GetActiveClientsForPeriod(endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to get active clients: %w", err)
	}

	// Get overdue clients
	overdueClients, err := s.repo.GetOverdueClientsForPeriod(endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue clients: %w", err)
	}

	return &models.WeeklyReportData{
		WeeklyPaymentTotal:   payments,
		WeeklyReleaseTotal:   releases,
		TotalClients:         activeClients + overdueClients,
		ActiveClients:        activeClients,
		OverdueClients:       overdueClients,
		ActivePaymentTotal:   payments,
		TotalPaymentThisWeek: payments,
	}, nil
}

// GetMonthlyReport returns monthly report data for the current month
func (s *ReportService) GetMonthlyReport() (*models.WeeklyReportData, error) {
	now := time.Now().UTC()
	
	// Get start of current month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	
	// Get end of current month
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	endOfMonth = time.Date(endOfMonth.Year(), endOfMonth.Month(), endOfMonth.Day(), 23, 59, 59, 0, endOfMonth.Location())

	// Get payment total for this month
	payments, err := s.repo.GetPaymentTotalForPeriod(startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment total: %w", err)
	}

	// Get release total for this month
	releases, err := s.repo.GetReleaseTotalForPeriod(startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get release total: %w", err)
	}

	// Get active clients
	activeClients, err := s.repo.GetActiveClientsForPeriod(endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get active clients: %w", err)
	}

	// Get overdue clients
	overdueClients, err := s.repo.GetOverdueClientsForPeriod(endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue clients: %w", err)
	}

	return &models.WeeklyReportData{
		WeeklyPaymentTotal:   payments,
		WeeklyReleaseTotal:   releases,
		TotalClients:         activeClients + overdueClients,
		ActiveClients:        activeClients,
		OverdueClients:       overdueClients,
		ActivePaymentTotal:   payments,
		TotalPaymentThisWeek: payments,
	}, nil
}

// GetHistoricalReport returns historical data for multiple periods
func (s *ReportService) GetHistoricalReport(periodType string, periods int) (*models.HistoricalReportResponse, error) {
	var records []models.HistoricalRecord
	now := time.Now().UTC()
	var totalPayments, totalReleases float64

	for i := 0; i < periods; i++ {
		var startDate, endDate time.Time
		var periodLabel string

		if periodType == "weekly" {
			// Calculate weekly period (Sunday to Saturday)
			daysBack := i * 7
			baseDate := now.AddDate(0, 0, -daysBack)
			
			// Find Sunday (start of week)
			daysUntilSunday := int(baseDate.Weekday())
			startDate = baseDate.AddDate(0, 0, -daysUntilSunday)
			startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
			
			// Saturday is 6 days after Sunday
			endDate = startDate.AddDate(0, 0, 6)
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

			periodLabel = startDate.Format("Jan 02") + " - " + endDate.Format("Jan 02, 2006")
		} else {
			// Calculate monthly period (first to last day of month)
			// Go back i months
			monthDate := now.AddDate(0, -i, 0)
			
			// Set to first day of that month
			startDate = time.Date(monthDate.Year(), monthDate.Month(), 1, 0, 0, 0, 0, monthDate.Location())
			
			// Set to last day of that month
			endDate = startDate.AddDate(0, 1, -1)
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

			periodLabel = endDate.Format("January 2006")
		}

		// Get data for this period using repository methods
		payments, err := s.repo.GetPaymentTotalForPeriod(startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get payments for period %s: %w", periodLabel, err)
		}

		releases, err := s.repo.GetReleaseTotalForPeriod(startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get releases for period %s: %w", periodLabel, err)
		}

		activeClients, err := s.repo.GetActiveClientsForPeriod(endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get active clients for period %s: %w", periodLabel, err)
		}

		overdueClients, err := s.repo.GetOverdueClientsForPeriod(endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get overdue clients for period %s: %w", periodLabel, err)
		}

		netFlow := payments - releases
		totalPayments += payments
		totalReleases += releases

		record := models.HistoricalRecord{
			Period:         periodLabel,
			StartDate:      startDate.Format("2006-01-02"),
			EndDate:        endDate.Format("2006-01-02"),
			Payments:       payments,
			Releases:       releases,
			ActiveClients:  activeClients,
			OverdueClients: overdueClients,
			NetFlow:        netFlow,
		}

		records = append(records, record)
	}

	// Calculate metadata
	avgPayment := 0.0
	avgRelease := 0.0
	if len(records) > 0 {
		avgPayment = totalPayments / float64(len(records))
		avgRelease = totalReleases / float64(len(records))
	}

	metadata := models.HistoryMetadata{
		PeriodType:     periodType,
		PeriodsCount:   len(records),
		AveragePayment: avgPayment,
		AverageRelease: avgRelease,
		TotalPayments:  totalPayments,
		TotalReleases:  totalReleases,
	}

	response := &models.HistoricalReportResponse{
		Data:      records,
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05"),
		Status:    "success",
		Metadata:  metadata,
	}

	return response, nil
}
