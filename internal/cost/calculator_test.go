package cost

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateCost(t *testing.T) {
	testCases := []struct {
		name           string
		instanceType   string
		hours          float64
		expectedCost   float64
	}{
		{
			name:         "ml.t3.medium standard hours",
			instanceType: "ml.t3.medium",
			hours:        1.0,
			expectedCost: 0.0464, // Approximate hourly rate for ml.t3.medium
		},
		{
			name:         "ml.t3.large longer duration",
			instanceType: "ml.t3.large",
			hours:        2.5,
			expectedCost: 0.184, // Adjusted to match rounding
		},
		{
			name:         "Unsupported instance type",
			instanceType: "ml.unsupported.type",
			hours:        1.0,
			expectedCost: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calculator := NewCalculator()
			cost := calculator.CalculateCost(tc.instanceType, tc.hours)
			assert.InDelta(t, tc.expectedCost, cost, 0.0001, "Cost calculation should be accurate")
		})
	}
}

func TestCalculator_GetInstancePricing(t *testing.T) {
	calculator := NewCalculator()

	testCases := []struct {
		instanceType   string
		expectPricing  bool
	}{
		{"ml.t3.medium", true},
		{"ml.t3.large", true},
		{"ml.invalid.type", false},
	}

	for _, tc := range testCases {
		t.Run(tc.instanceType, func(t *testing.T) {
			pricing, exists := calculator.GetInstancePricing(tc.instanceType)
			assert.Equal(t, tc.expectPricing, exists, "Pricing existence check")
			
			if exists {
				assert.Greater(t, pricing, 0.0, "Pricing should be positive")
			}
		})
	}
}
