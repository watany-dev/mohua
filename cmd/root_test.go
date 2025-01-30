package cmd

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

// resetCommand resets the root command and its flags to their initial state
func resetCommand() {
	rootCmd.ResetFlags()
	region = ""
	jsonOutput = false
}

// TestExecute tests the basic execution without any flags
func TestExecute(t *testing.T) {
	// Reset command before test
	resetCommand()

	t.Run("basic execution", func(t *testing.T) {
		err := Execute()
		assert.NoError(t, err)
	})
}

// TestExecuteWithFlags tests execution with various valid flag combinations
func TestExecuteWithFlags(t *testing.T) {
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

// TestExecuteWithInvalidFlags tests execution with invalid flag combinations
func TestExecuteWithInvalidFlags(t *testing.T) {
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
