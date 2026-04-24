package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// StackType represents the type of infrastructure stack.
type StackType string

const (
	StackTypeCloudFormation StackType = "cloudformation"
	StackTypeTerraform      StackType = "terraform"

	DefaultInterval = 60 * time.Second
	DefaultTimeout  = 30 * time.Second
)

// StackConfig holds configuration for a single stack to monitor.
type StackConfig struct {
	Name    string    `yaml:"name"`
	Type    StackType `yaml:"type"`
	Region  string    `yaml:"region,omitempty"`
	Profile string    `yaml:"profile,omitempty"`
	WorkDir string    `yaml:"work_dir,omitempty"`
}

// Config is the top-level configuration structure.
type Config struct {
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Stacks   []StackConfig `yaml:"stacks"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Stacks) == 0 {
		return errors.New("config must define at least one stack")
	}
	for _, s := range c.Stacks {
		if s.Name == "" {
			return errors.New("each stack must have a name")
		}
		if s.Type != StackTypeCloudFormation && s.Type != StackTypeTerraform {
			return errors.New("stack type must be \"cloudformation\" or \"terraform\"")
		}
		if s.Type == StackTypeCloudFormation && s.Region == "" {
			return errors.New("cloudformation stack requires a region")
		}
		if s.Type == StackTypeTerraform && s.WorkDir == "" {
			return errors.New("terraform stack requires a work_dir")
		}
	}
	return nil
}

func (c *Config) applyDefaults() {
	if c.Interval <= 0 {
		c.Interval = DefaultInterval
	}
	if c.Timeout <= 0 {
		c.Timeout = DefaultTimeout
	}
}
