package cost

import (
	"time"
)

// ResourceCost represents the cost information for a SageMaker resource
type ResourceCost struct {
	ResourceType     string
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
	pricing *PricingData
}

// NewCalculator creates a new cost calculator with the given pricing data
func NewCalculator(pricing *PricingData) *Calculator {
	return &Calculator{
		pricing: pricing,
	}
}

// CalculateEndpointCost calculates costs for a SageMaker endpoint
func (c *Calculator) CalculateEndpointCost(name, instanceType string, count int, startTime time.Time) *ResourceCost {
	runningTime := time.Since(startTime)
	hourlyRate := c.pricing.GetEndpointPrice(instanceType) * float64(count)
	
	return &ResourceCost{
		ResourceType:    "Endpoint",
		ResourceName:    name,
		InstanceType:    instanceType,
		RunningTime:     runningTime,
		HourlyCost:     hourlyRate,
		CurrentCost:    calculateCurrentCost(hourlyRate, runningTime),
		ProjectedCost:  calculateProjectedMonthlyCost(hourlyRate),
	}
}

// CalculateNotebookCost calculates costs for a SageMaker notebook instance
func (c *Calculator) CalculateNotebookCost(name, instanceType string, startTime time.Time, volumeSizeGB int) *ResourceCost {
	runningTime := time.Since(startTime)
	hourlyRate := c.pricing.GetNotebookPrice(instanceType)
	storageCost := calculateStorageCost(float64(volumeSizeGB), c.pricing.GetStoragePrice())
	
	return &ResourceCost{
		ResourceType:    "Notebook Instance",
		ResourceName:    name,
		InstanceType:    instanceType,
		RunningTime:     runningTime,
		HourlyCost:     hourlyRate,
		CurrentCost:    calculateCurrentCost(hourlyRate, runningTime),
		ProjectedCost:  calculateProjectedMonthlyCost(hourlyRate),
		StorageSizeGB:  float64(volumeSizeGB),
		StorageCost:    storageCost,
	}
}

// CalculateStudioCost calculates costs for a SageMaker Studio instance
func (c *Calculator) CalculateStudioCost(name, instanceType string, startTime time.Time) *ResourceCost {
	runningTime := time.Since(startTime)
	hourlyRate := c.pricing.GetStudioPrice(instanceType)
	
	return &ResourceCost{
		ResourceType:    "Studio",
		ResourceName:    name,
		InstanceType:    instanceType,
		RunningTime:     runningTime,
		HourlyCost:     hourlyRate,
		CurrentCost:    calculateCurrentCost(hourlyRate, runningTime),
		ProjectedCost:  calculateProjectedMonthlyCost(hourlyRate),
	}
}

// Helper functions

func calculateCurrentCost(hourlyRate float64, duration time.Duration) float64 {
	hours := duration.Hours()
	return hourlyRate * hours
}

func calculateProjectedMonthlyCost(hourlyRate float64) float64 {
	// Assuming 730 hours in a month (365 * 24 / 12)
	return hourlyRate * 730
}

func calculateStorageCost(sizeGB, pricePerGBMonth float64) float64 {
	return sizeGB * pricePerGBMonth
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return time.Duration(h).String() + time.Duration(m).String()
}
