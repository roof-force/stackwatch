package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func fixedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2024-01-15T10:00:00Z")
	return t
}

func TestPrint_InSync(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true)

	f.Print(DriftResult{
		StackName: "my-stack",
		Provider:  "cloudformation",
		Status:    StatusInSync,
		DriftedAt: fixedTime(),
	})

	out := buf.String()
	if !strings.Contains(out, "my-stack") {
		t.Errorf("expected stack name in output, got: %s", out)
	}
	if !strings.Contains(out, "IN_SYNC") {
		t.Errorf("expected IN_SYNC status in output, got: %s", out)
	}
}

func TestPrint_Drifted_WithDetails(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true)

	f.Print(DriftResult{
		StackName: "prod-stack",
		Provider:  "terraform",
		Status:    StatusDrifted,
		DriftedAt: fixedTime(),
		Details:   []string{"aws_s3_bucket.main: tags changed"},
	})

	out := buf.String()
	if !strings.Contains(out, "DRIFTED") {
		t.Errorf("expected DRIFTED in output, got: %s", out)
	}
	if !strings.Contains(out, "aws_s3_bucket.main") {
		t.Errorf("expected detail line in output, got: %s", out)
	}
}

func TestPrint_Error(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true)

	f.Print(DriftResult{
		StackName: "broken-stack",
		Provider:  "cloudformation",
		DriftedAt: fixedTime(),
		Error:     errors.New("access denied"),
	})

	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", out)
	}
	if !strings.Contains(out, "access denied") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestNewFormatter_NilWriter(t *testing.T) {
	f := NewFormatter(nil, true)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
	if f.out == nil {
		t.Error("expected fallback writer to be set")
	}
}
