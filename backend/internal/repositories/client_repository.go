package repositories

import (
    "micro-lending-platform/backend/internal/models"
    "gorm.io/gorm"
    "time"
    "fmt"
)

type ClientRepository struct {
    db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
    return &ClientRepository{db: db}
}

// CreateClientWithRelatedData creates a new client with all related data in a transaction
func (r *ClientRepository) CreateClientWithRelatedData(clientData *models.ClientWithRelatedData) (*models.ClientWithRelatedData, error) {
    // Start a transaction
    tx := r.db.Begin()
    if tx.Error != nil {
        return nil, tx.Error
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 1. Create Client
    if err := tx.Create(&clientData.Client).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to create client: %w", err)
    }

    // 2. Create Income Information (if provided)
    if clientData.Income != nil {
        clientData.Income.ClientID = clientData.Client.ID
        if err := tx.Create(&clientData.Income).Error; err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to create income information: %w", err)
        }
    }

    // 3. Create Loan (if provided)
    if clientData.Loan != nil {
        clientData.Loan.ClientID = clientData.Client.ID
        if err := tx.Create(&clientData.Loan).Error; err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to create loan: %w", err)
        }

        // 4. Create Co-makers (if provided and linked to the loan)
        for i := range clientData.CoMakers {
            clientData.CoMakers[i].LoanID = clientData.Loan.ID
            if err := tx.Create(&clientData.CoMakers[i]).Error; err != nil {
                tx.Rollback()
                return nil, fmt.Errorf("failed to create co-maker: %w", err)
            }
        }
    }

    // 5. Create Family Members (Parents)
    familyMembers := []models.FamilyMember{}

    // Add father if provided
    if clientData.Family.FatherName != "" {
        familyMembers = append(familyMembers, models.FamilyMember{
            ClientID:     clientData.Client.ID,
            Relationship: "father",
            Name:         clientData.Family.FatherName,
        })
    }

    // Add mother if provided
    if clientData.Family.MotherName != "" {
        familyMembers = append(familyMembers, models.FamilyMember{
            ClientID:     clientData.Client.ID,
            Relationship: "mother",
            Name:         clientData.Family.MotherName,
        })
    }

    // Add family address if provided
    if clientData.Family.Address != "" {
        familyMembers = append(familyMembers, models.FamilyMember{
            ClientID:     clientData.Client.ID,
            Relationship: "family_address",
            Address:      clientData.Family.Address,
        })
    }

    // Create family members
    for i := range familyMembers {
        if err := tx.Create(&familyMembers[i]).Error; err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to create family member: %w", err)
        }
    }

    // 6. Create Siblings
    for i := range clientData.Siblings {
        clientData.Siblings[i].ClientID = clientData.Client.ID
        clientData.Siblings[i].Relationship = "sibling"
        if err := tx.Create(&clientData.Siblings[i]).Error; err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to create sibling: %w", err)
        }
    }

    // 7. Create Spouse (if provided)
    if clientData.Spouse != nil && clientData.Spouse.Name != "" {
        clientData.Spouse.ClientID = clientData.Client.ID
        clientData.Spouse.Relationship = "spouse"
        if err := tx.Create(clientData.Spouse).Error; err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to create spouse: %w", err)
        }
    }

    // 8. Create Dependents
    for i := range clientData.Dependents {
        clientData.Dependents[i].ClientID = clientData.Client.ID
        clientData.Dependents[i].Relationship = "dependent"
        if err := tx.Create(&clientData.Dependents[i]).Error; err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to create dependent: %w", err)
        }
    }

    // Commit the transaction
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return clientData, nil
}

// FindByID finds a client by ID
func (r *ClientRepository) FindByID(id uint) (*models.Client, error) {
    var client models.Client
    result := r.db.First(&client, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &client, nil
}

// FindWithDetails finds a client with all related data
func (r *ClientRepository) FindWithDetails(id uint) (*models.ClientWithRelatedData, error) {
    var client models.Client
    result := r.db.Preload("IncomeInfo").
        Preload("Loans").
        Preload("Loans.Payments").
        Preload("Loans.CoMakers").
        Preload("FamilyMembers").
        Preload("Documents").
        First(&client, id)
    
    if result.Error != nil {
        return nil, result.Error
    }

    // Convert to ClientWithRelatedData structure
    clientData := &models.ClientWithRelatedData{
        Client: client,
    }

    // Extract income information
    if len(client.IncomeInfo) > 0 {
        clientData.Income = &client.IncomeInfo[0]
    }

    // Extract loans
    if len(client.Loans) > 0 {
        clientData.Loan = &client.Loans[0]
        clientData.CoMakers = client.Loans[0].CoMakers
    }

    // Extract family members and categorize them
    var siblings []models.FamilyMember
    var spouse *models.FamilyMember
    family := models.Family{}

    for _, member := range client.FamilyMembers {
        switch member.Relationship {
        case "father":
            family.FatherName = member.Name
        case "mother":
            family.MotherName = member.Name
        case "family_address":
            family.Address = member.Address
        case "sibling":
            siblings = append(siblings, member)
        case "spouse":
            spouse = &member
        case "dependent":
            clientData.Dependents = append(clientData.Dependents, member)
        }
    }

    clientData.Family = family
    clientData.Siblings = siblings
    clientData.Spouse = spouse

    return clientData, nil
}

// FindByControlNumber finds a client by control number
func (r *ClientRepository) FindByControlNumber(controlNumber string) (*models.Client, error) {
    var client models.Client
    result := r.db.Where("control_number = ?", controlNumber).First(&client)
    if result.Error != nil {
        return nil, result.Error
    }
    return &client, nil
}

// FindAll retrieves all clients with pagination and search
func (r *ClientRepository) FindAll(offset, limit int, search string) ([]models.Client, error) {
    var clients []models.Client
    
    query := r.db.Model(&models.Client{})
    
    // Add search filter if provided
    if search != "" {
        searchPattern := "%" + search + "%"
        query = query.Where("first_name LIKE ? OR last_name LIKE ? OR control_number LIKE ? OR contact_number LIKE ?", 
            searchPattern, searchPattern, searchPattern, searchPattern)
    }
    
    // Apply pagination if limits are specified
    if limit > 0 {
        query = query.Offset(offset).Limit(limit)
    }
    
    result := query.Order("created_at DESC").Find(&clients)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return clients, nil
}

// Search searches clients with more comprehensive criteria
func (r *ClientRepository) Search(query string, offset, limit int) ([]models.Client, error) {
    var clients []models.Client
    
    searchPattern := "%" + query + "%"
    
    result := r.db.Where(`
        first_name LIKE ? OR 
        last_name LIKE ? OR 
        control_number LIKE ? OR 
        contact_number LIKE ? OR
        home_address LIKE ? OR
        nickname LIKE ?
    `, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern).
        Offset(offset).Limit(limit).
        Order("created_at DESC").
        Find(&clients)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return clients, nil
}

// Update updates an existing client
func (r *ClientRepository) Update(client *models.Client) (*models.Client, error) {
    result := r.db.Save(client)
    if result.Error != nil {
        return nil, result.Error
    }
    return client, nil
}

// Delete soft deletes a client
func (r *ClientRepository) Delete(id uint) error {
    result := r.db.Delete(&models.Client{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("client not found")
    }
    return nil
}

// Count returns the total number of clients matching the search criteria
func (r *ClientRepository) Count(search string) (int64, error) {
    var count int64
    
    query := r.db.Model(&models.Client{})
    
    if search != "" {
        searchPattern := "%" + search + "%"
        query = query.Where("first_name LIKE ? OR last_name LIKE ? OR control_number LIKE ?", 
            searchPattern, searchPattern, searchPattern)
    }
    
    result := query.Count(&count)
    return count, result.Error
}

// CountSimilarClients counts clients with similar names
func (r *ClientRepository) CountSimilarClients(name string) (int64, error) {
    var count int64
    searchPattern := "%" + name + "%"
    
    result := r.db.Model(&models.Client{}).
        Where("first_name LIKE ? OR last_name LIKE ? OR first_name || ' ' || last_name LIKE ?", 
            searchPattern, searchPattern, searchPattern).
        Count(&count)
    
    return count, result.Error
}

// HasActiveLoans checks if a client has active loans
func (r *ClientRepository) HasActiveLoans(clientID uint) (bool, error) {
    var count int64
    result := r.db.Model(&models.Loan{}).
        Where("client_id = ? AND status IN ?", clientID, []string{"Active", "Overdue"}).
        Count(&count)
    
    return count > 0, result.Error
}

// CountActive counts active clients (clients with recent activity)
func (r *ClientRepository) CountActive() (int64, error) {
    var count int64
    result := r.db.Model(&models.Client{}).
        Where("updated_at > ?", time.Now().AddDate(0, -6, 0)). // Active in last 6 months
        Count(&count)
    
    return count, result.Error
}

// CountNewThisMonth counts clients created this month
func (r *ClientRepository) CountNewThisMonth() (int64, error) {
    var count int64
    now := time.Now()
    firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    
    result := r.db.Model(&models.Client{}).
        Where("created_at >= ?", firstOfMonth).
        Count(&count)
    
    return count, result.Error
}

// CountWithActiveLoans counts clients with active loans
func (r *ClientRepository) CountWithActiveLoans() (int64, error) {
    var count int64
    result := r.db.Model(&models.Client{}).
        Joins("INNER JOIN loans ON loans.client_id = clients.id").
        Where("loans.status IN ?", []string{"Active", "Overdue"}).
        Distinct("clients.id").
        Count(&count)
    
    return count, result.Error
}

// CreateSimple inserts a new client into the database (for backward compatibility)
func (r *ClientRepository) CreateSimple(client *models.Client) (*models.Client, error) {
    result := r.db.Create(client)
    if result.Error != nil {
        return nil, result.Error
    }
    return client, nil
}
