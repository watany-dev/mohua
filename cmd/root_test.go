//go:build !integration

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

// mockExecute is a helper function that executes the command with a mock client
func mockExecute(t *testing.T, args []string) error {
	// Reset command before test
	resetCommand()

	// Save original args
	oldArgs := os.Args
	// Set up new args for test
	os.Args = append([]string{"mohua"}, args...)

	// Reset args after test
	defer func() {
		os.Args = oldArgs
	}()

	return Execute()
}

func TestExecute_Unit(t *testing.T) {
	t.Run("with no resources", func(t *testing.T) {
		err := mockExecute(t, []string{})
		assert.NoError(t, err)
	})
}

func TestExecuteWithFlags_Unit(t *testing.T) {
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
			err := mockExecute(t, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExecuteWithInvalidFlags_Unit(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "with unknown flag",
			args:    []string{"--unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockExecute(t, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
