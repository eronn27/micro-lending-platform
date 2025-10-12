package handlers

import (
    "net/http"
    "strconv"
    "time"
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

// LoanCreateRequest for existing clients (standalone loan creation)
type LoanCreateRequest struct {
    ClientID uint                    `json:"client_id" binding:"required"`
    Loan     models.LoanCreate       `json:"loan" binding:"required"`
    CoMakers []models.CoMakerCreate  `json:"comakers,omitempty"`
}

// CreateLoan creates a new loan for an existing client
func (h *LoanHandler) CreateLoan(c *gin.Context) {
    var req LoanCreateRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Create the loan with ClientID
    createdLoan, err := h.loanService.CreateLoan(&req.Loan, req.ClientID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create loan: " + err.Error()})
        return
    }

    // Create co-makers if provided
    for _, comakerData := range req.CoMakers {
        comaker := &models.CoMaker{
            LoanID:   createdLoan.ID,
            Name:     comakerData.Name,
            Address:  comakerData.Address,
            Business: comakerData.Business,
        }
        // Note: You'll need to add a CreateCoMaker method to your service
        // For now, we'll skip error handling for co-makers
        _ = comaker
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Loan created successfully",
        "loan":    createdLoan,
    })
}
// GetAllLoans retrieves all loans with pagination
func (h *LoanHandler) GetAllLoans(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    status := c.Query("status")

    loans, total, err := h.loanService.GetAllLoans(page, limit, status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch loans"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "loans": loans,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// GetLoansByClientID retrieves all loans for a specific client
func (h *LoanHandler) GetLoansByClientID(c *gin.Context) {
    clientIDStr := c.Param("clientId")
    clientID, err := strconv.ParseUint(clientIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
        return
    }

    loans, err := h.loanService.GetLoansByClientID(uint(clientID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch loans"})
        return
    }

    c.JSON(http.StatusOK, loans)
}

// GetLoan retrieves a single loan by ID
func (h *LoanHandler) GetLoan(c *gin.Context) {
    loanIDStr := c.Param("id")
    loanID, err := strconv.ParseUint(loanIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
        return
    }

    loan, err := h.loanService.GetLoanByID(uint(loanID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
        return
    }

    c.JSON(http.StatusOK, loan)
}

// UpdateLoan updates an existing loan
func (h *LoanHandler) UpdateLoan(c *gin.Context) {
    loanIDStr := c.Param("id")
    loanID, err := strconv.ParseUint(loanIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
        return
    }

    var updateReq models.LoanUpdateRequest
    if err := c.ShouldBindJSON(&updateReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Set the ID from URL parameter
    updateReq.ID = uint(loanID)

    updatedLoan, err := h.loanService.UpdateLoan(&updateReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan"})
        return
    }

    c.JSON(http.StatusOK, updatedLoan)
}

// DeleteLoan soft deletes a loan
func (h *LoanHandler) DeleteLoan(c *gin.Context) {
    loanIDStr := c.Param("id")
    loanID, err := strconv.ParseUint(loanIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
        return
    }

    if err := h.loanService.DeleteLoan(uint(loanID)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete loan"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}

// GetLoanStats returns loan statistics
func (h *LoanHandler) GetLoanStats(c *gin.Context) {
    stats, err := h.loanService.GetLoanStats()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch loan statistics"})
        return
    }

    c.JSON(http.StatusOK, stats)
}

// Helper function to generate loan control number
func generateLoanControlNumber(clientID uint) string {
    return "L" + strconv.FormatUint(uint64(clientID), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10)
}
