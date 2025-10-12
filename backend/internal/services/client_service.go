package services

import (
    "fmt"
    "micro-lending-platform/backend/internal/models"
    "micro-lending-platform/backend/internal/repositories"
    "time"
)

type ClientService struct {
    clientRepo *repositories.ClientRepository
}

func NewClientService(clientRepo *repositories.ClientRepository) *ClientService {
    return &ClientService{clientRepo: clientRepo}
}

type DuplicateCheckResult struct {
    IsDuplicate   bool   `json:"is_duplicate"`
    SimilarCount  int64  `json:"similar_count"`
    SearchedName  string `json:"searched_name"`
}

type ClientStats struct {
    TotalClients    int64 `json:"total_clients"`
    ActiveClients   int64 `json:"active_clients"`
    NewThisMonth    int64 `json:"new_this_month"`
    WithActiveLoans int64 `json:"with_active_loans"`
}

// CheckDuplicate checks for potential duplicate clients
func (s *ClientService) CheckDuplicate(name string) (*DuplicateCheckResult, error) {
    count, err := s.clientRepo.CountSimilarClients(name)
    if err != nil {
        return nil, fmt.Errorf("failed to check duplicates: %w", err)
    }

    return &DuplicateCheckResult{
        IsDuplicate:  count > 0,
        SimilarCount: count,
        SearchedName: name,
    }, nil
}

// Update the CreateClientWithRelatedData method in client_service.go
func (s *ClientService) CreateClientWithRelatedData(req *models.ClientCreateRequest) (*models.ClientWithRelatedData, error) {
    // Validate the main client data
    if err := s.validateClientBusinessRules(&req.Client); err != nil {
        return nil, err
    }

    // Convert request to the complete data structure
    clientData, err := s.convertRequestToClientData(req)
    if err != nil {
        return nil, fmt.Errorf("failed to convert request data: %w", err)
    }

    // Generate control numbers if not provided
    if clientData.Client.ControlNumber == "" {
        clientData.Client.ControlNumber = s.generateControlNumber()
    }

    if clientData.Loan != nil && clientData.Loan.ControlNumber == "" {
        clientData.Loan.ControlNumber = s.generateLoanControlNumber()
    }

    // Check if client control number already exists
    if clientData.Client.ControlNumber != "" {
        existing, err := s.clientRepo.FindByControlNumber(clientData.Client.ControlNumber)
        if err == nil && existing != nil {
            return nil, fmt.Errorf("client with control number %s already exists", clientData.Client.ControlNumber)
        }
    }

    // Use the repository to create everything in a transaction
    createdClientData, err := s.clientRepo.CreateClientWithRelatedData(clientData)
    if err != nil {
        return nil, fmt.Errorf("failed to create client with related data: %w", err)
    }

    return createdClientData, nil
}

// Fix the CreateClient method in client_service.go
func (s *ClientService) CreateClient(req *models.ClientCreate) (*models.Client, error) {
    // Parse date of birth
    dob, err := s.parseDate(req.DateOfBirth)
    if err != nil {
        return nil, fmt.Errorf("invalid date of birth: %w", err)
    }

    // Convert to full Client model
    client := &models.Client{
        ControlNumber:    req.ControlNumber,
        FirstName:        req.FirstName,
        MiddleName:       req.MiddleName,
        LastName:         req.LastName,
        Nickname:         req.Nickname,
        DateOfBirth:      dob, // Use parsed time.Time
        Gender:           req.Gender,
        Religion:         req.Religion,
        CivilStatus:      req.CivilStatus,
        HomeAddress:      req.HomeAddress,
        YearsOfResidence: req.YearsOfResidence,
        FacebookAccount:  req.FacebookAccount,
        Age:              req.Age,
        ContactNumber:    req.ContactNumber,
    }

    // Generate control number if not provided
    if client.ControlNumber == "" {
        client.ControlNumber = s.generateControlNumber()
    }

    // Validate business rules
    if err := s.validateClientBusinessRules(req); err != nil {
        return nil, err
    }

    // Check if control number already exists
    if client.ControlNumber != "" {
        existing, err := s.clientRepo.FindByControlNumber(client.ControlNumber)
        if err == nil && existing != nil {
            return nil, fmt.Errorf("client with control number %s already exists", client.ControlNumber)
        }
    }

    createdClient, err := s.clientRepo.CreateSimple(client)
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }

    return createdClient, nil
}

// GetClientByID retrieves a client by ID
func (s *ClientService) GetClientByID(id uint) (*models.Client, error) {
    client, err := s.clientRepo.FindByID(id)
    if err != nil {
        if err.Error() == "record not found" {
            return nil, fmt.Errorf("client not found")
        }
        return nil, fmt.Errorf("failed to get client: %w", err)
    }
    return client, nil
}

// GetClientWithDetails retrieves a client with all related data
func (s *ClientService) GetClientWithDetails(id uint) (*models.ClientWithRelatedData, error) {
    clientData, err := s.clientRepo.FindWithDetails(id)
    if err != nil {
        if err.Error() == "record not found" {
            return nil, fmt.Errorf("client not found")
        }
        return nil, fmt.Errorf("failed to get client details: %w", err)
    }
    return clientData, nil
}

// GetAllClients retrieves all clients with pagination and search
func (s *ClientService) GetAllClients(page, limit int, search string) ([]models.Client, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 20
    }
    
    offset := (page - 1) * limit

    clients, err := s.clientRepo.FindAll(offset, limit, search)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get clients: %w", err)
    }

    // Debug: Check if loans are loaded
    for _, client := range clients {
        fmt.Printf("Client %d (%s %s) - Loans count: %d\n", 
            client.ID, client.FirstName, client.LastName, len(client.Loans))
        if len(client.Loans) > 0 {
            fmt.Printf("  Loan: %+v\n", client.Loans[0])
        }
    }

    total, err := s.clientRepo.Count(search)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count clients: %w", err)
    }

    return clients, total, nil
}

// UpdateClient updates client information
func (s *ClientService) UpdateClient(client *models.Client) (*models.Client, error) {
    // Check if client exists
    _, err := s.clientRepo.FindByID(client.ID)
    if err != nil {
        return nil, fmt.Errorf("client not found")
    }

    updatedClient, err := s.clientRepo.Update(client)
    if err != nil {
        return nil, fmt.Errorf("failed to update client: %w", err)
    }

    return updatedClient, nil
}

// DeleteClient soft deletes a client
func (s *ClientService) DeleteClient(id uint) error {
    // Check if client exists
    _, err := s.clientRepo.FindByID(id)
    if err != nil {
        if err.Error() == "record not found" {
            return fmt.Errorf("client not found")
        }
        return fmt.Errorf("failed to find client: %w", err)
    }

    // Check if client has active loans
    hasActiveLoans, err := s.clientRepo.HasActiveLoans(id)
    if err != nil {
        return fmt.Errorf("failed to check active loans: %w", err)
    }
    if hasActiveLoans {
        return fmt.Errorf("cannot delete client with active loans")
    }

    err = s.clientRepo.Delete(id)
    if err != nil {
        return fmt.Errorf("failed to delete client: %w", err)
    }

    return nil
}

// GetClientStats returns client statistics
func (s *ClientService) GetClientStats() (*ClientStats, error) {
    total, err := s.clientRepo.Count("")
    if err != nil {
        return nil, fmt.Errorf("failed to get total clients: %w", err)
    }

    active, err := s.clientRepo.CountActive()
    if err != nil {
        return nil, fmt.Errorf("failed to get active clients: %w", err)
    }

    newThisMonth, err := s.clientRepo.CountNewThisMonth()
    if err != nil {
        return nil, fmt.Errorf("failed to get new clients this month: %w", err)
    }

    withActiveLoans, err := s.clientRepo.CountWithActiveLoans()
    if err != nil {
        return nil, fmt.Errorf("failed to get clients with active loans: %w", err)
    }

    return &ClientStats{
        TotalClients:    total,
        ActiveClients:   active,
        NewThisMonth:    newThisMonth,
        WithActiveLoans: withActiveLoans,
    }, nil
}

// SearchClients searches clients by various criteria
func (s *ClientService) SearchClients(query string, page, limit int) ([]models.Client, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 20
    }
    
    offset := (page - 1) * limit

    clients, err := s.clientRepo.Search(query, offset, limit)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to search clients: %w", err)
    }

    total, err := s.clientRepo.Count(query)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count search results: %w", err)
    }

    return clients, total, nil
}

// GetClientByControlNumber finds a client by control number
func (s *ClientService) GetClientByControlNumber(controlNumber string) (*models.Client, error) {
    client, err := s.clientRepo.FindByControlNumber(controlNumber)
    if err != nil {
        if err.Error() == "record not found" {
            return nil, fmt.Errorf("client with control number %s not found", controlNumber)
        }
        return nil, fmt.Errorf("failed to get client: %w", err)
    }
    return client, nil
}

// BulkCreateClients creates multiple clients at once
func (s *ClientService) BulkCreateClients(clients []*models.Client) ([]*models.Client, error) {
    var createdClients []*models.Client

    for i, client := range clients {
        // Validate each client
        if err := s.validateClientBusinessRulesSimple(client); err != nil {
            return nil, fmt.Errorf("client %d validation failed: %w", i+1, err)
        }

        // Generate control number if not provided
        if client.ControlNumber == "" {
            client.ControlNumber = s.generateControlNumber()
        }

        createdClient, err := s.clientRepo.CreateSimple(client)
        if err != nil {
            return nil, fmt.Errorf("failed to create client %d: %w", i+1, err)
        }

        createdClients = append(createdClients, createdClient)
    }

    return createdClients, nil
}

// ExportClients returns all client data for export
func (s *ClientService) ExportClients() ([]models.Client, error) {
    clients, err := s.clientRepo.FindAll(0, 0, "") // 0,0 means no pagination
    if err != nil {
        return nil, fmt.Errorf("failed to export clients: %w", err)
    }
    return clients, nil
}

// validateClientBusinessRules validates business-specific rules for ClientCreate
func (s *ClientService) validateClientBusinessRules(client *models.ClientCreate) error {
    // Minimum age requirement
    if client.Age < 18 {
        return fmt.Errorf("client must be at least 18 years old")
    }
    
    // Valid contact number format
    if len(client.ContactNumber) < 10 {
        return fmt.Errorf("contact number must be at least 10 digits")
    }

    // Required fields
    if client.FirstName == "" {
        return fmt.Errorf("first name is required")
    }
    if client.LastName == "" {
        return fmt.Errorf("last name is required")
    }
    if client.HomeAddress == "" {
        return fmt.Errorf("home address is required")
    }

    return nil
}

// validateClientBusinessRulesSimple validates business rules for Client model
func (s *ClientService) validateClientBusinessRulesSimple(client *models.Client) error {
    // Minimum age requirement
    if client.Age < 18 {
        return fmt.Errorf("client must be at least 18 years old")
    }
    
    // Valid contact number format
    if len(client.ContactNumber) < 10 {
        return fmt.Errorf("contact number must be at least 10 digits")
    }

    // Required fields
    if client.FirstName == "" {
        return fmt.Errorf("first name is required")
    }
    if client.LastName == "" {
        return fmt.Errorf("last name is required")
    }
    if client.HomeAddress == "" {
        return fmt.Errorf("home address is required")
    }

    return nil
}

// generateControlNumber generates a unique client control number
func (s *ClientService) generateControlNumber() string {
    count, err := s.clientRepo.Count("")
    if err != nil {
        count = 0
    }
    // Format: MLP-YYYY-XXX (e.g., MLP-2024-001)
    year := time.Now().Year()
    return fmt.Sprintf("MLP-%d-%03d", year, count+1)
}

// generateLoanControlNumber generates a unique loan control number
func (s *ClientService) generateLoanControlNumber() string {
    // You might want to implement a separate counter for loans
    // For now, using timestamp-based approach
    timestamp := time.Now().Unix()
    return fmt.Sprintf("LOAN-%d", timestamp)
}

// Update the parseDate function in client_service.go
func (s *ClientService) parseDate(dateStr string) (time.Time, error) {
    if dateStr == "" {
        return time.Time{}, nil // Return zero time for empty dates
    }
    // Directly parse the Vue date format
    return time.Parse("2006-01-02", dateStr)    
}

// Update the convertRequestToClientData function in client_service.go
func (s *ClientService) convertRequestToClientData(req *models.ClientCreateRequest) (*models.ClientWithRelatedData, error) {
    // Parse dates for client
    dob, err := s.parseDate(req.Client.DateOfBirth)
    if err != nil {
        return nil, fmt.Errorf("invalid date of birth: %w", err)
    }

    // Convert main client
    client := models.Client{
        ControlNumber:    req.Client.ControlNumber,
        FirstName:        req.Client.FirstName,
        MiddleName:       req.Client.MiddleName,
        LastName:         req.Client.LastName,
        Nickname:         req.Client.Nickname,
        DateOfBirth:      dob,
        Gender:           req.Client.Gender,
        Religion:         req.Client.Religion,
        CivilStatus:      req.Client.CivilStatus,
        HomeAddress:      req.Client.HomeAddress,
        YearsOfResidence: req.Client.YearsOfResidence,
        FacebookAccount:  req.Client.FacebookAccount,
        Age:              req.Client.Age,
        ContactNumber:    req.Client.ContactNumber,
    }

    // Convert income info
    var income *models.IncomeInfo
    if req.Income != nil {
        income = &models.IncomeInfo{
            FamilyIncomeDaily:   req.Income.FamilyIncomeDaily,
            FamilyIncomeMonthly: req.Income.FamilyIncomeMonthly,
            TotalCostDaily:      req.Income.TotalCostDaily,
            TotalCostMonthly:    req.Income.TotalCostMonthly,
            NetIncomeDaily:      req.Income.NetIncomeDaily,
            NetIncomeMonthly:    req.Income.NetIncomeMonthly,
        }
    }

    // Parse dates for loan
    var loan *models.Loan
    if req.Loan != nil {
        dateOfRelease, err := s.parseDate(req.Loan.DateOfRelease)
        if err != nil {
            return nil, fmt.Errorf("invalid date of release: %w", err)
        }

        applicationDate, err := s.parseDate(req.Loan.ApplicationDate)
        if err != nil {
            return nil, fmt.Errorf("invalid application date: %w", err)
        }

        loan = &models.Loan{
            ControlNumber:         req.Loan.ControlNumber,
            DateOfRelease:         dateOfRelease,
            TotalAmount:           req.Loan.TotalAmount,
            Ammortization:         req.Loan.Ammortization,
            Terms:                 req.Loan.Terms,
            Mode:                  req.Loan.Mode,
            OutstandingBalance:    req.Loan.OutstandingBalance,
            Status:                models.LoanStatus(req.Loan.Status),
            DueDate:               req.Loan.DueDate,
            Deductions:            req.Loan.Deductions,
            AmountRelease:         req.Loan.AmountRelease,
            PaymentPeriodWeeks:    req.Loan.PaymentPeriodWeeks,
            MethodOfPayment:       req.Loan.MethodOfPayment,
            CreditHistory:         req.Loan.CreditHistory,
            RecommendedBy:         req.Loan.RecommendedBy,
            ApprovedBy:            req.Loan.ApprovedBy,
            LoanCycle:             req.Loan.LoanCycle,
            RecommendedLoanAmount: req.Loan.RecommendedLoanAmount,
            ApprovedLoanAmount:    req.Loan.ApprovedLoanAmount,
            CheckedBy:             req.Loan.CheckedBy,
            NameCI:                req.Loan.NameCI,
            NotedBy:               req.Loan.NotedBy,
            ApplicationDate:       applicationDate,
        }
    }

    // Handle spouse information if provided
    var spouse *models.FamilyMember
    if req.Spouse != nil && req.Spouse.Name != "" {
        spouse = &models.FamilyMember{
            Relationship:  "spouse",
            Name:          req.Spouse.Name,
            Age:           req.Spouse.Age,
            Nickname:      req.Spouse.Nickname,
            Work:          req.Spouse.Work,
            ContactNumber: req.Spouse.ContactNumber,
            Address:       req.Spouse.Address,
            Birthday:      req.Spouse.Birthday, // Already a time.Time, no parsing needed
        }
    }

    // Convert co-makers, siblings, and dependents
    var coMakers []models.CoMaker
    if req.CoMakers != nil {
        for _, cm := range req.CoMakers {
            coMakers = append(coMakers, models.CoMaker{
                Name:     cm.Name,
                Address:  cm.Address,
                Business: cm.Business,
            })
        }
    }

    var siblings []models.FamilyMember
    if req.Siblings != nil {
        for _, sibling := range req.Siblings {
            siblings = append(siblings, models.FamilyMember{
                Relationship: "sibling",
                Name:         sibling.Name,
                Age:          sibling.Age,
                Address:      sibling.Address,
            })
        }
    }

    var dependents []models.FamilyMember
    if req.Dependents != nil {
        for _, dependent := range req.Dependents {
            dependents = append(dependents, models.FamilyMember{
                Relationship: "dependent",
                Name:         dependent.Name,
                Surname:      dependent.Surname,
                Age:          dependent.Age,
                Address:      dependent.Address,
            })
        }
    }

    return &models.ClientWithRelatedData{
        Client:     client,
        Income:     income,
        Loan:       loan,
        CoMakers:   coMakers,
        Family:     req.Family,
        Siblings:   siblings,
        Spouse:     spouse,
        Dependents: dependents,
    }, nil
}

// GetClientsWithActiveLoans retrieves clients with their active loans for payment management
func (s *ClientService) GetClientsWithActiveLoans() ([]models.Client, error) {
    clients, err := s.clientRepo.FindWithActiveLoans()
    if err != nil {
        return nil, fmt.Errorf("failed to get clients with active loans: %w", err)
    }
    return clients, nil
}
