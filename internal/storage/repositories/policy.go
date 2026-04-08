// SPDX-License-Identifier: GPL-3.0-or-later
package repositories

import (
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"gorm.io/gorm"
)

type PolicyRepository interface {
	List() ([]models.Policy, error)
	GetByID(id string) (*models.Policy, error)
	Create(policy *models.Policy) error
	Update(policy *models.Policy) error
	Delete(id string) error
}

type PolicyRepo struct {
	db *gorm.DB
}

func NewPolicyRepo(db *gorm.DB) *PolicyRepo {
	return &PolicyRepo{db: db}
}

func (r *PolicyRepo) List() ([]models.Policy, error) {
	var policies []models.Policy
	err := r.db.Order("priority desc, id asc").Find(&policies).Error
	return policies, err
}

func (r *PolicyRepo) GetByID(id string) (*models.Policy, error) {
	var policy models.Policy
	err := r.db.First(&policy, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *PolicyRepo) Create(policy *models.Policy) error {
	return r.db.Create(policy).Error
}

func (r *PolicyRepo) Update(policy *models.Policy) error {
	return r.db.Save(policy).Error
}

func (r *PolicyRepo) Delete(id string) error {
	result := r.db.Delete(&models.Policy{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
