package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"mohua/internal/cost"
	"mohua/internal/display"
	"mohua/internal/sagemaker"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mohua",
	Short: "Monitor AWS SageMaker compute resources and their costs",
	Long: `A monitoring tool for AWS SageMaker that helps track running compute resources
and their associated costs. It provides information about:

- Endpoints
- Notebook Instances
- Studio Applications
- Running Jobs

The tool shows current status, running time, and cost projections for each resource.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMonitor()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region (required)")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	rootCmd.MarkPersistentFlagRequired("region")
}

func runMonitor() error {
	// Create SageMaker client
	client, err := sagemaker.NewClient(region)
	if err != nil {
		return fmt.Errorf("failed to create SageMaker client: %w", err)
	}

	// Load pricing data
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	pricingPath := filepath.Join(filepath.Dir(execPath), "configs", "pricing.yaml")
	pricing, err := cost.LoadPricing(pricingPath)
	if err != nil {
		return fmt.Errorf("failed to load pricing data: %w", err)
	}

	// Create cost calculator
	calculator := cost.NewCalculator(pricing)

	// Create printer
	printer := display.NewPrinter(jsonOutput)

	ctx := context.Background()

	// Get endpoints
	endpoints, err := client.ListEndpoints(ctx)
	if err != nil {
		fmt.Printf("Warning: Failed to list endpoints: %v\n", err)
	}
	for _, endpoint := range endpoints {
		cost := calculator.CalculateEndpointCost(
			endpoint.Name,
			endpoint.InstanceType,
			endpoint.InstanceCount,
			endpoint.CreationTime,
		)
		printer.AddResource(display.ResourceInfo{
			ResourceType:  "Endpoint",
			Name:         cost.ResourceName,
			Status:       endpoint.Status,
			InstanceType: cost.InstanceType,
			RunningTime:  cost.RunningTime.String(),
			HourlyCost:   cost.HourlyCost,
			CurrentCost:  cost.CurrentCost,
			ProjectedCost: cost.ProjectedCost,
			TotalCost:    cost.CurrentCost,
		})
	}

	// Get notebooks
	notebooks, err := client.ListNotebooks(ctx)
	if err != nil {
		fmt.Printf("Warning: Failed to list notebooks: %v\n", err)
	}
	for _, notebook := range notebooks {
		cost := calculator.CalculateNotebookCost(
			notebook.Name,
			notebook.InstanceType,
			notebook.CreationTime,
			notebook.VolumeSize,
		)
		printer.AddResource(display.ResourceInfo{
			ResourceType:  "Notebook",
			Name:         cost.ResourceName,
			Status:       notebook.Status,
			InstanceType: cost.InstanceType,
			RunningTime:  cost.RunningTime.String(),
			HourlyCost:   cost.HourlyCost,
			CurrentCost:  cost.CurrentCost,
			ProjectedCost: cost.ProjectedCost,
			StorageCost:  cost.StorageCost,
			TotalCost:    cost.CurrentCost + cost.StorageCost,
		})
	}

	// Get Studio apps
	apps, err := client.ListStudioApps(ctx)
	if err != nil {
		fmt.Printf("Warning: Failed to list Studio apps: %v\n", err)
	}
	for _, app := range apps {
		cost := calculator.CalculateStudioCost(
			fmt.Sprintf("%s/%s", app.UserProfile, app.AppType),
			app.InstanceType,
			app.CreationTime,
		)
		printer.AddResource(display.ResourceInfo{
			ResourceType:  "Studio",
			Name:         cost.ResourceName,
			Status:       app.Status,
			InstanceType: cost.InstanceType,
			RunningTime:  cost.RunningTime.String(),
			HourlyCost:   cost.HourlyCost,
			CurrentCost:  cost.CurrentCost,
			ProjectedCost: cost.ProjectedCost,
			TotalCost:    cost.CurrentCost,
		})
	}

	// Print results
	printer.Print()
	return nil
}
