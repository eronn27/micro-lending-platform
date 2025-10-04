package handlers

import (
    "fmt"
    "net/http"
    "strconv"
    "micro-lending-platform/backend/internal/models"
    "micro-lending-platform/backend/internal/services"
    "github.com/gin-gonic/gin"
)

type ClientHandler struct {
    clientService *services.ClientService
}

func NewClientHandler(clientService *services.ClientService) *ClientHandler {
    return &ClientHandler{clientService: clientService}
}

// CheckDuplicate checks if a client with similar name already exists
func (h *ClientHandler) CheckDuplicate(c *gin.Context) {
    name := c.Query("name")
    if name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter required"})
        return
    }

    result, err := h.clientService.CheckDuplicate(name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check duplicate"})
        return
    }

    c.JSON(http.StatusOK, result)
}

// CreateClient handles client creation with complete data structure
func (h *ClientHandler) CreateClient(c *gin.Context) {
    var req models.ClientCreateRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request data", 
            "details": err.Error(),
        })
        return
    }

    // Validate required fields
    if req.Client.FirstName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "First name is required"})
        return
    }
    if req.Client.LastName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Last name is required"})
        return
    }
    if req.Client.ContactNumber == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Contact number is required"})
        return
    }
    if req.Client.HomeAddress == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Home address is required"})
        return
    }

    createdClientData, err := h.clientService.CreateClientWithRelatedData(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create client", 
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Client created successfully with all related data",
        "client":  createdClientData,
    })
}

// CreateSimpleClient handles simple client creation (for backward compatibility)
func (h *ClientHandler) CreateSimpleClient(c *gin.Context) {
    var client models.ClientCreate
    
    if err := c.ShouldBindJSON(&client); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request data", 
            "details": err.Error(),
        })
        return
    }

    createdClient, err := h.clientService.CreateClient(&client)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create client", 
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Client created successfully",
        "client":  createdClient,
    })
}

// GetAllClients retrieves all clients with pagination
func (h *ClientHandler) GetAllClients(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    search := c.Query("search")

    clients, total, err := h.clientService.GetAllClients(page, limit, search)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve clients"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "clients": clients,
        "pagination": gin.H{
            "page":  page,
            "limit": limit,
            "total": total,
            "pages": (total + int64(limit) - 1) / int64(limit),
        },
    })
}

// GetClientByID retrieves a client by ID
func (h *ClientHandler) GetClientByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
        return
    }

    client, err := h.clientService.GetClientByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
        return
    }

    c.JSON(http.StatusOK, client)
}

// GetClientWithDetails retrieves a client with all related data
func (h *ClientHandler) GetClientWithDetails(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
        return
    }

    clientData, err := h.clientService.GetClientWithDetails(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
        return
    }

    c.JSON(http.StatusOK, clientData)
}

// UpdateClient updates client information
func (h *ClientHandler) UpdateClient(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
        return
    }

    var client models.Client
    if err := c.ShouldBindJSON(&client); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    client.ID = uint(id)
    updatedClient, err := h.clientService.UpdateClient(&client)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Client updated successfully",
        "client":  updatedClient,
    })
}

// DeleteClient soft deletes a client
func (h *ClientHandler) DeleteClient(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
        return
    }

    err = h.clientService.DeleteClient(uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}

// RestoreClient restores a soft-deleted client
func (h *ClientHandler) RestoreClient(c *gin.Context) {
    c.JSON(501, gin.H{"message": "Restore client - coming soon"})
}

// SearchClients searches clients by various criteria
func (h *ClientHandler) SearchClients(c *gin.Context) {
    query := c.Query("q")
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search query parameter 'q' is required"})
        return
    }

    clients, total, err := h.clientService.SearchClients(query, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search clients"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "clients": clients,
        "pagination": gin.H{
            "page":  page,
            "limit": limit,
            "total": total,
            "pages": (total + int64(limit) - 1) / int64(limit),
        },
        "query": query,
    })
}

// GetClientStats returns client statistics
func (h *ClientHandler) GetClientStats(c *gin.Context) {
    stats, err := h.clientService.GetClientStats()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get client statistics"})
        return
    }

    c.JSON(http.StatusOK, stats)
}

// ExportClients returns all client data for export
func (h *ClientHandler) ExportClients(c *gin.Context) {
    clients, err := h.clientService.ExportClients()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export clients"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "clients": clients,
        "exported_at": gin.H{
            "iso":  "", // Would be populated with actual timestamp
            "unix": 0,  // Would be populated with actual timestamp
        },
        "total": len(clients),
    })
}

// BulkCreateClients creates multiple clients at once
func (h *ClientHandler) BulkCreateClients(c *gin.Context) {
    var requests []models.ClientCreateRequest
    
    if err := c.ShouldBindJSON(&requests); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    var createdClients []*models.ClientWithRelatedData
    var errors []string

    for i, req := range requests {
        clientData, err := h.clientService.CreateClientWithRelatedData(&req)
        if err != nil {
            errors = append(errors, fmt.Sprintf("Client %d: %s", i+1, err.Error()))
        } else {
            createdClients = append(createdClients, clientData)
        }
    }

    response := gin.H{
        "created_count": len(createdClients),
        "error_count":   len(errors),
        "clients":       createdClients,
    }

    if len(errors) > 0 {
        response["errors"] = errors
        c.JSON(http.StatusMultiStatus, response)
    } else {
        c.JSON(http.StatusCreated, response)
    }
}

// GetClientByControlNumber finds a client by control number
func (h *ClientHandler) GetClientByControlNumber(c *gin.Context) {
    controlNumber := c.Param("controlNumber")
    if controlNumber == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Control number parameter is required"})
        return
    }

    client, err := h.clientService.GetClientByControlNumber(controlNumber)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
        return
    }

    c.JSON(http.StatusOK, client)
}
