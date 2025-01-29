package calculator

import "time"

// CalculateEndpointCost calculates costs for a SageMaker endpoint
func (c *Calculator) CalculateEndpointCost(name, instanceType string, count int, startTime time.Time) *ResourceCost {
	runningTime := time.Since(startTime)
	hourlyRate := c.pricing.GetEndpointPrice(instanceType) * float64(count)
	
	return &ResourceCost{
		ResourceType:   "Endpoint",
		ResourceName:   name,
		InstanceType:   instanceType,
		RunningTime:    runningTime,
		HourlyCost:    hourlyRate,
		CurrentCost:   calculateCurrentCost(hourlyRate, runningTime),
		ProjectedCost: calculateProjectedMonthlyCost(hourlyRate),
	}
}

// CalculateNotebookCost calculates costs for a SageMaker notebook instance
func (c *Calculator) CalculateNotebookCost(name, instanceType string, startTime time.Time, volumeSizeGB int) *ResourceCost {
	runningTime := time.Since(startTime)
	hourlyRate := c.pricing.GetNotebookPrice(instanceType)
	storageCost := calculateStorageCost(float64(volumeSizeGB), c.pricing.GetStoragePrice())
	
	return &ResourceCost{
		ResourceType:   "Notebook Instance",
		ResourceName:   name,
		InstanceType:   instanceType,
		RunningTime:    runningTime,
		HourlyCost:    hourlyRate,
		CurrentCost:   calculateCurrentCost(hourlyRate, runningTime),
		ProjectedCost: calculateProjectedMonthlyCost(hourlyRate),
		StorageSizeGB: float64(volumeSizeGB),
		StorageCost:   storageCost,
	}
}

// CalculateStudioCost calculates costs for a SageMaker Studio instance
func (c *Calculator) CalculateStudioCost(name, instanceType string, startTime time.Time) *ResourceCost {
	runningTime := time.Since(startTime)
	hourlyRate := c.pricing.GetStudioPrice(instanceType)
	
	return &ResourceCost{
		ResourceType:   "Studio",
		ResourceName:   name,
		InstanceType:   instanceType,
		RunningTime:    runningTime,
		HourlyCost:    hourlyRate,
		CurrentCost:   calculateCurrentCost(hourlyRate, runningTime),
		ProjectedCost: calculateProjectedMonthlyCost(hourlyRate),
	}
}

// CalculateCanvasCost calculates costs for a SageMaker Canvas application
func (c *Calculator) CalculateCanvasCost(name, instanceType string, startTime time.Time) *ResourceCost {
	runningTime := time.Since(startTime)
	hourlyRate := c.pricing.GetCanvasPrice(instanceType)
	
	return &ResourceCost{
		ResourceType:   "Canvas",
		ResourceName:   name,
		InstanceType:   instanceType,
		RunningTime:    runningTime,
		HourlyCost:    hourlyRate,
		CurrentCost:   calculateCurrentCost(hourlyRate, runningTime),
		ProjectedCost: calculateProjectedMonthlyCost(hourlyRate),
	}
}
