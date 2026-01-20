// Handles upstream DNS resolvers with connection pooling, retry, and failover.
// SPDX-License-Identifier: GPL-3.0-or-later
package dnsengine

import (
	"net"
	"sort"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/lopster568/phantomDNS/internal/logger"
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"github.com/miekg/dns"
)

type ManagedResolver struct {
	meta  models.UpstreamResolver
	pool  *UpstreamPool
	state atomic.Value // healthy / degraded / down
}

type UpstreamManager struct {
	resolvers []*ManagedResolver
}

type ResolverState string

const (
	StateHealthy  ResolverState = "healthy"
	StateDegraded ResolverState = "degraded"
	StateDown     ResolverState = "down"
)

// NewUpstreamManager builds a pool for each configured resolver
func NewUpstreamManager(resolvers []models.UpstreamResolver, poolSize int) (*UpstreamManager, error) {
	m := &UpstreamManager{}

	for _, r := range resolvers {
		addr := net.JoinHostPort(r.Address, strconv.Itoa(r.Port))

		pool, err := NewUpstreamPool(addr, poolSize)
		if err != nil {
			return nil, err
		}

		mr := &ManagedResolver{
			meta: r,
			pool: pool,
		}
		mr.state.Store(StateHealthy)

		m.resolvers = append(m.resolvers, mr)
	}

	// IMPORTANT: enforce priority ordering once
	sort.Slice(m.resolvers, func(i, j int) bool {
		return m.resolvers[i].meta.Priority < m.resolvers[j].meta.Priority
	})

	return m, nil
}

func (m *UpstreamManager) Close() {
	for _, r := range m.resolvers {
		_ = r.pool.Close()
	}
}

// Exchange forwards query to resolvers with retry+failover
func (m *UpstreamManager) Exchange(q *dns.Msg, timeout time.Duration, maxRetries int) (*dns.Msg, error) {
	var lastErr error

	for _, r := range m.resolvers {
		state := r.state.Load().(ResolverState)
		if state == StateDown {
			continue
		}

		for attempt := 0; attempt < maxRetries; attempt++ {
			resp, err := r.pool.Exchange(q, timeout)
			if err == nil {
				return resp, nil
			}

			lastErr = err
			logger.Log.Warnf(
				"resolver %s (%s) failed (attempt %d): %v",
				r.meta.Name,
				r.meta.ID,
				attempt+1,
				err,
			)
		}

		// downgrade health after repeated failures
		r.state.Store(StateDegraded)
	}

	return nil, lastErr
}
