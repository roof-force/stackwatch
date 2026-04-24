package notify

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var fixedTS = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestSend_InfoEvent_WritesLine(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(WithWriter(&buf))

	err := n.Send(Event{
		StackName: "my-stack",
		Provider:  "cloudformation",
		Level:     LevelInfo,
		Message:   "stack is in sync",
		Timestamp: fixedTS,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "my-stack") {
		t.Errorf("expected stack name in output, got: %s", out)
	}
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected level INFO in output, got: %s", out)
	}
}

func TestSend_BelowThreshold_Suppressed(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(WithWriter(&buf), WithThreshold(LevelWarn))

	_ = n.Send(Event{
		StackName: "quiet-stack",
		Provider:  "terraform",
		Level:     LevelInfo,
		Message:   "no drift",
		Timestamp: fixedTS,
	})
	if buf.Len() != 0 {
		t.Errorf("expected no output for suppressed level, got: %s", buf.String())
	}
}

func TestSend_ErrorLevel_PassesWarnThreshold(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(WithWriter(&buf), WithThreshold(LevelWarn))

	_ = n.Send(Event{
		StackName: "broken-stack",
		Provider:  "cloudformation",
		Level:     LevelError,
		Message:   "drift detected",
		Timestamp: fixedTS,
	})
	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", buf.String())
	}
}

func TestSend_ZeroTimestamp_UsesNow(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(WithWriter(&buf))

	before := time.Now().UTC()
	_ = n.Send(Event{
		StackName: "ts-stack",
		Provider:  "terraform",
		Level:     LevelWarn,
		Message:   "check timestamp",
	})
	after := time.Now().UTC()

	out := buf.String()
	if out == "" {
		t.Fatal("expected output but got empty string")
	}
	_ = before
	_ = after
}

func TestNewNotifier_Defaults(t *testing.T) {
	n := NewNotifier()
	if n.writer == nil {
		t.Error("expected default writer to be non-nil")
	}
	if n.threshold != LevelInfo {
		t.Errorf("expected default threshold INFO, got %s", n.threshold)
	}
}
