package calculator

import "time"

// ResourceCost represents the cost information for a SageMaker resource
type ResourceCost struct {
	ResourceType    string
	ResourceName    string
	InstanceType    string
	RunningTime     time.Duration
	HourlyCost     float64
	CurrentCost    float64
	ProjectedCost  float64
	StorageSizeGB  float64
	StorageCost    float64
}

// Calculator handles cost calculations for SageMaker resources
type Calculator struct {
	pricing interface {
		GetEndpointPrice(instanceType string) float64
		GetNotebookPrice(instanceType string) float64
		GetStudioPrice(instanceType string) float64
		GetCanvasPrice(instanceType string) float64
		GetStoragePrice() float64
		GetInstancePrice(instanceType string) (float64, bool)
	}
}
