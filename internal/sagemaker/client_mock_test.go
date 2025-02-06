package sagemaker

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSageMakerClient is a mock implementation of the SageMakerClientInterface
type MockSageMakerClient struct {
	mock.Mock
}

func (m *MockSageMakerClient) ListApps(ctx context.Context, params *sagemaker.ListAppsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListAppsOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sagemaker.ListAppsOutput), args.Error(1)
}

func (m *MockSageMakerClient) ListEndpoints(ctx context.Context, params *sagemaker.ListEndpointsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListEndpointsOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sagemaker.ListEndpointsOutput), args.Error(1)
}

func (m *MockSageMakerClient) ListNotebookInstances(ctx context.Context, params *sagemaker.ListNotebookInstancesInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListNotebookInstancesOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sagemaker.ListNotebookInstancesOutput), args.Error(1)
}

func (m *MockSageMakerClient) ListDomains(ctx context.Context, params *sagemaker.ListDomainsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListDomainsOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sagemaker.ListDomainsOutput), args.Error(1)
}

// TestMockSageMakerClientBasic verifies that the mock client implements the interface correctly
func TestMockSageMakerClientBasic(t *testing.T) {
	mockClient := new(MockSageMakerClient)
	
	// Test that mockClient implements SageMakerClientInterface
	var _ SageMakerClientInterface = mockClient

	// Setup test data
	ctx := context.Background()
	now := time.Now()

	// Test ListDomains
	mockClient.On("ListDomains", ctx, &sagemaker.ListDomainsInput{MaxResults: aws.Int32(1)}, mock.Anything).
		Return(&sagemaker.ListDomainsOutput{}, nil)

	// Test ListApps
	mockClient.On("ListApps", ctx, &sagemaker.ListAppsInput{}, mock.Anything).
		Return(&sagemaker.ListAppsOutput{
			Apps: []types.AppDetails{
				{
					AppName:         aws.String("TestApp"),
					Status:         types.AppStatusInService,
					CreationTime:   aws.Time(now),
					UserProfileName: aws.String("TestUser"),
					AppType:        types.AppTypeJupyterServer,
				},
			},
		}, nil)

	// Execute tests
	client := &clientImpl{client: mockClient}

	// Test ValidateConfiguration
	hasResources, err := client.ValidateConfiguration(ctx)
	assert.NoError(t, err)
	assert.True(t, hasResources)

	// Test ListStudioApps
	apps, err := client.ListStudioApps(ctx)
	assert.NoError(t, err)
	assert.Len(t, apps, 1)
	assert.Equal(t, "TestApp", apps[0].Name)
	assert.Equal(t, "TestUser", apps[0].UserProfile)
	assert.Equal(t, "JupyterServer", apps[0].AppType)

	// Verify all expectations were met
	mockClient.AssertExpectations(t)
}
