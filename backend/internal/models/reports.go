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

// HistoricalRecord represents data for a single historical period
type HistoricalRecord struct {
	Period         string  `json:"period"`
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
	Payments       float64 `json:"payments"`
	Releases       float64 `json:"releases"`
	ActiveClients  int64   `json:"active_clients"`
	OverdueClients int64   `json:"overdue_clients"`
	NetFlow        float64 `json:"net_flow"`
}

// HistoricalReportResponse represents the response for historical reports
type HistoricalReportResponse struct {
	Data      []HistoricalRecord `json:"data"`
	Timestamp string             `json:"timestamp"`
	Status    string             `json:"status"`
	Metadata  HistoryMetadata    `json:"metadata"`
}

// HistoryMetadata provides summary information about the historical data
type HistoryMetadata struct {
	PeriodType    string  `json:"period_type"`
	PeriodsCount  int     `json:"periods_count"`
	AveragePayment float64 `json:"average_payment"`
	AverageRelease float64 `json:"average_release"`
	TotalPayments  float64 `json:"total_payments"`
	TotalReleases  float64 `json:"total_releases"`
}
