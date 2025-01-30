package sagemaker

import (
	"errors"
	"fmt"
	"net"
	"syscall"
	"testing"

	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/assert"
)

func TestRetryableError(t *testing.T) {
	baseErr := errors.New("base error")
	retryableErr := &RetryableError{Err: baseErr}

	// Test Error() method
	assert.Equal(t, baseErr.Error(), retryableErr.Error())

	// Test IsRetryable() method
	assert.True(t, retryableErr.IsRetryable())
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name           string
		err           error
		expectedType  error
		isRetryable   bool
	}{
		{
			name: "network error - connection reset",
			err:  &net.OpError{Err: syscall.ECONNRESET},
			expectedType: &RetryableError{},
			isRetryable: true,
		},
		{
			name: "network error - connection refused",
			err:  &net.OpError{Err: syscall.ECONNREFUSED},
			expectedType: &RetryableError{},
			isRetryable: true,
		},
		{
			name: "throttling error",
			err:  &smithy.GenericAPIError{Code: "ThrottlingException"},
			expectedType: &RetryableError{},
			isRetryable: true,
		},
		{
			name: "non-retryable error",
			err:  errors.New("random error"),
			expectedType: &NonRetryableError{},
			isRetryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedErr := WrapError(tt.err)
			
			if tt.isRetryable {
				var retryable interface{ IsRetryable() bool }
				assert.True(t, errors.As(wrappedErr, &retryable))
				assert.True(t, retryable.IsRetryable())
				assert.IsType(t, tt.expectedType, wrappedErr)
			} else {
				var nonRetryable interface{ IsRetryable() bool }
				assert.True(t, errors.As(wrappedErr, &nonRetryable))
				assert.False(t, nonRetryable.IsRetryable())
				assert.IsType(t, tt.expectedType, wrappedErr)
			}
		})
	}
}

func TestIsNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "connection reset",
			err:      &net.OpError{Err: syscall.ECONNRESET},
			expected: true,
		},
		{
			name:     "connection refused",
			err:      &net.OpError{Err: syscall.ECONNREFUSED},
			expected: true,
		},
		{
			name:     "timeout error",
			err:      &net.OpError{Err: syscall.ETIMEDOUT},
			expected: true,
		},
		{
			name:     "non-network error",
			err:      fmt.Errorf("random error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isNetworkError(tt.err))
		})
	}
}
