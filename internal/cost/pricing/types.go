package pricing

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

// ResourceType represents different SageMaker resource types
type ResourceType string

const (
	ResourceTypeEndpoint  ResourceType = "endpoint"
	ResourceTypeNotebook ResourceType = "notebook"
	ResourceTypeStudio   ResourceType = "studio"
	ResourceTypeCanvas   ResourceType = "canvas"
)
