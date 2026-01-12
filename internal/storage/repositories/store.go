// SPDX-License-Identifier: GPL-3.0-or-later
package repositories

import (
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"gorm.io/gorm"
)

type Store struct {
	QueryLogs   QueryLogRepository
	Blocklist   BlocklistRepository
	Statistics  StatisticsRepository
	SystemState SystemStateRepository
	// Policies  PolicyRepository
	// Stats     StatsRepository
	// add more repos here...
}

func NewStore(db *gorm.DB) *Store {
	// Auto-migrate all models here (central place)
	_ = db.AutoMigrate(
		&models.DNSQuery{},
		&models.BlocklistSource{},
		&models.BlocklistSnapshot{},
		&models.BlocklistEntry{},
		&models.Statistics{},
		&models.SystemState{},
		// &models.Policy{},
		// &models.Statistic{},
		// &models.BlockedDomain{},
		// &models.SystemConfig{},
	)

	return &Store{
		QueryLogs:   NewGormQueryLogRepo(db),
		Blocklist:   NewBlocklistRepo(db),
		Statistics:  NewGormStatisticsRepo(db),
		SystemState: NewSystemStateRepo(db),
		// Policies:  NewGormPolicyRepo(db),
		// Stats:     NewGormStatsRepo(db),
	}
}
