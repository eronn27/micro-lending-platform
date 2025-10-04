package models

import (
    "time"
)

type LoanStatus string

const (
    LoanStatusActive  LoanStatus = "Active"
    LoanStatusPaid    LoanStatus = "Paid"
    LoanStatusOverdue LoanStatus = "Overdue"
    LoanStatusDefault LoanStatus = "Default"
)

type Loan struct {
    BaseModel
    ClientID              uint      `gorm:"not null;index" json:"client_id"`
    ControlNumber         string    `gorm:"uniqueIndex;not null;size:20" json:"control_number"`
    DateOfRelease         time.Time `gorm:"not null" json:"date_of_release"`
    TotalAmount           float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
    Ammortization         float64   `gorm:"type:decimal(10,2);not null" json:"ammortization"`
    Terms                 int       `gorm:"not null" json:"terms"`
    Mode                  string    `gorm:"size:20;default:'Weekly'" json:"mode"`
    OutstandingBalance    float64   `gorm:"type:decimal(10,2);not null" json:"outstanding_balance"`
    Status                LoanStatus `gorm:"size:20;default:'Active'" json:"status"`
    DueDate               string    `gorm:"size:20" json:"due_date"`
    Deductions            string    `gorm:"size:100" json:"deductions"`
    AmountRelease         float64   `gorm:"type:decimal(10,2);not null" json:"amount_release"`
    PaymentPeriodWeeks    int       `json:"payment_period_weeks"`
    MethodOfPayment       string    `gorm:"size:50" json:"method_of_payment"`
    CreditHistory         string    `gorm:"size:50" json:"credit_history"`
    RecommendedBy         string    `gorm:"size:100" json:"recommended_by"`
    ApprovedBy            string    `gorm:"size:100" json:"approved_by"`
    LoanCycle             int       `json:"loan_cycle"`
    RecommendedLoanAmount float64   `gorm:"type:decimal(10,2)" json:"recommended_loan_amount"`
    ApprovedLoanAmount    float64   `gorm:"type:decimal(10,2)" json:"approved_loan_amount"`
    CheckedBy             string    `gorm:"size:100" json:"checked_by"`
    NameCI                string    `gorm:"size:100" json:"name_ci"` // ADD THIS MISSING FIELD
    NotedBy               string    `gorm:"size:100" json:"noted_by"`
    ApplicationDate       time.Time `json:"application_date"`
    
    // Relationships
    Client     Client      `gorm:"foreignKey:ClientID" json:"client,omitempty"`
    Payments   []Payment   `gorm:"foreignKey:LoanID" json:"payments,omitempty"`
    CoMakers   []CoMaker   `gorm:"foreignKey:LoanID" json:"co_makers,omitempty"`
}

func (Loan) TableName() string {
    return "loans"
}
