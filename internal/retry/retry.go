package retry

import (
	"context"
	"math/rand"
	"time"
)

// Config holds the retry configuration parameters
type Config struct {
	MaxAttempts         int           // Maximum number of retry attempts
	InitialInterval     time.Duration // Initial backoff interval
	MaxInterval         time.Duration // Maximum backoff interval
	Multiplier         float64       // Backoff multiplier
	RandomizationFactor float64       // Randomization factor for jitter
}

// DefaultConfig provides reasonable default values for retry configuration
var DefaultConfig = Config{
	MaxAttempts:         3,
	InitialInterval:     1 * time.Second,
	MaxInterval:         30 * time.Second,
	Multiplier:         2.0,
	RandomizationFactor: 0.1,
}

// Retrier handles the retry logic with exponential backoff
type Retrier struct {
	config Config
}

// NewRetrier creates a new Retrier with the given configuration
func NewRetrier(config Config) *Retrier {
	return &Retrier{config: config}
}

// Do executes the given operation with retry logic
func (r *Retrier) Do(ctx context.Context, operation func() error) error {
	var err error
	currentInterval := r.config.InitialInterval

	for attempt := 0; attempt <= r.config.MaxAttempts; attempt++ {
		// Execute the operation
		err = operation()
		if err == nil {
			return nil
		}

		// Check if error is retryable
		if retryable, ok := err.(interface{ IsRetryable() bool }); ok && !retryable.IsRetryable() {
			return err
		}

		// Check if we've exhausted all attempts
		if attempt == r.config.MaxAttempts {
			return err
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Calculate next backoff duration with jitter
			jitter := 1.0 + (rand.Float64()*2-1.0)*r.config.RandomizationFactor
			backoff := time.Duration(float64(currentInterval) * jitter)
			if backoff > r.config.MaxInterval {
				backoff = r.config.MaxInterval
			}

			// Wait for backoff duration
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}

			// Update interval for next iteration
			currentInterval = time.Duration(float64(currentInterval) * r.config.Multiplier)
		}
	}

	return err
}
