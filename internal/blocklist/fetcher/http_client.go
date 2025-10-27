package fetcher

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lopster568/phantomDNS/internal/blocklist"
	"github.com/lopster568/phantomDNS/internal/logger"
)

type BlocklistFetcher interface {
	Fetch(source blocklist.SourceConfig) ([]byte, error)
}

var registry = map[string]BlocklistFetcher{}

func RegisterFetcher(name string, f BlocklistFetcher) {
	if _, exists := registry[name]; exists {
		panic("fetcher already registered: " + name)
	}
	registry[name] = f
}

// defaulting to http fetcher as fallback
func NewFetcher(name string) BlocklistFetcher {
	if f, exists := registry[name]; exists {
		return f
	}
	return &HTTPFetcher{Client: &http.Client{Timeout: 30 * time.Second}}
}

type HTTPFetcher struct {
	Client *http.Client
}

// Fetch downloads the raw blocklist data for the given source.
func (h *HTTPFetcher) Fetch(source blocklist.SourceConfig) ([]byte, error) {
	if source.URL == "" {
		return nil, errors.New("source URL is empty")
	}
	start := time.Now()
	resp, err := h.Client.Get(source.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Errorf("unexpected HTTP status %d for %s", resp.StatusCode, source.Name)
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body for %s: %w", source.Name, err)
	}

	duration := time.Since(start)
	logger.Log.Infof("[fetcher] fetched %s (%d bytes) in %v\n", source.Name, len(data), duration)
	return data, nil
}
