package models

import "time"

type AdminCredential struct {
	ID           uint   `gorm:"primaryKey"`
	PasswordHash string `gorm:"not null"`
	APIKey       string `gorm:"uniqueIndex;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
