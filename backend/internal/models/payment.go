package models

import (
    "time"
)

type PaymentStatus string

const (
    PaymentStatusPending PaymentStatus = "Pending"
    PaymentStatusPaid    PaymentStatus = "Paid"
    PaymentStatusOverdue PaymentStatus = "Overdue"
    PaymentStatusPartial PaymentStatus = "Partial"
)

type Payment struct {
    BaseModel
    LoanID       uint          `gorm:"not null;index" json:"loan_id"`
    WeekNumber   int           `gorm:"not null" json:"week_number"`
    PaymentDate  time.Time     `json:"payment_date"`
    AmountDue    float64       `gorm:"type:decimal(10,2);not null" json:"amount_due"`
    AmountPaid   float64       `gorm:"type:decimal(10,2)" json:"amount_paid"`
    Status       PaymentStatus `gorm:"size:20;default:'Pending'" json:"status"`
    PaymentMethod string       `gorm:"size:50" json:"payment_method"`
    
    // Relationship
    Loan Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
}

func (Payment) TableName() string {
    return "payments"
}

// Unique constraint for loan_id and week_number
func (Payment) Constraints() []string {
    return []string{"UNIQUE(loan_id, week_number)"}
}
