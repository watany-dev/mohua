package calculator

import (
	"fmt"
	"time"
)

// calculateCurrentCost calculates the cost based on hourly rate and duration
func calculateCurrentCost(hourlyRate float64, duration time.Duration) float64 {
	hours := duration.Hours()
	return hourlyRate * hours
}

// calculateProjectedMonthlyCost calculates the projected monthly cost based on hourly rate
func calculateProjectedMonthlyCost(hourlyRate float64) float64 {
	// Assuming 730 hours in a month (365 * 24 / 12)
	return hourlyRate * 730
}

// calculateStorageCost calculates the storage cost based on size and price per GB
func calculateStorageCost(sizeGB, pricePerGBMonth float64) float64 {
	return sizeGB * pricePerGBMonth
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%dh%dm", h, m)
}
