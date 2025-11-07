// SPDX-License-Identifier: GPL-3.0-or-later
package repositories

import (
	"strings"
	"time"

	"github.com/lopster568/phantomDNS/internal/logger"
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"gorm.io/gorm"
)

// Interface (clean, mockable)
type BlocklistRepository interface {
	SaveSnapshotWithEntries(src models.BlocklistSource, checksum string, entries []models.BlocklistEntry) (models.BlocklistSnapshot, error)
	GetAll() ([]string, error)
	IsBlocked(domain string) (bool, error)
}

// Implementation
type BlocklistRepo struct {
	db *gorm.DB
}

func NewBlocklistRepo(db *gorm.DB) *BlocklistRepo {
	return &BlocklistRepo{db: db}
}

func (r *BlocklistRepo) IsBlocked(domain string) (bool, error) {
	logger.Log.Infof("Checking if domain %s is blocklisted", domain)
	var count int64
	// Normalize domain (lowercase, remove trailing dot)
	d := strings.TrimSuffix(strings.ToLower(domain), ".")
	err := r.db.Model(&models.BlocklistEntry{}).
		Where("domain = ?", d).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *BlocklistRepo) GetAll() ([]string, error) {
	var domains []string
	if err := r.db.Model(&models.BlocklistEntry{}).Pluck("domain", &domains).Error; err != nil {
		return nil, err
	}
	return domains, nil
}

func (r *BlocklistRepo) SaveSnapshotWithEntries(src models.BlocklistSource, checksum string, entries []models.BlocklistEntry) (models.BlocklistSnapshot, error) {
	tx := r.db.Begin()
	snapshot := models.BlocklistSnapshot{
		SourceID: src.ID, CreatedAt: time.Now(), Size: len(entries), Checksum: checksum,
	}
	if err := tx.Create(&snapshot).Error; err != nil {
		tx.Rollback()
		return snapshot, err
	}
	for i := range entries {
		entries[i].SnapshotID = snapshot.ID
		if err := tx.Create(&entries[i]).Error; err != nil {
			tx.Rollback()
			return snapshot, err
		}
	}
	// update source metadata
	src.UpdatedAt = time.Now()
	if err := tx.Save(&src).Error; err != nil {
		tx.Rollback()
		return snapshot, err
	}
	if err := tx.Commit().Error; err != nil {
		return snapshot, err
	}
	return snapshot, nil
}
