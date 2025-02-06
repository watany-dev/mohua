package cmd

import (
	"os"
	"testing"

	"mohua/internal/sagemaker"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// resetCommand resets the root command and its flags to their initial state
func resetCommand() {
	rootCmd.ResetFlags()
	region = ""
	jsonOutput = false
}

// mockExecute is a helper function that executes the command with a mock client
func mockExecute(t *testing.T, args []string, client sagemaker.Client) error {
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

	// Store the original NewClient function
	origNewClient := sagemaker.NewClient
	// Replace it with our mock
	sagemaker.NewClient = func(region string) (sagemaker.Client, error) {
		return client, nil
	}
	// Restore the original function after the test
	defer func() {
		sagemaker.NewClient = origNewClient
	}()

	return Execute()
}

func TestExecute_Unit(t *testing.T) {
	mockClient := new(MockSageMakerClient)

	// Setup mock expectations
	mockClient.On("GetRegion").Return("us-west-2")
	mockClient.On("ValidateConfiguration", mock.Anything).Return(true, nil)
	mockClient.On("ListEndpoints", mock.Anything).Return([]sagemaker.ResourceInfo{}, nil)
	mockClient.On("ListNotebooks", mock.Anything).Return([]sagemaker.ResourceInfo{}, nil)
	mockClient.On("ListStudioApps", mock.Anything).Return([]sagemaker.ResourceInfo{}, nil)

	err := mockExecute(t, []string{}, mockClient)
	assert.NoError(t, err)

	// Assert that all mock expectations were met
	mockClient.AssertExpectations(t)
}

func TestExecuteWithFlags_Unit(t *testing.T) {
	mockClient := new(MockSageMakerClient)

	// Setup mock expectations
	mockClient.On("GetRegion").Return("us-west-2")
	mockClient.On("ValidateConfiguration", mock.Anything).Return(true, nil)
	mockClient.On("ListEndpoints", mock.Anything).Return([]sagemaker.ResourceInfo{}, nil)
	mockClient.On("ListNotebooks", mock.Anything).Return([]sagemaker.ResourceInfo{}, nil)
	mockClient.On("ListStudioApps", mock.Anything).Return([]sagemaker.ResourceInfo{}, nil)

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
			err := mockExecute(t, tt.args, mockClient)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// Assert that all mock expectations were met
	mockClient.AssertExpectations(t)
}

// func TestExecuteWithInvalidFlags_Unit(t *testing.T) {
// 	mockClient := new(MockSageMakerClient)

// 	// Setup mock expectations
// 	mockClient.On("GetRegion").Return("us-west-2")
// 	mockClient.On("ValidateConfiguration", mock.Anything).Return(true, nil)
// 	mockClient.On("ListEndpoints", mock.Anything).Return([]sagemaker.ResourceInfo{}, nil)
// 	mockClient.On("ListNotebooks", mock.Anything).Return([]sagemaker.ResourceInfo{}, nil)
// 	mockClient.On("ListStudioApps", mock.Anything).Return([]sagemaker.ResourceInfo{}, nil)

// 	tests := []struct {
// 		name    string
// 		args    []string
// 		wantErr bool
// 	}{
// 		{
// 			name:    "with unknown flag",
// 			args:    []string{"--unknown"},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := mockExecute(t, tt.args, mockClient)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}

// 	// Assert that all mock expectations were met
// 	mockClient.AssertExpectations(t)
// }
