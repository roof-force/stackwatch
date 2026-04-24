package output

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	colWidth = 30
	sepLine  = "-"
)

// TablePrinter renders drift results as an ASCII table.
type TablePrinter struct {
	out io.Writer
}

// NewTablePrinter creates a TablePrinter writing to out.
func NewTablePrinter(out io.Writer) *TablePrinter {
	return &TablePrinter{out: out}
}

// PrintHeader writes the table header row.
func (tp *TablePrinter) PrintHeader() {
	sep := strings.Repeat(sepLine, 85)
	fmt.Fprintln(tp.out, sep)
	fmt.Fprintf(tp.out, "%-20s %-12s %-30s %-12s\n",
		"TIMESTAMP", "PROVIDER", "STACK", "STATUS")
	fmt.Fprintln(tp.out, sep)
}

// PrintRow writes a single result row.
func (tp *TablePrinter) PrintRow(r DriftResult) {
	ts := r.DriftedAt.Format(time.Kitchen)
	if r.Error != nil {
		fmt.Fprintf(tp.out, "%-20s %-12s %-30s %-12s\n",
			ts, r.Provider, truncate(r.StackName, colWidth), "ERROR")
		return
	}
	fmt.Fprintf(tp.out, "%-20s %-12s %-30s %-12s\n",
		ts, r.Provider, truncate(r.StackName, colWidth), string(r.Status))
}

// PrintSummary writes a summary line after all rows.
func (tp *TablePrinter) PrintSummary(results []DriftResult) {
	drifted := 0
	for _, r := range results {
		if r.Status == StatusDrifted {
			drifted++
		}
	}
	sep := strings.Repeat(sepLine, 85)
	fmt.Fprintln(tp.out, sep)
	fmt.Fprintf(tp.out, "Total: %d stacks checked, %d drifted\n", len(results), drifted)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
