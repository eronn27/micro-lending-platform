package repositories

import (
    "micro-lending-platform/backend/internal/models"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
    var user models.User
    result := r.db.Where("username = ?", username).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    result := r.db.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) (*models.User, error) {
    result := r.db.Create(user)
    if result.Error != nil {
        return nil, result.Error
    }
    return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) (*models.User, error) {
    result := r.db.Save(user)
    if result.Error != nil {
        return nil, result.Error
    }
    return user, nil
}

// ListAll returns all users (for admin purposes)
func (r *UserRepository) ListAll() ([]models.User, error) {
    var users []models.User
    result := r.db.Find(&users)
    if result.Error != nil {
        return nil, result.Error
    }
    return users, nil
}
