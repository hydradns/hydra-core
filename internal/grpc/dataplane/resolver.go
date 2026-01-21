// SPDX-License-Identifier: GPL-3.0-or-later
package dataplane

import (
	"context"
	"time"

	"github.com/lopster568/phantomDNS/internal/dnsengine"
	pb "github.com/lopster568/phantomDNS/internal/gen/proto/phantomdns/v1"
	"github.com/lopster568/phantomDNS/internal/logger"
)

type ResolverReporter struct {
	client      pb.DataPlaneReporterServiceClient
	engine      *dnsengine.Engine
	dataplaneID string
	interval    time.Duration
}

func NewResolverReporter(
	client pb.DataPlaneReporterServiceClient,
	engine *dnsengine.Engine,
	dataplaneID string,
) *ResolverReporter {
	return &ResolverReporter{
		client:      client,
		engine:      engine,
		dataplaneID: dataplaneID,
		interval:    10 * time.Second,
	}
}

func (r *ResolverReporter) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.reportOnce(ctx)
		}
	}
}

func (r *ResolverReporter) reportOnce(parent context.Context) {
	snapshots := r.engine.ResolverSnapshot()

	if len(snapshots) == 0 {
		return
	}

	resolvers := make([]*pb.ResolverStatus, 0, len(snapshots))
	for _, s := range snapshots {
		var lastSuccess int64
		if !s.LastSuccess.IsZero() {
			lastSuccess = s.LastSuccess.Unix()
		}

		resolvers = append(resolvers, &pb.ResolverStatus{
			ResolverId:      s.ID,
			Healthy:         s.Healthy,
			AvgLatencyMs:    s.AvgLatencyMs,
			LastError:       s.LastError,
			LastSuccessUnix: lastSuccess,
		})
	}

	req := &pb.ResolverStatusReport{
		DataplaneId:   r.dataplaneID,
		TimestampUnix: time.Now().Unix(),
		Resolvers:     resolvers,
	}

	ctx, cancel := context.WithTimeout(parent, 2*time.Second)
	defer cancel()

	if _, err := r.client.ReportResolverStatus(ctx, req); err != nil {
		logger.Log.Warnf("failed to report resolver status: %v", err)
	}
}
