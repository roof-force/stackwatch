package output

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

// DriftStatus represents the drift state of a stack resource.
type DriftStatus string

const (
	StatusDrifted    DriftStatus = "DRIFTED"
	StatusInSync     DriftStatus = "IN_SYNC"
	StatusNotChecked DriftStatus = "NOT_CHECKED"
	StatusUnknown    DriftStatus = "UNKNOWN"
)

// DriftResult holds the result of a drift detection run.
type DriftResult struct {
	StackName  string
	Provider   string
	Status     DriftStatus
	DriftedAt  time.Time
	Details    []string
	Error      error
}

// Formatter writes drift results to an output destination.
type Formatter struct {
	out    io.Writer
	noColor bool
}

// NewFormatter creates a Formatter writing to the given writer.
func NewFormatter(out io.Writer, noColor bool) *Formatter {
	if noColor {
		color.NoColor = true
	}
	if out == nil {
		out = os.Stdout
	}
	return &Formatter{out: out, noColor: noColor}
}

// Print writes a formatted drift result to the output.
func (f *Formatter) Print(r DriftResult) {
	timestamp := r.DriftedAt.Format(time.RFC3339)

	if r.Error != nil {
		fmt.Fprintf(f.out, "[%s] %-12s %-30s %s\n",
			timestamp,
			r.Provider,
			r.StackName,
			color.RedString("ERROR: %v", r.Error),
		)
		return
	}

	var statusStr string
	switch r.Status {
	case StatusDrifted:
		statusStr = color.YellowString(string(r.Status))
	case StatusInSync:
		statusStr = color.GreenString(string(r.Status))
	default:
		statusStr = color.CyanString(string(r.Status))
	}

	fmt.Fprintf(f.out, "[%s] %-12s %-30s %s\n",
		timestamp, r.Provider, r.StackName, statusStr)

	for _, detail := range r.Details {
		fmt.Fprintf(f.out, "  → %s\n", detail)
	}
}
