// SPDX-License-Identifier: GPL-3.0-or-later
package parser

import "time"

type SourceConfig struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	Format         string `json:"format"` // "hosts","adblock","domains"
	Category       string `json:"category"`
	UpdateInterval int    `json:"update_interval_seconds"`
	Enabled        bool   `json:"enabled"`
	Priority       int    `json:"priority"`
}

type ParsedEntry struct {
	Domain   string
	SourceID string
	Category string
	Fetched  time.Time
}
