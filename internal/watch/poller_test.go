package watch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stackwatch/internal/config"
	"github.com/stackwatch/internal/watch"
)

type mockDetector struct {
	name    string
	dtype   string
	drifted bool
	details []string
	err     error
}

func (m *mockDetector) Detect(_ context.Context) (bool, []string, error) {
	return m.drifted, m.details, m.err
}
func (m *mockDetector) Name() string { return m.name }
func (m *mockDetector) Type() string { return m.dtype }

func newTestConfig(interval time.Duration) *config.Config {
	return &config.Config{Interval: interval}
}

func TestPoller_ReceivesResults(t *testing.T) {
	detectors := []watch.Detector{
		&mockDetector{name: "stack-a", dtype: "cloudformation", drifted: true, details: []string{"sg changed"}},
		&mockDetector{name: "stack-b", dtype: "terraform", drifted: false},
	}

	cfg := newTestConfig(10 * time.Second)
	poller := watch.NewPoller(cfg, detectors)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go poller.Run(ctx)

	var results []watch.DriftResult
	for r := range poller.Results() {
		results = append(results, r)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestPoller_PropagatesError(t *testing.T) {
	detectors := []watch.Detector{
		&mockDetector{name: "broken", dtype: "terraform", err: errors.New("exec failed")},
	}

	poller := watch.NewPoller(newTestConfig(10*time.Second), detectors)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go poller.Run(ctx)

	var got watch.DriftResult
	for r := range poller.Results() {
		got = r
	}

	if got.Err == nil {
		t.Fatal("expected error in result, got nil")
	}
}

func TestPoller_ContextCancellation(t *testing.T) {
	detectors := []watch.Detector{
		&mockDetector{name: "stack-x", dtype: "cloudformation"},
	}

	poller := watch.NewPoller(newTestConfig(50*time.Millisecond), detectors)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		poller.Run(ctx)
		close(done)
	}()

	time.Sleep(120 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// success
	case <-time.After(time.Second):
		t.Fatal("poller did not stop after context cancellation")
	}
}
