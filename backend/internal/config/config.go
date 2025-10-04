package config

import (
    "os"
    "strconv"
)

type Config struct {
    DBPath       string
    ServerPort   string
    JWTSecret    string
    Environment  string
}

func Load() *Config {
    return &Config{
        DBPath:      getEnv("DB_PATH", "./data/micro_lending.db"),
        ServerPort:  getEnv("SERVER_PORT", "8080"),
        JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
        Environment: getEnv("ENVIRONMENT", "development"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}
