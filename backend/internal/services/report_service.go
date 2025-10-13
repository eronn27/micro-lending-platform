// backend/internal/services/report_service.go
package services

import (
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

// GetWeeklyReport returns weekly report data
func (s *ReportService) GetWeeklyReport() (*models.WeeklyReportData, error) {
	// Implementation for weekly report
	// This should call repository methods to get the actual data
	return &models.WeeklyReportData{
		WeeklyPaymentTotal:   0,
		WeeklyReleaseTotal:   0,
		TotalClients:         0,
		ActiveClients:        0,
		OverdueClients:       0,
		ActivePaymentTotal:   0,
		TotalPaymentThisWeek: 0,
	}, nil
}

// GetMonthlyReport returns monthly report data
func (s *ReportService) GetMonthlyReport() (*models.WeeklyReportData, error) {
	// Implementation for monthly report
	// This should call repository methods to get the actual data
	return &models.WeeklyReportData{
		WeeklyPaymentTotal:   0,
		WeeklyReleaseTotal:   0,
		TotalClients:         0,
		ActiveClients:        0,
		OverdueClients:       0,
		ActivePaymentTotal:   0,
		TotalPaymentThisWeek: 0,
	}, nil
}

// GetHistoricalReport returns historical data for multiple periods
func (s *ReportService) GetHistoricalReport(periodType string, periods int) (*models.HistoricalReportResponse, error) {
	var records []models.HistoricalRecord
	now := time.Now()
	var totalPayments, totalReleases float64

	for i := 0; i < periods; i++ {
		var startDate, endDate time.Time
		var periodLabel string

		if periodType == "weekly" {
			// Calculate weekly period (Sunday to Saturday)
			endDate = now.AddDate(0, 0, -(i * 7))
			// Adjust to Saturday (end of week)
			daysUntilSaturday := (6 - int(endDate.Weekday()) + 7) % 7
			endDate = endDate.AddDate(0, 0, daysUntilSaturday)
			startDate = endDate.AddDate(0, 0, -6)

			periodLabel = startDate.Format("Jan 02") + " - " + endDate.Format("Jan 02, 2006")
		} else {
			// Calculate monthly period (first to last day of month)
			endDate = now.AddDate(0, -i, 0)
			endDate = time.Date(endDate.Year(), endDate.Month()+1, 0, 23, 59, 59, 0, endDate.Location())
			startDate = time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, endDate.Location())

			periodLabel = endDate.Format("January 2006")
		}

		// Get data for this period using repository methods
		payments, _ := s.repo.GetPaymentTotalForPeriod(startDate, endDate)
		releases, _ := s.repo.GetReleaseTotalForPeriod(startDate, endDate)
		activeClients, _ := s.repo.GetActiveClientsForPeriod(endDate)
		overdueClients, _ := s.repo.GetOverdueClientsForPeriod(endDate)

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
		PeriodType:    periodType,
		PeriodsCount:  len(records),
		AveragePayment: avgPayment,
		AverageRelease: avgRelease,
		TotalPayments:  totalPayments,
		TotalReleases:  totalReleases,
	}

	response := &models.HistoricalReportResponse{
		Data:      records,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Status:    "success",
		Metadata:  metadata,
	}

	return response, nil
}
