package sagemaker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/aws/smithy-go"
	"mohua/internal/retry"
)

// SageMakerClientInterface defines the methods used by Client
type SageMakerClientInterface interface {
	ListApps(ctx context.Context, params *sagemaker.ListAppsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListAppsOutput, error)
	ListEndpoints(ctx context.Context, params *sagemaker.ListEndpointsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListEndpointsOutput, error)
	ListNotebookInstances(ctx context.Context, params *sagemaker.ListNotebookInstancesInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListNotebookInstancesOutput, error)
	ListDomains(ctx context.Context, params *sagemaker.ListDomainsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListDomainsOutput, error)
}

// Client implements only the necessary SageMaker API operations
type Client struct {
	client SageMakerClientInterface
}

// NewClient creates a new SageMaker client
func NewClient(region string) (*Client, error) {
	var opts []func(*config.LoadOptions) error
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}
	
	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &Client{
		client: sagemaker.NewFromConfig(cfg),
	}, nil
}

// ValidateConfiguration checks if the AWS configuration is valid and resources are likely to exist
func (c *Client) ValidateConfiguration(ctx context.Context) (bool, error) {
	// Check if we can list domains as a lightweight way to validate configuration
	input := &sagemaker.ListDomainsInput{
		MaxResults: aws.Int32(1), // We only need to check if we can list
	}

	_, err := c.client.ListDomains(ctx, input)
	if err != nil {
		// If it's an authorization or configuration error, return false
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "AccessDeniedException", 
				 "InvalidClientTokenId", 
				 "SignatureDoesNotMatch", 
				 "ExpiredToken":
				return false, nil
			}
		}
		// For other errors, return the error
		return false, err
	}

	return true, nil
}

// ListEndpoints returns only active endpoints
func (c *Client) ListEndpoints(ctx context.Context) ([]ResourceInfo, error) {
	var resources []ResourceInfo
	
	retrier := retry.NewRetrier(retry.DefaultConfig)
	err := retrier.Do(ctx, func() error {
		input := &sagemaker.ListEndpointsInput{}
		output, err := c.client.ListEndpoints(ctx, input)
		if err != nil {
			return WrapError(err)
		}

		resources = make([]ResourceInfo, 0, len(output.Endpoints))
		for _, endpoint := range output.Endpoints {
			if endpoint.EndpointStatus == types.EndpointStatusInService {
				// we'll skip detailed endpoint config
				resources = append(resources, ResourceInfo{
					Name:         *endpoint.EndpointName,
					Status:       string(endpoint.EndpointStatus),
					InstanceType: "unknown", // Simplified version doesn't fetch detailed config
					InstanceCount: 1,        // Default to 1 for simplified version
					CreationTime: *endpoint.CreationTime,
				})
			}
		}

		return nil
	})

	return resources, err
}

// ListNotebooks returns only running notebook instances
func (c *Client) ListNotebooks(ctx context.Context) ([]ResourceInfo, error) {
	var resources []ResourceInfo

	retrier := retry.NewRetrier(retry.DefaultConfig)
	err := retrier.Do(ctx, func() error {
		input := &sagemaker.ListNotebookInstancesInput{}
		output, err := c.client.ListNotebookInstances(ctx, input)
		if err != nil {
			return WrapError(err)
		}

		resources = make([]ResourceInfo, 0, len(output.NotebookInstances))
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

		return nil
	})

	return resources, err
}

// ListStudioApps returns only running studio applications
func (c *Client) ListStudioApps(ctx context.Context) ([]ResourceInfo, error) {
	var resources []ResourceInfo

	retrier := retry.NewRetrier(retry.DefaultConfig)
	err := retrier.Do(ctx, func() error {
		input := &sagemaker.ListAppsInput{}
		output, err := c.client.ListApps(ctx, input)
		if err != nil {
			return WrapError(err)
		}

		resources = make([]ResourceInfo, 0, len(output.Apps))
		for _, app := range output.Apps {
			if app.Status == types.AppStatusInService {
				// Defensive nil checks
				var name, userProfile, appType, instanceType string
				var creationTime time.Time

				if app.AppName != nil {
					name = *app.AppName
				}

				if app.UserProfileName != nil {
					userProfile = *app.UserProfileName
				}

				if app.CreationTime != nil {
					creationTime = *app.CreationTime
				}

				// Handle potential nil ResourceSpec
				if app.ResourceSpec != nil {
					instanceType = string(app.ResourceSpec.InstanceType)
				}

				appType = string(app.AppType)

				// Only add resource if we have a meaningful name
				if name != "" {
					resources = append(resources, ResourceInfo{
						Name:         name,
						Status:       string(app.Status),
						InstanceType: instanceType,
						CreationTime: creationTime,
						UserProfile:  userProfile,
						AppType:      appType,
					})
				}
			}
		}

		return nil
	})

	return resources, err
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
