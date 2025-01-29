package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Custom retryable error type
type testRetryableError struct {
	retryable bool
	message   string
}

func (e *testRetryableError) Error() string {
	return e.message
}

func (e *testRetryableError) IsRetryable() bool {
	return e.retryable
}

func TestRetrier_SuccessOnFirstAttempt(t *testing.T) {
	ctx := context.Background()
	retrier := NewRetrier(DefaultConfig)

	attempts := 0
	err := retrier.Do(ctx, func() error {
		attempts++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, attempts)
}

func TestRetrier_SuccessOnSubsequentAttempt(t *testing.T) {
	ctx := context.Background()
	retrier := NewRetrier(DefaultConfig)

	attempts := 0
	err := retrier.Do(ctx, func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
}

func TestRetrier_MaxAttemptsExceeded(t *testing.T) {
	ctx := context.Background()
	config := Config{
		MaxAttempts:         2,
		InitialInterval:     10 * time.Millisecond,
		MaxInterval:         100 * time.Millisecond,
		Multiplier:          2.0,
		RandomizationFactor: 0.1,
	}
	retrier := NewRetrier(config)

	attempts := 0
	err := retrier.Do(ctx, func() error {
		attempts++
		return errors.New("persistent error")
	})

	assert.Error(t, err)
	assert.Equal(t, 3, attempts) // Note: attempts is MaxAttempts + 1
}

func TestRetrier_NonRetryableError(t *testing.T) {
	ctx := context.Background()
	retrier := NewRetrier(DefaultConfig)

	attempts := 0
	nonRetryableErr := &testRetryableError{retryable: false, message: "non-retryable error"}
	err := retrier.Do(ctx, func() error {
		attempts++
		return nonRetryableErr
	})

	assert.Error(t, err)
	assert.Equal(t, 1, attempts)
	assert.Equal(t, nonRetryableErr, err)
}

func TestRetrier_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	retrier := NewRetrier(DefaultConfig)

	var attempts int32 = 0
	errChan := make(chan error, 1)

	go func() {
		errChan <- retrier.Do(ctx, func() error {
			atomic.AddInt32(&attempts, 1)
			time.Sleep(100 * time.Millisecond)
			return errors.New("temporary error")
		})
	}()

	// Cancel context quickly
	time.Sleep(10 * time.Millisecond)
	cancel()

	err := <-errChan
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.LessOrEqual(t, atomic.LoadInt32(&attempts), int32(2))
}

func TestRetrier_BackoffTiming(t *testing.T) {
	ctx := context.Background()
	config := Config{
		MaxAttempts:         3,
		InitialInterval:     50 * time.Millisecond,
		MaxInterval:         500 * time.Millisecond,
		Multiplier:          2.0,
		RandomizationFactor: 0.1,
	}
	retrier := NewRetrier(config)

	var attempts []time.Time
	err := retrier.Do(ctx, func() error {
		attempts = append(attempts, time.Now())
		if len(attempts) < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Len(t, attempts, 3)

	// Check intervals between attempts
	for i := 1; i < len(attempts); i++ {
		interval := attempts[i].Sub(attempts[i-1])
		
		// First interval should be close to initial interval (with jitter)
		if i == 1 {
			assert.InDelta(t, config.InitialInterval.Seconds(), interval.Seconds(), 0.03, 
				"First retry interval should be close to initial interval")
		} else {
			// Subsequent intervals should increase exponentially
			expectedMinInterval := time.Duration(float64(config.InitialInterval) * math.Pow(config.Multiplier, float64(i-1)))
			expectedMinIntervalMs := expectedMinInterval.Milliseconds()
			
			// Add a small buffer to account for potential slight variations
			assert.GreaterOrEqual(t, interval.Milliseconds(), expectedMinIntervalMs, 
				fmt.Sprintf("Retry interval should increase exponentially. Expected at least %d ms, got %d ms", 
					expectedMinIntervalMs, interval.Milliseconds()))
		}

		// Ensure no interval exceeds max interval
		assert.LessOrEqual(t, interval.Milliseconds(), config.MaxInterval.Milliseconds(), 
			"Retry interval should not exceed max interval")
	}
}

func TestRetrier_Jitter(t *testing.T) {
	ctx := context.Background()
	config := Config{
		MaxAttempts:         5,
		InitialInterval:     50 * time.Millisecond,
		MaxInterval:         500 * time.Millisecond,
		Multiplier:          2.0,
		RandomizationFactor: 0.5, // High jitter for more noticeable variation
	}
	retrier := NewRetrier(config)

	var attempts []time.Time
	err := retrier.Do(ctx, func() error {
		attempts = append(attempts, time.Now())
		if len(attempts) < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Len(t, attempts, 3)

	// Check jitter variation
	for i := 1; i < len(attempts); i++ {
		interval := attempts[i].Sub(attempts[i-1])
		
		// Calculate expected base interval
		baseInterval := time.Duration(float64(config.InitialInterval) * math.Pow(config.Multiplier, float64(i-1)))
		
		// Jitter should create variation around the base interval
		assert.InDelta(t, baseInterval.Seconds(), interval.Seconds(), baseInterval.Seconds()*config.RandomizationFactor, 
			"Retry interval should have jitter")
	}
}
