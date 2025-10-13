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
func (r *ReportRepository) GetPaymentTotalForPeriod(startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.Table("payments").
		Where("status = ? AND created_at BETWEEN ? AND ?", "Paid", startDate, endDate).
		Where("deleted_at IS NULL").
		Select("COALESCE(SUM(amount_paid), 0)").
		Scan(&total).Error
	
	return total, err
}

// GetReleaseTotalForPeriod returns total loan releases within a date range
func (r *ReportRepository) GetReleaseTotalForPeriod(startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.Table("loans").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("deleted_at IS NULL").
		Select("COALESCE(SUM(amount_release), 0)").
		Scan(&total).Error
	
	return total, err
}

// GetActiveClientsForPeriod returns count of active clients at a specific date
func (r *ReportRepository) GetActiveClientsForPeriod(endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Table("clients").
		Joins("INNER JOIN loans ON loans.client_id = clients.id").
		Where("loans.status IN (?)", []string{"Active", "Pending"}).
		Where("loans.created_at <= ?", endDate).
		Where("clients.deleted_at IS NULL AND loans.deleted_at IS NULL").
		Distinct("clients.id").
		Count(&count).Error
	
	return count, err
}

// GetOverdueClientsForPeriod returns count of overdue clients at a specific date
func (r *ReportRepository) GetOverdueClientsForPeriod(endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Table("clients").
		Joins("INNER JOIN loans ON loans.client_id = clients.id").
		Where("loans.status IN (?)", []string{"Overdue", "Default"}).
		Where("loans.created_at <= ?", endDate).
		Where("clients.deleted_at IS NULL AND loans.deleted_at IS NULL").
		Distinct("clients.id").
		Count(&count).Error
	
	return count, err
}
