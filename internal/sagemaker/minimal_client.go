package sagemaker

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

// MinimalClient implements only the necessary SageMaker API operations
type MinimalClient struct {
	client *sagemaker.Client
}

// NewMinimalClient creates a new minimal SageMaker client
func NewMinimalClient(region string) (*MinimalClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &MinimalClient{
		client: sagemaker.NewFromConfig(cfg),
	}, nil
}

// ListEndpoints returns only active endpoints
func (c *MinimalClient) ListEndpoints(ctx context.Context) ([]ResourceInfo, error) {
	var resources []ResourceInfo
	
	input := &sagemaker.ListEndpointsInput{}
	output, err := c.client.ListEndpoints(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, endpoint := range output.Endpoints {
		if endpoint.EndpointStatus == types.EndpointStatusInService {
			// For minimal version, we'll skip detailed endpoint config
			resources = append(resources, ResourceInfo{
				Name:         *endpoint.EndpointName,
				Status:       string(endpoint.EndpointStatus),
				InstanceType: "unknown", // Simplified version doesn't fetch detailed config
				InstanceCount: 1,        // Default to 1 for simplified version
				CreationTime: *endpoint.CreationTime,
			})
		}
	}

	return resources, nil
}

// ListNotebooks returns only running notebook instances
func (c *MinimalClient) ListNotebooks(ctx context.Context) ([]ResourceInfo, error) {
	var resources []ResourceInfo

	input := &sagemaker.ListNotebookInstancesInput{}
	output, err := c.client.ListNotebookInstances(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, notebook := range output.NotebookInstances {
		if notebook.NotebookInstanceStatus == types.NotebookInstanceStatusInService {
			resources = append(resources, ResourceInfo{
				Name:         *notebook.NotebookInstanceName,
				Status:       string(notebook.NotebookInstanceStatus),
				InstanceType: string(notebook.InstanceType),
				CreationTime: *notebook.CreationTime,
				VolumeSize:   0, // Simplified version doesn't fetch volume size
			})
		}
	}

	return resources, nil
}

// ListStudioApps returns only running studio applications
func (c *MinimalClient) ListStudioApps(ctx context.Context) ([]ResourceInfo, error) {
	var resources []ResourceInfo

	input := &sagemaker.ListAppsInput{}
	output, err := c.client.ListApps(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, app := range output.Apps {
		if app.Status == types.AppStatusInService {
			resources = append(resources, ResourceInfo{
				Name:         *app.AppName,
				Status:       string(app.Status),
				InstanceType: string(app.ResourceSpec.InstanceType),
				CreationTime: *app.CreationTime,
				UserProfile:  *app.UserProfileName,
				AppType:     string(app.AppType),
			})
		}
	}

	return resources, nil
}

// ResourceInfo contains common fields for SageMaker resources
type ResourceInfo struct {
	Name          string
	Status        string
	InstanceType  string
	InstanceCount int
	CreationTime  time.Time
	VolumeSize    int
	UserProfile   string
	AppType       string
}
