package pricing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPricing(t *testing.T) {
	pricing, err := LoadPricing("")
	assert.NoError(t, err)
	assert.NotNil(t, pricing)

	// Test default values are loaded
	assert.Equal(t, 0.05, pricing.Endpoints["ml.t3.medium"])
	assert.Equal(t, 0.10, pricing.Storage.EBS)
}

func TestGetPrice(t *testing.T) {
	pricing, _ := LoadPricing("")

	testCases := []struct {
		name         string
		resourceType ResourceType
		instanceType string
		expectError  bool
		expectedPrice float64
	}{
		{
			name:         "valid endpoint instance",
			resourceType: ResourceTypeEndpoint,
			instanceType: "ml.t3.medium",
			expectError:  false,
			expectedPrice: 0.05,
		},
		{
			name:         "valid notebook instance",
			resourceType: ResourceTypeNotebook,
			instanceType: "ml.t3.medium",
			expectError:  false,
			expectedPrice: 0.05,
		},
		{
			name:         "invalid resource type",
			resourceType: "invalid",
			instanceType: "ml.t3.medium",
			expectError:  true,
			expectedPrice: 0,
		},
		{
			name:         "invalid instance type",
			resourceType: ResourceTypeEndpoint,
			instanceType: "invalid",
			expectError:  true,
			expectedPrice: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			price, err := pricing.GetPrice(tc.resourceType, tc.instanceType)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPrice, price)
			}
		})
	}
}

func TestResourceSpecificPricing(t *testing.T) {
	pricing, _ := LoadPricing("")

	t.Run("endpoint pricing", func(t *testing.T) {
		price := pricing.GetEndpointPrice("ml.t3.medium")
		assert.Equal(t, 0.05, price)
	})

	t.Run("notebook pricing", func(t *testing.T) {
		price := pricing.GetNotebookPrice("ml.t3.medium")
		assert.Equal(t, 0.05, price)
	})

	t.Run("studio pricing", func(t *testing.T) {
		price := pricing.GetStudioPrice("ml.t3.medium")
		assert.Equal(t, 0.05, price)
	})

	t.Run("canvas pricing", func(t *testing.T) {
		price := pricing.GetCanvasPrice("ml.t3.medium")
		assert.Equal(t, 0.05, price)
	})

	t.Run("storage pricing", func(t *testing.T) {
		price := pricing.GetStoragePrice()
		assert.Equal(t, 0.10, price)
	})

	t.Run("instance pricing", func(t *testing.T) {
		price, exists := pricing.GetInstancePrice("ml.t3.medium")
		assert.True(t, exists)
		assert.Equal(t, 0.0464, price)
	})
}
