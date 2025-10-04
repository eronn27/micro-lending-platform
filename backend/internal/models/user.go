package models

import (
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    BaseModel
    Username     string `gorm:"uniqueIndex;not null;size:50" json:"username"`
    PasswordHash string `gorm:"not null" json:"-"`
    IsAdmin      bool   `gorm:"default:false" json:"is_admin"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}

func (User) TableName() string {
    return "users"
}

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPasswordHash compares a plain text password with a hash
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
