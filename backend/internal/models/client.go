package models

import (
    "time"
)

type Client struct {
    BaseModel
    ControlNumber     string    `gorm:"uniqueIndex;not null;size:20" json:"control_number"`
    FirstName         string    `gorm:"not null;size:100" json:"first_name"`
    MiddleName        string    `gorm:"size:100" json:"middle_name"`
    LastName          string    `gorm:"not null;size:100" json:"last_name"`
    Nickname          string    `gorm:"size:50" json:"nickname"`
    DateOfBirth       time.Time `json:"date_of_birth"`
    Gender            string    `gorm:"size:10" json:"gender"`
    Religion          string    `gorm:"size:50" json:"religion"`
    CivilStatus       string    `gorm:"size:20" json:"civil_status"`
    HomeAddress       string    `gorm:"type:text" json:"home_address"`
    YearsOfResidence  int       `json:"years_of_residence"`
    FacebookAccount   string    `gorm:"size:100" json:"facebook_account"`
    Age               int       `json:"age"`
    ContactNumber     string    `gorm:"size:20" json:"contact_number"`
    
    // Relationships - use slices instead of pointers
    IncomeInfo    []IncomeInfo    `gorm:"foreignKey:ClientID" json:"income_info,omitempty"`
    Loans         []Loan          `gorm:"foreignKey:ClientID" json:"loans,omitempty"`
    Documents     []Document      `gorm:"foreignKey:ClientID" json:"documents,omitempty"`
    FamilyMembers []FamilyMember  `gorm:"foreignKey:ClientID" json:"family_members,omitempty"`
}

// TableName specifies the table name for GORM
func (Client) TableName() string {
    return "clients"
}
