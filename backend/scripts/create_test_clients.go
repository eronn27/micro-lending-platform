package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    "io"
    "bytes"
)

// Test client data structure matching your API
type TestClient struct {
    Client     ClientData     `json:"client"`
    Income     *IncomeData    `json:"income,omitempty"`
    Loan       *LoanData      `json:"loan,omitempty"`
    CoMakers   []CoMakerData  `json:"comakers,omitempty"`
    Family     FamilyData     `json:"family"`
    Siblings   []SiblingData  `json:"siblings,omitempty"`
    Spouse     *SpouseData    `json:"spouse,omitempty"`
    Dependents []DependentData `json:"dependents,omitempty"`
}

type ClientData struct {
    FirstName        string `json:"first_name"`
    MiddleName       string `json:"middle_name"`
    LastName         string `json:"last_name"`
    Nickname         string `json:"nickname"`
    DateOfBirth      string `json:"date_of_birth"`
    Gender           string `json:"gender"`
    Religion         string `json:"religion"`
    CivilStatus      string `json:"civil_status"`
    HomeAddress      string `json:"home_address"`
    YearsOfResidence int    `json:"years_of_residence"`
    FacebookAccount  string `json:"facebook_account"`
    Age              int    `json:"age"`
    ContactNumber    string `json:"contact_number"`
}

type IncomeData struct {
    FamilyIncomeDaily   float64 `json:"family_income_daily"`
    FamilyIncomeMonthly float64 `json:"family_income_monthly"`
    TotalCostDaily      float64 `json:"total_cost_daily"`
    TotalCostMonthly    float64 `json:"total_cost_monthly"`
    NetIncomeDaily      float64 `json:"net_income_daily"`
    NetIncomeMonthly    float64 `json:"net_income_monthly"`
}

type LoanData struct {
    TotalAmount           float64 `json:"total_amount"`
    Terms                 int     `json:"terms"`
    Ammortization         float64 `json:"ammortization"`
    OutstandingBalance    float64 `json:"outstanding_balance"`
    Mode                  string  `json:"mode"`
    Deductions            string  `json:"deductions"`
    DueDate               string  `json:"due_date"`
    AmountRelease         float64 `json:"amount_release"`
    DateOfRelease         string  `json:"date_of_release"`
    CreditHistory         string  `json:"credit_history"`
    MethodOfPayment       string  `json:"method_of_payment"`
    RecommendedBy         string  `json:"recommended_by"`
    ApprovedBy            string  `json:"approved_by"`
    LoanCycle             int     `json:"loan_cycle"`
    RecommendedLoanAmount float64 `json:"recommended_loan_amount"`
    ApprovedLoanAmount    float64 `json:"approved_loan_amount"`
    CheckedBy             string  `json:"checked_by"`
    NameCI                string  `json:"name_ci"`
    NotedBy               string  `json:"noted_by"`
    ApplicationDate       string  `json:"application_date"`
    Status                string  `json:"status"`
}

type CoMakerData struct {
    Name     string `json:"name"`
    Address  string `json:"address"`
    Business string `json:"business"`
}

type FamilyData struct {
    FatherName string `json:"father_name"`
    MotherName string `json:"mother_name"`
    Address    string `json:"address"`
}

type SiblingData struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Address string `json:"address"`
}

type SpouseData struct {
    Name          string `json:"name"`
    Age           int    `json:"age"`
    Nickname      string `json:"nickname"`
    Birthday      string `json:"birthday"`
    Work          string `json:"work"`
    ContactNumber string `json:"contact_number"`
}

type DependentData struct {
    Name    string `json:"name"`
    Surname string `json:"surname"`
    Age     int    `json:"age"`
    Spouse  string `json:"spouse"`
    Address string `json:"address"`
}

func main() {
    baseURL := "http://localhost:8080/api"
    jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJ1c2VybmFtZSI6InN0YWZmMSIsImlzX2FkbWluIjpmYWxzZSwiZXhwIjoxNzU5NTIwMzA1fQ.regxexcT8wU8k_bh7oS6Dc3_8pkFRbaw6AUUGxvdV8I" // Replace with actual JWT token from login
    
    testClients := generateTestClients()
    
    fmt.Printf("Creating %d test clients...\n", len(testClients))
    
    for i, client := range testClients {
        fmt.Printf("Creating client %d: %s %s\n", i+1, client.Client.FirstName, client.Client.LastName)
        
        err := createClient(baseURL, jwtToken, client)
        if err != nil {
            log.Printf("Failed to create client %d: %v\n", i+1, err)
        } else {
            fmt.Printf("Successfully created client %d\n", i+1)
        }
        
        // Small delay to avoid overwhelming the server
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("Test client creation completed!")
}

func createClient(baseURL, jwtToken string, client TestClient) error {
    // Convert client data to JSON
    jsonData, err := json.Marshal(client)
    if err != nil {
        return fmt.Errorf("failed to marshal client data: %w", err)
    }
    
    // Create HTTP request
    req, err := http.NewRequest("POST", baseURL+"/clients", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    
    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+jwtToken)
    
    // Send request
    clientHTTP := &http.Client{Timeout: 30 * time.Second}
    resp, err := clientHTTP.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    // Check response
    if resp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
    }
    
    return nil
}

func generateTestClients() []TestClient {
    return []TestClient{
        {
            Client: ClientData{
                FirstName:        "Juan",
                MiddleName:       "Dela",
                LastName:         "Cruz",
                Nickname:         "John",
                DateOfBirth:      "1985-05-15",
                Gender:           "Male",
                Religion:         "Catholic",
                CivilStatus:      "Married",
                HomeAddress:      "123 Main Street, Manila, Metro Manila",
                YearsOfResidence: 5,
                FacebookAccount:  "juan.delacruz",
                Age:              38,
                ContactNumber:    "0917-123-4567",
            },
            Income: &IncomeData{
                FamilyIncomeDaily:   2500.00,
                FamilyIncomeMonthly: 75000.00,
                TotalCostDaily:      1500.00,
                TotalCostMonthly:    45000.00,
                NetIncomeDaily:      1000.00,
                NetIncomeMonthly:    30000.00,
            },
            Loan: &LoanData{
                TotalAmount:           10000.00,
                Terms:                 4,
                Ammortization:         625.00,
                OutstandingBalance:    10000.00,
                Mode:                  "Weekly",
                Deductions:            "Without",
                DueDate:               "Monday",
                AmountRelease:         8800.00,
                DateOfRelease:         "2024-01-15",
                CreditHistory:         "New",
                MethodOfPayment:       "Cash",
                RecommendedBy:         "Maria Santos",
                ApprovedBy:            "Manager Reyes",
                LoanCycle:             1,
                RecommendedLoanAmount: 10000.00,
                ApprovedLoanAmount:    10000.00,
                CheckedBy:             "AO Garcia",
                NameCI:                "CI Officer",
                NotedBy:               "Branch Manager",
                ApplicationDate:       "2024-01-10",
                Status:                "Active",
            },
            CoMakers: []CoMakerData{
                {
                    Name:     "Pedro Santos",
                    Address:  "456 Oak Street, Quezon City",
                    Business: "Sari-sari Store",
                },
                {
                    Name:     "Maria Gonzales",
                    Address:  "789 Pine Street, Makati",
                    Business: "Food Stall",
                },
            },
            Family: FamilyData{
                FatherName: "Carlos Cruz",
                MotherName: "Elena Cruz",
                Address:    "123 Main Street, Manila, Metro Manila",
            },
            Siblings: []SiblingData{
                {
                    Name:    "Ana Cruz",
                    Age:     35,
                    Address: "321 Elm Street, Pasig",
                },
                {
                    Name:    "Luis Cruz",
                    Age:     32,
                    Address: "654 Maple Street, Mandaluyong",
                },
            },
            Spouse: &SpouseData{
                Name:          "Maria Cruz",
                Age:           36,
                Nickname:      "Maya",
                Birthday:      "1987-08-20",
                Work:          "Teacher",
                ContactNumber: "0918-765-4321",
            },
            Dependents: []DependentData{
                {
                    Name:    "Sophia",
                    Surname: "Cruz",
                    Age:     12,
                    Spouse:  "",
                    Address: "123 Main Street, Manila, Metro Manila",
                },
                {
                    Name:    "Lucas",
                    Surname: "Cruz",
                    Age:     8,
                    Spouse:  "",
                    Address: "123 Main Street, Manila, Metro Manila",
                },
            },
        },
        {
            Client: ClientData{
                FirstName:        "Maria",
                MiddleName:       "Santos",
                LastName:         "Reyes",
                Nickname:         "Mary",
                DateOfBirth:      "1990-12-10",
                Gender:           "Female",
                Religion:         "Protestant",
                CivilStatus:      "Single",
                HomeAddress:      "456 Oak Avenue, Quezon City",
                YearsOfResidence: 3,
                FacebookAccount:  "maria.reyes",
                Age:              33,
                ContactNumber:    "0919-888-9999",
            },
            Income: &IncomeData{
                FamilyIncomeDaily:   1800.00,
                FamilyIncomeMonthly: 54000.00,
                TotalCostDaily:      1200.00,
                TotalCostMonthly:    36000.00,
                NetIncomeDaily:      600.00,
                NetIncomeMonthly:    18000.00,
            },
            Loan: &LoanData{
                TotalAmount:           8000.00,
                Terms:                 3,
                Ammortization:         667.00,
                OutstandingBalance:    8000.00,
                Mode:                  "Weekly",
                Deductions:            "With",
                DueDate:               "Friday",
                AmountRelease:         6800.00,
                DateOfRelease:         "2024-01-20",
                CreditHistory:         "Existing",
                MethodOfPayment:       "GCash",
                RecommendedBy:         "Juan Dela Cruz",
                ApprovedBy:            "Supervisor Lim",
                LoanCycle:             2,
                RecommendedLoanAmount: 8000.00,
                ApprovedLoanAmount:    8000.00,
                CheckedBy:             "AO Tan",
                NameCI:                "CI Officer",
                NotedBy:               "Area Manager",
                ApplicationDate:       "2024-01-18",
                Status:                "Active",
            },
            CoMakers: []CoMakerData{
                {
                    Name:     "Roberto Lim",
                    Address:  "123 Pine Street, Makati",
                    Business: "Computer Shop",
                },
            },
            Family: FamilyData{
                FatherName: "Antonio Reyes",
                MotherName: "Carmen Reyes",
                Address:    "456 Oak Avenue, Quezon City",
            },
            Siblings: []SiblingData{
                {
                    Name:    "Carlos Reyes",
                    Age:     28,
                    Address: "789 Cedar Street, Taguig",
                },
            },
            // No spouse for single client
            Dependents: []DependentData{
                {
                    Name:    "Andrea",
                    Surname: "Reyes",
                    Age:     5,
                    Spouse:  "",
                    Address: "456 Oak Avenue, Quezon City",
                },
            },
        },
        {
            Client: ClientData{
                FirstName:        "Pedro",
                MiddleName:       "Gonzales",
                LastName:         "Lim",
                Nickname:         "Pete",
                DateOfBirth:      "1978-03-25",
                Gender:           "Male",
                Religion:         "Catholic",
                CivilStatus:      "Married",
                HomeAddress:      "789 Pine Road, Makati City",
                YearsOfResidence: 8,
                FacebookAccount:  "pedro.lim",
                Age:              45,
                ContactNumber:    "0920-111-2222",
            },
            Income: &IncomeData{
                FamilyIncomeDaily:   3500.00,
                FamilyIncomeMonthly: 105000.00,
                TotalCostDaily:      2000.00,
                TotalCostMonthly:    60000.00,
                NetIncomeDaily:      1500.00,
                NetIncomeMonthly:    45000.00,
            },
            Loan: &LoanData{
                TotalAmount:           15000.00,
                Terms:                 6,
                Ammortization:         625.00,
                OutstandingBalance:    15000.00,
                Mode:                  "Weekly",
                Deductions:            "Without",
                DueDate:               "Wednesday",
                AmountRelease:         13800.00,
                DateOfRelease:         "2024-01-25",
                CreditHistory:         "New W/ Existing",
                MethodOfPayment:       "Bank Transfer",
                RecommendedBy:         "Ana Torres",
                ApprovedBy:            "Manager Chen",
                LoanCycle:             3,
                RecommendedLoanAmount: 15000.00,
                ApprovedLoanAmount:    15000.00,
                CheckedBy:             "AO Wong",
                NameCI:                "CI Officer",
                NotedBy:               "Regional Head",
                ApplicationDate:       "2024-01-22",
                Status:                "Active",
            },
            CoMakers: []CoMakerData{
                {
                    Name:     "Miguel Tan",
                    Address:  "321 Elm Street, Pasig",
                    Business: "Auto Repair",
                },
                {
                    Name:     "Elena Sy",
                    Address:  "654 Maple Street, Mandaluyong",
                    Business: "Grocery Store",
                },
            },
            Family: FamilyData{
                FatherName: "Jose Lim",
                MotherName: "Teresita Lim",
                Address:    "789 Pine Road, Makati City",
            },
            Siblings: []SiblingData{
                {
                    Name:    "Roberto Lim",
                    Age:     42,
                    Address: "123 Cedar Lane, Paranaque",
                },
                {
                    Name:    "Susan Lim",
                    Age:     40,
                    Address: "456 Birch Street, Las Pinas",
                },
                {
                    Name:    "Michael Lim",
                    Age:     38,
                    Address: "789 Walnut Road, Muntinlupa",
                },
            },
            Spouse: &SpouseData{
                Name:          "Catherine Lim",
                Age:           43,
                Nickname:      "Cathy",
                Birthday:      "1980-07-15",
                Work:          "Accountant",
                ContactNumber: "0921-333-4444",
            },
            Dependents: []DependentData{
                {
                    Name:    "Daniel",
                    Surname: "Lim",
                    Age:     16,
                    Spouse:  "",
                    Address: "789 Pine Road, Makati City",
                },
                {
                    Name:    "Sarah",
                    Surname: "Lim",
                    Age:     14,
                    Spouse:  "",
                    Address: "789 Pine Road, Makati City",
                },
                {
                    Name:    "David",
                    Surname: "Lim",
                    Age:     10,
                    Spouse:  "",
                    Address: "789 Pine Road, Makati City",
                },
            },
        },
        // Add more test clients as needed
    }
}


