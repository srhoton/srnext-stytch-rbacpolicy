package main

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		wantErr     bool
	}{
		{
			name:        "Development logger",
			environment: "development",
			wantErr:     false,
		},
		{
			name:        "Production logger",
			environment: "production",
			wantErr:     false,
		},
		{
			name:        "Default logger (empty environment)",
			environment: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			os.Setenv("ENVIRONMENT", tt.environment)
			defer os.Unsetenv("ENVIRONMENT")

			logger, err := initLogger()

			if tt.wantErr {
				if err == nil {
					t.Errorf("initLogger() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("initLogger() unexpected error: %v", err)
				}
				if logger == nil {
					t.Errorf("initLogger() returned nil logger")
				} else {
					// Test that logger is functional
					logger.Info("test message", zap.String("test", "value"))
				}
			}
		})
	}
}
