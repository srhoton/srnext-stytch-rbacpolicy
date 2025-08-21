package config

import (
	"errors"
	"os"
)

type Config struct {
	WorkspaceKeyID     string
	WorkspaceKeySecret string
	ProjectID          string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		WorkspaceKeyID:     os.Getenv("STYTCH_WORKSPACE_KEY_ID"),
		WorkspaceKeySecret: os.Getenv("STYTCH_WORKSPACE_KEY_SECRET"),
		ProjectID:          os.Getenv("STYTCH_PROJECT_ID"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.WorkspaceKeyID == "" {
		return errors.New("STYTCH_WORKSPACE_KEY_ID environment variable is required")
	}
	if c.WorkspaceKeySecret == "" {
		return errors.New("STYTCH_WORKSPACE_KEY_SECRET environment variable is required")
	}
	if c.ProjectID == "" {
		return errors.New("STYTCH_PROJECT_ID environment variable is required")
	}
	return nil
}
