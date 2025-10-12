package services

import (
    "fmt"
    "micro-lending-platform/backend/internal/models"
    "micro-lending-platform/backend/internal/repositories"
    "time"
)

type LoanService struct {
    loanRepo   *repositories.LoanRepository
    clientRepo *repositories.ClientRepository  // Add clientRepo
}

func NewLoanService(loanRepo *repositories.LoanRepository, clientRepo *repositories.ClientRepository) *LoanService {
    return &LoanService{
        loanRepo:   loanRepo,
        clientRepo: clientRepo,
    }
}

type LoanStats struct {
    TotalLoans       int64   `json:"total_loans"`
    ActiveLoans      int64   `json:"active_loans"`
    PaidLoans        int64   `json:"paid_loans"`
    OverdueLoans     int64   `json:"overdue_loans"`
    TotalDisbursed   float64 `json:"total_disbursed"`
    TotalOutstanding float64 `json:"total_outstanding"`
}

// CreateLoan creates a new loan - UPDATED to accept ClientID
func (s *LoanService) CreateLoan(req *models.LoanCreate, clientID uint) (*models.Loan, error) {
    // Parse dates
    dateOfRelease, err := s.parseDate(req.DateOfRelease)
    if err != nil {
        return nil, fmt.Errorf("invalid date of release: %w", err)
    }

    applicationDate, err := s.parseDate(req.ApplicationDate)
    if err != nil {
        applicationDate = time.Now()
    }

    // Create loan object - NOW WITH CLIENT ID
    loan := &models.Loan{
        ClientID:              clientID, // Use the provided client ID
        ControlNumber:         req.ControlNumber,
        DateOfRelease:         dateOfRelease,
        TotalAmount:           req.TotalAmount,
        Ammortization:         req.Ammortization,
        Terms:                 req.Terms,
        Mode:                  req.Mode,
        OutstandingBalance:    req.OutstandingBalance,
        Status:                models.LoanStatus(req.Status),
        DueDate:               req.DueDate,
        Deductions:            req.Deductions,
        AmountRelease:         req.AmountRelease,
        PaymentPeriodWeeks:    req.PaymentPeriodWeeks,
        PaidWeeks:             0, // New loan starts at 0 paid weeks
        MethodOfPayment:       req.MethodOfPayment,
        CreditHistory:         req.CreditHistory,
        RecommendedBy:         req.RecommendedBy,
        ApprovedBy:            req.ApprovedBy,
        LoanCycle:             req.LoanCycle,
        RecommendedLoanAmount: req.RecommendedLoanAmount,
        ApprovedLoanAmount:    req.ApprovedLoanAmount,
        CheckedBy:             req.CheckedBy,
        NameCI:                req.NameCI,
        NotedBy:               req.NotedBy,
        ApplicationDate:       applicationDate,
    }

    // Generate control number if not provided
    if loan.ControlNumber == "" {
        loan.ControlNumber = s.generateLoanControlNumber()
    }

    // Set default mode if empty
    if loan.Mode == "" {
        loan.Mode = "Weekly"
    }

    // Set default status if empty
    if loan.Status == "" {
        loan.Status = models.LoanStatusActive
    }

    createdLoan, err := s.loanRepo.Create(loan)
    if err != nil {
        return nil, fmt.Errorf("failed to create loan: %w", err)
    }

    return createdLoan, nil
}

// UpdateLoanFromHandler updates loan information from handler (for ClientID updates)
func (s *LoanService) UpdateLoanFromHandler(loan *models.Loan) (*models.Loan, error) {
    updatedLoan, err := s.loanRepo.Update(loan)
    if err != nil {
        return nil, fmt.Errorf("failed to update loan: %w", err)
    }

    return updatedLoan, nil
}

// GetAllLoans retrieves all loans with pagination and filtering
func (s *LoanService) GetAllLoans(page, limit int, status string) ([]models.Loan, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 20
    }
    
    offset := (page - 1) * limit

    loans, err := s.loanRepo.FindAll(offset, limit, status)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get loans: %w", err)
    }

    total, err := s.loanRepo.Count(status)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count loans: %w", err)
    }

    return loans, total, nil
}

// GetLoanByID retrieves a loan by ID
func (s *LoanService) GetLoanByID(id uint) (*models.Loan, error) {
    loan, err := s.loanRepo.FindByID(id)
    if err != nil {
        if err.Error() == "record not found" {
            return nil, fmt.Errorf("loan not found")
        }
        return nil, fmt.Errorf("failed to get loan: %w", err)
    }
    return loan, nil
}

// GetLoansByClientID retrieves all loans for a specific client
func (s *LoanService) GetLoansByClientID(clientID uint) ([]models.Loan, error) {
    loans, err := s.loanRepo.FindByClientID(clientID)
    if err != nil {
        return nil, fmt.Errorf("failed to get loans: %w", err)
    }
    return loans, nil
}

// UpdateLoan updates loan information
func (s *LoanService) UpdateLoan(req *models.LoanUpdateRequest) (*models.Loan, error) {
    // Check if loan exists
    loan, err := s.loanRepo.FindByID(req.ID)
    if err != nil {
        return nil, fmt.Errorf("loan not found")
    }

    // Update only provided fields
    if req.OutstandingBalance != nil {
        loan.OutstandingBalance = *req.OutstandingBalance
    }
    if req.Status != "" {
        loan.Status = models.LoanStatus(req.Status)
    }
    if req.PaymentPeriodWeeks != nil {
        loan.PaymentPeriodWeeks = *req.PaymentPeriodWeeks
    }
    if req.DueDate != "" {
        loan.DueDate = req.DueDate
    }

    updatedLoan, err := s.loanRepo.Update(loan)
    if err != nil {
        return nil, fmt.Errorf("failed to update loan: %w", err)
    }

    return updatedLoan, nil
}

// DeleteLoan soft deletes a loan
func (s *LoanService) DeleteLoan(id uint) error {
    // Check if loan exists
    _, err := s.loanRepo.FindByID(id)
    if err != nil {
        if err.Error() == "record not found" {
            return fmt.Errorf("loan not found")
        }
        return fmt.Errorf("failed to find loan: %w", err)
    }

    err = s.loanRepo.Delete(id)
    if err != nil {
        return fmt.Errorf("failed to delete loan: %w", err)
    }

    return nil
}

// GetLoanStats returns loan statistics
func (s *LoanService) GetLoanStats() (*LoanStats, error) {
    total, err := s.loanRepo.Count("")
    if err != nil {
        return nil, fmt.Errorf("failed to get total loans: %w", err)
    }

    active, err := s.loanRepo.CountByStatus(models.LoanStatusActive)
    if err != nil {
        return nil, fmt.Errorf("failed to get active loans: %w", err)
    }

    paid, err := s.loanRepo.CountByStatus(models.LoanStatusPaid)
    if err != nil {
        return nil, fmt.Errorf("failed to get paid loans: %w", err)
    }

    overdue, err := s.loanRepo.CountByStatus(models.LoanStatusOverdue)
    if err != nil {
        return nil, fmt.Errorf("failed to get overdue loans: %w", err)
    }

    totalDisbursed, err := s.loanRepo.SumTotalAmount()
    if err != nil {
        return nil, fmt.Errorf("failed to get total disbursed: %w", err)
    }

    totalOutstanding, err := s.loanRepo.SumOutstandingBalance()
    if err != nil {
        return nil, fmt.Errorf("failed to get total outstanding: %w", err)
    }

    return &LoanStats{
        TotalLoans:       total,
        ActiveLoans:      active,
        PaidLoans:        paid,
        OverdueLoans:     overdue,
        TotalDisbursed:   totalDisbursed,
        TotalOutstanding: totalOutstanding,
    }, nil
}

// generateLoanControlNumber generates a unique loan control number
func (s *LoanService) generateLoanControlNumber() string {
    timestamp := time.Now().Unix()
    return fmt.Sprintf("LOAN-%d", timestamp)
}

// parseDate parses date string
func (s *LoanService) parseDate(dateStr string) (time.Time, error) {
    if dateStr == "" {
        return time.Time{}, fmt.Errorf("empty date string")
    }
    return time.Parse("2006-01-02", dateStr)
}

// GetClientsForPayments retrieves clients that need payment attention
func (s *LoanService) GetClientsForPayments() ([]models.Client, error) {
    loans, err := s.loanRepo.GetLoansForPayments()
    if err != nil {
        return nil, err
    }

    // Group loans by client and calculate payment status
    clientsMap := make(map[uint]*models.Client)
    for _, loan := range loans {
        if client, exists := clientsMap[loan.ClientID]; exists {
            client.Loans = append(client.Loans, loan)
        } else {
            clientsMap[loan.ClientID] = &models.Client{
                BaseModel:      loan.Client.BaseModel, // Use BaseModel instead of ID
                ControlNumber:  loan.Client.ControlNumber,
                FirstName:      loan.Client.FirstName,
                LastName:       loan.Client.LastName,
                Loans:          []models.Loan{loan},
            }
        }
    }

    // Convert map to slice
    clients := make([]models.Client, 0, len(clientsMap))
    for _, client := range clientsMap {
        clients = append(clients, *client)
    }

    return clients, nil
}

// GetAllClientsWithLoans retrieves all clients with their loans
func (s *LoanService) GetAllClientsWithLoans() ([]models.Client, error) {
    clients, err := s.clientRepo.FindAll(0, 0, "") // 0,0 means no pagination
    if err != nil {
        return nil, err
    }

    // Load loans for each client
    for i := range clients {
        loans, err := s.loanRepo.FindByClientID(clients[i].ID)
        if err == nil {
            clients[i].Loans = loans
        }
    }

    return clients, nil
}
