package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)

var JWTSecret = []byte("your-secret-key-change-in-production") // Use env var in production

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    IsAdmin  bool   `json:"is_admin"`
    jwt.RegisteredClaims
}

func GenerateJWT(userID uint, username string, isAdmin bool) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    
    claims := &Claims{
        UserID:   userID,
        Username: username,
        IsAdmin:  isAdmin,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(JWTSecret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return JWTSecret, nil
    })
    
    if err != nil || !token.Valid {
        return nil, err
    }
    
    return claims, nil
}
