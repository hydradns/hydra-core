// SPDX-License-Identifier: GPL-3.0-or-later
package repositories

import (
	"time"

	"github.com/lopster568/phantomDNS/internal/storage/models"
	"gorm.io/gorm"
)

// Interface (optional but recommended for tests)
type SystemStateRepository interface {
	Get() (*models.SystemState, error)
	SetDNSEnabled(enabled bool) error
	SetPolicyEnabled(enabled bool) error
}

// Implementation
type SystemStateRepo struct {
	db *gorm.DB
}

func NewSystemStateRepo(db *gorm.DB) *SystemStateRepo {
	return &SystemStateRepo{db: db}
}

// Get returns the single system state row (ID = 1).
// If it does not exist, it is created.
func (r *SystemStateRepo) Get() (*models.SystemState, error) {
	var state models.SystemState
	err := r.db.First(&state, 1).Error
	if err == gorm.ErrRecordNotFound {
		state = models.SystemState{
			ID:            1,
			DNSEnabled:    false,
			PolicyEnabled: false,
			UpdatedAt:     time.Now(),
		}
		if err := r.db.Create(&state).Error; err != nil {
			return nil, err
		}
		return &state, nil
	}
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *SystemStateRepo) SetDNSEnabled(enabled bool) error {
	return r.db.Model(&models.SystemState{}).
		Where("id = ?", 1).
		Updates(map[string]interface{}{
			"dns_enabled": enabled,
			"updated_at":  time.Now(),
		}).Error
}

func (r *SystemStateRepo) SetPolicyEnabled(enabled bool) error {
	return r.db.Model(&models.SystemState{}).
		Where("id = ?", 1).
		Updates(map[string]interface{}{
			"policy_enabled": enabled,
			"updated_at":     time.Now(),
		}).Error
}
