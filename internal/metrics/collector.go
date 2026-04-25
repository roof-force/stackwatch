package metrics

import (
	"sync"
	"time"
)

// Result represents the outcome of a single stack drift check.
type Result struct {
	Drifted   bool
	Errored   bool
	Latency   time.Duration
	CheckedAt time.Time
}

// Summary is a point-in-time view of all recorded results.
type Summary struct {
	Total      int
	Drifted    int
	Errors     int
	CheckedAt  time.Time
	AvgLatency time.Duration
}

// Collector accumulates drift check results.
type Collector struct {
	mu         sync.Mutex
	total      int
	drifted    int
	errors     int
	checkedAt  time.Time
	latencySum time.Duration
	now        func() time.Time
}

// NewCollector returns an initialised Collector.
func NewCollector() *Collector {
	return &Collector{now: time.Now}
}

// Record incorporates a single check result into the collector.
func (c *Collector) Record(r Result) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.total++
	if r.Drifted {
		c.drifted++
	}
	if r.Errored {
		c.errors++
	}
	c.latencySum += r.Latency
	if r.CheckedAt.IsZero() {
		c.checkedAt = c.now()
	} else {
		c.checkedAt = r.CheckedAt
	}
}

// Summary returns a point-in-time copy of the current metrics.
func (c *Collector) Summary() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()
	var avg time.Duration
	if c.total > 0 {
		avg = c.latencySum / time.Duration(c.total)
	}
	return Summary{
		Total:      c.total,
		Drifted:    c.drifted,
		Errors:     c.errors,
		CheckedAt:  c.checkedAt,
		AvgLatency: avg,
	}
}

// Snapshot converts the current Summary into a Snapshot suitable for
// storage in a SnapshotStore.
func (c *Collector) Snapshot() Snapshot {
	s := c.Summary()
	return Snapshot{
		Total:      s.Total,
		Drifted:    s.Drifted,
		Errors:     s.Errors,
		CheckedAt:  s.CheckedAt,
		AvgLatency: s.AvgLatency,
	}
}
