package metrics

import (
	"sync/atomic"
	"time"
)

type QueryMetrics struct {
	total  atomic.Uint64
	errors atomic.Uint64

	latency []atomic.Uint64
}

func NewQueryMetrics() *QueryMetrics {
	return &QueryMetrics{}
}

func (qm *QueryMetrics) Record(elapsed time.Duration, success bool) {}
