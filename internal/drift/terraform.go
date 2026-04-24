package drift

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// TerraformDetector detects drift in Terraform-managed infrastructure.
type TerraformDetector struct {
	WorkingDir string
	Timeout    time.Duration
}

// DriftResult holds the result of a Terraform drift check.
type DriftResult struct {
	HasDrift  bool
	Changes   []string
	DetectedAt time.Time
	Error     error
}

// NewTerraformDetector creates a new TerraformDetector.
// workingDir must be the path to the Terraform root module.
func NewTerraformDetector(workingDir string, timeout time.Duration) (*TerraformDetector, error) {
	if workingDir == "" {
		return nil, fmt.Errorf("workingDir must not be empty")
	}
	if timeout <= 0 {
		timeout = 5 * time.Minute
	}
	return &TerraformDetector{
		WorkingDir: workingDir,
		Timeout:    timeout,
	}, nil
}

// Detect runs `terraform plan` and reports whether drift exists.
func (t *TerraformDetector) Detect(ctx context.Context) (*DriftResult, error) {
	ctx, cancel := context.WithTimeout(ctx, t.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "terraform", "plan", "-detailed-exitcode", "-no-color")
	cmd.Dir = t.WorkingDir

	out, err := cmd.CombinedOutput()
	result := &DriftResult{
		DetectedAt: time.Now(),
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("terraform plan timed out after %s", t.Timeout)
	}

	// terraform plan exits with code 2 when there are changes (drift).
	if exitErr, ok := err.(*exec.ExitError); ok {
		switch exitErr.ExitCode() {
		case 2:
			result.HasDrift = true
			result.Changes = parsePlanOutput(string(out))
			return result, nil
		default:
			return nil, fmt.Errorf("terraform plan failed (exit %d): %s", exitErr.ExitCode(), string(out))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to run terraform plan: %w", err)
	}

	// Exit code 0 means no changes.
	result.HasDrift = false
	return result, nil
}

// parsePlanOutput extracts changed resource lines from terraform plan output.
func parsePlanOutput(output string) []string {
	var changes []string
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			changes = append(changes, trimmed)
		}
	}
	return changes
}
