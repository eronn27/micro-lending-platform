package repositories

import (
    "micro-lending-platform/backend/internal/models"
    "gorm.io/gorm"
    "fmt"
)

type PaymentRepository struct {
    db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
    return &PaymentRepository{db: db}
}

// Create inserts a new payment into the database
func (r *PaymentRepository) Create(payment *models.Payment) (*models.Payment, error) {
    result := r.db.Create(payment)
    if result.Error != nil {
        return nil, result.Error
    }
    return payment, nil
}

// FindByID finds a payment by ID
func (r *PaymentRepository) FindByID(id uint) (*models.Payment, error) {
    var payment models.Payment
    result := r.db.Preload("Loan").First(&payment, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &payment, nil
}

// FindByLoanID retrieves all payments for a specific loan
func (r *PaymentRepository) FindByLoanID(loanID uint) ([]models.Payment, error) {
    var payments []models.Payment
    result := r.db.Where("loan_id = ?", loanID).
        Order("week_number ASC").
        Find(&payments)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return payments, nil
}


// FindAll retrieves all payments with pagination
func (r *PaymentRepository) FindAll(offset, limit int) ([]models.Payment, error) {
    var payments []models.Payment
    
    query := r.db.Preload("Loan")
    
    if limit > 0 {
        query = query.Offset(offset).Limit(limit)
    }
    
    result := query.Order("created_at DESC").Find(&payments)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return payments, nil
}

// Update updates an existing payment
func (r *PaymentRepository) Update(payment *models.Payment) (*models.Payment, error) {
    result := r.db.Save(payment)
    if result.Error != nil {
        return nil, result.Error
    }
    return payment, nil
}

// Delete soft deletes a payment
func (r *PaymentRepository) Delete(id uint) error {
    result := r.db.Delete(&models.Payment{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("payment not found")
    }
    return nil
}

// Count returns the total number of payments
func (r *PaymentRepository) Count() (int64, error) {
    var count int64
    result := r.db.Model(&models.Payment{}).Count(&count)
    return count, result.Error
}

// CountPaidByLoanID counts paid payments for a specific loan
func (r *PaymentRepository) CountPaidByLoanID(loanID uint) (int64, error) {
    var count int64
    result := r.db.Model(&models.Payment{}).
        Where("loan_id = ? AND status = ?", loanID, "Paid").
        Count(&count)
    return count, result.Error
}


// FindByLoanAndWeek finds a payment by loan ID and week number
func (r *PaymentRepository) FindByLoanAndWeek(loanID uint, weekNumber int) (*models.Payment, error) {
    var payment models.Payment
    result := r.db.Where("loan_id = ? AND week_number = ?", loanID, weekNumber).First(&payment)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error
    }
    return &payment, nil
}

// NEW: Find full payment (non-partial) for loan and week
func (r *PaymentRepository) FindFullPaymentByLoanAndWeek(loanID uint, weekNumber int) (*models.Payment, error) {
    var payment models.Payment
    result := r.db.Where("loan_id = ? AND week_number = ? AND is_partial = ?",
        loanID, weekNumber, false).First(&payment)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error
    }
    return &payment, nil
}

// FindPartialsByLoanAndWeek finds partial payments for a specific loan and week
func (r *PaymentRepository) FindPartialsByLoanAndWeek(loanID uint, weekNumber int) ([]models.Payment, error) {
    var payments []models.Payment
    result := r.db.Where("loan_id = ? AND week_number = ? AND is_partial = ?",
        loanID, weekNumber, true).Find(&payments)
    if result.Error != nil {
        return nil, result.Error
    }
    return payments, nil
}

// FindByLoanAndWeekWithStatus finds a payment by loan ID, week number, and status
func (r *PaymentRepository) FindByLoanAndWeekWithStatus(loanID uint, weekNumber int, status models.PaymentStatus) (*models.Payment, error) {
    var payment models.Payment
    result := r.db.Where("loan_id = ? AND week_number = ? AND status = ?",
        loanID, weekNumber, status).First(&payment)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error
    }
    return &payment, nil
}

// FindCurrentWeekPartialPayments finds partial payments for the current week (paid_weeks + 1)
func (r *PaymentRepository) FindCurrentWeekPartialPayments(loanID uint, currentWeek int) ([]models.Payment, error) {
    var payments []models.Payment
    result := r.db.Where("loan_id = ? AND week_number = ? AND is_partial = ?",
        loanID, currentWeek, true).Order("created_at ASC").Find(&payments)
    if result.Error != nil {
        return nil, result.Error
    }
    return payments, nil
}

// GetCurrentWeekRemainingBalance calculates remaining balance for current week
func (r *PaymentRepository) GetCurrentWeekRemainingBalance(loanID uint, currentWeek int) (float64, error) {
    var partialPayments []models.Payment
    result := r.db.Where("loan_id = ? AND week_number = ? AND is_partial = ?",
        loanID, currentWeek, true).Find(&partialPayments)
    if result.Error != nil {
        return 0, result.Error
    }

    totalPaid := 0.0
    for _, payment := range partialPayments {
        totalPaid += payment.AmountPaid
    }

    // Get loan to know the amortization amount
    var loan models.Loan
    if err := r.db.First(&loan, loanID).Error; err != nil {
        return 0, err
    }

    remaining := loan.Ammortization - totalPaid
    if remaining < 0 {
        remaining = 0
    }

    return remaining, nil
}

// FindLatestPartialPayment finds the most recent partial payment for a loan
func (r *PaymentRepository) FindLatestPartialPayment(loanID uint, weekNumber int) (*models.Payment, error) {
    var payment models.Payment
    result := r.db.Where("loan_id = ? AND week_number = ? AND is_partial = ?",
        loanID, weekNumber, true).Order("created_at DESC").First(&payment)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error
    }
    return &payment, nil
}
