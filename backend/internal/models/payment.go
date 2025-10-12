package models

import (
    "time"
    "gorm.io/gorm"
)

type PaymentStatus string

const (
    PaymentStatusPending  PaymentStatus = "Pending"
    PaymentStatusPaid     PaymentStatus = "Paid"
    PaymentStatusPartial  PaymentStatus = "Partial"
    PaymentStatusOverdue  PaymentStatus = "Overdue"
)

type Payment struct {
    ID              uint          `json:"id" gorm:"primaryKey"`
    LoanID          uint          `json:"loan_id"`
    WeekNumber      int           `json:"week_number"`
    PaymentDate     time.Time     `json:"payment_date"`
    AmountDue       float64       `json:"amount_due" gorm:"type:decimal(10,2)"`
    AmountPaid      float64       `json:"amount_paid" gorm:"type:decimal(10,2)"`
    RemainingBalance float64      `json:"remaining_balance" gorm:"type:decimal(10,2);default:0"`
    Status          PaymentStatus `json:"status" gorm:"type:varchar(20)"`
    PaymentMethod   string        `json:"payment_method" gorm:"type:varchar(50)"`
    IsPartial       bool          `json:"is_partial" gorm:"default:false"`
    CompletesWeek   bool          `json:"completes_week" gorm:"default:false"`
    CreatedAt       time.Time     `json:"created_at"`
    UpdatedAt       time.Time     `json:"updated_at"`
    DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
    
    // Relations
    Loan *Loan `json:"loan,omitempty" gorm:"foreignKey:LoanID"`
}

type PaymentCreateRequest struct {
    LoanID          uint    `json:"loan_id" binding:"required"`
    WeekNumber      int     `json:"week_number" binding:"required"`
    PaymentDate     string  `json:"payment_date,omitempty"`
    AmountDue       float64 `json:"amount_due" binding:"required"`
    AmountPaid      float64 `json:"amount_paid" binding:"required"`
    RemainingBalance float64 `json:"remaining_balance,omitempty"`
    Status          string  `json:"status" binding:"required"`
    PaymentMethod   string  `json:"payment_method" binding:"required"`
    IsPartial       bool    `json:"is_partial,omitempty"`
    CompletesWeek   bool    `json:"completes_week,omitempty"`
}
type LoanUpdateRequest struct {
    ID                 uint     `json:"id"`
    OutstandingBalance *float64 `json:"outstanding_balance,omitempty"`
    Status             string   `json:"status,omitempty"`
    PaymentPeriodWeeks *int     `json:"payment_period_weeks,omitempty"`
    DueDate            string   `json:"due_date,omitempty"`
}


