package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stackwatch/stackwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
interval: 30s
timeout: 10s
stacks:
  - name: my-cf-stack
    type: cloudformation
    region: us-east-1
  - name: my-tf-stack
    type: terraform
    work_dir: /tmp/tf
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected interval 30s, got %v", cfg.Interval)
	}
	if len(cfg.Stacks) != 2 {
		t.Errorf("expected 2 stacks, got %d", len(cfg.Stacks))
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	yaml := `
stacks:
  - name: my-stack
    type: cloudformation
    region: eu-west-1
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != config.DefaultInterval {
		t.Errorf("expected default interval %v, got %v", config.DefaultInterval, cfg.Interval)
	}
	if cfg.Timeout != config.DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", config.DefaultTimeout, cfg.Timeout)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoStacks(t *testing.T) {
	path := writeTemp(t, "stacks: []\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty stacks")
	}
}

func TestLoad_MissingRegionForCF(t *testing.T) {
	yaml := `
stacks:
  - name: bad-stack
    type: cloudformation
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestLoad_MissingWorkDirForTerraform(t *testing.T) {
	yaml := `
stacks:
  - name: bad-tf
    type: terraform
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing work_dir")
	}
}
