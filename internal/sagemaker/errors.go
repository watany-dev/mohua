package sagemaker

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/smithy-go"
)

// RetryableError represents an error that can be retried
type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

func (e *RetryableError) IsRetryable() bool {
	return true
}

// NonRetryableError represents an error that should not be retried
type NonRetryableError struct {
	Err error
}

func (e *NonRetryableError) Error() string {
	return e.Err.Error()
}

func (e *NonRetryableError) IsRetryable() bool {
	return false
}

// ErrorTracker manages error suppression and tracking
type ErrorTracker struct {
	mu            sync.Mutex
	lastError     error
	suppressCount int
	maxSuppress   int
}

// NewErrorTracker creates a new ErrorTracker with default max suppress count
func NewErrorTracker() *ErrorTracker {
	return &ErrorTracker{
		maxSuppress: 1, // Suppress after first duplicate
	}
}

// Track manages error suppression logic
func (et *ErrorTracker) Track(err error) error {
	et.mu.Lock()
	defer et.mu.Unlock()

	if err == nil {
		return nil
	}

	// Extract the core error message
	var errMsg string
	if authErrorMessage := extractAuthenticationErrorMessage(err); authErrorMessage != "" {
		errMsg = fmt.Sprintf("AWS Authentication Error: %s", authErrorMessage)
	} else {
		errMsg = err.Error()
	}

	// Compare with last error and manage suppression
	if et.lastError != nil && et.lastError.Error() == errMsg {
		et.suppressCount++
		if et.suppressCount > et.maxSuppress {
			return nil
		}
	} else {
		et.suppressCount = 0
		et.lastError = fmt.Errorf(errMsg)
	}

	return et.lastError
}

// extractAuthenticationErrorMessage provides user-friendly authentication error messages
func extractAuthenticationErrorMessage(err error) string {
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "no EC2 IMDS role found"):
		return "No AWS role configured. Set credentials using AWS CLI or environment variables."
	case strings.Contains(errStr, "failed to refresh cached credentials"):
		return "Credential refresh failed. Verify AWS configuration and permissions."
	case strings.Contains(errStr, "failed to get API token"):
		return "Unable to obtain AWS API token. Check network and authentication settings."
	default:
		return ""
	}
}

// WrapError wraps AWS errors and determines if they are retryable
func WrapError(err error) error {
	if err == nil {
		return nil
	}

	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch ae.ErrorCode() {
		// Retryable errors
		case "RequestTimeout", 
			 "ThrottlingException", 
			 "ProvisionedThroughputExceededException", 
			 "TransactionInProgressException":
			return &RetryableError{Err: err}

		// Non-retryable errors
		case "ValidationError", 
			 "AccessDeniedException", 
			 "InvalidParameterException":
			return &NonRetryableError{Err: err}
		}
	}

	// Network errors or context cancellations are typically retryable
	if isNetworkError(err) {
		return &RetryableError{Err: err}
	}

	// Default to non-retryable for unknown errors
	return &NonRetryableError{Err: err}
}

// isNetworkError checks if the error is a network-related error
func isNetworkError(err error) bool {
	networkErrorMessages := []string{
		"connection refused",
		"connection reset",
		"network is unreachable",
		"timeout",
		"i/o timeout",
	}

	errStr := err.Error()
	for _, msg := range networkErrorMessages {
		if strings.Contains(strings.ToLower(errStr), msg) {
			return true
		}
	}

	return false
}
