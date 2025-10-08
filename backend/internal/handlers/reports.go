package handlers

import (
	"net/http"
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
