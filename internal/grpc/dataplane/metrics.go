// SPDX-License-Identifier: GPL-3.0-or-later
package dataplane

import (
	"context"

	"github.com/lopster568/phantomDNS/internal/dnsengine"
	pb "github.com/lopster568/phantomDNS/internal/gen/proto/phantomdns/v1"
	"github.com/lopster568/phantomDNS/internal/metrics"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetricsService struct {
	engine *dnsengine.Engine
}

func NewMetricsService(e *dnsengine.Engine) *MetricsService {
	return &MetricsService{engine: e}
}

func (s *MetricsService) GetLiveQueryMetrics(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.LiveQueryMetrics, error) {

	agg := s.engine.Metrics().Aggregate()

	p50 := metrics.EstimatePercentile(agg.Buckets, 0.50)
	p95 := metrics.EstimatePercentile(agg.Buckets, 0.95)
	p99 := metrics.EstimatePercentile(agg.Buckets, 0.99)

	return &pb.LiveQueryMetrics{
		WindowSizeSeconds: uint64(5 * 60),

		TotalQueries: agg.Total,
		ErrorQueries: agg.Errors,

		P50Ms: uint64(p50.Milliseconds()),
		P95Ms: uint64(p95.Milliseconds()),
		P99Ms: uint64(p99.Milliseconds()),
	}, nil
}
