package services

import (
	"fmt"
	"micro-lending-platform/backend/internal/models"
	"micro-lending-platform/backend/internal/repositories"
)

type ReportService struct {
	reportRepo *repositories.ReportRepository
}

func NewReportService(reportRepo *repositories.ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

// GetWeeklyReport generates weekly report data
func (s *ReportService) GetWeeklyReport() (*models.WeeklyReportData, error) {
	// Get weekly payment total
	weeklyPayment, err := s.reportRepo.GetWeeklyPaymentTotal()
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly payment total: %w", err)
	}

	// Get weekly release total
	weeklyRelease, err := s.reportRepo.GetWeeklyReleaseTotal()
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly release total: %w", err)
	}

	// Get total clients
	totalClients, err := s.reportRepo.GetTotalClients()
	if err != nil {
		return nil, fmt.Errorf("failed to get total clients: %w", err)
	}

	// Get active clients
	activeClients, err := s.reportRepo.GetActiveClients()
	if err != nil {
		return nil, fmt.Errorf("failed to get active clients: %w", err)
	}

	// Get overdue clients
	overdueClients, err := s.reportRepo.GetOverdueClients()
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue clients: %w", err)
	}

	// Get active payment total
	activePayment, err := s.reportRepo.GetActivePaymentTotal()
	if err != nil {
		return nil, fmt.Errorf("failed to get active payment total: %w", err)
	}

	return &models.WeeklyReportData{
		WeeklyPaymentTotal:   weeklyPayment,
		WeeklyReleaseTotal:   weeklyRelease,
		TotalClients:         totalClients,
		ActiveClients:        activeClients,
		OverdueClients:       overdueClients,
		ActivePaymentTotal:   activePayment,
		TotalPaymentThisWeek: weeklyPayment,
	}, nil
}

// GetMonthlyReport generates monthly report data
func (s *ReportService) GetMonthlyReport() (*models.WeeklyReportData, error) {
	// Get monthly payment total
	monthlyPayment, err := s.reportRepo.GetMonthlyPaymentTotal()
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly payment total: %w", err)
	}

	// Get monthly release total
	monthlyRelease, err := s.reportRepo.GetMonthlyReleaseTotal()
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly release total: %w", err)
	}

	// Get total clients
	totalClients, err := s.reportRepo.GetTotalClients()
	if err != nil {
		return nil, fmt.Errorf("failed to get total clients: %w", err)
	}

	// Get active clients
	activeClients, err := s.reportRepo.GetActiveClients()
	if err != nil {
		return nil, fmt.Errorf("failed to get active clients: %w", err)
	}

	// Get overdue clients
	overdueClients, err := s.reportRepo.GetOverdueClients()
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue clients: %w", err)
	}

	// Get active payment total
	activePayment, err := s.reportRepo.GetActivePaymentTotal()
	if err != nil {
		return nil, fmt.Errorf("failed to get active payment total: %w", err)
	}

	return &models.WeeklyReportData{
		WeeklyPaymentTotal:   monthlyPayment,
		WeeklyReleaseTotal:   monthlyRelease,
		TotalClients:         totalClients,
		ActiveClients:        activeClients,
		OverdueClients:       overdueClients,
		ActivePaymentTotal:   activePayment,
		TotalPaymentThisWeek: monthlyPayment,
	}, nil
}
