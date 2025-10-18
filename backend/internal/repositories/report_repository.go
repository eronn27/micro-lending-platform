// backend/internal/repositories/report_repository.go
package repositories

import (
	"gorm.io/gorm"
	"time"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetPaymentTotalForPeriod returns total payments within a date range
// Filters for 'Paid' status payments and excludes soft-deleted records
func (r *ReportRepository) GetPaymentTotalForPeriod(startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.Table("payments").
		Where("status = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL", "Paid", startDate, endDate).
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&total).Error
	
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetReleaseTotalForPeriod returns total loan releases within a date range
// Sums all loan amounts released in the period
func (r *ReportRepository) GetReleaseTotalForPeriod(startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.Table("loans").
		Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).
		Select("COALESCE(SUM(amount_release), 0)").
		Scan(&total).Error
	
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetActiveClientsForPeriod returns count of distinct clients with active loans at a specific date
// Active loans are those with status 'Active' or 'Pending'
func (r *ReportRepository) GetActiveClientsForPeriod(endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Table("clients").
		Joins("INNER JOIN loans ON loans.client_id = clients.id").
		Where("loans.status IN (?, ?)", "Active", "Pending").
		Where("loans.created_at <= ?", endDate).
		Where("clients.deleted_at IS NULL AND loans.deleted_at IS NULL").
		Distinct("clients.id").
		Count(&count).Error
	
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetOverdueClientsForPeriod returns count of distinct clients with overdue/default loans at a specific date
// Overdue loans are those with status 'Overdue' or 'Default'
func (r *ReportRepository) GetOverdueClientsForPeriod(endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Table("clients").
		Joins("INNER JOIN loans ON loans.client_id = clients.id").
		Where("loans.status IN (?, ?)", "Overdue", "Default").
		Where("loans.created_at <= ?", endDate).
		Where("clients.deleted_at IS NULL AND loans.deleted_at IS NULL").
		Distinct("clients.id").
		Count(&count).Error
	
	if err != nil {
		return 0, err
	}
	return count, nil
}
