package cmd

import (
	"context"
	"fmt"
	"time"
	"github.com/spf13/cobra"
	"mohua/internal/display"
	"mohua/internal/sagemaker"
)

var (
	region    string
	jsonOutput bool
)

// minimalRootCmd represents the base command when called without any subcommands
var minimalRootCmd = &cobra.Command{
	Use:   "mohua",
	Short: "Monitor AWS SageMaker compute resources and their costs",
	Long: `A monitoring tool for AWS SageMaker that helps track running compute resources
and their associated costs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMinimalMonitor()
	},
}

// ExecuteMinimal adds all child commands to the root command and sets flags appropriately.
func ExecuteMinimal() error {
	minimalRootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region (optional, defaults to AWS_REGION env var or us-east-1)")
	minimalRootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	
	return minimalRootCmd.Execute()
}

func runMinimalMonitor() error {
	// Create minimal SageMaker client
	client, err := sagemaker.NewMinimalClient(region)
	if err != nil {
		return fmt.Errorf("failed to create SageMaker client: %w", err)
	}

	// Create printer
	printer := display.NewPrinter(jsonOutput)

	ctx := context.Background()

	// Get endpoints
	endpoints, err := client.ListEndpoints(ctx)
	if err != nil {
		fmt.Printf("Warning: Failed to list endpoints: %v\n", err)
	}
	for _, endpoint := range endpoints {
		printer.AddResource(display.ResourceInfo{
			ResourceType:  "Endpoint",
			Name:         endpoint.Name,
			Status:       endpoint.Status,
			InstanceType: endpoint.InstanceType,
			RunningTime:  time.Since(endpoint.CreationTime).String(),
		})
	}

	// Get notebooks
	notebooks, err := client.ListNotebooks(ctx)
	if err != nil {
		fmt.Printf("Warning: Failed to list notebooks: %v\n", err)
	}
	for _, notebook := range notebooks {
		printer.AddResource(display.ResourceInfo{
			ResourceType:  "Notebook",
			Name:         notebook.Name,
			Status:       notebook.Status,
			InstanceType: notebook.InstanceType,
			RunningTime:  time.Since(notebook.CreationTime).String(),
		})
	}

	// Get Studio apps
	apps, err := client.ListStudioApps(ctx)
	if err != nil {
		fmt.Printf("Warning: Failed to list Studio apps: %v\n", err)
	}
	for _, app := range apps {	
		printer.AddResource(display.ResourceInfo{
			ResourceType:  "Studio",
			Name:         fmt.Sprintf("%s/%s", app.UserProfile, app.AppType),
			Status:       app.Status,
			InstanceType: app.InstanceType,
			RunningTime:  time.Since(app.CreationTime).String(),
		})
	}

	// Print results
	printer.Print()
	return nil
}
