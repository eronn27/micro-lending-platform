package models

// PaymentCreateRequest represents the request to create a payment
type PaymentCreateRequest struct {
    LoanID        uint    `json:"loan_id" binding:"required"`
    WeekNumber    int     `json:"week_number" binding:"required"`
    PaymentDate   string  `json:"payment_date"`
    AmountDue     float64 `json:"amount_due" binding:"required"`
    AmountPaid    float64 `json:"amount_paid" binding:"required"`
    Status        string  `json:"status"`
    PaymentMethod string  `json:"payment_method"`
}

// LoanUpdateRequest represents the request to update a loan
type LoanUpdateRequest struct {
    ID                 uint     `json:"id"`
    OutstandingBalance *float64 `json:"outstanding_balance,omitempty"`
    Status             string   `json:"status,omitempty"`
    PaymentPeriodWeeks *int     `json:"payment_period_weeks,omitempty"`
    DueDate            string   `json:"due_date,omitempty"`
}
