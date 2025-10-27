package models

import "time"

type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null;"` // e.g., "adult", "phishing"
	Description string `gorm:"type:text;"`            // optional description
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
