package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "bytes"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  struct {
        ID       uint   `json:"id"`
        Username string `json:"username"`
        IsAdmin  bool   `json:"is_admin"`
    } `json:"user"`
}

func main() {
    baseURL := "http://localhost:8080/api"
    
    // Default admin credentials
    loginReq := LoginRequest{
        Username: "admin",
        Password: "admin123",
    }
    
    token, err := getJWTToken(baseURL, loginReq)
    if err != nil {
        fmt.Printf("Failed to get JWT token: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("JWT Token: %s\n", token)
    fmt.Println("Copy this token and replace 'YOUR_JWT_TOKEN_HERE' in create_test_clients.go")
}

func getJWTToken(baseURL string, loginReq LoginRequest) (string, error) {
    jsonData, err := json.Marshal(loginReq)
    if err != nil {
        return "", fmt.Errorf("failed to marshal login request: %w", err)
    }
    
    resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to send login request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
    }
    
    var loginResp LoginResponse
    if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
        return "", fmt.Errorf("failed to decode login response: %w", err)
    }
    
    return loginResp.Token, nil
}
