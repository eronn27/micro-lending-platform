package models


type IncomeInfo struct {
    BaseModel
    ClientID             uint    `gorm:"not null;uniqueIndex" json:"client_id"`
    FamilyIncomeDaily    float64 `gorm:"type:decimal(10,2)" json:"family_income_daily"`
    FamilyIncomeMonthly  float64 `gorm:"type:decimal(10,2)" json:"family_income_monthly"`
    TotalCostDaily       float64 `gorm:"type:decimal(10,2)" json:"total_cost_daily"`
    TotalCostMonthly     float64 `gorm:"type:decimal(10,2)" json:"total_cost_monthly"`
    NetIncomeDaily       float64 `gorm:"type:decimal(10,2)" json:"net_income_daily"`
    NetIncomeMonthly     float64 `gorm:"type:decimal(10,2)" json:"net_income_monthly"`
    
    Client Client `gorm:"foreignKey:ClientID" json:"client,omitempty"`
}

func (IncomeInfo) TableName() string {
    return "income_info"
}

type CoMaker struct {
    BaseModel
    LoanID        uint   `gorm:"not null;index" json:"loan_id"`
    Name          string `gorm:"not null;size:100" json:"name"`
    Address       string `gorm:"type:text" json:"address"`
    Business      string `gorm:"size:100" json:"business"`
    SignaturePath string `gorm:"size:255" json:"signature_path"`
    IDCopyPath    string `gorm:"size:255" json:"id_copy_path"`
    
    Loan Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
}

func (CoMaker) TableName() string {
    return "co_makers"
}

// Update the FamilyMember struct in models/related.go
type FamilyMember struct {
    BaseModel
    ClientID     uint      `gorm:"not null;index" json:"client_id"`
    Relationship string    `gorm:"not null;size:50" json:"relationship"` // father, mother, sibling, spouse, dependent
    Name         string    `gorm:"not null;size:100" json:"name"`
    Surname      string    `gorm:"size:100" json:"surname"`
    Age          int       `json:"age"`
    Work         string    `gorm:"size:100" json:"work"`
    ContactNumber string   `gorm:"size:20" json:"contact_number"`
    Address      string    `gorm:"type:text" json:"address"`
    Birthday     string `json:"birthday"` // ADD THIS FIELD
    Nickname     string    `gorm:"size:50" json:"nickname"` // ADD THIS FIELD for spouse nickname
    
    Client Client `gorm:"foreignKey:ClientID" json:"client,omitempty"`
}

func (FamilyMember) TableName() string {
    return "family_members"
}
type Document struct {
    BaseModel
    ClientID     uint   `gorm:"not null;index" json:"client_id"`
    DocumentType string `gorm:"not null;size:50" json:"document_type"` // valid_id, barangay_clearance, etc.
    FilePath     string `gorm:"not null;size:255" json:"file_path"`
    FileName     string `gorm:"not null;size:255" json:"file_name"`
    
    Client Client `gorm:"foreignKey:ClientID" json:"client,omitempty"`
}

func (Document) TableName() string {
    return "documents"
}


// Add these to the bottom of models/related.go

// ClientWithRelatedData represents the complete client data structure from the Vue form
type ClientWithRelatedData struct {
    Client     Client         `json:"client"`
    Income     *IncomeInfo    `json:"income,omitempty"`
    Loan       *Loan          `json:"loan,omitempty"`
    CoMakers   []CoMaker      `json:"comakers,omitempty"`
    Family     Family         `json:"family"`
    Siblings   []FamilyMember `json:"siblings,omitempty"`
    Spouse     *FamilyMember  `json:"spouse,omitempty"`
    Dependents []FamilyMember `json:"dependents,omitempty"`
}

// Family represents the family information structure
type Family struct {
    FatherName string `json:"father_name"`
    MotherName string `json:"mother_name"`
    Address    string `json:"address"`
}

// ClientCreateRequest represents the complete client creation request from frontend
type ClientCreateRequest struct {
    Client     ClientCreate     `json:"client"`
    Income     *IncomeInfo      `json:"income,omitempty"`
    Loan       *LoanCreate      `json:"loan,omitempty"`
    CoMakers   []CoMakerCreate  `json:"comakers,omitempty"`
    Family     Family           `json:"family"`
    Siblings   []FamilyMember   `json:"siblings,omitempty"`
    Spouse     *FamilyMember    `json:"spouse,omitempty"`
    Dependents []FamilyMember   `json:"dependents,omitempty"`
}

// ClientCreate represents client data for creation (without ID and relationships)
type ClientCreate struct {
    ControlNumber     string    `json:"control_number"`
    FirstName         string    `json:"first_name" binding:"required"`
    MiddleName        string    `json:"middle_name"`
    LastName          string    `json:"last_name" binding:"required"`
    Nickname          string    `json:"nickname"`
    DateOfBirth       string    `json:"date_of_birth"`
    Gender            string    `json:"gender"`
    Religion          string    `json:"religion"`
    CivilStatus       string    `json:"civil_status"`
    HomeAddress       string    `json:"home_address" binding:"required"`
    YearsOfResidence  int       `json:"years_of_residence"`
    FacebookAccount   string    `json:"facebook_account"`
    Age               int       `json:"age" binding:"required"`
    ContactNumber     string    `json:"contact_number" binding:"required"`
}

// LoanCreate represents loan data for creation
type LoanCreate struct {
    ControlNumber         string    `json:"control_number"`
    DateOfRelease         string    `json:"date_of_release"`
    TotalAmount           float64   `json:"total_amount" binding:"required"`
    Ammortization         float64   `json:"ammortization" binding:"required"`
    Terms                 int       `json:"terms" binding:"required"`
    Mode                  string    `json:"mode"`
    OutstandingBalance    float64   `json:"outstanding_balance" binding:"required"`
    Status                string    `json:"status"`
    DueDate               string    `json:"due_date"`
    Deductions            string    `json:"deductions"`
    AmountRelease         float64   `json:"amount_release" binding:"required"`
    PaymentPeriodWeeks    int       `json:"payment_period_weeks"`
    MethodOfPayment       string    `json:"method_of_payment"`
    CreditHistory         string    `json:"credit_history"`
    RecommendedBy         string    `json:"recommended_by"`
    ApprovedBy            string    `json:"approved_by"`
    LoanCycle             int       `json:"loan_cycle"`
    RecommendedLoanAmount float64   `json:"recommended_loan_amount"`
    ApprovedLoanAmount    float64   `json:"approved_loan_amount"`
    CheckedBy             string    `json:"checked_by"`
    NameCI                string    `json:"name_ci"` // ADDED THIS MISSING FIELD
    NotedBy               string    `json:"noted_by"`
    ApplicationDate       string    `json:"application_date"`
}

// CoMakerCreate represents co-maker data for creation
type CoMakerCreate struct {
    Name     string `json:"name"`
    Address  string `json:"address"`
    Business string `json:"business"`
}
