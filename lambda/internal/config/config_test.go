package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid configuration",
			envVars: map[string]string{
				"STYTCH_WORKSPACE_KEY_ID":     "test-key-id",
				"STYTCH_WORKSPACE_KEY_SECRET": "test-key-secret",
				"STYTCH_PROJECT_ID":           "test-project-id",
			},
			wantErr: false,
		},
		{
			name: "Missing workspace key ID",
			envVars: map[string]string{
				"STYTCH_WORKSPACE_KEY_SECRET": "test-key-secret",
				"STYTCH_PROJECT_ID":           "test-project-id",
			},
			wantErr: true,
			errMsg:  "STYTCH_WORKSPACE_KEY_ID environment variable is required",
		},
		{
			name: "Missing workspace key secret",
			envVars: map[string]string{
				"STYTCH_WORKSPACE_KEY_ID": "test-key-id",
				"STYTCH_PROJECT_ID":       "test-project-id",
			},
			wantErr: true,
			errMsg:  "STYTCH_WORKSPACE_KEY_SECRET environment variable is required",
		},
		{
			name: "Missing project ID",
			envVars: map[string]string{
				"STYTCH_WORKSPACE_KEY_ID":     "test-key-id",
				"STYTCH_WORKSPACE_KEY_SECRET": "test-key-secret",
			},
			wantErr: true,
			errMsg:  "STYTCH_PROJECT_ID environment variable is required",
		},
		{
			name:    "All environment variables missing",
			envVars: map[string]string{},
			wantErr: true,
			errMsg:  "STYTCH_WORKSPACE_KEY_ID environment variable is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := LoadConfig()

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadConfig() expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("LoadConfig() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("LoadConfig() unexpected error: %v", err)
				}
				if cfg == nil {
					t.Errorf("LoadConfig() returned nil config")
				} else {
					if cfg.WorkspaceKeyID != tt.envVars["STYTCH_WORKSPACE_KEY_ID"] {
						t.Errorf("WorkspaceKeyID = %v, want %v", cfg.WorkspaceKeyID, tt.envVars["STYTCH_WORKSPACE_KEY_ID"])
					}
					if cfg.WorkspaceKeySecret != tt.envVars["STYTCH_WORKSPACE_KEY_SECRET"] {
						t.Errorf("WorkspaceKeySecret = %v, want %v", cfg.WorkspaceKeySecret, tt.envVars["STYTCH_WORKSPACE_KEY_SECRET"])
					}
					if cfg.ProjectID != tt.envVars["STYTCH_PROJECT_ID"] {
						t.Errorf("ProjectID = %v, want %v", cfg.ProjectID, tt.envVars["STYTCH_PROJECT_ID"])
					}
				}
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid config",
			config: Config{
				WorkspaceKeyID:     "key-id",
				WorkspaceKeySecret: "key-secret",
				ProjectID:          "project-id",
			},
			wantErr: false,
		},
		{
			name: "Empty workspace key ID",
			config: Config{
				WorkspaceKeyID:     "",
				WorkspaceKeySecret: "key-secret",
				ProjectID:          "project-id",
			},
			wantErr: true,
			errMsg:  "STYTCH_WORKSPACE_KEY_ID environment variable is required",
		},
		{
			name: "Empty workspace key secret",
			config: Config{
				WorkspaceKeyID:     "key-id",
				WorkspaceKeySecret: "",
				ProjectID:          "project-id",
			},
			wantErr: true,
			errMsg:  "STYTCH_WORKSPACE_KEY_SECRET environment variable is required",
		},
		{
			name: "Empty project ID",
			config: Config{
				WorkspaceKeyID:     "key-id",
				WorkspaceKeySecret: "key-secret",
				ProjectID:          "",
			},
			wantErr: true,
			errMsg:  "STYTCH_PROJECT_ID environment variable is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("validate() expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validate() unexpected error: %v", err)
				}
			}
		})
	}
}
