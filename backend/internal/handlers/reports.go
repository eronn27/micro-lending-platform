package handlers

import (
	"net/http"
	"strconv"
	"time"
	"micro-lending-platform/backend/internal/models"
	"micro-lending-platform/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService *services.ReportService
}

func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// GetWeeklyReport returns weekly report data
// @Summary Get Weekly Report
// @Description Returns weekly statistics including payments, releases, and client metrics
// @Tags reports
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.ReportResponse
// @Failure 500 {object} gin.H
// @Router /api/v1/reports/weekly [get]
func (h *ReportHandler) GetWeeklyReport(c *gin.Context) {
	data, err := h.reportService.GetWeeklyReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate report",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ReportResponse{
		Data:      *data,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Status:    "success",
	})
}

// GetMonthlyReport returns monthly report data
// @Summary Get Monthly Report
// @Description Returns monthly statistics including payments, releases, and client metrics
// @Tags reports
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.ReportResponse
// @Failure 500 {object} gin.H
// @Router /api/v1/reports/monthly [get]
func (h *ReportHandler) GetMonthlyReport(c *gin.Context) {
	data, err := h.reportService.GetMonthlyReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate report",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ReportResponse{
		Data:      *data,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Status:    "success",
	})
}

// GetHistoricalReport returns historical report data for multiple periods
// @Summary Get Historical Report
// @Description Returns historical statistics for multiple time periods
// @Tags reports
// @Security BearerAuth
// @Produce json
// @Param period query string false "Period type: weekly or monthly" default(weekly)
// @Param periods query int false "Number of periods to retrieve" default(4)
// @Success 200 {object} models.HistoricalReportResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/reports/history [get]
func (h *ReportHandler) GetHistoricalReport(c *gin.Context) {
	// Get query parameters
	periodType := c.DefaultQuery("period", "weekly")
	periodsStr := c.DefaultQuery("periods", "4")
	
	// Validate period type
	if periodType != "weekly" && periodType != "monthly" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid period type. Must be 'weekly' or 'monthly'",
		})
		return
	}
	
	// Parse and validate periods count
	periods, err := strconv.Atoi(periodsStr)
	if err != nil || periods < 1 || periods > 52 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid periods count. Must be between 1 and 52",
		})
		return
	}
	
	// Generate historical report
	data, err := h.reportService.GetHistoricalReport(periodType, periods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate historical report",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, data)
}
