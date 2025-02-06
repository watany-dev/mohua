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
	"github.com/stretchr/testify/mock"
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

func TestGetRegion(t *testing.T) {
	// Test with explicit region
	client, err := NewClient("us-west-2")
	assert.NoError(t, err)
	assert.Equal(t, "us-west-2", client.GetRegion())

	// Test with AWS_REGION environment variable
	os.Setenv("AWS_REGION", "us-east-1")
	client, err = NewClient("")
	assert.NoError(t, err)
	assert.Equal(t, "us-east-1", client.GetRegion())
	os.Unsetenv("AWS_REGION")
}

func TestListStudioApps_NilFields(t *testing.T) {
	// Prepare a context
	ctx := context.Background()

	// Create a mock client
	mockClient := new(MockSageMakerClient)
	now := time.Now()

	// Setup mock expectations
	mockClient.On("ListApps", ctx, &sagemaker.ListAppsInput{}, mock.Anything).
		Return(&sagemaker.ListAppsOutput{
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
					// Old Studio app with some fields populated
					Status:     types.AppStatusInService,
					AppType:    types.AppTypeJupyterServer,
					AppName:    aws.String("TestApp"),
					CreationTime: aws.Time(now),
					UserProfileName: aws.String("TestUser"),
					ResourceSpec: &types.ResourceSpec{
						InstanceType: types.AppInstanceType("ml.t3.medium"),
					},
				},
				{
					// New Studio app with space name
					Status:     types.AppStatusInService,
					AppType:    types.AppTypeJupyterLab,
					AppName:    aws.String("NewTestApp"),
					CreationTime: aws.Time(now),
					UserProfileName: aws.String("NewTestUser"),
					SpaceName:   aws.String("TestSpace"),
					ResourceSpec: &types.ResourceSpec{
						InstanceType: types.AppInstanceType("ml.t3.large"),
					},
				},
			},
		}, nil)

	// Create a Client with the mock
	client := &clientImpl{
		client: mockClient,
	}

	// Call the method
	resources, err := client.ListStudioApps(ctx)

	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, resources, 2, "Should include apps with non-nil names")
	
	// Verify the old Studio app details
	oldStudioApp := resources[0]
	assert.Equal(t, "TestApp", oldStudioApp.Name)
	assert.Equal(t, "TestUser", oldStudioApp.UserProfile)
	assert.Equal(t, "ml.t3.medium", oldStudioApp.InstanceType)
	assert.Equal(t, "JupyterServer", oldStudioApp.AppType)
	assert.Equal(t, "Old Studio (JupyterServer)", oldStudioApp.StudioType)
	assert.Empty(t, oldStudioApp.SpaceName)

	// Verify the new Studio app details
	newStudioApp := resources[1]
	assert.Equal(t, "NewTestApp", newStudioApp.Name)
	assert.Equal(t, "NewTestUser", newStudioApp.UserProfile)
	assert.Equal(t, "ml.t3.large", newStudioApp.InstanceType)
	assert.Equal(t, "JupyterLab", newStudioApp.AppType)
	assert.Equal(t, "New Studio (JupyterLab)", newStudioApp.StudioType)
	assert.Equal(t, "TestSpace", newStudioApp.SpaceName)
}

func TestListStudioApps_StatusHandling(t *testing.T) {
	// Prepare a context
	ctx := context.Background()

	// Create a mock client with mixed statuses
	mockClient := new(MockSageMakerClient)
	now := time.Now()

	// Setup mock expectations
	mockClient.On("ListApps", ctx, &sagemaker.ListAppsInput{}, mock.Anything).
		Return(&sagemaker.ListAppsOutput{
			Apps: []types.AppDetails{
				{
					// Running old Studio app
					Status:     types.AppStatusInService,
					AppType:    types.AppTypeJupyterServer,
					AppName:    aws.String("RunningOldApp"),
					CreationTime: aws.Time(now),
					UserProfileName: aws.String("OldUser"),
				},
				{
					// Stopped new Studio app
					Status:     types.AppStatusDeleted,
					AppType:    types.AppTypeJupyterLab,
					AppName:    aws.String("StoppedNewApp"),
					CreationTime: aws.Time(now.Add(-1 * time.Hour)),
					UserProfileName: aws.String("NewUser"),
					SpaceName:   aws.String("StoppedSpace"),
				},
			},
		}, nil)

	// Create a Client with the mock
	client := &clientImpl{
		client: mockClient,
	}

	// Call the method
	resources, err := client.ListStudioApps(ctx)

	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, resources, 1, "Should only include InService apps")
	
	// Verify the running old Studio app details
	runningOldApp := resources[0]
	assert.Equal(t, "RunningOldApp", runningOldApp.Name)
	assert.Equal(t, "InService", runningOldApp.Status)
	assert.Equal(t, "Old Studio (JupyterServer)", runningOldApp.StudioType)
}

func TestConcurrentResourceListing(t *testing.T) {
	// Prepare a context
	ctx := context.Background()

	// Create a mock client with simulated delays
	mockClient := new(MockSageMakerClient)
	now := time.Now()

	// Setup mock expectations with delays
	mockClient.On("ListEndpoints", ctx, &sagemaker.ListEndpointsInput{}, mock.Anything).
		Run(func(args mock.Arguments) {
			time.Sleep(100 * time.Millisecond) // Simulate some delay
		}).
		Return(&sagemaker.ListEndpointsOutput{
			Endpoints: []types.EndpointSummary{
				{
					EndpointName:   aws.String("Endpoint1"),
					EndpointStatus: types.EndpointStatusInService,
					CreationTime:   aws.Time(now.Add(-1 * time.Hour)),
				},
			},
		}, nil)

	mockClient.On("ListNotebookInstances", ctx, &sagemaker.ListNotebookInstancesInput{}, mock.Anything).
		Run(func(args mock.Arguments) {
			time.Sleep(50 * time.Millisecond) // Simulate some delay
		}).
		Return(&sagemaker.ListNotebookInstancesOutput{
			NotebookInstances: []types.NotebookInstanceSummary{
				{
					NotebookInstanceName:   aws.String("Notebook1"),
					NotebookInstanceStatus: types.NotebookInstanceStatusInService,
					CreationTime:           aws.Time(now.Add(-2 * time.Hour)),
				},
			},
		}, nil)

	mockClient.On("ListApps", ctx, &sagemaker.ListAppsInput{}, mock.Anything).
		Run(func(args mock.Arguments) {
			time.Sleep(75 * time.Millisecond) // Simulate some delay
		}).
		Return(&sagemaker.ListAppsOutput{
			Apps: []types.AppDetails{
				{
					AppName:      aws.String("App1"),
					Status:       types.AppStatusInService,
					CreationTime: aws.Time(now.Add(-3 * time.Hour)),
				},
			},
		}, nil)

	// Create a Client with the mock
	client := &clientImpl{
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
	mockClientSuccess := new(MockSageMakerClient)
	mockClientSuccess.On("ListDomains", ctx, &sagemaker.ListDomainsInput{MaxResults: aws.Int32(1)}, mock.Anything).
		Return(&sagemaker.ListDomainsOutput{}, nil)
	clientSuccess := &clientImpl{client: mockClientSuccess}
	hasResources, err := clientSuccess.ValidateConfiguration(ctx)
	assert.NoError(t, err)
	assert.True(t, hasResources)

	// Test case 2: Access Denied
	mockClientAccessDenied := new(MockSageMakerClient)
	mockClientAccessDenied.On("ListDomains", ctx, &sagemaker.ListDomainsInput{MaxResults: aws.Int32(1)}, mock.Anything).
		Return(nil, &smithy.GenericAPIError{
			Code:    "AccessDeniedException",
			Message: "Access Denied",
		})
	clientAccessDenied := &clientImpl{client: mockClientAccessDenied}
	hasResources, err = clientAccessDenied.ValidateConfiguration(ctx)
	assert.NoError(t, err)
	assert.False(t, hasResources)

	// Test case 3: Invalid Token
	mockClientInvalidToken := new(MockSageMakerClient)
	mockClientInvalidToken.On("ListDomains", ctx, &sagemaker.ListDomainsInput{MaxResults: aws.Int32(1)}, mock.Anything).
		Return(nil, &smithy.GenericAPIError{
			Code:    "InvalidClientTokenId",
			Message: "Invalid Token",
		})
	clientInvalidToken := &clientImpl{client: mockClientInvalidToken}
	hasResources, err = clientInvalidToken.ValidateConfiguration(ctx)
	assert.NoError(t, err)
	assert.False(t, hasResources)
}
