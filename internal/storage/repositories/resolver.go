// SPDX-License-Identifier: GPL-3.0-or-later
package repositories

import (
	"fmt"
	"net"
	"strconv"

	"github.com/google/uuid"
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"gorm.io/gorm"
)

// Interface (optional but recommended for tests)
type ResolverRepository interface {
	ListResolvers() ([]models.UpstreamResolver, error)
	GetResolver(id string) (*models.UpstreamResolver, error)
	CreateResolver(resolver *models.UpstreamResolver) error
	UpdateResolver(resolver *models.UpstreamResolver) error
	DisableResolver(id string) error
	BootstrapResolvers(cfgResolvers []string) error
}

// Implementation
type ResolverRepo struct {
	db *gorm.DB
}

func NewResolverRepo(db *gorm.DB) *ResolverRepo {
	return &ResolverRepo{db: db}
}

// ListResolvers returns all upstream resolvers ordered by priority
func (r *ResolverRepo) ListResolvers() ([]models.UpstreamResolver, error) {
	var resolvers []models.UpstreamResolver
	err := r.db.
		Where("enabled = ?", true).
		Order("priority ASC").
		Find(&resolvers).Error
	return resolvers, err
}

// GetResolver returns a specific resolver by ID
func (r *ResolverRepo) GetResolver(id string) (*models.UpstreamResolver, error) {
	var resolver models.UpstreamResolver
	err := r.db.Where("id = ?", id).First(&resolver).Error
	if err != nil {
		return nil, err
	}
	return &resolver, nil
}

// CreateResolver creates a new upstream resolver
func (r *ResolverRepo) CreateResolver(resolver *models.UpstreamResolver) error {
	return r.db.Create(resolver).Error
}

// UpdateResolver updates an existing resolver
func (r *ResolverRepo) UpdateResolver(resolver *models.UpstreamResolver) error {
	return r.db.Save(resolver).Error
}

// DisableResolver deletes a resolver by ID
func (r *ResolverRepo) DisableResolver(id string) error {
	return r.db.Model(&models.UpstreamResolver{}).
		Where("id = ?", id).
		Update("enabled", false).Error
}
func (r *ResolverRepo) BootstrapResolvers(
	cfgResolvers []string,
) error {
	existing, err := r.ListResolvers()
	if err != nil {
		return err
	}

	// DB already initialized → do nothing
	if len(existing) > 0 {
		return nil
	}

	for i, raw := range cfgResolvers {
		host, portStr, err := net.SplitHostPort(raw)
		if err != nil {
			return fmt.Errorf("invalid upstream resolver %q: %w", raw, err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid port in upstream resolver %q: %w", raw, err)
		}

		resolver := &models.UpstreamResolver{
			ID:       uuid.NewString(),
			Name:     "bootstrap-" + host,
			Address:  host,
			Port:     port,
			Scope:    models.ResolverPublic,
			Priority: i,
			Enabled:  true,
		}

		if err := r.CreateResolver(resolver); err != nil {
			return err
		}
	}

	return nil
}
