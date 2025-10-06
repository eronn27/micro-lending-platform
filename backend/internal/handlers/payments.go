package handlers

import (
    "net/http"
    "strconv"
    "micro-lending-platform/backend/internal/models"
    "micro-lending-platform/backend/internal/services"
    "github.com/gin-gonic/gin"
)

type PaymentHandler struct {
    paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
    return &PaymentHandler{paymentService: paymentService}
}

// CreatePayment handles payment creation
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
    var req models.PaymentCreateRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request data", 
            "details": err.Error(),
        })
        return
    }

    // Validate required fields
    if req.LoanID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Loan ID is required"})
        return
    }
    if req.WeekNumber == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Week number is required"})
        return
    }
    if req.AmountDue == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Amount due is required"})
        return
    }

    createdPayment, err := h.paymentService.CreatePayment(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create payment", 
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Payment created successfully",
        "payment": createdPayment,
    })
}

// GetPaymentsByLoanID retrieves all payments for a specific loan
func (h *PaymentHandler) GetPaymentsByLoanID(c *gin.Context) {
    idStr := c.Param("loanId")
    loanId, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
        return
    }

    payments, err := h.paymentService.GetPaymentsByLoanID(uint(loanId))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payments"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "payments": payments,
        "total": len(payments),
    })
}

// GetPaymentByID retrieves a single payment by ID
func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
        return
    }

    payment, err := h.paymentService.GetPaymentByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
        return
    }

    c.JSON(http.StatusOK, payment)
}

// UpdatePayment updates payment information
func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
        return
    }

    var payment models.Payment
    if err := c.ShouldBindJSON(&payment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    payment.ID = uint(id)
    updatedPayment, err := h.paymentService.UpdatePayment(&payment)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Payment updated successfully",
        "payment": updatedPayment,
    })
}

// DeletePayment soft deletes a payment
func (h *PaymentHandler) DeletePayment(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
        return
    }

    err = h.paymentService.DeletePayment(uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}

// GetAllPayments retrieves all payments with pagination
func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

    payments, total, err := h.paymentService.GetAllPayments(page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payments"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "payments": payments,
        "pagination": gin.H{
            "page":  page,
            "limit": limit,
            "total": total,
            "pages": (total + int64(limit) - 1) / int64(limit),
        },
    })
}
