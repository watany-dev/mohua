package cost

// PricingData represents the pricing data structure
type PricingData struct {
	Endpoints map[string]float64 `yaml:"endpoints"`
	Notebooks map[string]float64 `yaml:"notebooks"`
	Studio    map[string]float64 `yaml:"studio"`
	Canvas    map[string]float64 `yaml:"canvas"`
	Storage   struct {
		EBS float64 `yaml:"ebs"`
	} `yaml:"storage"`
	InstancePrices map[string]float64 `yaml:"instance_prices"`
}

// LoadPricing returns the hardcoded pricing data
func LoadPricing(filepath string) (*PricingData, error) {
	return &PricingData{
		Endpoints: map[string]float64{
			"ml.t2.medium":   0.05,
			"ml.t2.large":    0.10,
			"ml.t2.xlarge":   0.20,
			"ml.t3.medium":   0.05,
			"ml.t3.large":    0.10,
			"ml.t3.xlarge":   0.20,
			"ml.m4.xlarge":   0.28,
			"ml.m5.large":    0.13,
			"ml.m5.xlarge":   0.27,
			"ml.m5.2xlarge":  0.54,
			"ml.c5.large":    0.12,
			"ml.c5.xlarge":   0.24,
			"ml.c5.2xlarge":  0.48,
			"ml.p3.2xlarge":  3.825,
			"ml.g4dn.xlarge": 0.736,
		},
		Notebooks: map[string]float64{
			"ml.t2.medium":   0.05,
			"ml.t2.large":    0.10,
			"ml.t2.xlarge":   0.20,
			"ml.t3.medium":   0.05,
			"ml.t3.large":    0.10,
			"ml.t3.xlarge":   0.20,
			"ml.m4.xlarge":   0.28,
			"ml.m5.large":    0.13,
			"ml.m5.xlarge":   0.27,
			"ml.m5.2xlarge":  0.54,
			"ml.c5.large":    0.12,
			"ml.c5.xlarge":   0.24,
			"ml.c5.2xlarge":  0.48,
			"ml.p3.2xlarge":  3.825,
			"ml.g4dn.xlarge": 0.736,
		},
		Studio: map[string]float64{
			"ml.t3.medium":   0.05,
			"ml.m5.large":    0.13,
			"ml.m5.xlarge":   0.27,
			"ml.m5.2xlarge":  0.54,
			"ml.c5.large":    0.12,
			"ml.c5.xlarge":   0.24,
			"ml.c5.2xlarge":  0.48,
			"ml.g4dn.xlarge": 0.736,
			"ml.p3.2xlarge":  3.825,
		},
		Canvas: map[string]float64{
			"ml.t3.medium":   0.05,
			"ml.m5.large":    0.13,
			"ml.m5.xlarge":   0.27,
			"ml.m5.2xlarge":  0.54,
			"ml.c5.large":    0.12,
			"ml.c5.xlarge":   0.24,
			"ml.c5.2xlarge":  0.48,
			"ml.g4dn.xlarge": 0.736,
			"ml.p3.2xlarge":  3.825,
		},
		Storage: struct {
			EBS float64 `yaml:"ebs"`
		}{
			EBS: 0.10,
		},
		InstancePrices: map[string]float64{
			"ml.t3.medium": 0.0464,
			"ml.t3.large":  0.0736,
		},
	}, nil
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

// GetCanvasPrice returns the hourly price for a Canvas instance type
func (p *PricingData) GetCanvasPrice(instanceType string) float64 {
	if price, ok := p.Canvas[instanceType]; ok {
		return price
	}
	return 0
}

// GetStoragePrice returns the monthly price per GB for EBS storage
func (p *PricingData) GetStoragePrice() float64 {
	return p.Storage.EBS
}
