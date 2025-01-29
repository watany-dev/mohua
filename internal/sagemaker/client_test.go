package sagemaker

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
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
