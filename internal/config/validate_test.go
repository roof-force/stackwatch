package config

import (
	"strings"
	"testing"
	"time"
)

func baseConfig() *Config {
	return &Config{
		Interval: 30 * time.Second,
		Stacks: []Stack{
			{
				Name:   "my-stack",
				Type:   "cloudformation",
				Region: "us-east-1",
			},
		},
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	if err := Validate(baseConfig()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_NoStacks(t *testing.T) {
	cfg := baseConfig()
	cfg.Stacks = nil
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for empty stacks")
	}
	if !strings.Contains(err.Error(), "at least one stack") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MissingStackName(t *testing.T) {
	cfg := baseConfig()
	cfg.Stacks[0].Name = ""
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing stack name")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_UnknownType(t *testing.T) {
	cfg := baseConfig()
	cfg.Stacks[0].Type = "pulumi"
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
	if !strings.Contains(err.Error(), "unknown type") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_CloudFormationMissingRegion(t *testing.T) {
	cfg := baseConfig()
	cfg.Stacks[0].Region = ""
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing region")
	}
	if !strings.Contains(err.Error(), "region is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_TerraformMissingWorkingDir(t *testing.T) {
	cfg := baseConfig()
	cfg.Stacks[0].Type = "terraform"
	cfg.Stacks[0].Region = ""
	cfg.Stacks[0].WorkingDir = ""
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing working_dir")
	}
	if !strings.Contains(err.Error(), "working_dir is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_IntervalTooShort(t *testing.T) {
	cfg := baseConfig()
	cfg.Interval = 5 * time.Second
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for short interval")
	}
	if !strings.Contains(err.Error(), "interval must be at least") {
		t.Errorf("unexpected error message: %v", err)
	}
}
