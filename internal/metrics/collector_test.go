package metrics

import (
	"testing"
	"time"
)

func TestRecord_IncreasesTotal(t *testing.T) {
	c := NewCollector()
	c.Record(DriftEvent{StackName: "stack-a", StackType: "cloudformation"})
	c.Record(DriftEvent{StackName: "stack-b", StackType: "terraform"})

	s := c.Summary()
	if s.Total != 2 {
		t.Fatalf("expected Total=2, got %d", s.Total)
	}
}

func TestRecord_CountsDrifted(t *testing.T) {
	c := NewCollector()
	c.Record(DriftEvent{Drifted: true})
	c.Record(DriftEvent{Drifted: false})
	c.Record(DriftEvent{Drifted: true})

	s := c.Summary()
	if s.Drifted != 2 {
		t.Fatalf("expected Drifted=2, got %d", s.Drifted)
	}
}

func TestRecord_CountsErrors(t *testing.T) {
	c := NewCollector()
	c.Record(DriftEvent{Error: true})
	c.Record(DriftEvent{Error: false})

	s := c.Summary()
	if s.Errors != 1 {
		t.Fatalf("expected Errors=1, got %d", s.Errors)
	}
}

func TestRecord_SetsCheckedAtWhenZero(t *testing.T) {
	c := NewCollector()
	before := time.Now()
	c.Record(DriftEvent{StackName: "x"})
	after := time.Now()

	c.mu.Lock()
	ts := c.events[0].CheckedAt
	c.mu.Unlock()

	if ts.Before(before) || ts.After(after) {
		t.Fatalf("CheckedAt %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestSummary_AvgLatency(t *testing.T) {
	c := NewCollector()
	c.Record(DriftEvent{Latency: 100 * time.Millisecond})
	c.Record(DriftEvent{Latency: 200 * time.Millisecond})

	s := c.Summary()
	if s.AvgLatency != 150*time.Millisecond {
		t.Fatalf("expected AvgLatency=150ms, got %v", s.AvgLatency)
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	c := NewCollector()
	c.Record(DriftEvent{})
	c.Record(DriftEvent{})
	c.Reset()

	s := c.Summary()
	if s.Total != 0 {
		t.Fatalf("expected Total=0 after Reset, got %d", s.Total)
	}
}
