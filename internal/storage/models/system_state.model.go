// SPDX-License-Identifier: GPL-3.0-or-later
package models

import "time"

// Statistics will store aggregated DNS query metrics.
type SystemState struct {
	ID            uint `gorm:"primaryKey"`
	DNSEnabled    bool `gorm:"not null"`
	PolicyEnabled bool `gorm:"not null"`
	// UpdatedBy string
	UpdatedAt time.Time
}
