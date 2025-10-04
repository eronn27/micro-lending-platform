package handlers

import (
    "net/http"
    "micro-lending-platform/backend/internal/services"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
    var loginReq struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&loginReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // Use the auth service
    response, err := h.authService.Login(loginReq.Username, loginReq.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
    // Get user info from context (set by auth middleware)
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    username, _ := c.Get("username")
    isAdmin, _ := c.Get("is_admin")

    c.JSON(http.StatusOK, gin.H{
        "id":       userID,
        "username": username,
        "is_admin": isAdmin,
    })
}

// Placeholder methods for new routes
func (h *AuthHandler) Register(c *gin.Context) {
    c.JSON(501, gin.H{"message": "User registration - coming soon"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
    c.JSON(501, gin.H{"message": "Token refresh - coming soon"})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
    c.JSON(501, gin.H{"message": "Update profile - coming soon"})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
    c.JSON(501, gin.H{"message": "Change password - coming soon"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
    c.JSON(501, gin.H{"message": "Logout - coming soon"})
}

