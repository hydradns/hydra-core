package dataplane

import (
	"github.com/lopster568/phantomDNS/internal/dnsengine"
)

type MetricsService struct {
	engine *dnsengine.Engine
}

func NewMetricsService(e *dnsengine.Engine) *MetricsService {
	return &MetricsService{
		engine: e,
	}
}
