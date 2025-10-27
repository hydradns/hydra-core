package blocklist

// SPDX-License-Identifier: GPL-3.0-or-later

import "time"

// SourceConfig defines a single blocklist source as declared in configuration.
type SourceConfig struct {
	Name           string        `json:"name"`            // unique source name
	URL            string        `json:"url"`             // remote list URL
	Format         string        `json:"format"`          // e.g. "hosts", "adblock", "domains"
	Category       string        `json:"category"`        // malware, ads, tracking, etc.
	UpdateInterval time.Duration `json:"update_interval"` // interval for auto refresh
	Enabled        bool          `json:"enabled"`         // toggle source activation
}

// BlocklistEntry represents a single parsed and normalized domain record.
type BlocklistEntry struct {
	Domain   string    // normalized domain (e.g., "ads.example.com")
	Source   string    // which list it came from
	Category string    // category tag
	Updated  time.Time // when this record was last updated
}

// ParsedBlocklist holds the result of parsing a single source.
type ParsedBlocklist struct {
	SourceName string
	Entries    []BlocklistEntry
	FetchedAt  time.Time
}

// UpdateResult summarizes a blocklist update attempt.
type UpdateResult struct {
	Source     string
	Success    bool
	Added      int
	Removed    int
	Duration   time.Duration
	ErrMessage string
}

// SourceHealth tracks recent performance metrics for a blocklist source.
type SourceHealth struct {
	Name           string
	LastChecked    time.Time
	LastSuccess    time.Time
	FailureCount   int
	AverageLatency time.Duration
	LastError      string
}

// SupportedFormats enumerates recognized parser formats.
const (
	FormatHosts   = "hosts"
	FormatAdBlock = "adblock"
	FormatDomains = "domains"
)
