package metrics

import (
	"sync"
	"time"
)

// DriftEvent represents a single drift check result recorded for metrics.
type DriftEvent struct {
	StackName string
	StackType string
	Drifted   bool
	Error     bool
	CheckedAt time.Time
	Latency   time.Duration
}

// Collector accumulates drift check statistics across polling cycles.
type Collector struct {
	mu     sync.Mutex
	events []DriftEvent
}

// NewCollector returns an initialised Collector.
func NewCollector() *Collector {
	return &Collector{
		events: make([]DriftEvent, 0, 64),
	}
}

// Record appends a DriftEvent to the collector.
func (c *Collector) Record(e DriftEvent) {
	if e.CheckedAt.IsZero() {
		e.CheckedAt = time.Now()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, e)
}

// Summary returns an aggregated view of all recorded events.
func (c *Collector) Summary() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := Summary{Total: len(c.events)}
	var totalLatency time.Duration
	for _, e := range c.events {
		if e.Drifted {
			s.Drifted++
		}
		if e.Error {
			s.Errors++
		}
		totalLatency += e.Latency
	}
	if s.Total > 0 {
		s.AvgLatency = totalLatency / time.Duration(s.Total)
	}
	return s
}

// Reset clears all recorded events.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = c.events[:0]
}

// Summary holds aggregated metrics for a collection of drift events.
type Summary struct {
	Total      int
	Drifted    int
	Errors     int
	AvgLatency time.Duration
}
