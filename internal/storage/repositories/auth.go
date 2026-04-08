package repositories

import (
	"errors"

	"github.com/lopster568/phantomDNS/internal/storage/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) AuthRepository {
	return AuthRepository{db: db}
}

// IsSetup returns true if an admin credential has been created.
func (r AuthRepository) IsSetup() (bool, error) {
	var count int64
	if err := r.db.Model(&models.AdminCredential{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateAdmin creates the singleton admin credential. Returns error if already exists.
func (r AuthRepository) CreateAdmin(passwordHash, apiKey string) error {
	setup, err := r.IsSetup()
	if err != nil {
		return err
	}
	if setup {
		return errors.New("admin already exists")
	}

	return r.db.Create(&models.AdminCredential{
		ID:           1,
		PasswordHash: passwordHash,
		APIKey:       apiKey,
	}).Error
}

// GetAdmin returns the singleton admin credential.
func (r AuthRepository) GetAdmin() (*models.AdminCredential, error) {
	var admin models.AdminCredential
	if err := r.db.First(&admin, 1).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// ValidateAPIKey checks if the provided key matches the stored API key.
func (r AuthRepository) ValidateAPIKey(key string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.AdminCredential{}).Where("api_key = ?", key).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
