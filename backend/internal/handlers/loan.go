package handlers

import (
    "net/http"
    "strconv"
    "micro-lending-platform/backend/internal/models"
    "micro-lending-platform/backend/internal/services"
    "github.com/gin-gonic/gin"
)

type LoanHandler struct {
    loanService *services.LoanService
}

func NewLoanHandler(loanService *services.LoanService) *LoanHandler {
    return &LoanHandler{loanService: loanService}
}

// CreateLoan handles loan creation
func (h *LoanHandler) CreateLoan(c *gin.Context) {
    var req models.LoanCreate
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request data", 
            "details": err.Error(),
        })
        return
    }

    createdLoan, err := h.loanService.CreateLoan(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create loan", 
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Loan created successfully",
        "loan": createdLoan,
    })
}

// GetAllLoans retrieves all loans with pagination
func (h *LoanHandler) GetAllLoans(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    status := c.Query("status")

    loans, total, err := h.loanService.GetAllLoans(page, limit, status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve loans"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "loans": loans,
        "pagination": gin.H{
            "page":  page,
            "limit": limit,
            "total": total,
            "pages": (total + int64(limit) - 1) / int64(limit),
        },
    })
}

// GetLoanByID retrieves a loan by ID
func (h *LoanHandler) GetLoanByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
        return
    }

    loan, err := h.loanService.GetLoanByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
        return
    }

    c.JSON(http.StatusOK, loan)
}

// GetLoansByClientID retrieves all loans for a specific client
func (h *LoanHandler) GetLoansByClientID(c *gin.Context) {
    idStr := c.Param("clientId")
    clientId, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
        return
    }

    loans, err := h.loanService.GetLoansByClientID(uint(clientId))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve loans"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "loans": loans,
        "total": len(loans),
    })
}

// UpdateLoan updates loan information
func (h *LoanHandler) UpdateLoan(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
        return
    }

    var updateReq models.LoanUpdateRequest
    if err := c.ShouldBindJSON(&updateReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
        return
    }

    updateReq.ID = uint(id)
    updatedLoan, err := h.loanService.UpdateLoan(&updateReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Loan updated successfully",
        "loan": updatedLoan,
    })
}

// DeleteLoan soft deletes a loan
func (h *LoanHandler) DeleteLoan(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
        return
    }

    err = h.loanService.DeleteLoan(uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}

// GetLoanStats returns loan statistics
func (h *LoanHandler) GetLoanStats(c *gin.Context) {
    stats, err := h.loanService.GetLoanStats()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get loan statistics"})
        return
    }

    c.JSON(http.StatusOK, stats)
}
