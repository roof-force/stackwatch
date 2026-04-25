package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time copy of collector metrics.
type Snapshot struct {
	Total      int
	Drifted    int
	Errors     int
	CheckedAt  time.Time
	AvgLatency time.Duration
}

// SnapshotStore retains the last N snapshots for trend analysis.
type SnapshotStore struct {
	mu       sync.RWMutex
	capacity int
	items    []Snapshot
}

// NewSnapshotStore creates a store that retains up to capacity snapshots.
func NewSnapshotStore(capacity int) *SnapshotStore {
	if capacity <= 0 {
		capacity = 10
	}
	return &SnapshotStore{
		capacity: capacity,
		items:    make([]Snapshot, 0, capacity),
	}
}

// Add appends a snapshot, evicting the oldest when at capacity.
func (s *SnapshotStore) Add(snap Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) >= s.capacity {
		s.items = s.items[1:]
	}
	s.items = append(s.items, snap)
}

// Latest returns the most recent snapshot and a boolean indicating
// whether any snapshot exists.
func (s *SnapshotStore) Latest() (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		return Snapshot{}, false
	}
	return s.items[len(s.items)-1], true
}

// All returns a copy of all stored snapshots in insertion order.
func (s *SnapshotStore) All() []Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Snapshot, len(s.items))
	copy(out, s.items)
	return out
}

// DriftRate returns the fraction of checks that detected drift across
// all retained snapshots. Returns 0 if no snapshots are stored.
func (s *SnapshotStore) DriftRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		return 0
	}
	var totalChecks, totalDrifted int
	for _, snap := range s.items {
		totalChecks += snap.Total
		totalDrifted += snap.Drifted
	}
	if totalChecks == 0 {
		return 0
	}
	return float64(totalDrifted) / float64(totalChecks)
}
