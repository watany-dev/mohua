package cost

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// PricingData represents the structure of the pricing.yaml file
type PricingData struct {
	Endpoints map[string]float64 `yaml:"endpoints"`
	Notebooks map[string]float64 `yaml:"notebooks"`
	Studio    map[string]float64 `yaml:"studio"`
	Storage   struct {
		EBS float64 `yaml:"ebs"`
	} `yaml:"storage"`
}

// LoadPricing loads pricing data from the YAML file
func LoadPricing(filepath string) (*PricingData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pricing file: %w", err)
	}

	var pricing PricingData
	if err := yaml.Unmarshal(data, &pricing); err != nil {
		return nil, fmt.Errorf("failed to parse pricing data: %w", err)
	}

	return &pricing, nil
}

// GetEndpointPrice returns the hourly price for an endpoint instance type
func (p *PricingData) GetEndpointPrice(instanceType string) float64 {
	if price, ok := p.Endpoints[instanceType]; ok {
		return price
	}
	return 0
}

// GetNotebookPrice returns the hourly price for a notebook instance type
func (p *PricingData) GetNotebookPrice(instanceType string) float64 {
	if price, ok := p.Notebooks[instanceType]; ok {
		return price
	}
	return 0
}

// GetStudioPrice returns the hourly price for a studio instance type
func (p *PricingData) GetStudioPrice(instanceType string) float64 {
	if price, ok := p.Studio[instanceType]; ok {
		return price
	}
	return 0
}

// GetStoragePrice returns the monthly price per GB for EBS storage
func (p *PricingData) GetStoragePrice() float64 {
	return p.Storage.EBS
}
