package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ribice/twisk/pkg/config"
)

func TestLoad(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		wantData *config.Configuration
		wantErr  bool
	}{
		{
			name:    "does not exist",
			path:    "no",
			wantErr: true,
		},
		{
			name:    "invalid format",
			path:    "testdata/invalid.yaml",
			wantErr: true,
		},
		{
			name: "success",
			path: "testdata/success.yaml",
			wantData: &config.Configuration{
				Server: config.Server{
					Port:                ":8080",
					ReadTimeoutSeconds:  31,
					WriteTimeoutSeconds: 30,
				},
				DB: config.Database{
					PSN:            "postgre",
					LogQueries:     true,
					TimeoutSeconds: 10,
				},
				JWT: config.JWT{
					Secret:    "changedvalue",
					Duration:  15,
					Algorithm: "HS256",
				},
				App: config.Application{
					MinPasswordStrength: 1,
				},
				OpenAPI: config.OpenAPI{
					Username: "twisk",
					Password: "twisk",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := config.Load(tc.path)
			assert.Equal(t, tc.wantData, cfg)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
