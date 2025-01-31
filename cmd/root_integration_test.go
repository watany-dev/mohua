//go:build integration

package cmd

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestExecute_Integration(t *testing.T) {
	skipIfNoAWSCredentials(t)

	// Reset command before test
	resetCommand()

	t.Run("basic execution", func(t *testing.T) {
		err := Execute()
		assert.NoError(t, err)
	})
}

func TestExecuteWithFlags_Integration(t *testing.T) {
	skipIfNoAWSCredentials(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "with region flag",
			args:    []string{"-r", "us-west-2"},
			wantErr: false,
		},
		{
			name:    "with json flag",
			args:    []string{"-j"},
			wantErr: false,
		},
		{
			name:    "with both flags",
			args:    []string{"-r", "us-west-2", "-j"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset command before each test
			resetCommand()

			// Save original args
			oldArgs := os.Args
			// Set up new args for test
			os.Args = append([]string{"mohua"}, tt.args...)

			// Reset args after test
			defer func() {
				os.Args = oldArgs
			}()

			err := Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExecuteWithInvalidFlags_Integration(t *testing.T) {
	skipIfNoAWSCredentials(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "with invalid region",
			args:    []string{"-r", "invalid-region"},
			wantErr: true,
		},
		{
			name:    "with empty region",
			args:    []string{"-r", ""},
			wantErr: false, // Empty region should fall back to default
		},
		{
			name:    "with unknown flag",
			args:    []string{"--unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset command before each test
			resetCommand()

			// Save original args
			oldArgs := os.Args
			// Set up new args for test
			os.Args = append([]string{"mohua"}, tt.args...)

			// Reset args after test
			defer func() {
				os.Args = oldArgs
			}()

			err := Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
