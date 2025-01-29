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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mohua",
	Short: "Monitor AWS SageMaker compute resources and their costs",
	Long: `A monitoring tool for AWS SageMaker that helps track running compute resources
and their associated costs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMonitor()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region (optional, defaults to AWS_REGION env var or us-east-1)")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	
	return rootCmd.Execute()
}

func runMonitor() error {
	// Create SageMaker client
	client, err := sagemaker.NewClient(region)
	if err != nil {
		return fmt.Errorf("failed to create SageMaker client: %w", err)
	}

	ctx := context.Background()

	// Track if any resources were found
	resourceFound := false
	var printer *display.Printer

	// Get endpoints
	endpoints, err := client.ListEndpoints(ctx)
	if err != nil {
		return fmt.Errorf("failed to list SageMaker resources: %w", err)
	}
	if len(endpoints) > 0 {
		// Create printer only when first resource is found
		printer = display.NewPrinter(jsonOutput)
		printer.PrintHeader()
		resourceFound = true

		for _, endpoint := range endpoints {
			printer.PrintResource(display.ResourceInfo{
				ResourceType:  "Endpoint",
				Name:         endpoint.Name,
				Status:       endpoint.Status,
				InstanceType: endpoint.InstanceType,
				RunningTime:  time.Since(endpoint.CreationTime).String(),
			})
		}
	}

	// Get notebooks
	notebooks, err := client.ListNotebooks(ctx)
	if err != nil {
		return fmt.Errorf("failed to list SageMaker resources: %w", err)
	}
	if len(notebooks) > 0 {
		// Create printer only when first resource is found
		if !resourceFound {
			printer = display.NewPrinter(jsonOutput)
			printer.PrintHeader()
			resourceFound = true
		}

		for _, notebook := range notebooks {
			printer.PrintResource(display.ResourceInfo{
				ResourceType:  "Notebook",
				Name:         notebook.Name,
				Status:       notebook.Status,
				InstanceType: notebook.InstanceType,
				RunningTime:  time.Since(notebook.CreationTime).String(),
			})
		}
	}

	// Get Studio apps
	apps, err := client.ListStudioApps(ctx)
	if err != nil {
		return fmt.Errorf("failed to list SageMaker resources: %w", err)
	}
	if len(apps) > 0 {
		// Create printer only when first resource is found
		if !resourceFound {
			printer = display.NewPrinter(jsonOutput)
			printer.PrintHeader()
			resourceFound = true
		}

		for _, app := range apps {	
			printer.PrintResource(display.ResourceInfo{
				ResourceType:  "Studio",
				Name:         fmt.Sprintf("%s/%s", app.UserProfile, app.AppType),
				Status:       app.Status,
				InstanceType: app.InstanceType,
				RunningTime:  time.Since(app.CreationTime).String(),
			})
		}
	}

	// If no resources found, return an error
	if !resourceFound {
		return fmt.Errorf("no SageMaker resources found")
	}

	// Print footer
	printer.PrintFooter()
	return nil
}
