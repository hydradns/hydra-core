// SPDX-License-Identifier: GPL-3.0-or-later
package models

import "time"

// DomainPolicy links a domain to a category and an action.
type DomainPolicy struct {
	ID         uint   `gorm:"primaryKey"`
	Domain     string `gorm:"uniqueIndex;not null;"` // domain name, e.g., "example.com"
	CategoryID uint   `gorm:"index;"`                // reference to Category.ID
	ActionID   uint   `gorm:"index;"`                // reference to Action.ID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
