package sagemaker

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Test with explicit region
	client, err := NewClient("us-west-2")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Test with empty region (should use AWS SDK's default region resolution)
	client, err = NewClient("")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Test with AWS_REGION environment variable
	os.Setenv("AWS_REGION", "us-east-1")
	client, err = NewClient("")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	os.Unsetenv("AWS_REGION")
}

// MockSageMakerClient provides a mock implementation of the SageMaker client
type MockSageMakerClient struct {
	listAppsFunc         func(ctx context.Context, params *sagemaker.ListAppsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListAppsOutput, error)
	listEndpointsFunc    func(ctx context.Context, params *sagemaker.ListEndpointsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListEndpointsOutput, error)
	listNotebookFunc     func(ctx context.Context, params *sagemaker.ListNotebookInstancesInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListNotebookInstancesOutput, error)
	listDomainsFunc      func(ctx context.Context, params *sagemaker.ListDomainsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListDomainsOutput, error)
}

func (m *MockSageMakerClient) ListApps(ctx context.Context, params *sagemaker.ListAppsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListAppsOutput, error) {
	if m.listAppsFunc != nil {
		return m.listAppsFunc(ctx, params, optFns...)
	}
	return &sagemaker.ListAppsOutput{}, nil
}

func (m *MockSageMakerClient) ListEndpoints(ctx context.Context, params *sagemaker.ListEndpointsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListEndpointsOutput, error) {
	if m.listEndpointsFunc != nil {
		return m.listEndpointsFunc(ctx, params, optFns...)
	}
	return &sagemaker.ListEndpointsOutput{}, nil
}

func (m *MockSageMakerClient) ListNotebookInstances(ctx context.Context, params *sagemaker.ListNotebookInstancesInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListNotebookInstancesOutput, error) {
	if m.listNotebookFunc != nil {
		return m.listNotebookFunc(ctx, params, optFns...)
	}
	return &sagemaker.ListNotebookInstancesOutput{}, nil
}

func (m *MockSageMakerClient) ListDomains(ctx context.Context, params *sagemaker.ListDomainsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListDomainsOutput, error) {
	if m.listDomainsFunc != nil {
		return m.listDomainsFunc(ctx, params, optFns...)
	}
	return &sagemaker.ListDomainsOutput{}, nil
}

func TestListStudioApps_NilFields(t *testing.T) {
	// Prepare a context
	ctx := context.Background()

	// Create a mock client
	mockClient := &MockSageMakerClient{
		listAppsFunc: func(ctx context.Context, params *sagemaker.ListAppsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListAppsOutput, error) {
			return &sagemaker.ListAppsOutput{
				Apps: []types.AppDetails{
					{
						// Intentionally leave some fields nil
						Status:     types.AppStatusInService,
						AppType:    types.AppTypeJupyterServer,
						AppName:    nil,
						CreationTime: nil,
						UserProfileName: nil,
						ResourceSpec: nil,
					},
					{
						// Another app with some fields populated
						Status:     types.AppStatusInService,
						AppType:    types.AppTypeJupyterServer,
						AppName:    aws.String("TestApp"),
						CreationTime: aws.Time(time.Now()),
						UserProfileName: aws.String("TestUser"),
						ResourceSpec: &types.ResourceSpec{
							InstanceType: types.AppInstanceType("ml.t3.medium"),
						},
					},
				},
			}, nil
		},
	}

	// Create a Client with the mock
	client := &Client{
		client: mockClient,
	}

	// Call the method
	resources, err := client.ListStudioApps(ctx)

	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, resources, 1, "Should only include apps with non-nil names")
	
	// Verify the populated app's details
	if len(resources) > 0 {
		assert.Equal(t, "TestApp", resources[0].Name)
		assert.Equal(t, "TestUser", resources[0].UserProfile)
		assert.Equal(t, "ml.t3.medium", resources[0].InstanceType)
	}
}

func TestConcurrentResourceListing(t *testing.T) {
	// Prepare a context
	ctx := context.Background()

	// Create a mock client with simulated delays
	mockClient := &MockSageMakerClient{
		listEndpointsFunc: func(ctx context.Context, params *sagemaker.ListEndpointsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListEndpointsOutput, error) {
			time.Sleep(100 * time.Millisecond) // Simulate some delay
			return &sagemaker.ListEndpointsOutput{
				Endpoints: []types.EndpointSummary{
					{
						EndpointName:     aws.String("Endpoint1"),
						EndpointStatus:   types.EndpointStatusInService,
						CreationTime:     aws.Time(time.Now().Add(-1 * time.Hour)),
					},
				},
			}, nil
		},
		listNotebookFunc: func(ctx context.Context, params *sagemaker.ListNotebookInstancesInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListNotebookInstancesOutput, error) {
			time.Sleep(50 * time.Millisecond) // Simulate some delay
			return &sagemaker.ListNotebookInstancesOutput{
				NotebookInstances: []types.NotebookInstanceSummary{
					{
						NotebookInstanceName:     aws.String("Notebook1"),
						NotebookInstanceStatus:   types.NotebookInstanceStatusInService,
						CreationTime:             aws.Time(time.Now().Add(-2 * time.Hour)),
					},
				},
			}, nil
		},
		listAppsFunc: func(ctx context.Context, params *sagemaker.ListAppsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListAppsOutput, error) {
			time.Sleep(75 * time.Millisecond) // Simulate some delay
			return &sagemaker.ListAppsOutput{
				Apps: []types.AppDetails{
					{
						AppName:     aws.String("App1"),
						Status:      types.AppStatusInService,
						CreationTime: aws.Time(time.Now().Add(-3 * time.Hour)),
					},
				},
			}, nil
		},
	}

	// Create a Client with the mock
	client := &Client{
		client: mockClient,
	}

	// Measure total time for concurrent calls
	startTime := time.Now()
	
	// Perform concurrent resource listing
	var wg sync.WaitGroup
	wg.Add(3)

	var endpointResults, notebookResults, appResults []ResourceInfo
	var endpointErr, notebookErr, appErr error

	go func() {
		defer wg.Done()
		endpointResults, endpointErr = client.ListEndpoints(ctx)
	}()

	go func() {
		defer wg.Done()
		notebookResults, notebookErr = client.ListNotebooks(ctx)
	}()

	go func() {
		defer wg.Done()
		appResults, appErr = client.ListStudioApps(ctx)
	}()

	wg.Wait()

	// Calculate total time
	totalTime := time.Since(startTime)

	// Assert no errors
	assert.NoError(t, endpointErr)
	assert.NoError(t, notebookErr)
	assert.NoError(t, appErr)

	// Assert results
	assert.Len(t, endpointResults, 1)
	assert.Len(t, notebookResults, 1)
	assert.Len(t, appResults, 1)

	// Total time should be less than sequential calls (sum of delays)
	// Allowing some buffer for goroutine overhead
	assert.Less(t, totalTime.Milliseconds(), int64(250), "Concurrent calls should be faster than sequential")
}

func TestValidateConfiguration(t *testing.T) {
	ctx := context.Background()

	// Test case 1: Successful configuration
	mockClientSuccess := &MockSageMakerClient{
		listDomainsFunc: func(ctx context.Context, params *sagemaker.ListDomainsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListDomainsOutput, error) {
			return &sagemaker.ListDomainsOutput{}, nil
		},
	}
	clientSuccess := &Client{client: mockClientSuccess}
	hasResources, err := clientSuccess.ValidateConfiguration(ctx)
	assert.NoError(t, err)
	assert.True(t, hasResources)

	// Test case 2: Access Denied
	mockClientAccessDenied := &MockSageMakerClient{
		listDomainsFunc: func(ctx context.Context, params *sagemaker.ListDomainsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListDomainsOutput, error) {
			return nil, &smithy.GenericAPIError{
				Code:    "AccessDeniedException",
				Message: "Access Denied",
			}
		},
	}
	clientAccessDenied := &Client{client: mockClientAccessDenied}
	hasResources, err = clientAccessDenied.ValidateConfiguration(ctx)
	assert.NoError(t, err)
	assert.False(t, hasResources)

	// Test case 3: Invalid Token
	mockClientInvalidToken := &MockSageMakerClient{
		listDomainsFunc: func(ctx context.Context, params *sagemaker.ListDomainsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListDomainsOutput, error) {
			return nil, &smithy.GenericAPIError{
				Code:    "InvalidClientTokenId",
				Message: "Invalid Token",
			}
		},
	}
	clientInvalidToken := &Client{client: mockClientInvalidToken}
	hasResources, err = clientInvalidToken.ValidateConfiguration(ctx)
	assert.NoError(t, err)
	assert.False(t, hasResources)
}
