package metrics

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Reporter prints a human-readable metrics summary to a writer.
type Reporter struct {
	w io.Writer
}

// NewReporter creates a Reporter that writes to w.
// If w is nil, os.Stdout is used.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Print writes a formatted summary table to the reporter's writer.
func (r *Reporter) Print(s Summary) {
	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "METRIC\tVALUE")
	fmt.Fprintln(tw, "------\t-----")
	fmt.Fprintf(tw, "Total checks\t%d\n", s.Total)
	fmt.Fprintf(tw, "Drifted\t%d\n", s.Drifted)
	fmt.Fprintf(tw, "Errors\t%d\n", s.Errors)
	fmt.Fprintf(tw, "Avg latency\t%s\n", roundDuration(s.AvgLatency))
	_ = tw.Flush()
}

// roundDuration truncates a duration to millisecond precision for display.
func roundDuration(d time.Duration) string {
	if d == 0 {
		return "0ms"
	}
	return d.Round(time.Millisecond).String()
}
