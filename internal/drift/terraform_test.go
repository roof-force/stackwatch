package drift

import (
	"context"
	"testing"
	"time"
)

func TestNewTerraformDetector_EmptyWorkingDir(t *testing.T) {
	_, err := NewTerraformDetector("", 0)
	if err == nil {
		t.Fatal("expected error for empty workingDir, got nil")
	}
}

func TestNewTerraformDetector_DefaultTimeout(t *testing.T) {
	d, err := NewTerraformDetector("/some/path", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Timeout != 5*time.Minute {
		t.Errorf("expected default timeout of 5m, got %s", d.Timeout)
	}
}

func TestNewTerraformDetector_CustomTimeout(t *testing.T) {
	expected := 2 * time.Minute
	d, err := NewTerraformDetector("/some/path", expected)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Timeout != expected {
		t.Errorf("expected timeout %s, got %s", expected, d.Timeout)
	}
}

func TestDetect_ContextCancelled(t *testing.T) {
	d, err := NewTerraformDetector("/tmp", 10*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err = d.Detect(ctx)
	if err == nil {
		t.Fatal("expected error when context is cancelled, got nil")
	}
}

func TestParsePlanOutput(t *testing.T) {
	input := `Terraform will perform the following actions:

  # aws_instance.web will be updated in-place
  ~ resource "aws_instance" "web" {
      ~ ami = "ami-old" -> "ami-new"
    }

  # aws_s3_bucket.data must be replaced
  - resource "aws_s3_bucket" "data" {
    }
`
	changes := parsePlanOutput(input)
	if len(changes) != 2 {
		t.Errorf("expected 2 changes, got %d: %v", len(changes), changes)
	}
}
