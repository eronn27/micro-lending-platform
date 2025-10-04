package models

import (
    "time"
    "gorm.io/gorm"
)

// Base model with common fields
type BaseModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// JSONTime for custom time formatting
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
    stamp := time.Time(t).Format("2006-01-02 15:04:05")
    return []byte(`"` + stamp + `"`), nil
}
