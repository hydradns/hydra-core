// SPDX-License-Identifier: GPL-3.0-or-later
package repositories

import (
	"time"

	"github.com/lopster568/phantomDNS/internal/logger"
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"gorm.io/gorm"
)

// Interface
type StatisticsRepository interface {
	Save(stat *models.Statistics) error
	ListRecent(limit int) ([]models.Statistics, error)
	IncrementCounter(action string) error
}

// Implementation
type GormStatisticsRepo struct {
	db *gorm.DB
}

func NewGormStatisticsRepo(db *gorm.DB) *GormStatisticsRepo {
	return &GormStatisticsRepo{db: db}
}

func (r *GormStatisticsRepo) Save(stat *models.Statistics) error {
	stat.UpdatedAt = time.Now()
	logger.Log.Debug("Saving statistics record")
	logger.Log.Debug("stats", stat)
	return r.db.Save(stat).Error
}

func (r *GormStatisticsRepo) ListRecent(limit int) ([]models.Statistics, error) {
	var stats []models.Statistics
	err := r.db.Order("updated_at desc").Limit(limit).Find(&stats).Error
	return stats, err
}

// IncrementCounter increments the global counters (single-row statistics)
func (r *GormStatisticsRepo) IncrementCounter(action string) error {
	var stats models.Statistics

	// Use single-row stats; create row if not exists (ID = 1)
	if err := r.db.FirstOrCreate(&stats, models.Statistics{ID: 1}).Error; err != nil {
		logger.Log.Error("Failed to get or create statistics row: " + err.Error())
		return err
	}

	switch action {
	case "allow":
		stats.AllowedQueries++
	case "block":
		stats.BlockedQueries++
	case "redirect":
		stats.RedirectedQueries++
	default:
		// treat unknown as total only
	}

	stats.TotalQueries++
	stats.UpdatedAt = time.Now()

	if err := r.db.Save(&stats).Error; err != nil {
		logger.Log.Error("Failed to save statistics: " + err.Error())
		return err
	}

	return nil
}
