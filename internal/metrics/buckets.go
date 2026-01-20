// SPDX-License-Identifier: GPL-3.0-or-later
package metrics

import (
	"time"
)

// LatencyBucket represents a latency range index.
// Buckets are ordered from fastest to slowest.
type LatencyBucket int

const (
	BucketLT5ms LatencyBucket = iota
	BucketLT10ms
	BucketLT20ms
	BucketLT50ms
	BucketLT100ms
	BucketLT250ms
	BucketLT500ms
	BucketGTE500ms

	BucketCount
)

// bucketUpperBounds defines the exclusive upper bound
// for each latency bucket.
var bucketUpperBounds = [...]time.Duration{
	5 * time.Millisecond,
	10 * time.Millisecond,
	20 * time.Millisecond,
	50 * time.Millisecond,
	100 * time.Millisecond,
	250 * time.Millisecond,
	500 * time.Millisecond,
	5 * time.Second, // cap worst
}

func BucketForLatency(d time.Duration) LatencyBucket {
	switch {
	case d < 5*time.Millisecond:
		return BucketLT5ms
	case d < 10*time.Millisecond:
		return BucketLT10ms
	case d < 20*time.Millisecond:
		return BucketLT20ms
	case d < 50*time.Millisecond:
		return BucketLT50ms
	case d < 100*time.Millisecond:
		return BucketLT100ms
	case d < 250*time.Millisecond:
		return BucketLT250ms
	case d < 500*time.Millisecond:
		return BucketLT500ms
	default:
		return BucketGTE500ms
	}
}

func BucketUpperBound(b LatencyBucket) time.Duration {
	if b < 0 || b >= BucketCount {
		return bucketUpperBounds[BucketGTE500ms]
	}
	return bucketUpperBounds[b]
}
