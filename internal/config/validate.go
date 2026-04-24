package config

import (
	"errors"
	"fmt"
	"time"
)

// ValidationError holds a list of field-level validation issues.
type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	if len(v.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("config validation failed with %d error(s): %v", len(v.Errors), v.Errors)
}

func (v *ValidationError) add(msg string) {
	v.Errors = append(v.Errors, msg)
}

// Validate checks the loaded Config for semantic correctness.
func Validate(cfg *Config) error {
	ve := &ValidationError{}

	if len(cfg.Stacks) == 0 {
		ve.add("at least one stack must be defined")
	}

	for i, s := range cfg.Stacks {
		prefix := fmt.Sprintf("stacks[%d]", i)

		if s.Name == "" {
			ve.add(fmt.Sprintf("%s: name is required", prefix))
		}

		switch s.Type {
		case "cloudformation", "terraform":
			// valid
		case "":
			ve.add(fmt.Sprintf("%s: type is required", prefix))
		default:
			ve.add(fmt.Sprintf("%s: unknown type %q (must be cloudformation or terraform)", prefix, s.Type))
		}

		if s.Type == "cloudformation" && s.Region == "" {
			ve.add(fmt.Sprintf("%s: region is required for cloudformation stacks", prefix))
		}

		if s.Type == "terraform" && s.WorkingDir == "" {
			ve.add(fmt.Sprintf("%s: working_dir is required for terraform stacks", prefix))
		}
	}

	if cfg.Interval < 10*time.Second {
		ve.add(fmt.Sprintf("interval must be at least 10s, got %s", cfg.Interval))
	}

	if len(ve.Errors) > 0 {
		return ve
	}
	return nil
}

// ErrNoStacks is returned when a config has no stacks defined.
var ErrNoStacks = errors.New("no stacks defined in config")
