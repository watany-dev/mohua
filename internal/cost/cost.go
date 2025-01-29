package cost

import (
	"mohua/internal/cost/calculator"
	"mohua/internal/cost/pricing"
)

// LoadPricing delegates to the pricing package's LoadPricing function
func LoadPricing(filepath string) (*pricing.PricingData, error) {
	return pricing.LoadPricing(filepath)
}

// NewCalculator creates a new cost calculator
func NewCalculator(pricingData *pricing.PricingData) *calculator.Calculator {
	return calculator.NewCalculator(pricingData)
}
