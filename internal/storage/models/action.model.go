package models

import "time"

// Action defines what should happen to a domain — e.g., allow, block, log.
type Action struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null;"` // "allow", "block", "log"
	CreatedAt time.Time
	UpdatedAt time.Time
}
