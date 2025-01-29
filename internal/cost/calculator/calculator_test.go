package calculator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockPricing implements the pricing interface for testing
type mockPricing struct {
	endpointPrices map[string]float64
	notebookPrices map[string]float64
	studioPrices   map[string]float64
	canvasPrices   map[string]float64
	instancePrices map[string]float64
	storagePrice   float64
}

func (m *mockPricing) GetEndpointPrice(instanceType string) float64 {
	if price, ok := m.endpointPrices[instanceType]; ok {
		return price
	}
	return 0
}

func (m *mockPricing) GetNotebookPrice(instanceType string) float64 {
	if price, ok := m.notebookPrices[instanceType]; ok {
		return price
	}
	return 0
}

func (m *mockPricing) GetStudioPrice(instanceType string) float64 {
	if price, ok := m.studioPrices[instanceType]; ok {
		return price
	}
	return 0
}

func (m *mockPricing) GetCanvasPrice(instanceType string) float64 {
	if price, ok := m.canvasPrices[instanceType]; ok {
		return price
	}
	return 0
}

func (m *mockPricing) GetStoragePrice() float64 {
	return m.storagePrice
}

func (m *mockPricing) GetInstancePrice(instanceType string) (float64, bool) {
	price, exists := m.instancePrices[instanceType]
	return price, exists
}

func newMockPricing() *mockPricing {
	return &mockPricing{
		endpointPrices: map[string]float64{
			"ml.t3.medium": 0.05,
			"ml.t3.large":  0.10,
		},
		notebookPrices: map[string]float64{
			"ml.t3.medium": 0.05,
			"ml.t3.large":  0.10,
		},
		studioPrices: map[string]float64{
			"ml.t3.medium": 0.05,
			"ml.t3.large":  0.10,
		},
		canvasPrices: map[string]float64{
			"ml.t3.medium": 0.05,
			"ml.t3.large":  0.10,
		},
		instancePrices: map[string]float64{
			"ml.t3.medium": 0.0464,
			"ml.t3.large":  0.0736,
		},
		storagePrice: 0.10,
	}
}

func TestCalculateCost(t *testing.T) {
	testCases := []struct {
		name         string
		instanceType string
		hours        float64
		expectedCost float64
	}{
		{
			name:         "ml.t3.medium standard hours",
			instanceType: "ml.t3.medium",
			hours:        1.0,
			expectedCost: 0.0464,
		},
		{
			name:         "ml.t3.large longer duration",
			instanceType: "ml.t3.large",
			hours:        2.5,
			expectedCost: 0.184,
		},
		{
			name:         "unsupported instance type",
			instanceType: "ml.unsupported.type",
			hours:        1.0,
			expectedCost: 0.0,
		},
	}

	calculator := NewCalculator(newMockPricing())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cost := calculator.CalculateCost(tc.instanceType, tc.hours)
			assert.InDelta(t, tc.expectedCost, cost, 0.0001, "Cost calculation should be accurate")
		})
	}
}

func TestResourceCalculations(t *testing.T) {
	calculator := NewCalculator(newMockPricing())
	startTime := time.Now().Add(-1 * time.Hour) // 1 hour ago

	t.Run("endpoint cost calculation", func(t *testing.T) {
		cost := calculator.CalculateEndpointCost("test-endpoint", "ml.t3.medium", 2, startTime)
		assert.Equal(t, "Endpoint", cost.ResourceType)
		assert.Equal(t, "test-endpoint", cost.ResourceName)
		assert.Equal(t, "ml.t3.medium", cost.InstanceType)
		assert.InDelta(t, 0.10, cost.HourlyCost, 0.0001) // 0.05 * 2 instances
	})

	t.Run("notebook cost calculation", func(t *testing.T) {
		cost := calculator.CalculateNotebookCost("test-notebook", "ml.t3.medium", startTime, 100)
		assert.Equal(t, "Notebook Instance", cost.ResourceType)
		assert.Equal(t, float64(100), cost.StorageSizeGB)
		assert.InDelta(t, 10.0, cost.StorageCost, 0.0001) // 100GB * $0.10/GB
	})

	t.Run("studio cost calculation", func(t *testing.T) {
		cost := calculator.CalculateStudioCost("test-studio", "ml.t3.medium", startTime)
		assert.Equal(t, "Studio", cost.ResourceType)
		assert.InDelta(t, 0.05, cost.HourlyCost, 0.0001)
	})

	t.Run("canvas cost calculation", func(t *testing.T) {
		cost := calculator.CalculateCanvasCost("test-canvas", "ml.t3.medium", startTime)
		assert.Equal(t, "Canvas", cost.ResourceType)
		assert.InDelta(t, 0.05, cost.HourlyCost, 0.0001)
	})
}

func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "hours and minutes",
			duration: 2*time.Hour + 30*time.Minute,
			expected: "2h30m",
		},
		{
			name:     "only hours",
			duration: 5 * time.Hour,
			expected: "5h0m",
		},
		{
			name:     "only minutes",
			duration: 45 * time.Minute,
			expected: "0h45m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatDuration(tc.duration)
			assert.Equal(t, tc.expected, result)
		})
	}
}
