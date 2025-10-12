package repositories

import (
    "micro-lending-platform/backend/internal/models"
    "gorm.io/gorm"
    "fmt"
    "time"
)

type LoanRepository struct {
    db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) *LoanRepository {
    return &LoanRepository{db: db}
}

// Create inserts a new loan into the database
func (r *LoanRepository) Create(loan *models.Loan) (*models.Loan, error) {
    result := r.db.Create(loan)
    if result.Error != nil {
        return nil, result.Error
    }
    return loan, nil
}

// FindByID finds a loan by ID
func (r *LoanRepository) FindByID(id uint) (*models.Loan, error) {
    var loan models.Loan
    result := r.db.Preload("Payments").Preload("CoMakers").First(&loan, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &loan, nil
}

// FindAll retrieves all loans with pagination and status filtering
func (r *LoanRepository) FindAll(offset, limit int, status string) ([]models.Loan, error) {
    var loans []models.Loan
    
    query := r.db.Preload("Client").Preload("Payments")
    
    // Add status filter if provided
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Apply pagination if limits are specified
    if limit > 0 {
        query = query.Offset(offset).Limit(limit)
    }
    
    result := query.Order("created_at DESC").Find(&loans)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return loans, nil
}

// FindByClientID retrieves all loans for a specific client
func (r *LoanRepository) FindByClientID(clientID uint) ([]models.Loan, error) {
    var loans []models.Loan
    result := r.db.Preload("Payments").
        Where("client_id = ?", clientID).
        Order("created_at DESC").
        Find(&loans)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return loans, nil
}

// Update updates an existing loan
func (r *LoanRepository) Update(loan *models.Loan) (*models.Loan, error) {
    result := r.db.Save(loan)
    if result.Error != nil {
        return nil, result.Error
    }
    return loan, nil
}

// UpdateBalance updates the outstanding balance and status of a loan
func (r *LoanRepository) UpdateBalance(loanID uint, newBalance float64, newStatus models.LoanStatus) error {
    result := r.db.Model(&models.Loan{}).
        Where("id = ?", loanID).
        Updates(map[string]interface{}{
            "outstanding_balance": newBalance,
            "status":              newStatus,
        })
    
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return fmt.Errorf("loan not found")
    }
    
    return nil
}

// Delete soft deletes a loan
func (r *LoanRepository) Delete(id uint) error {
    result := r.db.Delete(&models.Loan{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("loan not found")
    }
    return nil
}

// Count returns the total number of loans matching the status filter
func (r *LoanRepository) Count(status string) (int64, error) {
    var count int64
    
    query := r.db.Model(&models.Loan{})
    
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    result := query.Count(&count)
    return count, result.Error
}

// CountByStatus counts loans by specific status
func (r *LoanRepository) CountByStatus(status models.LoanStatus) (int64, error) {
    var count int64
    result := r.db.Model(&models.Loan{}).
        Where("status = ?", status).
        Count(&count)
    return count, result.Error
}

// SumTotalAmount calculates the sum of all loan amounts
func (r *LoanRepository) SumTotalAmount() (float64, error) {
    var total float64
    result := r.db.Model(&models.Loan{}).
        Select("COALESCE(SUM(total_amount), 0)").
        Scan(&total)
    return total, result.Error
}

// SumOutstandingBalance calculates the sum of all outstanding balances
func (r *LoanRepository) SumOutstandingBalance() (float64, error) {
    var total float64
    result := r.db.Model(&models.Loan{}).
        Select("COALESCE(SUM(outstanding_balance), 0)").
        Scan(&total)
    return total, result.Error
}

// FindActiveLoans retrieves all active loans
func (r *LoanRepository) FindActiveLoans() ([]models.Loan, error) {
    var loans []models.Loan
    result := r.db.Preload("Client").
        Where("status = ?", models.LoanStatusActive).
        Order("due_date ASC").
        Find(&loans)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return loans, nil
}

// FindOverdueLoans retrieves all overdue loans
func (r *LoanRepository) FindOverdueLoans() ([]models.Loan, error) {
    var loans []models.Loan
    result := r.db.Preload("Client").
        Where("status = ?", models.LoanStatusOverdue).
        Order("due_date ASC").
        Find(&loans)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return loans, nil
}

// FindByControlNumber finds a loan by control number
func (r *LoanRepository) FindByControlNumber(controlNumber string) (*models.Loan, error) {
    var loan models.Loan
    result := r.db.Preload("Client").Preload("Payments").
        Where("control_number = ?", controlNumber).
        First(&loan)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return &loan, nil
}



// FindWithPartialPayments retrieves loans with their partial payments
func (r *LoanRepository) FindWithPartialPayments(offset, limit int) ([]models.Loan, error) {
    var loans []models.Loan

    query := r.db.Preload("Client").Preload("Payments", "status = ?", models.PaymentStatusPartial)

    if limit > 0 {
        query = query.Offset(offset).Limit(limit)
    }

    result := query.Order("created_at DESC").Find(&loans)

    if result.Error != nil {
        return nil, result.Error
    }
    return loans, nil
}

// UpdateBalanceAndProgress updates loan balance and paid weeks
func (r *LoanRepository) UpdateBalanceAndProgress(loanID uint, balance float64, paidWeeks int, status models.LoanStatus) error {
    result := r.db.Model(&models.Loan{}).
        Where("id = ?", loanID).
        Updates(map[string]interface{}{
            "outstanding_balance": balance,
            "paid_weeks": paidWeeks,
            "status": status,
            "updated_at": time.Now(),
        })

    return result.Error
}

// GetLoansForPayments retrieves loans that need payment attention
func (r *LoanRepository) GetLoansForPayments() ([]models.Loan, error) {
    var loans []models.Loan

    // Get loans that are active and not fully paid
    result := r.db.Preload("Client").
        Preload("Payments", func(db *gorm.DB) *gorm.DB {
            return db.Where("status IN ?", []models.PaymentStatus{models.PaymentStatusPartial, models.PaymentStatusPaid})
        }).
        Where("status = ? AND outstanding_balance > 0", models.LoanStatusActive).
        Order("created_at ASC").
        Find(&loans)

    if result.Error != nil {
        return nil, result.Error
    }
    return loans, nil
}
