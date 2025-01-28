package sagemaker

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
)

// Client represents a SageMaker client with methods to fetch resource information
type Client struct {
	sagemaker *sagemaker.Client
}

// NewClient creates a new SageMaker client for the specified region
func NewClient(region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &Client{
		sagemaker: sagemaker.NewFromConfig(cfg),
	}, nil
}

// EndpointInfo contains information about a SageMaker endpoint
type EndpointInfo struct {
	Name         string
	InstanceType string
	InstanceCount int
	Status       string
	CreationTime time.Time
}

// NotebookInfo contains information about a SageMaker notebook instance
type NotebookInfo struct {
	Name         string
	InstanceType string
	Status       string
	VolumeSize   int
	CreationTime time.Time
}

// StudioInfo contains information about a SageMaker Studio instance
type StudioInfo struct {
	DomainID     string
	UserProfile  string
	AppType      string
	InstanceType string
	Status       string
	CreationTime time.Time
}

// ListEndpoints returns information about all SageMaker endpoints
func (c *Client) ListEndpoints(ctx context.Context) ([]EndpointInfo, error) {
	var endpoints []EndpointInfo
	paginator := sagemaker.NewListEndpointsPaginator(c.sagemaker, &sagemaker.ListEndpointsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list endpoints: %w", err)
		}

		for _, endpoint := range output.Endpoints {
			desc, err := c.sagemaker.DescribeEndpoint(ctx, &sagemaker.DescribeEndpointInput{
				EndpointName: endpoint.EndpointName,
			})
			if err != nil {
				continue
			}

			config, err := c.sagemaker.DescribeEndpointConfig(ctx, &sagemaker.DescribeEndpointConfigInput{
				EndpointConfigName: desc.EndpointConfigName,
			})
			if err != nil {
				continue
			}

			if len(config.ProductionVariants) > 0 {
				endpoints = append(endpoints, EndpointInfo{
					Name:          *endpoint.EndpointName,
					InstanceType:  string(config.ProductionVariants[0].InstanceType),
					InstanceCount: int(*config.ProductionVariants[0].InitialInstanceCount),
					Status:        string(desc.EndpointStatus),
					CreationTime:  *endpoint.CreationTime,
				})
			}
		}
	}

	return endpoints, nil
}

// ListNotebooks returns information about all SageMaker notebook instances
func (c *Client) ListNotebooks(ctx context.Context) ([]NotebookInfo, error) {
	var notebooks []NotebookInfo
	paginator := sagemaker.NewListNotebookInstancesPaginator(c.sagemaker, &sagemaker.ListNotebookInstancesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list notebooks: %w", err)
		}

		for _, notebook := range output.NotebookInstances {
			notebooks = append(notebooks, NotebookInfo{
				Name:         *notebook.NotebookInstanceName,
				InstanceType: string(notebook.InstanceType),
				Status:      string(notebook.NotebookInstanceStatus),
				VolumeSize:  20, // デフォルト値として20GBを設定
				CreationTime: *notebook.CreationTime,
			})
		}
	}

	return notebooks, nil
}

// ListStudioApps returns information about all SageMaker Studio applications
func (c *Client) ListStudioApps(ctx context.Context) ([]StudioInfo, error) {
	var apps []StudioInfo
	paginator := sagemaker.NewListDomainsPaginator(c.sagemaker, &sagemaker.ListDomainsInput{})

	for paginator.HasMorePages() {
		domains, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list domains: %w", err)
		}

		for _, domain := range domains.Domains {
			userPaginator := sagemaker.NewListUserProfilesPaginator(c.sagemaker, &sagemaker.ListUserProfilesInput{
				DomainIdEquals: domain.DomainId,
			})

			for userPaginator.HasMorePages() {
				users, err := userPaginator.NextPage(ctx)
				if err != nil {
					continue
				}

				for _, user := range users.UserProfiles {
					appPaginator := sagemaker.NewListAppsPaginator(c.sagemaker, &sagemaker.ListAppsInput{
						DomainIdEquals:        domain.DomainId,
						UserProfileNameEquals: user.UserProfileName,
					})

					for appPaginator.HasMorePages() {
						appOutput, err := appPaginator.NextPage(ctx)
						if err != nil {
							continue
						}

						for _, app := range appOutput.Apps {
							var instanceType string
							if app.ResourceSpec != nil {
								instanceType = string(app.ResourceSpec.InstanceType)
							}

							apps = append(apps, StudioInfo{
								DomainID:     *domain.DomainId,
								UserProfile:  *user.UserProfileName,
								AppType:      string(app.AppType),
								InstanceType: instanceType,
								Status:       string(app.Status),
								CreationTime: *app.CreationTime,
							})
						}
					}
				}
			}
		}
	}

	return apps, nil
}
