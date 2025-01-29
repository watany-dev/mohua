package pricing

// Default pricing constants for different instance types
var (
	defaultEndpointPrices = map[string]float64{
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
	}

	defaultNotebookPrices = map[string]float64{
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
	}

	defaultStudioPrices = map[string]float64{
		"ml.t3.medium":   0.05,
		"ml.m5.large":    0.13,
		"ml.m5.xlarge":   0.27,
		"ml.m5.2xlarge":  0.54,
		"ml.c5.large":    0.12,
		"ml.c5.xlarge":   0.24,
		"ml.c5.2xlarge":  0.48,
		"ml.g4dn.xlarge": 0.736,
		"ml.p3.2xlarge":  3.825,
	}

	defaultCanvasPrices = map[string]float64{
		"ml.t3.medium":   0.05,
		"ml.m5.large":    0.13,
		"ml.m5.xlarge":   0.27,
		"ml.m5.2xlarge":  0.54,
		"ml.c5.large":    0.12,
		"ml.c5.xlarge":   0.24,
		"ml.c5.2xlarge":  0.48,
		"ml.g4dn.xlarge": 0.736,
		"ml.p3.2xlarge":  3.825,
	}

	defaultInstancePrices = map[string]float64{
		"ml.t3.medium": 0.0464,
		"ml.t3.large":  0.0736,
	}

	defaultStoragePrice = 0.10
)
