package sagemaker

import (
	"errors"
	"strings"

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
	// Add common network error types or error message patterns
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
