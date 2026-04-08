package metrics

import (
	"testing"
	"time"
)

func TestBucketForLatency(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want LatencyBucket
	}{
		{1 * time.Millisecond, BucketLT5ms},
		{4 * time.Millisecond, BucketLT5ms},
		{5 * time.Millisecond, BucketLT10ms},
		{9 * time.Millisecond, BucketLT10ms},
		{10 * time.Millisecond, BucketLT20ms},
		{49 * time.Millisecond, BucketLT50ms},
		{99 * time.Millisecond, BucketLT100ms},
		{249 * time.Millisecond, BucketLT250ms},
		{499 * time.Millisecond, BucketLT500ms},
		{500 * time.Millisecond, BucketGTE500ms},
		{2 * time.Second, BucketGTE500ms},
	}
	for _, tt := range tests {
		if got := BucketForLatency(tt.d); got != tt.want {
			t.Errorf("BucketForLatency(%v) = %d, want %d", tt.d, got, tt.want)
		}
	}
}

func TestBucketUpperBound(t *testing.T) {
	if got := BucketUpperBound(BucketLT5ms); got != 5*time.Millisecond {
		t.Errorf("BucketUpperBound(LT5ms) = %v, want 5ms", got)
	}
	if got := BucketUpperBound(BucketGTE500ms); got != 5*time.Second {
		t.Errorf("BucketUpperBound(GTE500ms) = %v, want 5s", got)
	}
	// Out of range returns max
	if got := BucketUpperBound(LatencyBucket(-1)); got != 5*time.Second {
		t.Errorf("BucketUpperBound(-1) = %v, want 5s", got)
	}
}

func TestRecord_And_Aggregate(t *testing.T) {
	m := NewQueryMetrics()

	m.Record(1*time.Millisecond, true)
	m.Record(50*time.Millisecond, true)
	m.Record(600*time.Millisecond, false)

	agg := m.Aggregate()
	if agg.Total != 3 {
		t.Errorf("expected 3 total, got %d", agg.Total)
	}
	if agg.Errors != 1 {
		t.Errorf("expected 1 error, got %d", agg.Errors)
	}
	if agg.Buckets[BucketLT5ms] != 1 {
		t.Errorf("expected 1 in <5ms bucket, got %d", agg.Buckets[BucketLT5ms])
	}
	if agg.Buckets[BucketLT100ms] != 1 {
		t.Errorf("expected 1 in <100ms bucket, got %d", agg.Buckets[BucketLT100ms])
	}
	if agg.Buckets[BucketGTE500ms] != 1 {
		t.Errorf("expected 1 in >=500ms bucket, got %d", agg.Buckets[BucketGTE500ms])
	}
}

func TestEstimatePercentile(t *testing.T) {
	var buckets [BucketCount]uint64
	buckets[BucketLT5ms] = 50
	buckets[BucketLT10ms] = 40
	buckets[BucketLT100ms] = 10

	p50 := EstimatePercentile(buckets, 0.50)
	if p50 != 5*time.Millisecond {
		t.Errorf("p50 = %v, want 5ms", p50)
	}

	// 50+40=90, total=100, p90 target=90 → cumulative hits 90 at LT10ms bucket
	p90 := EstimatePercentile(buckets, 0.90)
	if p90 != 10*time.Millisecond {
		t.Errorf("p90 = %v, want 10ms", p90)
	}

	// p95 target=95, cumulative 90 at LT10ms, 100 at LT100ms → LT100ms
	p95 := EstimatePercentile(buckets, 0.95)
	if p95 != 100*time.Millisecond {
		t.Errorf("p95 = %v, want 100ms", p95)
	}
}

func TestEstimatePercentile_Empty(t *testing.T) {
	var buckets [BucketCount]uint64
	if got := EstimatePercentile(buckets, 0.50); got != 0 {
		t.Errorf("expected 0 for empty buckets, got %v", got)
	}
}
