// SPDX-License-Identifier: GPL-3.0-or-later
package models

import "time"

type BlocklistSource struct {
	ID        string `gorm:"primaryKey;size:64"`
	Name      string
	URL       string
	Format    string
	Category  string
	Enabled   bool
	Priority  int
	UpdatedAt time.Time
	ETag      string
	LastHash  string
	CreatedAt time.Time
}

type BlocklistSnapshot struct {
	ID        uint   `gorm:"primaryKey"`
	SourceID  string `gorm:"index"`
	CreatedAt time.Time
	Size      int
	Checksum  string
	Path      string // optional file path if persisted to disk
}

type BlocklistEntry struct {
	ID         uint   `gorm:"primaryKey"`
	SnapshotID uint   `gorm:"index"`
	Domain     string `gorm:"index;size:255"`
	SourceID   string
	Category   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
