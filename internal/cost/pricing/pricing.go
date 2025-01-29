package pricing

import "fmt"

// LoadPricing returns the pricing data with default values
func LoadPricing(filepath string) (*PricingData, error) {
	return &PricingData{
		Endpoints:      defaultEndpointPrices,
		Notebooks:     defaultNotebookPrices,
		Studio:        defaultStudioPrices,
		Canvas:        defaultCanvasPrices,
		Storage: struct {
			EBS float64 `yaml:"ebs"`
		}{
			EBS: defaultStoragePrice,
		},
		InstancePrices: defaultInstancePrices,
	}, nil
}

// GetPrice returns the price for a given resource type and instance type
func (p *PricingData) GetPrice(resourceType ResourceType, instanceType string) (float64, error) {
	var priceMap map[string]float64

	switch resourceType {
	case ResourceTypeEndpoint:
		priceMap = p.Endpoints
	case ResourceTypeNotebook:
		priceMap = p.Notebooks
	case ResourceTypeStudio:
		priceMap = p.Studio
	case ResourceTypeCanvas:
		priceMap = p.Canvas
	default:
		return 0, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	if price, ok := priceMap[instanceType]; ok {
		return price, nil
	}
	return 0, fmt.Errorf("unsupported instance type: %s for resource type: %s", instanceType, resourceType)
}

// GetEndpointPrice returns the hourly price for an endpoint instance type
func (p *PricingData) GetEndpointPrice(instanceType string) float64 {
	price, _ := p.GetPrice(ResourceTypeEndpoint, instanceType)
	return price
}

// GetNotebookPrice returns the hourly price for a notebook instance type
func (p *PricingData) GetNotebookPrice(instanceType string) float64 {
	price, _ := p.GetPrice(ResourceTypeNotebook, instanceType)
	return price
}

// GetStudioPrice returns the hourly price for a studio instance type
func (p *PricingData) GetStudioPrice(instanceType string) float64 {
	price, _ := p.GetPrice(ResourceTypeStudio, instanceType)
	return price
}

// GetCanvasPrice returns the hourly price for a Canvas instance type
func (p *PricingData) GetCanvasPrice(instanceType string) float64 {
	price, _ := p.GetPrice(ResourceTypeCanvas, instanceType)
	return price
}

// GetStoragePrice returns the monthly price per GB for EBS storage
func (p *PricingData) GetStoragePrice() float64 {
	return p.Storage.EBS
}

// GetInstancePrice returns the price for a given instance type from the general instance prices
func (p *PricingData) GetInstancePrice(instanceType string) (float64, bool) {
	price, exists := p.InstancePrices[instanceType]
	return price, exists
}
