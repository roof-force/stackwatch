package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestTablePrinter_PrintHeader(t *testing.T) {
	var buf bytes.Buffer
	tp := NewTablePrinter(&buf)
	tp.PrintHeader()

	out := buf.String()
	if !strings.Contains(out, "PROVIDER") {
		t.Errorf("expected PROVIDER column in header, got: %s", out)
	}
	if !strings.Contains(out, "STACK") {
		t.Errorf("expected STACK column in header, got: %s", out)
	}
	if !strings.Contains(out, "STATUS") {
		t.Errorf("expected STATUS column in header, got: %s", out)
	}
}

func TestTablePrinter_PrintRow_Drifted(t *testing.T) {
	var buf bytes.Buffer
	tp := NewTablePrinter(&buf)
	tp.PrintRow(DriftResult{
		StackName: "infra-stack",
		Provider:  "terraform",
		Status:    StatusDrifted,
		DriftedAt: fixedTime(),
	})

	out := buf.String()
	if !strings.Contains(out, "infra-stack") {
		t.Errorf("expected stack name in row, got: %s", out)
	}
	if !strings.Contains(out, "DRIFTED") {
		t.Errorf("expected DRIFTED status in row, got: %s", out)
	}
}

func TestTablePrinter_PrintRow_Error(t *testing.T) {
	var buf bytes.Buffer
	tp := NewTablePrinter(&buf)
	tp.PrintRow(DriftResult{
		StackName: "bad-stack",
		Provider:  "cloudformation",
		DriftedAt: fixedTime(),
		Error:     errors.New("timeout"),
	})

	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in row output, got: %s", out)
	}
}

func TestTablePrinter_PrintSummary(t *testing.T) {
	var buf bytes.Buffer
	tp := NewTablePrinter(&buf)

	results := []DriftResult{
		{Status: StatusInSync},
		{Status: StatusDrifted},
		{Status: StatusDrifted},
	}
	tp.PrintSummary(results)

	out := buf.String()
	if !strings.Contains(out, "3 stacks checked") {
		t.Errorf("expected total count in summary, got: %s", out)
	}
	if !strings.Contains(out, "2 drifted") {
		t.Errorf("expected drifted count in summary, got: %s", out)
	}
}

func TestTruncate(t *testing.T) {
	if got := truncate("short", 10); got != "short" {
		t.Errorf("expected 'short', got %q", got)
	}
	long := "this-is-a-very-long-stack-name-that-exceeds-limit"
	if got := truncate(long, 20); len(got) > 20 {
		t.Errorf("truncated string too long: %q", got)
	}
	if !strings.HasSuffix(truncate(long, 20), "...") {
		t.Errorf("expected ellipsis suffix")
	}
}
