package cmd

import (
	"os"
	"testing"
)

// hasAWSCredentials checks if AWS credentials are available in the environment
func hasAWSCredentials() bool {
	// Check for explicit credentials
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		return true
	}
	
	// Check for AWS_PROFILE
	if os.Getenv("AWS_PROFILE") != "" {
		return true
	}
	
	// Check for AWS_ROLE_ARN
	if os.Getenv("AWS_ROLE_ARN") != "" {
		return true
	}
	
	return false
}

// skipIfNoAWSCredentials skips the test if no AWS credentials are available
func skipIfNoAWSCredentials(t *testing.T) {
	if !hasAWSCredentials() {
		t.Skip("Skipping test that requires AWS credentials")
	}
}
