package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event holds the data for a single drift notification.
type Event struct {
	StackName string
	Provider  string
	Level     Level
	Message   string
	Timestamp time.Time
}

// Notifier sends drift events to one or more sinks.
type Notifier struct {
	writer    io.Writer
	threshold Level
}

// Option is a functional option for Notifier.
type Option func(*Notifier)

// WithWriter overrides the default stdout writer.
func WithWriter(w io.Writer) Option {
	return func(n *Notifier) {
		n.writer = w
	}
}

// WithThreshold sets the minimum level that will be emitted.
func WithThreshold(l Level) Option {
	return func(n *Notifier) {
		n.threshold = l
	}
}

// NewNotifier constructs a Notifier with the given options.
func NewNotifier(opts ...Option) *Notifier {
	n := &Notifier{
		writer:    os.Stdout,
		threshold: LevelInfo,
	}
	for _, o := range opts {
		o(n)
	}
	return n
}

// Send emits the event if its level meets the configured threshold.
func (n *Notifier) Send(e Event) error {
	if !n.shouldEmit(e.Level) {
		return nil
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	_, err := fmt.Fprintf(
		n.writer,
		"[%s] %s %-12s %-20s %s\n",
		e.Timestamp.UTC().Format(time.RFC3339),
		e.Level,
		e.Provider,
		e.StackName,
		e.Message,
	)
	return err
}

func (n *Notifier) shouldEmit(l Level) bool {
	return levelRank(l) >= levelRank(n.threshold)
}

func levelRank(l Level) int {
	switch l {
	case LevelInfo:
		return 0
	case LevelWarn:
		return 1
	case LevelError:
		return 2
	default:
		return -1
	}
}
