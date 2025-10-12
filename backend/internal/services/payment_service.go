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

// GetPartialPaymentsByLoanAndWeek retrieves partial payments for a specific loan and week
func (s *PaymentService) GetPartialPaymentsByLoanAndWeek(loanID uint, weekNumber int) ([]models.Payment, error) {
    payments, err := s.paymentRepo.FindPartialsByLoanAndWeek(loanID, weekNumber)
    if err != nil {
        return nil, fmt.Errorf("failed to get partial payments: %w", err)
    }
    return payments, nil
}


// GetNextPaymentWeek calculates the next available week for payment
func (s *PaymentService) GetNextPaymentWeek(loanID uint) (int, error) {
    loan, err := s.loanRepo.FindByID(loanID)
    if err != nil {
        return 0, fmt.Errorf("loan not found: %w", err)
    }

    // Start from the week after the last paid week
    nextWeek := loan.PaidWeeks + 1

    // Check if there are already payments for this week
    for {
        existingPayment, _ := s.paymentRepo.FindByLoanAndWeekWithStatus(loanID, nextWeek, models.PaymentStatusPaid)
        if existingPayment == nil {
            break
        }
        nextWeek++
    }

    return nextWeek, nil
}

// parseDate parses date string
func (s *PaymentService) parseDate(dateStr string) (time.Time, error) {
    if dateStr == "" {
        return time.Time{}, fmt.Errorf("empty date string")
    }
    return time.Parse("2006-01-02", dateStr)
}


// CreatePayment creates a new payment and updates the loan's outstanding balance
func (s *PaymentService) CreatePayment(req *models.PaymentCreateRequest) (*models.Payment, error) {
    // Validate the loan exists
    loan, err := s.loanRepo.FindByID(req.LoanID)
    if err != nil {
        return nil, fmt.Errorf("loan not found: %w", err)
    }

    // Parse payment date
    paymentDate, err := s.parseDate(req.PaymentDate)
    if err != nil {
        paymentDate = time.Now()
    }

    // Determine the correct week number
    weekNumber := req.WeekNumber
    if weekNumber == 0 {
        weekNumber = loan.PaidWeeks + 1
    }

    // Calculate remaining balance for this payment
    remainingBalance := 0.0
    if req.IsPartial {
        // For partial payments, calculate remaining balance after this payment
        currentRemaining, err := s.CalculateRemainingBalance(req.LoanID, weekNumber)
        if err != nil {
            return nil, fmt.Errorf("failed to calculate remaining balance: %w", err)
        }
        remainingBalance = currentRemaining - req.AmountPaid
        if remainingBalance < 0 {
            remainingBalance = 0
        }
    } else {
        // For full payments, check if there's already a full payment for this week
        existingFullPayment, _ := s.paymentRepo.FindFullPaymentByLoanAndWeek(req.LoanID, weekNumber)
        if existingFullPayment != nil {
            return nil, fmt.Errorf("full payment already exists for week %d", weekNumber)
        }
    }

    // Create payment with calculated remaining balance
    payment := &models.Payment{
        LoanID:          req.LoanID,
        WeekNumber:      weekNumber,
        PaymentDate:     paymentDate,
        AmountDue:       req.AmountDue,
        AmountPaid:      req.AmountPaid,
        RemainingBalance: remainingBalance,
        Status:          models.PaymentStatus(req.Status),
        PaymentMethod:   req.PaymentMethod,
        IsPartial:       req.IsPartial,
        CompletesWeek:   req.CompletesWeek,
    }

    // Create payment in database
    createdPayment, err := s.paymentRepo.Create(payment)
    if err != nil {
        return nil, fmt.Errorf("failed to create payment: %w", err)
    }

    // Update loan balance and progress
    err = s.updateLoanAfterPayment(loan, payment)
    if err != nil {
        fmt.Printf("Warning: Failed to update loan: %v\n", err)
    }

    return createdPayment, nil
}

// updateLoanAfterPayment handles loan updates after payment creation
func (s *PaymentService) updateLoanAfterPayment(loan *models.Loan, payment *models.Payment) error {
    // Calculate new outstanding balance
    newBalance := loan.OutstandingBalance - payment.AmountPaid
    if newBalance < 0 {
        newBalance = 0
    }

    newPaidWeeks := loan.PaidWeeks
    newStatus := loan.Status

    // Check if this payment completes a week
    if payment.CompletesWeek || (!payment.IsPartial && payment.Status == models.PaymentStatusPaid) {
        // Payment completes the current week
        if payment.WeekNumber > newPaidWeeks {
            newPaidWeeks = payment.WeekNumber
        } else {
            newPaidWeeks = loan.PaidWeeks + 1
        }
    } else if payment.IsPartial {
        // For partial payments, check if accumulated payments complete the week
        weekCompleted, err := s.checkIfWeekCompleted(loan.ID, payment.WeekNumber, loan.Ammortization)
        if err != nil {
            return err
        }
        if weekCompleted {
            newPaidWeeks = payment.WeekNumber
            // Mark all partial payments for this week as completing the week
            s.markWeekAsCompleted(loan.ID, payment.WeekNumber)
        }
    }

    // Update loan status if fully paid
    if newBalance == 0 || newPaidWeeks >= loan.PaymentPeriodWeeks {
        newStatus = models.LoanStatusPaid
    }

    // Update loan in database
    return s.loanRepo.UpdateBalanceAndProgress(loan.ID, newBalance, newPaidWeeks, newStatus)
}

// checkIfWeekCompleted checks if accumulated payments complete the week
func (s *PaymentService) checkIfWeekCompleted(loanID uint, weekNumber int, amortization float64) (bool, error) {
    partialPayments, err := s.paymentRepo.FindPartialsByLoanAndWeek(loanID, weekNumber)
    if err != nil {
        return false, err
    }

    totalPaid := 0.0
    for _, payment := range partialPayments {
        totalPaid += payment.AmountPaid
    }

    return totalPaid >= amortization, nil
}

// markWeekAsCompleted marks all partial payments for a week as completing the week
func (s *PaymentService) markWeekAsCompleted(loanID uint, weekNumber int) error {
    partialPayments, err := s.paymentRepo.FindPartialsByLoanAndWeek(loanID, weekNumber)
    if err != nil {
        return err
    }

    for _, payment := range partialPayments {
        payment.CompletesWeek = true
        _, err := s.paymentRepo.Update(&payment)
        if err != nil {
            return err
        }
    }

    return nil
}

// CalculateRemainingBalance calculates remaining balance for current week
func (s *PaymentService) CalculateRemainingBalance(loanID uint, weekNumber int) (float64, error) {
    loan, err := s.loanRepo.FindByID(loanID)
    if err != nil {
        return 0, fmt.Errorf("loan not found: %w", err)
    }

    partialPayments, err := s.paymentRepo.FindPartialsByLoanAndWeek(loanID, weekNumber)
    if err != nil {
        return 0, err
    }

    totalPaid := 0.0
    for _, payment := range partialPayments {
        totalPaid += payment.AmountPaid
    }

    remaining := loan.Ammortization - totalPaid
    if remaining < 0 {
        remaining = 0
    }

    return remaining, nil
}

// GetPaymentProgress returns detailed payment progress for a loan
func (s *PaymentService) GetPaymentProgress(loanID uint) (*PaymentProgress, error) {
    loan, err := s.loanRepo.FindByID(loanID)
    if err != nil {
        return nil, fmt.Errorf("loan not found: %w", err)
    }

    currentWeek := loan.PaidWeeks + 1
    remainingBalance, err := s.CalculateRemainingBalance(loanID, currentWeek)
    if err != nil {
        return nil, err
    }

    partialPayments, err := s.paymentRepo.FindPartialsByLoanAndWeek(loanID, currentWeek)
    if err != nil {
        return nil, err
    }

    return &PaymentProgress{
        LoanID:           loanID,
        CurrentWeek:      currentWeek,
        PaidWeeks:        loan.PaidWeeks,
        TotalWeeks:       loan.PaymentPeriodWeeks,
        Amortization:     loan.Ammortization,
        RemainingBalance: remainingBalance,
        PartialPayments:  partialPayments,
        IsWeekCompleted:  remainingBalance == 0,
    }, nil
}

type PaymentProgress struct {
    LoanID           uint            `json:"loan_id"`
    CurrentWeek      int             `json:"current_week"`
    PaidWeeks        int             `json:"paid_weeks"`
    TotalWeeks       int             `json:"total_weeks"`
    Amortization     float64         `json:"amortization"`
    RemainingBalance float64         `json:"remaining_balance"`
    PartialPayments  []models.Payment `json:"partial_payments"`
    IsWeekCompleted  bool            `json:"is_week_completed"`
}
