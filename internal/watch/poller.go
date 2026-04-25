package watch

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/stackwatch/internal/config"
)

// DriftResult holds the outcome of a single drift detection run.
type DriftResult struct {
	StackName string
	StackType string
	Drifted   bool
	Details   []string
	Err       error
	Timestamp time.Time
}

// Detector is the interface that drift detectors must implement.
type Detector interface {
	Detect(ctx context.Context) (bool, []string, error)
	Name() string
	Type() string
}

// Poller periodically runs drift detection for all configured stacks.
type Poller struct {
	detectors []Detector
	interval  time.Duration
	results   chan DriftResult
}

// NewPoller creates a Poller from a config and a set of detectors.
func NewPoller(cfg *config.Config, detectors []Detector) *Poller {
	return &Poller{
		detectors: detectors,
		interval:  cfg.Interval,
		results:   make(chan DriftResult, len(detectors)*2),
	}
}

// Results returns the read-only results channel.
func (p *Poller) Results() <-chan DriftResult {
	return p.results
}

// Run starts polling all detectors at the configured interval until ctx is cancelled.
func (p *Poller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	defer close(p.results)

	p.poll(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *Poller) poll(ctx context.Context) {
	var wg sync.WaitGroup
	for _, d := range p.detectors {
		wg.Add(1)
		go func(det Detector) {
			defer wg.Done()
			drifted, details, err := det.Detect(ctx)
			if ctx.Err() != nil {
				return
			}
			select {
			case p.results <- DriftResult{
				StackName: det.Name(),
				StackType: det.Type(),
				Drifted:   drifted,
				Details:   details,
				Err:       err,
				Timestamp: time.Now(),
			}:
			case <-ctx.Done():
				log.Printf("context cancelled while sending result for %s", det.Name())
			}
		}(d)
	}
	wg.Wait()
}
