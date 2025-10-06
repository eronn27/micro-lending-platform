package services

import (
    "fmt"
    "micro-lending-platform/backend/internal/models"
    "micro-lending-platform/backend/internal/repositories"
    "time"
)

type PaymentService struct {
    paymentRepo *repositories.PaymentRepository
    loanRepo    *repositories.LoanRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository, loanRepo *repositories.LoanRepository) *PaymentService {
    return &PaymentService{
        paymentRepo: paymentRepo,
        loanRepo:    loanRepo,
    }
}

// CreatePayment creates a new payment and updates the loan's outstanding balance
func (s *PaymentService) CreatePayment(req *models.PaymentCreateRequest) (*models.Payment, error) {
    // Validate the loan exists
    loan, err := s.loanRepo.FindByID(req.LoanID)
    if err != nil {
        return nil, fmt.Errorf("loan not found: %w", err)
    }

    // Check if payment for this week already exists
    existingPayment, _ := s.paymentRepo.FindByLoanAndWeek(req.LoanID, req.WeekNumber)
    if existingPayment != nil {
        return nil, fmt.Errorf("payment for week %d already exists", req.WeekNumber)
    }

    // Parse payment date
    paymentDate, err := s.parseDate(req.PaymentDate)
    if err != nil {
        paymentDate = time.Now()
    }

    // Create payment
    payment := &models.Payment{
        LoanID:        req.LoanID,
        WeekNumber:    req.WeekNumber,
        PaymentDate:   paymentDate,
        AmountDue:     req.AmountDue,
        AmountPaid:    req.AmountPaid,
        Status:        models.PaymentStatus(req.Status),
        PaymentMethod: req.PaymentMethod,
    }

    // Create payment in database
    createdPayment, err := s.paymentRepo.Create(payment)
    if err != nil {
        return nil, fmt.Errorf("failed to create payment: %w", err)
    }
// Update loan's outstanding balance and progress
newBalance := loan.OutstandingBalance - req.AmountPaid
if newBalance < 0 {
    newBalance = 0
}

// ✅ INCREMENT PAID WEEKS (use existing paid_weeks + 1)
newPaidWeeks := loan.PaidWeeks + 1

// Determine new loan status
newStatus := loan.Status
if newBalance == 0 || newPaidWeeks >= loan.PaymentPeriodWeeks {
    newStatus = models.LoanStatusPaid
}

// ✅ UPDATE WITH NEW METHOD that saves paid_weeks
err = s.loanRepo.UpdateBalanceAndProgress(req.LoanID, newBalance, newPaidWeeks, newStatus)
    if err != nil {
        // Log error but don't fail the payment creation
        fmt.Printf("Warning: Failed to update loan balance: %v\n", err)
    }

    return createdPayment, nil
}

// GetPaymentsByLoanID retrieves all payments for a specific loan
func (s *PaymentService) GetPaymentsByLoanID(loanID uint) ([]models.Payment, error) {
    payments, err := s.paymentRepo.FindByLoanID(loanID)
    if err != nil {
        return nil, fmt.Errorf("failed to get payments: %w", err)
    }
    return payments, nil
}

// GetPaymentByID retrieves a payment by ID
func (s *PaymentService) GetPaymentByID(id uint) (*models.Payment, error) {
    payment, err := s.paymentRepo.FindByID(id)
    if err != nil {
        if err.Error() == "record not found" {
            return nil, fmt.Errorf("payment not found")
        }
        return nil, fmt.Errorf("failed to get payment: %w", err)
    }
    return payment, nil
}

// UpdatePayment updates payment information
func (s *PaymentService) UpdatePayment(payment *models.Payment) (*models.Payment, error) {
    // Check if payment exists
    _, err := s.paymentRepo.FindByID(payment.ID)
    if err != nil {
        return nil, fmt.Errorf("payment not found")
    }

    updatedPayment, err := s.paymentRepo.Update(payment)
    if err != nil {
        return nil, fmt.Errorf("failed to update payment: %w", err)
    }

    return updatedPayment, nil
}

// DeletePayment soft deletes a payment
func (s *PaymentService) DeletePayment(id uint) error {
    // Check if payment exists
    _, err := s.paymentRepo.FindByID(id)
    if err != nil {
        if err.Error() == "record not found" {
            return fmt.Errorf("payment not found")
        }
        return fmt.Errorf("failed to find payment: %w", err)
    }

    err = s.paymentRepo.Delete(id)
    if err != nil {
        return fmt.Errorf("failed to delete payment: %w", err)
    }

    return nil
}

// GetAllPayments retrieves all payments with pagination
func (s *PaymentService) GetAllPayments(page, limit int) ([]models.Payment, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 50
    }
    
    offset := (page - 1) * limit

    payments, err := s.paymentRepo.FindAll(offset, limit)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get payments: %w", err)
    }

    total, err := s.paymentRepo.Count()
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count payments: %w", err)
    }

    return payments, total, nil
}

// parseDate parses date string
func (s *PaymentService) parseDate(dateStr string) (time.Time, error) {
    if dateStr == "" {
        return time.Time{}, fmt.Errorf("empty date string")
    }
    return time.Parse("2006-01-02", dateStr)
}
