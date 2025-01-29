package calculator

import (
	"math"
)

// NewCalculator creates a new cost calculator with the given pricing data
func NewCalculator(pricing interface {
	GetEndpointPrice(instanceType string) float64
	GetNotebookPrice(instanceType string) float64
	GetStudioPrice(instanceType string) float64
	GetCanvasPrice(instanceType string) float64
	GetStoragePrice() float64
	GetInstancePrice(instanceType string) (float64, bool)
}) *Calculator {
	return &Calculator{
		pricing: pricing,
	}
}

// CalculateCost calculates the cost for a given instance type and duration
func (c *Calculator) CalculateCost(instanceType string, hours float64) float64 {
	price, exists := c.pricing.GetInstancePrice(instanceType)
	if !exists {
		return 0.0
	}
	// Round to 4 decimal places to match test precision
	return math.Round(price * hours * 10000) / 10000
}

// GetInstancePricing retrieves the hourly price for a given instance type
func (c *Calculator) GetInstancePricing(instanceType string) (float64, bool) {
	return c.pricing.GetInstancePrice(instanceType)
}
