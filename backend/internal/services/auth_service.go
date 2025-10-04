package services

import (
    "fmt"
    "micro-lending-platform/backend/internal/models"
    "micro-lending-platform/backend/internal/repositories"
    "micro-lending-platform/backend/internal/auth"
)

type AuthService struct {
    userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
    return &AuthService{userRepo: userRepo}
}

// Login authenticates a user and returns JWT token
func (s *AuthService) Login(username, password string) (*models.LoginResponse, error) {
    // Validate input
    if username == "" || password == "" {
        return nil, fmt.Errorf("username and password are required")
    }

    // Find user by username
    user, err := s.userRepo.FindByUsername(username)
    if err != nil {
        // Return generic error for security
        return nil, fmt.Errorf("invalid credentials")
    }

    // Check password
    if !auth.CheckPasswordHash(password, user.PasswordHash) {
        return nil, fmt.Errorf("invalid credentials")
    }

    // Generate JWT token
    token, err := auth.GenerateJWT(user.ID, user.Username, user.IsAdmin)
    if err != nil {
        return nil, fmt.Errorf("failed to generate authentication token")
    }

    // Return login response
    response := &models.LoginResponse{
        Token: token,
        User:  *user,
    }

    return response, nil
}
