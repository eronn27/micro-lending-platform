package models

type WeeklyReportData struct {
    WeeklyPaymentTotal   float64 `json:"weekly_payment_total"`
    WeeklyReleaseTotal   float64 `json:"weekly_release_total"`
    TotalClients         int64   `json:"total_clients"`
    ActiveClients        int64   `json:"active_clients"`
    OverdueClients       int64   `json:"overdue_clients"`
    ActivePaymentTotal   float64 `json:"active_payment_total"`
    TotalPaymentThisWeek float64 `json:"total_payment_this_week"`
}

type ReportResponse struct {
    Data      WeeklyReportData `json:"data"`
    Timestamp string           `json:"timestamp"`
    Status    string           `json:"status"`
}

