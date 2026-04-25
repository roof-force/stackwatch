package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrint_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Print(Summary{})

	out := buf.String()
	if !strings.Contains(out, "METRIC") || !strings.Contains(out, "VALUE") {
		t.Fatalf("expected header row in output, got:\n%s", out)
	}
}

func TestPrint_ShowsTotals(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Print(Summary{Total: 10, Drifted: 3, Errors: 1})

	out := buf.String()
	for _, want := range []string{"10", "3", "1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestPrint_ShowsAvgLatency(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Print(Summary{Total: 1, AvgLatency: 250 * time.Millisecond})

	out := buf.String()
	if !strings.Contains(out, "250ms") {
		t.Fatalf("expected latency in output, got:\n%s", out)
	}
}

func TestPrint_ZeroLatency(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Print(Summary{})

	out := buf.String()
	if !strings.Contains(out, "0ms") {
		t.Fatalf("expected '0ms' for zero latency, got:\n%s", out)
	}
}

func TestNewReporter_NilWriter_UsesStdout(t *testing.T) {
	// Should not panic when w is nil.
	r := NewReporter(nil)
	if r.w == nil {
		t.Fatal("expected non-nil writer after NewReporter(nil)")
	}
}
