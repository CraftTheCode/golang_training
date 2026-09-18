package config

import (
	"errors"
	"testing"
	"time"
)

func TestLoadWithLookup(t *testing.T) {
	tests := []struct {
		name        string
		envMap      map[string]string
		wantErr     bool
		errField    string
		verifyFunc  func(t *testing.T, cfg *Config)
	}{
		{
			name:   "Default values applied when env is empty",
			envMap: map[string]string{},
			wantErr: false,
			verifyFunc: func(t *testing.T, cfg *Config) {
				if cfg.AppName != "GoApp" {
					t.Errorf("expected default AppName 'GoApp', got %q", cfg.AppName)
				}
				if cfg.Port != 8080 {
					t.Errorf("expected default Port 8080, got %d", cfg.Port)
				}
				if cfg.Environment != "development" {
					t.Errorf("expected default Environment 'development', got %q", cfg.Environment)
				}
				if cfg.ShutdownTimeout != 10*time.Second {
					t.Errorf("expected default timeout 10s, got %v", cfg.ShutdownTimeout)
				}
				if !cfg.EnableMetrics {
					t.Errorf("expected default EnableMetrics true, got %t", cfg.EnableMetrics)
				}
			},
		},
		{
			name: "Valid custom environment variables",
			envMap: map[string]string{
				"APP_NAME":                 "PaymentService",
				"PORT":                     "9090",
				"ENVIRONMENT":              "production",
				"SHUTDOWN_TIMEOUT_SECONDS": "30",
				"ENABLE_METRICS":           "false",
			},
			wantErr: false,
			verifyFunc: func(t *testing.T, cfg *Config) {
				if cfg.AppName != "PaymentService" {
					t.Errorf("expected AppName 'PaymentService', got %q", cfg.AppName)
				}
				if cfg.Port != 9090 {
					t.Errorf("expected Port 9090, got %d", cfg.Port)
				}
				if cfg.Environment != "production" {
					t.Errorf("expected Environment 'production', got %q", cfg.Environment)
				}
				if cfg.ShutdownTimeout != 30*time.Second {
					t.Errorf("expected timeout 30s, got %v", cfg.ShutdownTimeout)
				}
				if cfg.EnableMetrics {
					t.Errorf("expected EnableMetrics false, got %t", cfg.EnableMetrics)
				}
			},
		},
		{
			name: "Invalid port number out of range",
			envMap: map[string]string{
				"PORT": "70000",
			},
			wantErr:  true,
			errField: "PORT",
		},
		{
			name: "Invalid non-integer port",
			envMap: map[string]string{
				"PORT": "eight-zero",
			},
			wantErr:  true,
			errField: "PORT",
		},
		{
			name: "Invalid environment name",
			envMap: map[string]string{
				"ENVIRONMENT": "invalid-env",
			},
			wantErr:  true,
			errField: "ENVIRONMENT",
		},
		{
			name: "Invalid boolean for ENABLE_METRICS",
			envMap: map[string]string{
				"ENABLE_METRICS": "not-a-bool",
			},
			wantErr:  true,
			errField: "ENABLE_METRICS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLookup := func(key string) (string, bool) {
				val, ok := tt.envMap[key]
				return val, ok
			}

			cfg, err := LoadWithLookup(mockLookup)

			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadWithLookup() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				var cfgErr *ConfigError
				if !errors.As(err, &cfgErr) {
					t.Fatalf("expected error of type *ConfigError, got %T (%v)", err, err)
				}
				if cfgErr.Field != tt.errField {
					t.Errorf("expected error field %q, got %q", tt.errField, cfgErr.Field)
				}
				return
			}

			if tt.verifyFunc != nil {
				tt.verifyFunc(t, cfg)
			}
		})
	}
}
