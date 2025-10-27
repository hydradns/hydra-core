// SPDX-License-Identifier: GPL-3.0-or-later
package repositories

import (
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"gorm.io/gorm"
)

type Store struct {
	QueryLogs  QueryLogRepository
	Statistics StatisticsRepository
	// Policies  PolicyRepository
	// Stats     StatsRepository
	// add more repos here...
}

func NewStore(db *gorm.DB) *Store {
	// Auto-migrate all models here (central place)
	_ = db.AutoMigrate(
		&models.DNSQuery{},
		&models.Statistics{},
		// &models.Policy{},
		// &models.Statistic{},
		// &models.BlockedDomain{},
		// &models.SystemConfig{},
	)

	return &Store{
		QueryLogs:  NewGormQueryLogRepo(db),
		Statistics: NewGormStatisticsRepo(db),
		// Policies:  NewGormPolicyRepo(db),
		// Stats:     NewGormStatsRepo(db),
	}
}
