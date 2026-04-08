// SPDX-License-Identifier: GPL-3.0-or-later
package repositories

import (
	"gorm.io/gorm"
)

type Store struct {
	QueryLogs   QueryLogRepository
	Blocklist   BlocklistRepository
	Statistics  StatisticsRepository
	SystemState SystemStateRepository
	Policies    PolicyRepository
	Auth        AuthRepository
}

func NewStore(db *gorm.DB) *Store {
	// Migrations run in db.InitDB() — single source of truth.
	return &Store{
		QueryLogs:   NewGormQueryLogRepo(db),
		Blocklist:   NewBlocklistRepo(db),
		Statistics:  NewGormStatisticsRepo(db),
		SystemState: NewSystemStateRepo(db),
		Policies:    NewPolicyRepo(db),
		Auth:        NewAuthRepo(db),
	}
}
