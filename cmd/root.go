package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
	"github.com/spf13/cobra"
	"mohua/internal/display"
	"mohua/internal/sagemaker"
)

type NoResourcesError struct {
	Region      string
	DefaultUsed bool
}

func (e *NoResourcesError) Error() string {
	if e.DefaultUsed {
		return fmt.Sprintf("No SageMaker resources found in default region %s.", e.Region)
	}
	return fmt.Sprintf("No SageMaker resources found in specified region %s.", e.Region)
}

// ResourceResult holds the results and errors from API calls
type ResourceResult struct {
	Resources []sagemaker.ResourceInfo
	Error     error
}

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
		cmd.SilenceUsage = true
		return runMonitor()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region (optional, defaults to AWS CLI configuration)")
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

	// Validate AWS configuration and check if resources are likely to exist
	hasConfiguredResources, err := client.ValidateConfiguration(ctx)
	if err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	if !hasConfiguredResources {
		return &NoResourcesError{}
	}

	// Create channels for each resource type
	endpointsChan := make(chan ResourceResult, 1)
	notebooksChan := make(chan ResourceResult, 1)
	appsChan := make(chan ResourceResult, 1)

	// Launch goroutines for each API call
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		endpoints, err := client.ListEndpoints(ctx)
		endpointsChan <- ResourceResult{Resources: endpoints, Error: err}
	}()

	go func() {
		defer wg.Done()
		notebooks, err := client.ListNotebooks(ctx)
		notebooksChan <- ResourceResult{Resources: notebooks, Error: err}
	}()

	go func() {
		defer wg.Done()
		apps, err := client.ListStudioApps(ctx)
		appsChan <- ResourceResult{Resources: apps, Error: err}
	}()

	// Close channels after all goroutines complete
	go func() {
		wg.Wait()
		close(endpointsChan)
		close(notebooksChan)
		close(appsChan)
	}()

	// Track if any resources were found and collect errors
	resourceFound := false
	var printer *display.Printer
	var firstError error

	// Process endpoints
	if result := <-endpointsChan; result.Error != nil {
		// Check if the error is retryable
		if retryableErr, ok := result.Error.(*sagemaker.RetryableError); ok {
			// Log the retryable error, but don't stop execution
			fmt.Fprintf(os.Stderr, "Retryable error listing endpoints: %v\n", retryableErr)
		} else {
			if firstError == nil {
				firstError = fmt.Errorf("failed to list endpoints: %w", result.Error)
			}
		}
	} else if len(result.Resources) > 0 {
		if !resourceFound {
			printer = display.NewPrinter(jsonOutput)
			printer.PrintHeader()
			resourceFound = true
		}
		for _, endpoint := range result.Resources {
			printer.PrintResource(display.ResourceInfo{
				ResourceType:  "Endpoint",
				Name:         endpoint.Name,
				Status:       endpoint.Status,
				InstanceType: endpoint.InstanceType,
				RunningTime:  time.Since(endpoint.CreationTime).String(),
			})
		}
	}

	// Process notebooks
	if result := <-notebooksChan; result.Error != nil {
		// Check if the error is retryable
		if retryableErr, ok := result.Error.(*sagemaker.RetryableError); ok {
			// Log the retryable error, but don't stop execution
			fmt.Fprintf(os.Stderr, "Retryable error listing notebooks: %v\n", retryableErr)
		} else {
			if firstError == nil {
				firstError = fmt.Errorf("failed to list notebooks: %w", result.Error)
			}
		}
	} else if len(result.Resources) > 0 {
		if !resourceFound {
			printer = display.NewPrinter(jsonOutput)
			printer.PrintHeader()
			resourceFound = true
		}
		for _, notebook := range result.Resources {
			printer.PrintResource(display.ResourceInfo{
				ResourceType:  "Notebook",
				Name:         notebook.Name,
				Status:       notebook.Status,
				InstanceType: notebook.InstanceType,
				RunningTime:  time.Since(notebook.CreationTime).String(),
			})
		}
	}

	// Process Studio apps
	if result := <-appsChan; result.Error != nil {
		// Check if the error is retryable
		if retryableErr, ok := result.Error.(*sagemaker.RetryableError); ok {
			// Log the retryable error, but don't stop execution
			fmt.Fprintf(os.Stderr, "Retryable error listing studio apps: %v\n", retryableErr)
		} else {
			if firstError == nil {
				firstError = fmt.Errorf("failed to list studio apps: %w", result.Error)
			}
		}
	} else if len(result.Resources) > 0 {
		if !resourceFound {
			printer = display.NewPrinter(jsonOutput)
			printer.PrintHeader()
			resourceFound = true
		}
		for _, app := range result.Resources {
			printer.PrintResource(display.ResourceInfo{
				ResourceType:  "Studio",
				Name:         fmt.Sprintf("%s/%s", app.UserProfile, app.AppType),
				Status:       app.Status,
				InstanceType: app.InstanceType,
				RunningTime:  time.Since(app.CreationTime).String(),
			})
		}
	}

	// Return first error encountered if any
	if firstError != nil {
		return firstError
	}

	// If no resources found, return a NoResourcesError
	if !resourceFound {
		// Get the effective region from the client
		effectiveRegion := client.GetRegion()
		if effectiveRegion == "" {
			effectiveRegion = "configured region" // Fallback message
		}
		return &NoResourcesError{Region: effectiveRegion, DefaultUsed: region == ""}
	}

	// Print footer
	printer.PrintFooter()
	return nil
}
