// SPDX-License-Identifier: GPL-3.0-or-later
package blocklist

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/lopster568/phantomDNS/internal/blocklist/fetcher"
	"github.com/lopster568/phantomDNS/internal/blocklist/parser"
	"github.com/lopster568/phantomDNS/internal/logger"
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"github.com/lopster568/phantomDNS/internal/storage/repositories"
)

// Engine orchestrates fetching, parsing, and persisting blocklist snapshots.
type Engine struct {
	fetcher *fetcher.HTTPFetcher
	repo    repositories.BlocklistRepository
}

// NewEngine constructs a new blocklist engine.
func NewEngine(repo repositories.BlocklistRepository) *Engine {
	return &Engine{
		fetcher: fetcher.NewHTTPFetcher(),
		repo:    repo,
	}
}

// UpdateSource fetches and updates a specific blocklist source.
// It skips if the source is disabled or unchanged (via ETag/If-None-Match).
func (m *Engine) UpdateSource(ctx context.Context, src models.BlocklistSource, knownETag string) error {
	if !src.Enabled {
		logger.Log.Infof("source %s is disabled; skipping", src.Name)
		return nil
	}

	body, etag, err := m.fetcher.Fetch(ctx, toSourceConfig(src), knownETag)
	if err != nil {
		return fmt.Errorf("fetch failed for %s: %w", src.ID, err)
	}

	if body == nil {
		logger.Log.Infof("source %s not modified (etag=%s)", src.ID, etag)
		return nil
	}

	p, ok := parser.Get(src.Format)
	if !ok {
		return fmt.Errorf("no parser available for format %q (source %s)", src.Format, src.ID)
	}

	parsedEntries, err := p.Parse(body)
	if err != nil {
		return fmt.Errorf("parse failed for %s: %w", src.ID, err)
	}

	// Transform parsed entries into DB model objects
	entries := make([]models.BlocklistEntry, len(parsedEntries))
	now := time.Now()
	for i, pe := range parsedEntries {
		entries[i] = models.BlocklistEntry{
			Domain:    pe.Domain,
			SourceID:  src.ID,
			Category:  src.Category,
			CreatedAt: now,
		}
	}

	// Compute checksum for version tracking
	checksum := hash(body)

	// Save atomically to DB (snapshot + entries)
	snapshot, err := m.repo.SaveSnapshotWithEntries(src, checksum, entries)
	if err != nil {
		return fmt.Errorf("save snapshot failed for %s: %w", src.ID, err)
	}

	logger.Log.Infof("updated blocklist source=%s snapshotID=%d entries=%d etag=%s", src.ID, snapshot.ID, len(entries), etag)
	return nil
}

func (e *Engine) List() ([]string, error) {
	hosts, err := e.repo.GetAll() // assuming repo has GetAll or similar
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

// For a quick debug utility:
func (e *Engine) PrintAll() error {
	hosts, err := e.List()
	if err != nil {
		return err
	}
	for _, h := range hosts {
		fmt.Println(h)
	}
	return nil
}

// hash computes a hex SHA-256 of the given content.
func hash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// toSourceConfig converts a DB source model into a fetcher/parser-friendly struct.
func toSourceConfig(s models.BlocklistSource) parser.SourceConfig {
	return parser.SourceConfig{
		ID:       s.ID,
		Name:     s.Name,
		URL:      s.URL,
		Format:   s.Format,
		Category: s.Category,
	}
}
