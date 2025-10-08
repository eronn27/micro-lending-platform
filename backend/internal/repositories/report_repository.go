package repositories

import (
	"micro-lending-platform/backend/internal/models"
	"gorm.io/gorm"
	"time"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetWeeklyPaymentTotal calculates total payments for current week
func (r *ReportRepository) GetWeeklyPaymentTotal() (float64, error) {
	var total float64
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	
	result := r.db.Model(&models.Payment{}).
		Where("created_at >= ? AND status = ?", weekStart, "Paid").
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&total)
	
	return total, result.Error
}

// GetWeeklyReleaseTotal calculates total loan releases for current week
func (r *ReportRepository) GetWeeklyReleaseTotal() (float64, error) {
	var total float64
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	
	result := r.db.Model(&models.Loan{}).
		Where("created_at >= ?", weekStart).
		Select("COALESCE(SUM(amount_release), 0)").
		Scan(&total)
	
	return total, result.Error
}

// GetTotalClients returns total count of clients
func (r *ReportRepository) GetTotalClients() (int64, error) {
	var count int64
	result := r.db.Model(&models.Client{}).
		Where("deleted_at IS NULL").
		Count(&count)
	return count, result.Error
}

// GetActiveClients returns count of clients with active loans
func (r *ReportRepository) GetActiveClients() (int64, error) {
	var count int64
	result := r.db.Model(&models.Client{}).
		Joins("INNER JOIN loans ON loans.client_id = clients.id").
		Where("loans.status IN ? AND loans.deleted_at IS NULL", []string{"Active"}).
		Distinct("clients.id").
		Count(&count)
	return count, result.Error
}

// GetOverdueClients returns count of clients with overdue loans
func (r *ReportRepository) GetOverdueClients() (int64, error) {
	var count int64
	result := r.db.Model(&models.Client{}).
		Joins("INNER JOIN loans ON loans.client_id = clients.id").
		Where("loans.status IN ? AND loans.deleted_at IS NULL", []string{"Overdue", "Default"}).
		Distinct("clients.id").
		Count(&count)
	return count, result.Error
}

// GetActivePaymentTotal calculates total active payments (pending/partial)
func (r *ReportRepository) GetActivePaymentTotal() (float64, error) {
	var total float64
	result := r.db.Model(&models.Payment{}).
		Where("status IN ? AND deleted_at IS NULL", []string{"Pending", "Partial"}).
		Select("COALESCE(SUM(amount_due), 0)").
		Scan(&total)
	
	return total, result.Error
}

// GetMonthlyPaymentTotal calculates total payments for current month
func (r *ReportRepository) GetMonthlyPaymentTotal() (float64, error) {
	var total float64
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	
	result := r.db.Model(&models.Payment{}).
		Where("created_at >= ? AND status = ?", monthStart, "Paid").
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&total)
	
	return total, result.Error
}

// GetMonthlyReleaseTotal calculates total loan releases for current month
func (r *ReportRepository) GetMonthlyReleaseTotal() (float64, error) {
	var total float64
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	
	result := r.db.Model(&models.Loan{}).
		Where("created_at >= ?", monthStart).
		Select("COALESCE(SUM(amount_release), 0)").
		Scan(&total)
	
	return total, result.Error
}

// GetPaymentsByStatus returns count of payments by status
func (r *ReportRepository) GetPaymentsByStatus(status string) (int64, error) {
	var count int64
	result := r.db.Model(&models.Payment{}).
		Where("status = ? AND deleted_at IS NULL", status).
		Count(&count)
	return count, result.Error
}

// GetLoansByStatus returns count of loans by status
func (r *ReportRepository) GetLoansByStatus(status string) (int64, error) {
	var count int64
	result := r.db.Model(&models.Loan{}).
		Where("status = ? AND deleted_at IS NULL", status).
		Count(&count)
	return count, result.Error
}

// GetTotalOutstandingBalance returns total outstanding balance across all loans
func (r *ReportRepository) GetTotalOutstandingBalance() (float64, error) {
	var total float64
	result := r.db.Model(&models.Loan{}).
		Where("deleted_at IS NULL").
		Select("COALESCE(SUM(outstanding_balance), 0)").
		Scan(&total)
	return total, result.Error
}

// GetTotalLoanAmount returns total amount of all loans
func (r *ReportRepository) GetTotalLoanAmount() (float64, error) {
	var total float64
	result := r.db.Model(&models.Loan{}).
		Where("deleted_at IS NULL").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total)
	return total, result.Error
}
