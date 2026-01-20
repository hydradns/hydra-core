package models

import "time"

type ResolverScope string

const (
	ResolverPublic  ResolverScope = "public"
	ResolverPrivate ResolverScope = "private"
)

type UpstreamResolver struct {
	ID       string        `gorm:"primaryKey"`
	Name     string        `gorm:"not null"`
	Address  string        `gorm:"not null"` // IP only (no hostnames)
	Port     int           `gorm:"not null"`
	Scope    ResolverScope `gorm:"type:text;not null"`
	Priority int           `gorm:"not null"` // lower = higher priority
	Enabled  bool          `gorm:"not null;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
