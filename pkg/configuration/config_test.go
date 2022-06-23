// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package configuration

import (
	"context"
	"testing"
)

func TestValidateGLConfig(t *testing.T) {
	tCases := []struct {
		name   string
		config *Config
		expErr bool
	}{
		{
			name: "Valid GL Token Config",
			config: &Config{
				restURL:    "https://my-portal",
				useGLToken: true,
				token:      "havetoken",
			},
		},
		{
			name: "Valid GL Token Config with trf",
			config: &Config{
				restURL: "https://my-portal",
				trf:     func(ctx context.Context) (string, error) { return "", nil },
			},
		},
		{
			name: "Valid with gl_token and trf set",
			config: &Config{
				restURL:    "https://my-portal",
				trf:        func(ctx context.Context) (string, error) { return "", nil },
				useGLToken: true,
			},
		},
		{
			name: "Invalid GL Token Config without token",
			config: &Config{
				restURL:    "https://my-portal",
				useGLToken: true,
			},
			expErr: true,
		},
		{
			name: "Invalid GL Token Config without rest_url",
			config: &Config{
				useGLToken: true,
				trf:        func(ctx context.Context) (string, error) { return "", nil },
			},
			expErr: true,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateGLConfig(*tc.config)

			if got, want := (err != nil), tc.expErr; got != want {
				t.Fatalf("validateGLConfig return didn't match gotErr: %v, wantErr: %v", got, want)
			}
		})
	}
}

func TestValidateMetalConfig(t *testing.T) {
	tCases := []struct {
		name   string
		config *Config
		expErr bool
	}{
		{
			name: "Valid Metal Token Config",
			config: &Config{
				restURL: "https://my-portal",
				token:   "havetoken",
				user:    "B901063C-DB35-4FF",
			},
		},
		{
			name: "Invalid Metal Token Config - without rest_url",
			config: &Config{
				token: "havetoken",
				user:  "B901063C-DB35-4FF",
			},
			expErr: true,
		},
		{
			name: "Invalid Metal Token Config - without user",
			config: &Config{
				restURL: "https://my-portal",
				token:   "havetoken",
			},
			expErr: true,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateMetalConfig(*tc.config)

			if got, want := (err != nil), tc.expErr; got != want {
				t.Fatalf("validateConfig return didn't match gotErr: %v, wantErr: %v", got, want)
			}
		})
	}
}
